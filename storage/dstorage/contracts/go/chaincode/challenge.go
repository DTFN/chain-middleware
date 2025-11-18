// SPDX-License-Identifier: Apache-2.0
package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// ====================
// 数据结构定义
// ====================

type ChallengeContract struct {
	contractapi.Contract
}

// 必须设置合约名称
func (c *ChallengeContract) GetName() string {
    return "ChallengeContract"
}

type Challenge struct {
	ChallengeID uint64 `json:"challenge_id"`
	FileHash    string `json:"file_hash"`
	Nonce       string `json:"nonce"`
	CID         string `json:"cid"`
	Prover      string `json:"prover"`  // address
	Issuer      string `json:"issuer"`  // msg.sender
    ProofSubmitted  bool   `json:"proof_submitted"`  // 是否已提交证明
	ProofVerified   bool   `json:"proof_verified"`   // 验证是否通过
    ProofHasVerified bool   `json:"proof_has_verified"` // 这条挑战是否已验证
}

type Proof struct {
	ChallengeID      uint64 `json:"challenge_id"`
	ProofCID         string `json:"proof_cid"`
	ProofURL         string `json:"proof_url"`
	ProofHash        string `json:"proof_hash"`
	PublicCID        string `json:"public_cid"`
	PublicURL        string `json:"public_url"`
	PublicHash       string `json:"public_hash"`
	VerifyKeyCID     string `json:"vk_cid"`
	VerifyKeyURL     string `json:"vk_url"`
	VerifyKeyHash    string `json:"vk_hash"`
}

// ====================
// 合约方法实现
// ====================

func (s *ChallengeContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// 发起挑战
func (s *ChallengeContract) CreateChallenge(ctx contractapi.TransactionContextInterface,
	fileHash, nonce, cid, prover string) (string, error) {

	issuer, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to get client identity: %v", err)
	}

	challengeCounterKey := "CHALLENGE_COUNTER"
	challengeIDBytes, err := ctx.GetStub().GetState(challengeCounterKey)
	var challengeID uint64 = 1
	if err == nil && challengeIDBytes != nil {
		challengeID, _ = strconv.ParseUint(string(challengeIDBytes), 10, 64)
		challengeID++
	}

	newChallenge := Challenge{
		ChallengeID: challengeID,
		FileHash:    fileHash,
		Nonce:       nonce,
		CID:         cid,
		Prover:      prover,
		Issuer:      issuer,
        ProofSubmitted: false,
        ProofVerified:  false,
        ProofHasVerified: false,
	}

	challengeBytes, _ := json.Marshal(newChallenge)
	challengeKey := fmt.Sprintf("CHALLENGE_%d", challengeID)

	if err := ctx.GetStub().PutState(challengeKey, challengeBytes); err != nil {
		return "", err
	}
	if err := ctx.GetStub().PutState(challengeCounterKey, []byte(strconv.FormatUint(challengeID, 10))); err != nil {
		return "", err
	}

	// 维护用户挑战索引（issuer 或 prover 均可索引）
	userKey := fmt.Sprintf("USER_CHALLENGES_%s", prover)
	userListBytes, _ := ctx.GetStub().GetState(userKey)
	var userList []uint64
	json.Unmarshal(userListBytes, &userList)
	userList = append(userList, challengeID)
	userListBytesNew, _ := json.Marshal(userList)
	ctx.GetStub().PutState(userKey, userListBytesNew)

	return fmt.Sprintf("%d", challengeID), nil
}

