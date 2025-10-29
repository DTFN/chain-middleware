use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::collections::{BTreeMap, HashMap};

/**
 * 解析DID
 */
pub fn decode_did(content: &str) -> Result<DidDocument, Box<dyn std::error::Error>> {
    let did: DidDocument = serde_json::from_str(content)?;
    return Ok(did);
}

/**
 * 解析VC
 */
pub fn decode_vc(content: &str) -> Result<VerifiableCredential, Box<dyn std::error::Error>> {
    let vc: VerifiableCredential = serde_json::from_str(content)?;
    return Ok(vc);
}

// 跨链消息
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CrossChainMessage {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub origin: Option<VerifiableCredential>,

    #[serde(skip_serializing_if = "Option::is_none")]
    pub target: Option<VerifiableCredential>,
}

// DID
/// DID文档根结构体
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DidDocument {
    /// DID标识
    pub id: String,

    /// 验证方法列表
    #[serde(rename = "verificationMethod")]
    pub verification_methods: Vec<VerificationMethod>,

    /// 断言方法引用列表
    #[serde(rename = "assertionMethod")]
    pub assertion_methods: Option<Vec<String>>,

    /// 认证方法引用列表
    pub authentication: Option<Vec<String>>,

    /// 能力委派方法引用列表
    #[serde(rename = "capabilityDelegation")]
    pub capability_delegation: Option<Vec<String>>,

    /// 能力调用方法引用列表
    #[serde(rename = "capabilityInvocation")]
    pub capability_invocation: Option<Vec<String>>,
}

/// JWK公钥上下文定义
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PublicKeyJwkContext {
    /// 标识
    #[serde(rename = "@id")]
    pub id: String,
    /// 类型（指定为JSON）
    #[serde(rename = "@type")]
    pub type_: String,
}

/// 验证方法结构体（支持多种密钥类型）
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VerificationMethod {
    /// 验证方法标识
    pub id: String,

    /// 验证方法类型
    #[serde(rename = "type")]
    pub type_: String,

    /// 控制器DID
    pub controller: String,

    /// JWK格式公钥（可选，根据类型存在）
    #[serde(rename = "publicKeyJwk", skip_serializing_if = "Option::is_none")]
    pub public_key_jwk: Option<PublicKeyJwk>,

    /// base58格式公钥（可选，根据类型存在）
    #[serde(rename = "publicKeyBase58", skip_serializing_if = "Option::is_none")]
    pub public_key_base58: Option<String>,

    /// hex格式公钥（可选，根据类型存在）
    #[serde(rename = "publicKeyHex", skip_serializing_if = "Option::is_none")]
    pub public_key_hex: Option<String>,

    /// hex格式公钥（可选，根据类型存在）
    #[serde(rename = "address", skip_serializing_if = "Option::is_none")]
    pub address: Option<String>,
}

/// JWK格式公钥结构体
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PublicKeyJwk {
    /// 密钥类型
    pub kty: String,

    /// RSA模数（仅RSA密钥存在）
    pub n: Option<String>,

    /// RSA公开指数（仅RSA密钥存在）
    pub e: Option<String>,
    // 可以根据需要添加其他JWK字段（如crv, x, y等用于椭圆曲线密钥）
}

// 下面是VC
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct VerifiableCredential {
    #[serde(rename = "@context")]
    pub context: Vec<String>,

    pub id: Option<String>,

    #[serde(rename = "type")]
    pub type_: Option<Vec<String>>,

    pub issuer: String,

    #[serde(rename = "issuanceDate")]
    pub issuance_date: String,

    #[serde(rename = "credentialSubject")]
    pub credential_subject: CredentialSubject,

    #[serde(skip_serializing_if = "Option::is_none")]
    pub proof: Option<Proof>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct CredentialSubject {
    pub chain_rid: String,
    pub contract_address: String,
    pub contract_func: String,
    pub contract_type: String,
    pub func_params: HashMap<String, String>,
    pub param_name: String,
    pub gateway_id: String,
    pub resource_name: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Proof {
    // ed25519
    #[serde(rename = "type")]
    pub type_: Option<String>,
    pub created: Option<String>,
    #[serde(rename = "verificationMethod")]
    pub verification_method: Option<String>,
    #[serde(rename = "proofPurpose")]
    pub proof_purpose: Option<String>,
    pub jws: Option<String>,
    pub signature: Option<String>,

    // eth
    pub address: Option<String>,
    #[serde(rename = "contentHash")]
    pub content_hash: Option<String>,
    pub r: Option<String>,
    pub s: Option<String>,
    pub v: Option<u32>,
}

#[test]
pub fn testDecodeVcs() {
    let vcs_str = "{\"vcs\":\"{\\\"target\\\":{\\\"@context\\\":[\\\"https://www.w3.org/2018/credentials/v1\\\"],\\\"credentialSubject\\\":{\\\"chain_rid\\\":\\\"chainmaker001\\\",\\\"contract_address\\\":\\\"busiWasmCenter13120\\\",\\\"contract_func\\\":\\\"erc20MintVcs\\\",\\\"contract_type\\\":\\\"4\\\",\\\"func_params\\\":{\\\"amount\\\":\\\"10\\\",\\\"to\\\":\\\"did:1\\\"},\\\"gateway_id\\\":\\\"0\\\",\\\"param_name\\\":\\\"content\\\",\\\"resource_name\\\":\\\"chain1:contract01\\\"},\\\"id\\\":\\\"https://uni.example.com/credentials/12345\\\",\\\"issuanceDate\\\":\\\"2025-05-20\\\",\\\"issuer\\\":\\\"https://uni.example.com/organization/1\\\",\\\"proof\\\":{\\\"address\\\":\\\"000000000000000000000000bd5d9f6b1b70cfd36c90aa11a8830a43062ab4ea\\\",\\\"contentHash\\\":\\\"ba44b74e24968241245b9c91fdb5cc73a6b187c366c70ec115861dc96c9352c5\\\",\\\"created\\\":null,\\\"jws\\\":null,\\\"proofPurpose\\\":null,\\\"r\\\":\\\"286df6b6299e7f83439391518b56ceee6711f5eca2fafb1c7b406d49d78dbbf6\\\",\\\"s\\\":\\\"1305c091c3569c69c01ae576a9e56e798bbc9542b6f94de9cbddbf82d88f0a69\\\",\\\"signature\\\":null,\\\"type\\\":null,\\\"v\\\":28,\\\"verificationMethod\\\":\\\"did:example:foo#key4\\\"},\\\"type\\\":[\\\"VerifiableCredential\\\",\\\"UniversityDegreeCredential\\\"]}}\"}";
    let ccm: CrossChainMessage =
        serde_json::from_slice(vcs_str.as_bytes()).expect("vcs decode fail");
}
