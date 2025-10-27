package com.lingshu.server.core.business;

import cn.hutool.core.date.DateUtil;
import cn.hutool.json.JSONObject;
import cn.hutool.json.JSONUtil;
import com.lingshu.server.common.api.ApiException;
import com.lingshu.server.common.metrics.CrossChainTxMetric;
import com.lingshu.server.core.enums.ContractTypeEnum;
import com.lingshu.server.core.web3j.service.ContractService;
import com.lingshu.server.core.web3j.service.ResourceDomainService;
import com.lingshu.server.dto.*;
import com.lingshu.server.dto.resp.busi.TxHashResp;
import com.lingshu.server.utils.EthVCSigner;
import com.lingshu.server.utils.JsonOrderProcessor;
import com.lingshu.server.utils.Vc0Processor;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.web3j.crypto.Sign;
import org.web3j.utils.Numeric;

import java.util.Objects;
import java.util.Optional;

import static com.lingshu.server.core.enums.ContractTypeEnum.CHAINMAKER_WASM;
import static com.lingshu.server.core.enums.ContractTypeEnum.FABRIC_WASM;

/**
 * @author: derrick
 * @since: 2025-09-08
 */
@Slf4j
@Service
public class BusinessService {
    @Autowired
    private ResourceDomainService resourceDomainService;
    @Autowired
    private ContractService contractService;
    @Autowired
    private AccountService accountService;
    @Autowired
    private CrossChainTxMetric crossChainTxMetric;

    public TxHashResp crossChain(CrossChainRequest request) throws Exception {
        // 校验密钥
        boolean checkResult = accountService.checkAccount(request.getDid(), request.getPasswd(), request.getPrivateKeyHex());
        if (checkResult == false) {
            crossChainTxMetric.incrementSignFail(
                    request.getDid(),
                    Optional.ofNullable(request)
                            .map(CrossChainRequest::getOrigin)
                            .map(CredentialSubject::getResourceName)
                            .orElse(null),
                    Optional.ofNullable(request)
                            .map(CrossChainRequest::getTarget)
                            .map(CredentialSubject::getResourceName)
                            .orElse(null)
                    );
            throw new ApiException("账号密码不匹配");
        }

        String paramsStr = genCrossChainVcListStr(request);
        log.info("vc: {}", paramsStr);
        //调用origin合约方法
        TxHashResp result = contractService.call(paramsStr);
        return result;
    }

    public TxHashResp crossChainDirect(VerifiableCredentialListDto vcl) throws Exception {
        log.info("vc: {}", JsonOrderProcessor.convert(vcl));

        // 验证签名
        try {
            TxHashResp verifyResult = contractService.callVerifyFunc(JsonOrderProcessor.convert(vcl));
        } catch (Exception e) {
            crossChainTxMetric.incrementSignFail(
                    vcl.getOrigin().getIssuer(),
                    Optional.ofNullable(vcl)
                            .map(VerifiableCredentialListDto::getOrigin)
                            .map(VerifiableCredentialDto::getCredentialSubject)
                            .map(VerifiableCredentialDto.CredentialSubject::getResourceName)
                            .orElse(null),
                    Optional.ofNullable(vcl)
                            .map(VerifiableCredentialListDto::getTarget)
                            .map(VerifiableCredentialDto::getCredentialSubject)
                            .map(VerifiableCredentialDto.CredentialSubject::getResourceName)
                            .orElse(null)
            );
            throw e;
        }

        //调用origin合约方法
        TxHashResp result = contractService.call(JsonOrderProcessor.convert(vcl));
        return result;
    }

    public String genCrossChainVcListStr(CrossChainRequest request) throws Exception {
        //根据resource_name,获取ResourceDomain合约上的信息（orign/target）
        ResourceDomainInfo originResourceDomainInfo = getResourceDomainInfo(request.getOrigin().getResourceName());
        ResourceDomainInfo targetResourceDomainInfo = getResourceDomainInfo(request.getTarget().getResourceName());
        //将ResourceDomain合约上的信息，写入origin/target vc中
        VerifiableCredentialDto originObject = dealVc(originResourceDomainInfo, request.getPrivateKeyHex(), request.getOrigin(), request.getDid());
        VerifiableCredentialDto targetObject = dealVc(targetResourceDomainInfo, request.getPrivateKeyHex(), request.getTarget(), request.getDid());

        VerifiableCredentialListDto verifiableCredentialListDto = new VerifiableCredentialListDto();
        verifiableCredentialListDto
                .setOrigin(originObject)
                .setTarget(targetObject);

        String paramsStr = JsonOrderProcessor.convert(verifiableCredentialListDto);
        return paramsStr;
    }

