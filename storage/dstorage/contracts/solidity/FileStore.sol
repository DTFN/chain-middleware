// SPDX-License-Identifier: MIT
pragma solidity ^0.8.11;

contract FileStore {
    struct FileMeta {
        string rawFileId;
        string pk;
        address uploader;
        string fileName;
        string author;
        string fileId;
        uint256 fileSize;
        uint256 totalShards;
        uint256 dataShards;
        string[] ipfsUrls;
        string[] shardHashes;
        string[] poseidonHashes;
        uint256 timestamp;
    }

    struct FileInput {
        string rawFileId; // 记录文件原始哈希
        string pk;
        string fileId;
        string fileName;
        string author;
        uint256 fileSize;
        uint256 totalShards;
        uint256 dataShards;
        string[] ipfsUrls;
        string[] shardHashes;
        string[] poseidonHashes;
    }

    struct FileAccessMeta {
        bool exists;
        // pk -> fileId             // 加密后文件唯一标识
        mapping(string => string) pkFileIds;
        // pk->bool
        mapping(string => bool) accessList; // 授权访问列表
    }

    constructor() {
        // 空的构造函数，不做任何操作
    }

    // 文件列表
    FileMeta[] public files;
    mapping(string => FileMeta) public filesById;
    mapping(string => FileAccessMeta) public filesAccessById;
    mapping(string => mapping(string => bool)) public rawFileToFileId;

    event FileStored(
        string rawFileId,
        address indexed uploader,
        string fileName,
        string author,
        uint256 fileSize,
        uint256 totalShards,
        uint256 dataShards,
        uint256 timestamp
    );

    function storeFile(FileInput calldata input) external {
        _validateFileInput(input);

        filesById[input.fileId] = FileMeta({
            rawFileId: input.rawFileId,
            pk: input.pk,
            uploader: msg.sender,
            fileName: input.fileName,
            author: input.author,
            fileId: input.fileId,
            fileSize: input.fileSize,
            totalShards: input.totalShards,
            dataShards: input.dataShards,
            ipfsUrls: input.ipfsUrls,
            shardHashes: input.shardHashes,
            poseidonHashes: input.poseidonHashes,
            timestamp: block.timestamp
        });
        rawFileToFileId[input.rawFileId][input.fileId] = true;

        files.push(filesById[input.fileId]);

        filesAccessById[input.rawFileId].exists = true;
        filesAccessById[input.rawFileId].pkFileIds[input.pk] = input.fileId;
        filesAccessById[input.rawFileId].accessList[input.pk] = true;

        emit FileStored(
            input.rawFileId,
            msg.sender,
            input.fileName,
            input.author,
            input.fileSize,
            input.totalShards,
            input.dataShards,
            block.timestamp
        );
    }

    function _validateFileInput(FileInput calldata input) internal pure {
        require(input.totalShards >= input.dataShards, "Invalid erasure config");
        require(input.ipfsUrls.length == input.totalShards, "ipfs count mismatch");
        require(input.shardHashes.length == input.totalShards, "shard count mismatch");
        require(input.poseidonHashes.length == input.totalShards, "poseidon count mismatch");
    }

    function hasFile(string calldata fileId, string calldata fileName, string calldata author, string calldata pk) public view returns (string[] memory, bool[] memory) { 
        require(bytes(fileId).length != 0 || bytes(fileName).length != 0 || bytes(author).length != 0, "At least one condition must be non-empty");

        uint256 count = 0;
        string[] memory temp = new string[](files.length);

        // 第一次遍历，筛选出匹配项
        for (uint256 i = 0; i < files.length; i++) {
            if (
                (bytes(fileId).length == 0 || keccak256(bytes(files[i].rawFileId)) == keccak256(bytes(fileId))) &&
                (bytes(fileName).length == 0 || keccak256(bytes(files[i].fileName)) == keccak256(bytes(fileName))) &&
                (bytes(author).length == 0 || keccak256(bytes(files[i].author)) == keccak256(bytes(author)))
            ) {
                // 检查是否重复
                bool exists = false;
                for (uint256 j = 0; j < count; j++) {
                    if (keccak256(bytes(temp[j])) == keccak256(bytes(files[i].rawFileId))) {
                        exists = true;
                        break;
                    }
                }

                if (!exists) {
                    temp[count] = files[i].rawFileId;
                    count++;
                }
            }
        }

        // 返回定长数组
        string[] memory result = new string[](count);
        bool[] memory resultBool = new bool[](count);
        for (uint256 i = 0; i < count; i++) {
            result[i] = temp[i];
            resultBool[i] = filesAccessById[result[i]].accessList[pk];
        }

        return (result, resultBool);
    }

    function getShardHashes(string calldata rawFileId, string calldata pk) external view returns (string[] memory) {
        string memory fileId = _getFileId(rawFileId, pk);
        return filesById[fileId].shardHashes;
    }

    function getErasureConfig(string calldata rawFileId, string calldata pk) external view returns (uint256 total, uint256 data) {
        string memory fileId = _getFileId(rawFileId, pk);
        FileMeta storage file = filesById[fileId];
        return (file.totalShards, file.dataShards);
    }

    function getUploader(string calldata rawFileId, string calldata pk) external view returns (address) {
        string memory fileId = _getFileId(rawFileId, pk);
        return filesById[fileId].uploader;
    }

    function getFileMeta(string calldata rawFileId, string calldata pk) external view returns (
        string memory fileId,
        string memory fileName,
        uint256 fileSize,
        uint256 totalShards,
        uint256 dataShards,
        string[] memory ipfsUrls,
        string[] memory shardHashes,
        string[] memory poseidonHashes,
        uint256 timestamp,
        string memory author
    ) {
        string memory newFileId = _getFileId(rawFileId, pk);
        FileMeta storage file = filesById[newFileId];
        return (
            newFileId,
            file.fileName,
            file.fileSize,
            file.totalShards,
            file.dataShards,
            file.ipfsUrls,
            file.shardHashes,
            file.poseidonHashes,
            file.timestamp,
            file.author
        );
    }

    function _getFileId(string calldata rawFileId, string calldata pk) private view returns (string memory) {
        if (!filesAccessById[rawFileId].exists) {
            return "";
        }

        if (filesAccessById[rawFileId].accessList[pk]) {
            return filesAccessById[rawFileId].pkFileIds[pk];
        } else {
            return "";
        }
    }
}
