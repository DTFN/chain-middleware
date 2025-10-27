package com.lingshu.server.core.business;

import cn.hutool.core.io.FileUtil;
import com.lingshu.server.common.api.ApiErrorCode;
import com.lingshu.server.common.api.OpenAPIRespBuilder;
import com.lingshu.server.core.enums.NationalityEnum;
import com.lingshu.server.core.web3j.service.DIDManagerService;
import com.lingshu.server.core.web3j.service.ResourceDomainService;
import com.lingshu.server.dto.*;
import com.lingshu.server.dto.resp.account.CreateAccountResp;
import com.lingshu.server.dto.resp.account.DidDocumentResp;
import com.lingshu.server.key.PkeyByMnemonicService;
import com.lingshu.server.utils.*;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.web3j.crypto.Credentials;
import org.web3j.crypto.ECKeyPair;
import org.web3j.crypto.Keys;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.utils.Numeric;

import javax.annotation.Resource;
import java.io.File;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Map;

@Slf4j
@Service
public class AccountService {
    @Autowired
    private ResourceDomainService resourceDomainService;

    @Autowired
    private DIDManagerService didManagerService;
    @Resource
    private PkeyByMnemonicService pkeyByMnemonicService;

    @Value("${file-store-path}")
    private String fileStorePath;

    public CreateAccountResp create(String password, NationalityEnum nationality) throws Exception {
        Map<String, String> keyPairMap = generateEthKeyPairByMnemonic(password);
        String privateKey = keyPairMap.get("privateKey");
        log.info("privateKey: {}", privateKey);
        String publicKey = Ed25519Signer.derivePublicKey(privateKey);
        log.info("publicKey: {}", publicKey);
        String address = EthVCSigner.getAddressFromPrivateKey(privateKey);
        log.info("address: {}", address);
        String did = generateUniqueDID(publicKey);
        log.info("did: {}", did);

        String didDoc = DIDProcessor.fillAndStandardizeDIDObj(did, publicKey, address, nationality.getCode());
        log.info("didDoc: {}", didDoc);

        DIDCreateRequest didCreateRequest = new DIDCreateRequest();
        didCreateRequest.setDid(did);
        didCreateRequest.setDidDocument(didDoc);
        didManagerService.createDID(didCreateRequest);

        CreateAccountResp createAccountResp = new CreateAccountResp();
        createAccountResp.setDid(did);
        createAccountResp.setPrivateKey(privateKey);
        createAccountResp.setDidDocument(didDoc);
        createAccountResp.setMnemonic(keyPairMap.get("mnemonic"));

        // 保存账户信息
        storeAccount(did, privateKey, password);

        return createAccountResp;
    }

    private void storeAccount(String did, String privateKeyHex, String passwd) throws Exception {
        PkeyByMnemonicService service = new PkeyByMnemonicService();
        String privateKeyByP12 = service.encryptPrivateKeyByP12(passwd, Numeric.hexStringToByteArray(privateKeyHex));
        FileUtil.writeString(privateKeyByP12, fileStorePath + File.separator + did, StandardCharsets.UTF_8);
    }

    public boolean checkAccount(String did, String passwd, String privateKeyHex) throws Exception {
        PkeyByMnemonicService service = new PkeyByMnemonicService();
        String privateKeyByP12 = FileUtil.readString(fileStorePath + File.separator + did, StandardCharsets.UTF_8);

        // 私钥不是必须的,如果密码错误这一步会报错
        String decryptPrivateKeyByP12 = service.decryptPrivateKeyByP12(passwd, privateKeyByP12);

        return decryptPrivateKeyByP12.equals(privateKeyHex);
    }

    // todo did前缀
    public String generateUniqueDID(String publicKey) {
        return "did:lsid:" + publicKey;
    }

    public Map<String, String> generateEthKeyPair() throws Exception {
        ECKeyPair keyPair = Keys.createEcKeyPair();
        String privateKey = Numeric.toHexStringNoPrefix(keyPair.getPrivateKey());
        String publicKey = Numeric.toHexStringNoPrefix(keyPair.getPublicKey());
        String address = Credentials.create(keyPair).getAddress();

        Map<String, String> keyMap = new HashMap<>();
        keyMap.put("privateKey", privateKey);
        keyMap.put("publicKey", publicKey);
        keyMap.put("address", address);

        log.info("Generated ETH key pair - Address: {}", address);
        return keyMap;
    }

    public Map<String, String> generateEthKeyPairByMnemonic(String password) throws Exception {
        String mnemonic = pkeyByMnemonicService.createMnemonic();
        ECKeyPair keyPair = pkeyByMnemonicService.generatePrivateKeyByMnemonic(mnemonic, password);
        String privateKey = Numeric.toHexStringNoPrefix(keyPair.getPrivateKey());
        String publicKey = Numeric.toHexStringNoPrefix(keyPair.getPublicKey());
        String address = Credentials.create(keyPair).getAddress();

        Map<String, String> keyMap = new HashMap<>();
        keyMap.put("privateKey", privateKey);
        keyMap.put("publicKey", publicKey);
        keyMap.put("address", address);
        keyMap.put("mnemonic", mnemonic);

        log.info("Generated ETH key pair - Address: {}", address);
        return keyMap;
    }

    public Boolean update(UpdateAccountRequest request) {
        DIDCreateRequest didUpdateRequest = new DIDCreateRequest();
        didUpdateRequest.setDid(request.getDid());
        didUpdateRequest.setDidDocument(request.getDidDocument());
        TransactionReceipt receipt = didManagerService.updateDID(didUpdateRequest);
        return true;
    }

    public DidDocumentResp getDidDocument(String did) throws Exception {
        String didDetails = didManagerService.getDIDDetails(did);
        DidDocumentResp resp = new DidDocumentResp();
        resp.setDidDocument(didDetails);
        return resp;
    }

    public String backupAccount(BackupAccountRequest request) throws Exception {
        ECKeyPair ecKeyPair = pkeyByMnemonicService.generatePrivateKeyByMnemonic(request.getMnemonic(), request.getPassword());
        String privateKey = Numeric.toHexStringNoPrefix(ecKeyPair.getPrivateKey());

        // 权限检查
        checkAccount(request.getDid(), request.getPassword(), privateKey);

        String privateKeyByP12 = FileUtil.readString(fileStorePath + File.separator + request.getDid(), StandardCharsets.UTF_8);

        return privateKeyByP12;
    }

    public CreateAccountResp importAccount(RestoreAccountRequest request) throws Exception {
        String privateKey = pkeyByMnemonicService.decryptPrivateKeyByP12(request.getPassword(), request.getEncrypt());

        // 生成did
        String publicKey = Ed25519Signer.derivePublicKey(privateKey);
        String did = generateUniqueDID(publicKey);
        String address = EthVCSigner.getAddressFromPrivateKey(privateKey);
        String didDoc = DIDProcessor.fillAndStandardizeDIDObj(did, publicKey, address, request.getNationality().getCode());
        log.info("didDoc: {}", didDoc);

        DIDCreateRequest didCreateRequest = new DIDCreateRequest();
        didCreateRequest.setDid(did);
        didCreateRequest.setDidDocument(didDoc);
        didManagerService.updateDID(didCreateRequest);

        CreateAccountResp createAccountResp = new CreateAccountResp();
        createAccountResp.setDid(did);
        createAccountResp.setPrivateKey(privateKey);
        createAccountResp.setDidDocument(didDoc);
        createAccountResp.setMnemonic(null);    // 私钥无法推导出助记词

        // 保存账户信息
        storeAccount(did, privateKey, request.getPassword());

        return createAccountResp;
    }
}
