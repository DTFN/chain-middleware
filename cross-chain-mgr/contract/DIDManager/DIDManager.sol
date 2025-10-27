// SPDX-License-Identifier: MIT
pragma solidity ^0.4.24;

contract DIDManager {
    // 存储DID地址到详情的映射
    mapping(string => string) private _didDocuments;
    // 记录已存在的DID地址
    mapping(string => bool) private _exists;

    // 当DID被创建时触发的事件
    event DIDModify(string did, string didDocument);

    // 返回值
    event Result(string did);

    /**
     * @dev 根据DID地址查询详情
     * @param did DID地址
     * @return 详情序列化的字符串
     */
    function getDIDDetails(string did) external returns (string) {
        emit Result(_didDocuments[did]);
        return _didDocuments[did];
    }

    /**
     * @dev 创建新的DID记录，已存在则报错
     * @param did DID地址
     * @param didDocument 详情序列化的字符串
     */
    function createDID(string did, string didDocument) external {
        // 检查DID是否已存在
        require(!_exists[did], "DIDManager: DID already exists");
        
        // 保存DID详情
        _didDocuments[did] = didDocument;
        _exists[did] = true;
        
        // 触发创建事件，包含时间戳
        emit DIDModify(did, didDocument);
    }

    /**
     * @dev 更新已存在的DID记录
     * @param did DID地址
     * @param newDidDocument 新的详情序列化字符串
     */
    function updateDID(string did, string newDidDocument) external {
        // 检查DID是否存在
        // require(_exists[did], "DIDManager: DID does not exist");
        
        // 更新DID详情
        _didDocuments[did] = newDidDocument;
        _exists[did] = true;
        
        // 触发创建事件，包含时间戳
        emit DIDModify(did, newDidDocument);
    }

    /**
     * @dev echo测试函数
     * @param did DID地址
     */
    function echo(string did) external returns (string) {
        emit Result(did);
        return did;
    }

    /**
     * @dev 检查DID是否已存在
     * @param did DID地址
     * @return 是否存在的布尔值
     */
    function doesDIDExist(string did) external view returns (bool) {
        return _exists[did];
    }
}