// 查询挑战详情
func (s *ChallengeContract) GetChallenge(ctx contractapi.TransactionContextInterface, challengeID uint64) (*Challenge, error) {
	challengeKey := fmt.Sprintf("CHALLENGE_%d", challengeID)
	challengeBytes, err := ctx.GetStub().GetState(challengeKey)
	if err != nil || challengeBytes == nil {
		return nil, fmt.Errorf("challenge not found")
	}
	var c Challenge
	if err := json.Unmarshal(challengeBytes, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// 查询当前用户参与的所有挑战（作为 issuer）
func (s *ChallengeContract) GetMyChallenges(ctx contractapi.TransactionContextInterface) (*Challenge, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("cannot get client ID: %v", err)
	}
	userKey := fmt.Sprintf("USER_CHALLENGES_%s", clientID)
	userListBytes, err := ctx.GetStub().GetState(userKey)
	if err != nil || userListBytes == nil {
		return nil, nil
	}
	var userList []uint64
	if err := json.Unmarshal(userListBytes, &userList); err != nil {
		return nil, err
	}

    // 倒序遍历，找到第一个未提交证明的挑战
	for i := len(userList) - 1; i >= 0; i-- {
		chID := userList[i]
		challengeKey := fmt.Sprintf("CHALLENGE_%d", chID)

		chBytes, err := ctx.GetStub().GetState(challengeKey)
		if err != nil || chBytes == nil {
			continue // 跳过错误或缺失的数据
		}

		var ch Challenge
		if err := json.Unmarshal(chBytes, &ch); err != nil {
			continue
		}

		// 判断是否未提交证明
		if ch.ProofSubmitted == false {
			return &ch, nil
		}
	}
	return nil, nil
}

// 提交证明数据
func (s *ChallengeContract) SubmitProof(ctx contractapi.TransactionContextInterface,
	challengeID uint64,
	proofCID, proofURL, proofHash string,
	publicCID, publicURL, publicHash string,
	vkCID, vkURL, vkHash string,
) error {
	proof := Proof{
		ChallengeID:   challengeID,
		ProofCID:      proofCID,
		ProofURL:      proofURL,
		ProofHash:     proofHash,
		PublicCID:     publicCID,
		PublicURL:     publicURL,
		PublicHash:    publicHash,
		VerifyKeyCID:  vkCID,
		VerifyKeyURL:  vkURL,
		VerifyKeyHash: vkHash,
	}
	proofBytes, _ := json.Marshal(proof)
	proofKey := fmt.Sprintf("PROOF_%d", challengeID)
	if err := ctx.GetStub().PutState(proofKey, proofBytes); err != nil {
        return err
	}
    
    // 2. 获取对应 Challenge 并更新 ProofSubmitted = true
	challengeKey := fmt.Sprintf("CHALLENGE_%d", challengeID)
	challengeBytes, err := ctx.GetStub().GetState(challengeKey)
	if err != nil || challengeBytes == nil {
		return fmt.Errorf("challenge %d not found", challengeID)
	}

	var challenge Challenge
	if err := json.Unmarshal(challengeBytes, &challenge); err != nil {
		return fmt.Errorf("failed to unmarshal challenge: %v", err)
	}

	challenge.ProofSubmitted = true

	updatedChallengeBytes, err := json.Marshal(challenge)
	if err != nil {
		return fmt.Errorf("failed to marshal updated challenge: %v", err)
	}

	if err := ctx.GetStub().PutState(challengeKey, updatedChallengeBytes); err != nil {
		return fmt.Errorf("failed to update challenge state: %v", err)
	}

	return nil
}

// 查询某个挑战的证明数据
func (s *ChallengeContract) GetProof(ctx contractapi.TransactionContextInterface, challengeID uint64) (*Proof, error) {
	proofKey := fmt.Sprintf("PROOF_%d", challengeID)
	proofBytes, err := ctx.GetStub().GetState(proofKey)
	if err != nil || proofBytes == nil {
		return nil, fmt.Errorf("proof not found")
	}
	var p Proof
	if err := json.Unmarshal(proofBytes, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *ChallengeContract) UpdateChallenge(ctx contractapi.TransactionContextInterface, challengeID uint64, ProofVerified bool) error {
    challengeKey := fmt.Sprintf("CHALLENGE_%d", challengeID)
	challengeBytes, err := ctx.GetStub().GetState(challengeKey)
	if err != nil || challengeBytes == nil {
		return fmt.Errorf("challenge %d not found", challengeID)
	}

	var challenge Challenge
	if err := json.Unmarshal(challengeBytes, &challenge); err != nil {
		return fmt.Errorf("failed to unmarshal challenge: %v", err)
	}

	challenge.ProofHasVerified = true
    challenge.ProofVerified = ProofVerified

	updatedChallengeBytes, err := json.Marshal(challenge)
	if err != nil {
		return fmt.Errorf("failed to marshal updated challenge: %v", err)
	}

	if err := ctx.GetStub().PutState(challengeKey, updatedChallengeBytes); err != nil {
		return fmt.Errorf("failed to update challenge state: %v", err)
	}

	return nil
}

