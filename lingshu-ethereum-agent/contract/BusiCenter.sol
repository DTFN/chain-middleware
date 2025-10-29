// SPDX-License-Identifier: MIT
pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "./JsmnSolLib.sol";

// 签名测试用例
// "{\"origin\": {\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain1\",\"contract_address\":\"contract01\",\"contract_func\":\"pay\",\"func_params\":\"[{\\\"string\\\":\\\"c1\\\"},{\\\"string\\\":\\\"121.11\\\"}]\",\"gateway_id\":\"0\",\"reource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"address\":\"00000000000000000000000094d6dd6b20628043de7fc036eefaf425ac32c87c\",\"contentHash\":\"7e5f28475d086414322f86a7e5d0c9d0d9015382f9048daa7937d6e4d4672ec8\",\"r\":\"e451bf696df6979d8d3cb56386337c0c93c131b3f1e2818ea7fec994d9e308e9\",\"s\":\"3443c0cee85fb040fa823116afc614280ee4e2529887e44c3ec857874ae9a9cf\",\"v\":27},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"]},\"target\": {\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain1\",\"contract_address\":\"contract01\",\"contract_func\":\"stock\",\"func_params\":\"[{\\\"string\\\":\\\"c1\\\"},{\\\"string\\\":\\\"12.1\\\"}]\",\"gateway_id\":\"0\",\"reource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"address\":\"00000000000000000000000094d6dd6b20628043de7fc036eefaf425ac32c87c\",\"contentHash\":\"2387acb4f91b2fe0a77d446e2b36b0fa6d0bfa0ac8793d6bc21f2353f6e04c46\",\"r\":\"d996c09255403893e6e71797243f158399fa0da79e89d48fffc85b5e01d2a760\",\"s\":\"5d4b09326348b0bde0605a8dd0bb9614dfbf938aa8803614c17c84420c229765\",\"v\":27},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"],\"z\":\"1\"}}"

