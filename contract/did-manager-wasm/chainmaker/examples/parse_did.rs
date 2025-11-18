use did_demo::did_tool::{DidDocument, VerificationMethod};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let did_bytes = include_str!("../assets/did.jsonld");
    let did: DidDocument = serde_json::from_str(did_bytes)?;

    // 查找验证方法
    let verification_method: &VerificationMethod = did.verification_methods.iter()
    .filter(|item|item.id == "did:example:foo#key4".to_string())
    .nth(0).expect("get verification_method fail");

    println!("verification_method:\n{}", verification_method.public_key_hex.clone().expect("not have public key"));

    Ok(())
}
