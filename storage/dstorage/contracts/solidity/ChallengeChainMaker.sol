// SPDX-License-Identifier: MIT
pragma solidity ^0.8.11;

interface IVerifier {
    function verifyProof(
        uint[2] calldata a,
        uint[2][2] calldata b,
        uint[2] calldata c,
        uint[3] calldata input
    ) external view returns (bool);
}

contract Challenge {

    IVerifier public verifier;

    uint public challengeCounter;

    struct ChallengeInfo {
        address challenger;
        string fileHash;     // 被挑战的文件哈希或分片哈希
        string nonce;        // 挑战随机数
        string cid;
        uint timestamp;
        bool proofSubmitted;
        bool proofValid;
        address prover;
    }

    mapping(uint => ChallengeInfo) public challenges;

    event ChallengeCreated(uint indexed challengeId, address indexed challenger, string fileHash, string nonce);
    event ProofSubmitted(uint indexed challengeId, address indexed prover, bool valid);

    constructor() {
    }

    function init(address verifierAddress) external {
        verifier = IVerifier(verifierAddress);
    }

    // 发起挑战，challenger指定要挑战的文件哈希和nonce（随机数）
    function createChallenge(string memory fileHash, string memory nonce, string memory cid, address prover) external returns (uint) {
        challengeCounter++;
        challenges[challengeCounter] = ChallengeInfo({
            challenger: msg.sender,
            fileHash: fileHash,
            nonce: nonce,
            cid: cid,
            timestamp: block.timestamp,
            proofSubmitted: false,
            proofValid: false,
            prover: prover
        });

        emit ChallengeCreated(challengeCounter, msg.sender, fileHash, nonce);
        return challengeCounter;
    }

    // 矿工提交零知识证明
    function submitProof(
        uint challengeId, uint[2] calldata _pA, uint[2] calldata _pB0, uint[2] calldata _pB1, uint[2] calldata _pC, uint[3] calldata _pubSignals
    ) external {
        ChallengeInfo storage ch = challenges[challengeId];
        require(ch.challenger != address(0), "no challenge");
        require(ch.prover == msg.sender, "not your challenge");
        require(!ch.proofSubmitted, "Proof is submit");
        
        uint[2] memory a = [_pA[0], _pA[1]];
        uint[2][2] memory b = [_pB0, _pB1];
        uint[2] memory c = [_pC[0], _pC[1]];
        uint[3] memory input = [_pubSignals[0], _pubSignals[1], _pubSignals[2]];
        bool valid = verifier.verifyProof(a, b, c, input);

        ch.proofSubmitted = true;
        ch.proofValid = valid;
        ch.prover = msg.sender;

        emit ProofSubmitted(challengeId, msg.sender, valid);
    }

    // 查询挑战详情
    function getChallengeInfo(uint challengeId) external view returns (
        address challenger,
        string memory fileHash,
        string memory nonce,
        string memory cid,
        uint timestamp,
        bool proofSubmitted,
        bool proofValid,
        address prover
    ) {
        ChallengeInfo storage ch = challenges[challengeId];
        return (
            ch.challenger,
            ch.fileHash,
            ch.nonce,
            ch.cid,
            ch.timestamp,
            ch.proofSubmitted,
            ch.proofValid,
            ch.prover
        );
    }

    function getChallengesByUser()  external view returns (
        uint challengeId,
        address challenger,
        string memory fileHash,
        string memory nonce,
        string memory cid,
        uint timestamp,
        bool proofSubmitted,
        bool proofValid,
        address prover
    ) {
        for (uint i = 1; i <= challengeCounter; i++) {
            if (challenges[i].prover == msg.sender && !challenges[i].proofSubmitted) {
                ChallengeInfo storage ch = challenges[i];
                return (
                    i,
                    ch.challenger,
                    ch.fileHash,
                    ch.nonce,
                    ch.cid,
                    ch.timestamp,
                    ch.proofSubmitted,
                    ch.proofValid,
                    ch.prover
            );
            }
        }

        // 如果没找到，返回默认值（如全 0）
        return (0, address(0), "", "", "", 0, false, false, address(0));
    }

}