contract BusiCenter {
    event DEBUG(string);
    event DEBUG_ADDRESS(address);
    event DEBUG_BYTES32(bytes32);

    event CROSS_CHAIN_VC(string);


    function echo(string calldata json) external {
        emit DEBUG(json);
    }

    function keyFromString(string memory key) public returns (bytes32) {
        return keccak256(abi.encodePacked(key));
    }

    // 余额映射
    mapping(bytes32 => uint256) public balanceOf;

    // 转账事件
    event Erc20Transfer(string from, string to, uint256 value);

    // 铸造新代币并分配给指定地址(VC)
    function erc20MintVcs(string memory vcs) public {
        string memory usedVc = getUsedVc(vcs);

        string memory credentialSubject = valueByKey(usedVc, "credentialSubject");
        string memory func_params = valueByKey(credentialSubject, "func_params");
        string memory to = valueByKey(func_params, "to");
        string memory amountStr = valueByKey(func_params, "amount");
        uint amount = uint(JsmnSolLib.parseInt(amountStr));
        erc20Mint(to, amount);

        // 4 发送跨链事件
        sendVcIfNeed(vcs);
    }

    /**
     * @dev 铸造新代币并分配给指定地址
     * @param to 接收者地址
     * @param amount 铸造数量
     */
    function erc20Mint(string memory to, uint256 amount) public {
        balanceOf[keyFromString(to)] += amount;
        emit Erc20Transfer("", to, amount);
    }

    // 转账功能(VC)
    function erc20TransferVcs(string memory vcs) public {
        string memory usedVc = getUsedVc(vcs);
        string memory credentialSubject = valueByKey(usedVc, "credentialSubject");
        string memory func_params = valueByKey(credentialSubject, "func_params");
        
        string memory from = valueByKey(func_params, "from");
        string memory to = valueByKey(func_params, "to");
        string memory amountStr = valueByKey(func_params, "amount");
        uint amount = uint(JsmnSolLib.parseInt(amountStr));
        erc20Transfer(from, to, amount);

        // 4 发送跨链事件
        sendVcIfNeed(vcs);
    }

    /**
     * @dev 转账功能
     * @param to 接收者地址
     * @param amount 转账数量
     * @return 转账是否成功
     */
    function erc20Transfer(string memory from, string memory to, uint256 amount) public returns (bool) {
        require(balanceOf[keyFromString(from)] >= amount, "not ");
        
        balanceOf[keyFromString(from)] -= amount;
        balanceOf[keyFromString(to)] += amount;
        emit Erc20Transfer(from, to, amount);
        return true;
    }

    /**
    查看余额
     */
    function erc20GetBalance(string memory account) public returns (uint256) {
        return balanceOf[keyFromString(account)];
    }

    // 记录每个NFT的所有者（bytes32类型）
    mapping(uint256 => bytes32) private _owners;
    // 下一个要铸造的tokenId（从1开始自增）
    uint256 private _nextTokenId = 1;
    
    // 转移事件
    event Erc721Transfer(string from, string to, uint256 tokenId);

    // 铸造新NFT，自动分配递增的tokenId(VC)
    function erc721MintVcs(string memory vcs) public {
        string memory usedVc = getUsedVc(vcs);
        string memory credentialSubject = valueByKey(usedVc, "credentialSubject");
        string memory func_params = valueByKey(credentialSubject, "func_params");
        
        string memory to = valueByKey(func_params, "to");
        erc721Mint(to);

        // 4 发送跨链事件
        sendVcIfNeed(vcs);
    }
    
    /**
     * @dev 铸造新NFT，自动分配递增的tokenId:
     * @param to 接收者（bytes32类型）
     * @return 新铸造NFT的tokenId
     */
    function erc721Mint(string memory to) public returns (uint256) {
        uint256 tokenId = _nextTokenId;
        require(_owners[tokenId] == bytes32(0), "tokenId已存在");
        
        // 仅记录所有者，不维护余额统计
        _owners[tokenId] = keyFromString(to);
        _nextTokenId++;
        
        emit Erc721Transfer("", to, tokenId);
        return tokenId;
    }

    /**
     * @dev 查看NFT所有者
     * @param tokenId NFT的标识
     * @return 该NFT的所有者（bytes32类型）
     */
    function erc721OwnerOf(string memory from, uint256 tokenId) public returns (bool) {
        bytes32 owner = _owners[tokenId];
        if (owner == keyFromString(from)) {
            return true;
        } else {
            return false;
        }
    }

    // 铸造新NFT，自动分配递增的tokenId(VC)
    function erc721TransferVcs(string memory vcs) public {
        string memory usedVc = getUsedVc(vcs);
        string memory credentialSubject = valueByKey(usedVc, "credentialSubject");
        string memory func_params = valueByKey(credentialSubject, "func_params");
        
        string memory from = valueByKey(func_params, "from");
        string memory to = valueByKey(func_params, "to");
        string memory tokenIdStr = valueByKey(func_params, "tokenId");
        uint tokenId = uint(JsmnSolLib.parseInt(tokenIdStr));
        erc721Transfer(from, to, tokenId);

        // 4 发送跨链事件
        sendVcIfNeed(vcs);
    }
    
    /**
     * @dev 转账NFT
     * @param from 发送者（bytes32类型）
     * @param to 接收者（bytes32类型）
     * @param tokenId 要转移的NFT标识
     */
    function erc721Transfer(string memory from, string memory to, uint256 tokenId) public {
        require(_owners[tokenId] == keyFromString(from), "发送者不是NFT所有者");
        
        // 仅更新所有者，不处理余额
        _owners[tokenId] = keyFromString(to);
        
        emit Erc721Transfer(from, to, tokenId);
    }

    // 如果需要则发送VC
    function sendVcIfNeed(string memory json) public returns (string memory) {
        string memory originVc = valueByKey(json, "origin");
        string memory targetVc = valueByKey(json, "target");

        if (JsmnSolLib.strCompare(originVc, "") != 0 && JsmnSolLib.strCompare(targetVc, "") != 0) {
            verifyVc(targetVc);
            emit CROSS_CHAIN_VC(targetVc);
        }
    }

    // 先选择vc0,否则选择vc1
    function getUsedVc(string memory json) public returns (string memory) {
        string memory usedVc = valueByKey(json, "origin");

        if (JsmnSolLib.strCompare(usedVc, "") != 0) {
            verifyVc(usedVc);
        } else {
            usedVc = valueByKey(json, "target");
            verifyVc(usedVc);
        }

        return usedVc;
    }

    // 下面是示例
    function run(string calldata json) external {
        // 1 提取两个VC
        string memory vc0 = valueByKey(json, "origin");
        string memory vc1 = valueByKey(json, "target");

        // 2 验证VC
        verifyVc(vc0);
        verifyVc(vc1);

        // 3 todo 业务代码

        // 4 todo 发送跨链事件
    }

    function verifyVc(string memory vcJson) public {
        // 1 提取proof
        string memory proof = valueByKey(vcJson, "proof");

        // 2 提取 hash r s v
        string memory contentHash = valueByKey(proof, "contentHash");
        string memory r = valueByKey(proof, "r");
        string memory s = valueByKey(proof, "s");
        string memory v = valueByKey(proof, "v");
        bytes32 contentHashBytes = hexStringToBytes32(contentHash);
        bytes32 rBytes = hexStringToBytes32(r);
        bytes32 sBytes = hexStringToBytes32(s);
        uint8 vuint = uint8(JsmnSolLib.parseInt(v));

        // 3 todo 查询did合约

        // 4 验签（验证地址是否一致） test:address 000000000000000000000000bd5d9f6b1b70cfd36c90aa11a8830a43062ab4ea
        address calAddress = ecrecover(contentHashBytes, vuint, rBytes, sBytes);
        require(
            hexStringToBytes32("0x000000000000000000000000bd5d9f6b1b70cfd36c90aa11a8830a43062ab4ea")
            ==
            bytes32(uint256(uint160(calAddress)))
        , string(abi.encodePacked("address not true,contentHash,", contentHash)));

        // 5 序列化VC并hash，判断hash是否一致
        // 5.1 获取签名的VC
        (uint start, uint end) = getProofRange(vcJson);
        string memory vcStrNoProof = removestring(vcJson, start, end);

        // 5.2 对VC做hash,需要和VC.proof的hash一致
            // 测试hash函数,下面两个一样
            // emit DEBUG_BYTES32(hexStringToBytes32("0x2912723b3ed60c075b271f075d881d82fa5de112b6c25f7dfa4cab85de25045a"));
            // emit DEBUG_BYTES32(getEthSignedMessageHash("123456"));
        bytes32 vcStrNoProofHash = getEthSignedMessageHash(vcStrNoProof);
        require(vcStrNoProofHash == contentHashBytes, string(abi.encodePacked("vc hash not equal,", vcStrNoProof)));
    }

    function getProofRange(string memory json) public returns (uint, uint) {
        uint returnValue;
        JsmnSolLib.Token[] memory tokens;
        uint actualNum;

        (returnValue, tokens, actualNum) = JsmnSolLib.parse(json, 200);
        
        for (uint i = 0; i < tokens.length; i++) {
            JsmnSolLib.Token memory t = tokens[i];
            string memory key = JsmnSolLib.getBytes(json, t.start, t.end);
            if (JsmnSolLib.strCompare(key, "proof") == 0) {
                JsmnSolLib.Token memory vt = tokens[i+1];

                // 1. 字符串在Token不包含 "
                // 2. 删除proof前面的,
                // todo ,和 "proof" 之间可能有空格
                return (t.start - 2, vt.end);
            }
        }

        return (0, 0);
    }

    function showAll(string memory json) public returns (string memory) {
        uint returnValue;
        JsmnSolLib.Token[] memory tokens;
        uint actualNum;

        (returnValue, tokens, actualNum) = JsmnSolLib.parse(json, 200);
        
        for (uint i = 0; i < tokens.length; i++) {
            JsmnSolLib.Token memory t = tokens[i];
            string memory key = JsmnSolLib.getBytes(json, t.start, t.end);
            emit DEBUG(key);
        }
    }

    function valueByIndex(string memory json, uint key_index) public returns (string memory) {
        uint returnValue;
        JsmnSolLib.Token[] memory tokens;
        uint actualNum;

        (returnValue, tokens, actualNum) = JsmnSolLib.parse(json, 200);
        
        JsmnSolLib.Token memory t = tokens[key_index];
        string memory key = JsmnSolLib.getBytes(json, t.start, t.end);
        return key;
    }

    function valueByKey(string memory json, string memory needKey) public returns (string memory) {
        uint returnValue;
        JsmnSolLib.Token[] memory tokens;
        uint actualNum;

        (returnValue, tokens, actualNum) = JsmnSolLib.parse(json, 200);
        
        for (uint i = 0; i < tokens.length; i++) {
            JsmnSolLib.Token memory t = tokens[i];
            string memory key = JsmnSolLib.getBytes(json, t.start, t.end);
            if (JsmnSolLib.strCompare(key, needKey) == 0) {
                JsmnSolLib.Token memory vt = tokens[i+1];
                string memory valueStr = JsmnSolLib.getBytes(json, vt.start, vt.end);
                return valueStr;
            }
        }
        return "";
    }

    function removestring(string memory original, uint256 start, uint256 end) public pure returns (string memory) {
        bytes memory originalBytes = bytes(original);
        
        // 验证输入的有效性
        require(start < originalBytes.length, "Start index out of bounds");
        require(end <= originalBytes.length, "End index out of bounds");
        require(start <= end, "Start index must be less than or equal to end index");
        
        // 计算结果字符串的长度
        uint256 resultLength = originalBytes.length - (end - start);
        bytes memory result = new bytes(resultLength);
        
        uint256 resultIndex = 0;
        
        // 复制start之前的字符
        for (uint256 i = 0; i < start; i++) {
            result[resultIndex] = originalBytes[i];
            resultIndex++;
        }
        
        // 复制end之后的字符
        for (uint256 i = end; i < originalBytes.length; i++) {
            result[resultIndex] = originalBytes[i];
            resultIndex++;
        }
        
        return string(result);
    }

    function hexStringToBytes32(string memory hexString) public returns (bytes32) {
        if (bytes(hexString).length == 66) {
            return hexStringToBytes32(JsmnSolLib.getBytes(hexString, 2, 66));
        }
        require(bytes(hexString).length == 64, "Invalid hex string length");
        
        bytes memory bytesString = bytes(hexString);
        uint256 bytesLength = bytesString.length;
        bytes32 result;
        
        for (uint256 i = 0; i < bytesLength; i += 2) {
            uint8 b1 = uint8(bytesString[i]);
            uint8 b2 = uint8(bytesString[i + 1]);
            
            // 转换字符到对应的数值
            uint8 byteValue = (hexCharToByte(b1) << 4) | hexCharToByte(b2);
            result |= bytes32(uint256(byteValue) << (248 - (i / 2) * 8));
        }
        
        return result;
    }
    function hexCharToByte(uint8 c) private pure returns (uint8) {
        if (c >= 48 && c <= 57) { // 0-9
            return c - 48;
        } else if (c >= 65 && c <= 70) { // A-F
            return c - 55;
        } else if (c >= 97 && c <= 102) { // a-f
            return c - 87;
        } else {
            revert("Invalid hex character");
        }
    }

    function getEthSignedMessageHash(
        string memory _message
    ) public pure returns (bytes32) {
        uint stringLength = bytes(_message).length;
        string memory stringLengthStr = uintToString(stringLength);
        return
            keccak256(
                abi.encodePacked(
                    "\x19Ethereum Signed Message:\n",
                    stringLengthStr,
                    _message
                )
            );
    }
    function uintToString(uint256 _number) public pure returns (string memory) {
        // 处理0的特殊情况
        if (_number == 0) {
            return "0";
        }
        
        // 计算数字的位数
        uint256 temp = _number;
        uint256 digits;
        while (temp != 0) {
            digits++;
            temp /= 10;
        }
        
        // 分配字符串内存
        bytes memory buffer = new bytes(digits);
        
        // 填充字符串
        uint256 index = digits - 1;
        temp = _number;
        while (temp != 0) {
            buffer[index--] = bytes1(uint8(48 + temp % 10));
            temp /= 10;
        }
        
        return string(buffer);
    }

    function test_verify() public pure returns (address) {
        bytes32 _ethSignedMessageHash = hex"478ec551ff5eb934d827f6e3d2827ba63e58885bbeaf76820090c8e16d836bc6";
        bytes32 r = hex"815186624fffdc410a3d97e0762b6d2405883836a5e55349be2a4e9fe8ee7246";
        bytes32 s = hex"2243ca645c8397cd0f389d5fc721f586bc76407a8d8d70375ab4eedf7b60ef33";
        uint8 v = 28;

        //  0x94d6dd6b20628043de7fc036eefaf425ac32c87c
        return ecrecover(_ethSignedMessageHash, v, r, s);
    }

}
