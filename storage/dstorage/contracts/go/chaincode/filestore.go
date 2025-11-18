package chaincode

import (
	"encoding/json"
	"fmt"
	"time"
    "strconv"


	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type FileStoreContract struct {
	contractapi.Contract
}

// 必须设置合约名称
func (c *FileStoreContract) GetName() string {
    return "FileStoreContract"
}

func (f *FileStoreContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (f *FileStoreContract) GetUser(ctx contractapi.TransactionContextInterface) (string, error) {
	issuer, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}

	return fmt.Sprintf("%s", issuer), nil
}

// storeFile stores a new file's metadata
func (f *FileStoreContract) StoreFile(
    ctx contractapi.TransactionContextInterface,
    rawFileId, pk, fileId, fileName, author, fileSizeStr, totalShardsStr, dataShardsStr, ipfsUrlsJson, cidsJson string) error {

    fileSize, err := strconv.ParseUint(fileSizeStr, 10, 64)
    if err != nil {
        return fmt.Errorf("invalid fileSize: %v", err)
    }

    totalShards, err := strconv.ParseUint(totalShardsStr, 10, 64)
    if err != nil {
        return fmt.Errorf("invalid totalShards: %v", err)
    }

    dataShards, err := strconv.ParseUint(dataShardsStr, 10, 64)
    if err != nil {
        return fmt.Errorf("invalid dataShards: %v", err)
    }

    var ipfsUrls []string
    if err := json.Unmarshal([]byte(ipfsUrlsJson), &ipfsUrls); err != nil {
        return fmt.Errorf("invalid ipfsUrls JSON: %v", err)
    }

    var cids []string
    if err := json.Unmarshal([]byte(cidsJson), &cids); err != nil {
        return fmt.Errorf("invalid cids JSON: %v", err)
    }
	
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("cannot get client identity: %v", err)
	}

    // 构造 FileMeta
    meta := FileMeta{
		RawFileId:	 rawFileId,
		Pk:	 		 pk,
        Uploader:    clientID,
        FileName:    fileName,
        Author:      author,
        FileSize:    fileSize,
        TotalShards: totalShards,
        DataShards:  dataShards,
        IpfsUrls:    ipfsUrls,
        Cids:        cids,
        Timestamp:   time.Now().Unix(),
    }

    data, err := json.Marshal(meta)
    if err != nil {
        return fmt.Errorf("failed to marshal file meta: %v", err)
    }

    key := rawFileId + "_" + pk
    return ctx.GetStub().PutState(key, data)
}

// storeFile stores a new file's metadata
func (f *FileStoreContract) HasFile(ctx contractapi.TransactionContextInterface, fileId, fileName, author, pk string) (string, error) {
    fmt.Printf("[QueryFiles] Start query, params => fileId=%s, fileName=%s, author=%s, pk=%s\n", fileId, fileName, author, pk)

    var ids []string
    var access []bool
    accessMap := make(map[string]bool)
    seenOrder := make([]string, 0)

    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return "", fmt.Errorf("failed to get state range: %v", err)
    }
    defer resultsIterator.Close()

    count := 0
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return "", err
        }

        var fileMeta FileMeta
        if err := json.Unmarshal(queryResponse.Value, &fileMeta); err != nil {
            fmt.Printf("[WARN] Failed to unmarshal record: %v\n", err)
            continue
        }

        count++
        fmt.Printf("[DEBUG] #%d FileMeta => RawFileId=%s, FileName=%s, Author=%s, Pk=%s\n",
            count, fileMeta.RawFileId, fileMeta.FileName, fileMeta.Author, fileMeta.Pk)

        // 条件匹配
        if (fileId == "" || fileMeta.RawFileId == fileId) &&
            (fileName == "" || fileMeta.FileName == fileName) &&
            (author == "" || fileMeta.Author == author) {

            hasAccess := true
            if pk != "" && fileMeta.Pk != pk {
                hasAccess = false
            }

            if old, exists := accessMap[fileMeta.RawFileId]; exists {
                accessMap[fileMeta.RawFileId] = old || hasAccess
                fmt.Printf("[INFO] Update accessMap[%s] from %v to %v\n", fileMeta.RawFileId, old, accessMap[fileMeta.RawFileId])
            } else {
                accessMap[fileMeta.RawFileId] = hasAccess
                seenOrder = append(seenOrder, fileMeta.RawFileId)
                fmt.Printf("[INFO] New entry => fileId=%s, hasAccess=%v\n", fileMeta.RawFileId, hasAccess)
            }
        }
    }

    // 依序构造结果
    for _, fid := range seenOrder {
        ids = append(ids, fid)
        access = append(access, accessMap[fid])
    }

    result := map[string]interface{}{
        "ids":    ids,
        "access": access,
    }

    jsonBytes, err := json.Marshal(result)
    if err != nil {
        return "", fmt.Errorf("failed to marshal result: %v", err)
    }

    fmt.Printf("[QueryFiles] Finished, matched %d files\n", len(ids))
    fmt.Printf("[QueryFiles] Return => %s\n", string(jsonBytes))

    return string(jsonBytes), nil
}

// getUploader returns uploader's identity for a file
func (f *FileStoreContract) GetUploader(ctx contractapi.TransactionContextInterface, fileId, pk string) (string, error) {
	meta, err := f.GetFileMeta(ctx, fileId, pk)
	if err != nil {
		return "", err
	}
	return meta.Uploader, nil
}

// getErasureConfig returns total and data shards
func (f *FileStoreContract) GetTotalShards(ctx contractapi.TransactionContextInterface, fileId, pk string) (uint64, error) {
	meta, err := f.GetFileMeta(ctx, fileId, pk)
	if err != nil {
		return 0, err
	}
	return meta.TotalShards, nil
}

// getErasureConfig returns total and data shards
func (f *FileStoreContract) GetDataShards(ctx contractapi.TransactionContextInterface, fileId, pk string) (uint64, error) {
	meta, err := f.GetFileMeta(ctx, fileId, pk)
	if err != nil {
		return 0, err
	}
	return meta.DataShards, nil
}

// getFileMetaBasic returns key metadata
func (f *FileStoreContract) GetFileMetaBasic(ctx contractapi.TransactionContextInterface, fileId, pk string) (*FileMeta, error) {
	meta, err := f.GetFileMeta(ctx, fileId, pk)
	if err != nil {
		return nil, err
	}
	// strip urls/hashes to reduce payload
	meta.IpfsUrls = nil
	meta.Cids = nil
	return meta, nil
}

// getFileMetaIpfsUrls returns all IPFS urls
func (f *FileStoreContract) GetFileMetaIpfsUrls(ctx contractapi.TransactionContextInterface, fileId, pk string) ([]string, error) {
	meta, err := f.GetFileMeta(ctx, fileId, pk)
	if err != nil {
		return nil, err
	}
	return meta.IpfsUrls, nil
}

// getFileMetaPoseidonHashes returns all Poseidon hashes
func (f *FileStoreContract) GetFileMetaPoseidonHashes(ctx contractapi.TransactionContextInterface, fileId, pk string) ([]string, error) {
	meta, err := f.GetFileMeta(ctx, fileId, pk)
	if err != nil {
		return nil, err
	}
	return meta.Cids, nil
}

// Helper: get file meta from state
func (f *FileStoreContract) GetFileMeta(ctx contractapi.TransactionContextInterface, fileId, pk string) (*FileMeta, error) {
	key := fileId + "_" + pk
    bytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read state: %v", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("file not found: %s", key)
	}

	var meta FileMeta
	if err := json.Unmarshal(bytes, &meta); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	return &meta, nil
}
