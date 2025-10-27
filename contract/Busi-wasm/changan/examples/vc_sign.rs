use busi_left::did_tool::{Proof, VerifiableCredential};
use chrono::Local;
use ed25519_dalek::ed25519::signature::{Keypair, SignerMut};
use ed25519_dalek::{Signature, SigningKey, VerifyingKey};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    sign_test();

    // 构造签名字符串
    let message = include_bytes!("../assets/vc-raw.jsonld");
    let mut vc: VerifiableCredential = serde_json::from_slice(message)?;
    let vc_content = serde_json::to_string(&vc).expect("");

    let secret_vec =
        hex::decode("162642a20c4175dac1fe26b66d111716ab023d5a445de6de1a83474c8b0b310c").unwrap();
    let verify_vec =
        hex::decode("702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a").unwrap();
    let secret_bytes: &[u8; 32] = &secret_vec[..32].try_into().expect("msg");
    let verify_bytes: &[u8; 32] = &verify_vec[..32].try_into().expect("msg");
    let mut signing_key = SigningKey::from_bytes(secret_bytes);
    let signature: Signature = signing_key.sign(vc_content.as_bytes());
    println!("signature:\n{}", signature);

    // 第一次校验
    let verify_key: VerifyingKey = VerifyingKey::from_bytes(&verify_bytes).unwrap();
    let verify_result = verify_key.verify_strict(vc_content.as_bytes(), &signature);
    println!(
        "verify1-conent:\n{:?}",
        String::from_utf8(vc_content.as_bytes().to_vec()).expect("")
    );
    println!("verify1:{:?}", verify_result);

    // 构造VC
    let now = Local::now();
    let formated_time = now.to_rfc3339();
    let proof = Proof {
        type_: "JsonWebSignature2020".to_string(),
        created: formated_time,
        verification_method: "did:example:foo#key5".to_string(),
        proof_purpose: "assertionMethod".to_string(),
        jws: Option::None,
        signature: Option::Some(hex::encode(signature.to_bytes())),
    };
    vc.proof = Option::Some(proof);

    let vc_str = serde_json::to_string(&vc)?;
    println!("vc:\n{}", vc_str);

    // 校验VC
    let mut vc: VerifiableCredential =
        serde_json::from_slice(vc_str.as_bytes()).expect("vc decode fail");
    let proof = vc.proof.expect("not have proof");
    vc.proof = Option::None;
    let signature_hex = proof.signature.expect("not have signature");
    let signature_bytes = hex::decode(signature_hex.as_str()).expect("signature format fail");
    let signature = Signature::from_bytes(&signature_bytes.try_into().expect(""));
    let verify_vec =
        hex::decode("702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a").unwrap();
    let verify_bytes: &[u8; 32] = &verify_vec[..32].try_into().expect("msg");
    let verify_key: VerifyingKey = VerifyingKey::from_bytes(&verify_bytes).unwrap();
    let verify_result =
        verify_key.verify_strict(serde_json::to_string(&vc).expect("").as_bytes(), &signature);
    println!(
        "verify2-conent:\n{:?}",
        serde_json::to_string(&vc).expect("")
    );
    print!("verify2:{:?}", verify_result);

    Ok(())
}

fn sign_test() {
    let s = "{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain1\",\"contract_address\":\"contract01\",\"contract_func\":\"pay\",\"contract_type\":\"1\",\"func_params\":\"[{\\\"string\\\":\\\"c1\\\"},{\\\"string\\\":\\\"121.11\\\"}]\",\"gateway_id\":\"0\",\"param_name\":\"content\",\"reource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"]}";
    let verify_vec =
        hex::decode("702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a").unwrap();
    let verify_bytes: &[u8; 32] = &verify_vec[..32].try_into().expect("msg");
    let verify_key: VerifyingKey = VerifyingKey::from_bytes(&verify_bytes).unwrap();

    let signature_hex = "7112e69aadb1c24c9633ba5515da788dc3074df2f8e87b95bd1dbb2dd65269ee411bf0f659d87acccf9b83207a4998204229dc025b45de7662741d3b5d55ec09";
    let signature_bytes = hex::decode(signature_hex).expect("signature format fail");
    let signature = Signature::from_bytes(&signature_bytes.try_into().expect(""));

    let verify_result = verify_key.verify_strict(s.as_bytes(), &signature);
    println!("{}", s.len());
    println!("verify_result {:?}", verify_result);
    println!("{}", hex::encode(s));
}