    public VerifiableCredentialListDto genCrossChainVcList(CrossChainRequest request) throws Exception {
        //根据resource_name,获取ResourceDomain合约上的信息（orign/target）
        ResourceDomainInfo originResourceDomainInfo = getResourceDomainInfo(request.getOrigin().getResourceName());
        ResourceDomainInfo targetResourceDomainInfo = getResourceDomainInfo(request.getTarget().getResourceName());
        //将ResourceDomain合约上的信息，写入origin/target vc中
        VerifiableCredentialDto originObject = dealVc(originResourceDomainInfo, request.getPrivateKeyHex(), request.getOrigin(), request.getDid());
        VerifiableCredentialDto targetObject = dealVc(targetResourceDomainInfo, request.getPrivateKeyHex(), request.getTarget(), request.getDid());

        VerifiableCredentialListDto verifiableCredentialListDto = new VerifiableCredentialListDto();
        verifiableCredentialListDto
                .setOrigin(originObject)
                .setTarget(targetObject);

        return verifiableCredentialListDto;
    }

    private VerifiableCredentialDto dealVc(ResourceDomainInfo resourceDomainInfo, String privateKeyHex,
                                           CredentialSubject credentialSubject, String did) throws Exception {
        VerifiableCredentialDto vcObject = Vc0Processor.fillAndStandardizeVc(resourceDomainInfo.getChainRid(),
                resourceDomainInfo.getContractAddress(),
                credentialSubject.getContractFunc(),
                resourceDomainInfo.getContractType(),
                resourceDomainInfo.getGatewayId(),
                resourceDomainInfo.getResourceName(),
                DateUtil.today());

        // 填充issuer
        vcObject.setIssuer(did);

        VerifiableCredentialDto.CredentialSubject credentialSubjectObject = vcObject.getCredentialSubject();
        credentialSubjectObject.setFuncParams(credentialSubject.getFuncParams());

        if (StringUtils.isNotBlank(credentialSubject.getKey())) {
            credentialSubjectObject.setKey(credentialSubject.getKey());
        }
        if (StringUtils.isNotBlank(credentialSubject.getParamName())) {
            credentialSubjectObject.setParamName(credentialSubject.getParamName());
        }

        String contractType = resourceDomainInfo.getContractType();
        ContractTypeEnum contractTypeEnum = ContractTypeEnum.getByCode(contractType);
        if (Objects.isNull(contractTypeEnum)) {
            throw new ApiException("非法合约类型:" + contractType);
        }

        //加签之前需要对参数按字典排序
        String paramsStr = JsonOrderProcessor.convert(JSONUtil.toJsonStr(vcObject));
        log.info("加签:{}", paramsStr);
        // eth 方式签名
        // 计算address
        String address = EthVCSigner.getAddressFromPrivateKey(privateKeyHex);

        // 计算contentHash
        String contentHash = EthVCSigner.calculateContentHash(paramsStr);

        // 签名生成r/s/v
        Sign.SignatureData signatureData = EthVCSigner.signVC(privateKeyHex, paramsStr);
        String r = Numeric.toHexStringNoPrefix(signatureData.getR());
        String s = Numeric.toHexStringNoPrefix(signatureData.getS());
        int v = signatureData.getV()[0] & 0xFF; // 转换为无符号整数

        VerifiableCredentialDto.ProofEth proof = new VerifiableCredentialDto.ProofEth();
        proof
                .setContentHash(contentHash)
                .setR(r)
                .setS(s)
                .setV(v)
                .setVerificationMethod(did + "#key1");
        vcObject.setProof(proof);

        return vcObject;
    }

    private ResourceDomainInfo getResourceDomainInfo(String resourceName) {
        String details = null;
        try {
            details = resourceDomainService.getDomainDetails(resourceName);
        } catch (Exception e) {
            log.error("getDomainDetails error: {}", e.getMessage());
            throw new ApiException("获取链上资源失败");
        }
        ResourceDomainInfo resourceDomainInfo = null;
        if (StringUtils.isNotBlank(details)) {
            try {
                if (-1 != details.indexOf("{") && -1 != details.lastIndexOf("}")) {
                    details = details.substring(details.indexOf("{"), details.indexOf("}") + 1);
                }
                if (StringUtils.isNotBlank(details)) {
                    resourceDomainInfo = JSONUtil.toBean(details, ResourceDomainInfo.class);
                }
            } catch (Exception e) {
                log.error("getDomainDetails convert error: {}", e.getMessage());
                throw new ApiException("解析链上资源失败");
            }
        } else {
            throw new ApiException("链上资源不存在");
        }
        return resourceDomainInfo;
    }

    public TxHashResp call(CrossChainRequest request) throws Exception {
        //根据resource_name,获取ResourceDomain合约上的信息（orign/target）
        ResourceDomainInfo originResourceDomainInfo = getResourceDomainInfo(request.getOrigin().getResourceName());
        //将ResourceDomain合约上的信息，写入origin/target vc中
        VerifiableCredentialDto originObject = dealVc(originResourceDomainInfo, request.getPrivateKeyHex(), request.getOrigin(), request.getDid());

        VerifiableCredentialListDto verifiableCredentialListDto = new VerifiableCredentialListDto();
        verifiableCredentialListDto.setOrigin(originObject);

        String paramsStr = JsonOrderProcessor.convert(verifiableCredentialListDto);
        log.info("vc: {}", paramsStr);
        //调用origin合约方法
        TxHashResp result = contractService.call(paramsStr);
        return result;
    }

