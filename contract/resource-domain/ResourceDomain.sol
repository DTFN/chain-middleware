// SPDX-License-Identifier: MIT
pragma solidity ^0.4.24;

contract ResourceDomain {
    // 存储资源域名地址到详情的映射
    mapping(string => string) private _domainDetails;
    // 记录已存在的资源域名地址
    mapping(string => bool) private _exists;

    // 当资源域名被注册时触发
    event DomainRegistered(string domainAddress, string details);

    /**
     * @dev 根据资源域名地址查询详情
     * @param domainAddress 资源域名地址
     * @return 详情序列化的字符串
     */
    function getDomainDetails(string domainAddress) external view returns (string) {
        return _domainDetails[domainAddress];
    }

    /**
     * @dev 保存资源域名详情，已存在则报错
     * @param domainAddress 资源域名地址
     * @param details 详情序列化的字符串
     */
    function saveDomainDetails(string domainAddress, string details) external {
        // 检查资源域名是否已存在
        require(!_exists[domainAddress], "ResourceDomain: domain already exists");
        
        // 保存详情
        _domainDetails[domainAddress] = details;
        _exists[domainAddress] = true;
        
        // 触发事件
        emit DomainRegistered(domainAddress, details);
    }

    /**
     * @dev 保存资源域名详情，覆盖旧数据
     * @param domainAddress 资源域名地址
     * @param details 详情序列化的字符串
     */
    function updateDomainDetails(string domainAddress, string details) external {
        // 保存详情
        _domainDetails[domainAddress] = details;
        _exists[domainAddress] = true;
        
        // 触发事件
        emit DomainRegistered(domainAddress, details);
    }

    /**
     * @dev 检查资源域名是否已存在
     * @param domainAddress 资源域名地址
     * @return 是否存在的布尔值
     */
    function doesDomainExist(string domainAddress) external view returns (bool) {
        return _exists[domainAddress];
    }
}
