use busi_left::did_tool::{CrossChainMessage, DidDocument, VerificationMethod};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let vcs_bytes = include_str!("../assets/crossChainMessage.json");
    let vcs: CrossChainMessage = serde_json::from_str(vcs_bytes)?;

    println!("{:?}", vcs);

    Ok(())
}