    public String get(CrossChainRequest request) throws Exception {
        //根据resource_name,获取ResourceDomain合约上的信息（orign/target）
        ResourceDomainInfo originResourceDomainInfo = getResourceDomainInfo(request.getOrigin().getResourceName());
        //将ResourceDomain合约上的信息，写入origin/target vc中
        VerifiableCredentialDto originObject = dealVc(originResourceDomainInfo, request.getPrivateKeyHex(), request.getOrigin(), request.getDid());

        VerifiableCredentialListDto verifiableCredentialListDto = new VerifiableCredentialListDto();
        verifiableCredentialListDto.setOrigin(originObject);

        String paramsStr = JsonOrderProcessor.convert(verifiableCredentialListDto);
        log.info("vc: {}", paramsStr);
        //调用origin合约方法
        String result = contractService.get(paramsStr);
        return result;
    }

    public String getDidDocument(GetBusinessDidDocumentRequest request) throws Exception {

        //根据resource_name,获取ResourceDomain合约上的信息（orign/target）
        ResourceDomainInfo originResourceDomainInfo = getResourceDomainInfo(request.getResourceName());
        //将ResourceDomain合约上的信息，写入origin/target vc中
        CredentialSubject origin = new CredentialSubject();
        origin.setResourceName(request.getResourceName());
        origin.setContractFunc("getDIDDetails");
        origin.setFuncParams(request.getDid());
        origin.setKey("useFuncParams");
        if (CHAINMAKER_WASM.getCode().equals(originResourceDomainInfo.getContractType())) {
            origin.setParamName("did");
        }

//        origin.setDid(request.getDid());
        //origin.setParamName("vcs");
        VerifiableCredentialDto originObject = dealVc(originResourceDomainInfo, request.getPrivateKeyHex(), origin, request.getDid());

        VerifiableCredentialListDto verifiableCredentialListDto = new VerifiableCredentialListDto();
        verifiableCredentialListDto.setOrigin(originObject);


        String paramsStr = JsonOrderProcessor.convert(verifiableCredentialListDto);
        log.info("vc: {}", paramsStr);
        //调用origin合约方法
        String result = contractService.get(paramsStr);

        return result;
    }

    public String getBalance(GetBalanceRequest request) throws Exception {
        //根据resource_name,获取ResourceDomain合约上的信息（orign/target）
        ResourceDomainInfo originResourceDomainInfo = getResourceDomainInfo(request.getResourceName());
        //将ResourceDomain合约上的信息，写入origin/target vc中
        CredentialSubject origin = new CredentialSubject();
        origin.setResourceName(request.getResourceName());

        // 设置获取用户身份的方法
        if (FABRIC_WASM.getCode().equals(originResourceDomainInfo.getContractType())) {
            origin.setContractFunc("Erc20GetBalanceVcs");
        } else {
            origin.setContractFunc("erc20GetBalanceVcs");
        }

        JSONObject obj = new JSONObject();
        obj.put("account", request.getDid());
        origin.setFuncParams(obj);
        if (CHAINMAKER_WASM.getCode().equals(originResourceDomainInfo.getContractType())) {
            origin.setParamName("vcs");
        }
        VerifiableCredentialDto originObject = dealVc(originResourceDomainInfo, request.getPrivateKeyHex(), origin, request.getDid());
        VerifiableCredentialListDto verifiableCredentialListDto = new VerifiableCredentialListDto();
        verifiableCredentialListDto.setOrigin(originObject);

        String paramsStr = JsonOrderProcessor.convert(verifiableCredentialListDto);
        log.info("vc: {}", paramsStr);
        //调用origin合约方法
        String result = contractService.get(paramsStr);
        return result;
    }

    public TxHashResp mint(MintRequest request) throws Exception {
        //根据resource_name,获取ResourceDomain合约上的信息（orign/target）
        ResourceDomainInfo originResourceDomainInfo = getResourceDomainInfo(request.getResourceName());
        //将ResourceDomain合约上的信息，写入origin/target vc中
        CredentialSubject origin = new CredentialSubject();
        origin.setResourceName(request.getResourceName());
        origin.setContractFunc("erc20MintVcs");
        JSONObject obj = new JSONObject();
        obj.put("amount", request.getAmount());
        obj.put("to", request.getDid());
        origin.setFuncParams(obj);
        if (CHAINMAKER_WASM.getCode().equals(originResourceDomainInfo.getContractType())) {
            origin.setParamName("vcs");
        }
        VerifiableCredentialDto originObject = dealVc(originResourceDomainInfo, request.getPrivateKeyHex(), origin, request.getDid());
        VerifiableCredentialListDto verifiableCredentialListDto = new VerifiableCredentialListDto();
        verifiableCredentialListDto.setOrigin(originObject);

        String paramsStr = JsonOrderProcessor.convert(verifiableCredentialListDto);
        log.info("vc: {}", paramsStr);
        //调用origin合约方法
        TxHashResp result = contractService.call(paramsStr);
        return result;
    }
}


