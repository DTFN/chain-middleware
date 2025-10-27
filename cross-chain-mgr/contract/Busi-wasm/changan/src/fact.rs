use std::collections::BTreeMap;
use std::error::Error;

use contract_sdk_rust::easycodec::*;
use contract_sdk_rust::sim_context;
use contract_sdk_rust::sim_context::SimContext;
use ed25519_dalek::Signature;
use ed25519_dalek::SigningKey;
use ed25519_dalek::VerifyingKey;
use ed25519_dalek::ed25519::signature::{Keypair, SignerMut};
use hex_literal::hex;
use k256::EncodedPoint;
use k256::ecdsa::RecoveryId;
use k256::sha2::Sha256;
use serde_json::Serializer;
use serde_json::Value;
use sha3::Digest;
use sha3::Keccak256;

use crate::did_tool::CrossChainMessage;
use crate::did_tool::DidDocument;
use crate::did_tool::VerifiableCredential;
use crate::did_tool::VerificationMethod;

// 安装合约时会执行此方法，必须
#[unsafe(no_mangle)]
pub extern "C" fn init_contract() {
    // 安装时的业务逻辑，内容可为空
    sim_context::log("init_contract");
}

// 升级合约时会执行此方法，必须
#[unsafe(no_mangle)]
pub extern "C" fn upgrade() {
    // 升级时的业务逻辑，内容可为空
    sim_context::log("upgrade");
    let ctx = &mut sim_context::get_sim_context();
    ctx.ok("upgrade success".as_bytes());
}

struct Fact {
    file_hash: String,
    file_name: String,
    time: i32,
    ec: EasyCodec,
}

impl Fact {
    fn new_fact(file_hash: String, file_name: String, time: i32) -> Fact {
        let mut ec = EasyCodec::new();
        ec.add_string("file_hash", file_hash.as_str());
        ec.add_string("file_name", file_name.as_str());
        ec.add_i32("time", time);
        Fact {
            file_hash,
            file_name,
            time,
            ec,
        }
    }

    fn get_emit_event_data(&self) -> Vec<String> {
        let mut arr: Vec<String> = Vec::new();
        arr.push(self.file_hash.clone());
        arr.push(self.file_name.clone());
        arr.push(self.time.to_string());
        arr
    }

    fn to_json(&self) -> String {
        self.ec.to_json()
    }

    fn marshal(&self) -> Vec<u8> {
        self.ec.marshal()
    }

    fn unmarshal(data: &Vec<u8>) -> Fact {
        let ec = EasyCodec::new_with_bytes(data);
        Fact {
            file_hash: ec.get_string("file_hash").unwrap(),
            file_name: ec.get_string("file_name").unwrap(),
            time: ec.get_i32("time").unwrap(),
            ec,
        }
    }
}

// echo 测试方法
#[unsafe(no_mangle)]
pub extern "C" fn echo() {
    let ctx = &mut sim_context::get_sim_context();
    let vcs_str = ctx.arg_as_utf8_str("content");
    ctx.emit_event("result", &vec![format!("{}", vcs_str)]);

    // 返回查询结果
    ctx.ok(vcs_str.as_bytes());
    // ctx.log(&json_str);
}

// 测试验签的方法
#[unsafe(no_mangle)]
pub extern "C" fn test_eth_verify_wasm() {
    let ctx = &mut sim_context::get_sim_context();

    let content = "123456";
    let contentHash = "2912723b3ed60c075b271f075d881d82fa5de112b6c25f7dfa4cab85de25045a";
    let address = "bd5d9f6b1b70cfd36c90aa11a8830a43062ab4ea";
    let r = "68c33515e77b9f80404e57512ea112156a1ff058a364f7e96671bf62ca8a63eb";
    let s = "4ccc92170b52985d35be0478786b08a92d4a81bb5186d69ebe4075a30c48b775";
    let v = 27;

    // 测试
    let mut digest = Keccak256::new();
    digest.update(format!(
        "\x19Ethereum Signed Message:\n{}{}",
        content.len(),
        content
    ));
    let sig = hex::decode(format!("{r}{s}")).expect("gen sig fail");
    let sig = k256::ecdsa::Signature::try_from(sig.as_slice()).unwrap();
    let recid: RecoveryId = RecoveryId::new(v % 2 == 0, false);
    let pk = k256::ecdsa::VerifyingKey::recover_from_digest(digest, &sig, recid).unwrap();
    let pk = pk.to_encoded_point(false);
    let pk = pk.as_bytes();

    // 计算address
    let pubkey = &pk[1..];
    let mut hasher = Keccak256::new();
    hasher.update(pubkey);
    let pubkey_hash = hasher.finalize();
    let address = hex::encode(&pubkey_hash[12..]);

    // 返回查询结果
    ctx.ok(address.as_bytes());
}

// 测试方法
#[unsafe(no_mangle)]
pub extern "C" fn getDIDDetails() {
    let ctx = &mut sim_context::get_sim_context();
    let did = ctx.arg_as_utf8_str("did");
    ctx.emit_event("result", &vec![format!("{}", did)]);

    // 获取did文档
    let mut args = EasyCodec::new();
    args.add_string("did", did.as_str());
    let result = ctx.call_contract("DID_MANAGER", "getDIDDetails", args);

    // 解析did文档
    let did_doc = result
        .and_then(|item| String::from_utf8(item).map_err(|_| 1))
        .unwrap_or("".to_string());

    // 返回查询结果
    ctx.ok(format!("did_doc: {:?}", did_doc).as_bytes());
}

// verifyVc
#[unsafe(no_mangle)]
pub extern "C" fn verifyVcs() {
    // 获取上下文
    let ctx = &mut sim_context::get_sim_context();

    // 获取VC
    let vc = match get_used_vc(ctx) {
        Err(msg) => {
            ctx.error(msg.to_string().as_str());
            return;
        }
        Ok(vc) => {
            ctx.ok("verify success".as_bytes());
        }
    };
}

// ERC20
#[unsafe(no_mangle)]
pub extern "C" fn erc20MintVcs() {
    // 获取上下文
    let ctx = &mut sim_context::get_sim_context();

    // 获取VC
    let vc = match get_used_vc(ctx) {
        Err(msg) => {
            ctx.error(msg.to_string().as_str());
            return;
        }
        Ok(vc) => vc,
    };

    // 解析参数
    let credential_subject = &vc.credential_subject;
    let to: &String = credential_subject
        .func_params
        .get("to")
        .expect("[to] not in params");
    let amount_str = credential_subject
        .func_params
        .get("amount")
        .expect("[amount] not in params");
    let amount: u32 = amount_str.parse().expect("[amount] not a number");

    // 业务逻辑-铸造代币
    let erc20_field_name = "erc20";
    let old_amount = ctx
        .get_state(erc20_field_name, didEncode(to).as_str())
        .map(|item| String::from_utf8(item).expect("parse old amount to string fail"))
        .map(|item| {
            if item == "" {
                return 0;
            } else {
                return item.parse().expect("parse old amount to number fail");
            }
        })
        .unwrap_or(0);
    let new_amount = old_amount + amount;
    ctx.put_state(
        erc20_field_name,
        didEncode(to).as_str(),
        &new_amount.to_string().into_bytes(),
    );

    // 发送事件
    sendCrossChainEventIfNeed(ctx);

    // 返回结果
    ctx.ok(format!("{}", new_amount).as_bytes());
}

#[unsafe(no_mangle)]
pub extern "C" fn erc20TransferVcs() {
    // 获取上下文
    let ctx = &mut sim_context::get_sim_context();

    // 获取VC
    let vc = match get_used_vc(ctx) {
        Err(msg) => {
            ctx.error(msg.to_string().as_str());
            return;
        }
        Ok(vc) => vc,
    };

    // 解析参数
    let credential_subject = &vc.credential_subject;
    let from: &String = credential_subject
        .func_params
        .get("from")
        .expect("[from] not in params");
    let to: &String = credential_subject
        .func_params
        .get("to")
        .expect("[to] not in params");
    let amount_str = credential_subject
        .func_params
        .get("amount")
        .expect("[amount] not in params");
    let amount: u32 = amount_str.parse().expect("[amount] not a number");

    // 业务逻辑-转移代币
    let erc20_field_name = "erc20";
    let from_old_amount = ctx
        .get_state(erc20_field_name, didEncode(from).as_str())
        .map(|item| String::from_utf8(item).expect("parse old amount to string fail"))
        .map(|item| {
            if item == "" {
                return 0;
            } else {
                return item.parse().expect("parse old amount to number fail");
            }
        })
        .unwrap_or(0);
    let to_old_amount = ctx
        .get_state(erc20_field_name, didEncode(to).as_str())
        .map(|item| String::from_utf8(item).expect("parse old amount to string fail"))
        .map(|item| {
            if item == "" {
                return 0;
            } else {
                return item.parse().expect("parse old amount to number fail");
            }
        })
        .unwrap_or(0);

    if from_old_amount <= amount {
        panic!("from ot have enough ,{}", amount);
    }

    let from_amount = from_old_amount - amount;
    let to_amount = to_old_amount + amount;
    ctx.put_state(
        erc20_field_name,
        didEncode(from).as_str(),
        &from_amount.to_string().into_bytes(),
    );
    ctx.put_state(
        erc20_field_name,
        didEncode(to).as_str(),
        &to_amount.to_string().into_bytes(),
    );

    // 发送事件
    sendCrossChainEventIfNeed(ctx);

    // 返回结果
    ctx.ok(format!("from: {}, to: {}", from_amount, to_amount).as_bytes());
}

#[unsafe(no_mangle)]
pub extern "C" fn erc20GetBalanceVcs() {
    // 获取上下文
    let ctx = &mut sim_context::get_sim_context();

    // 获取VC
    let vc = match get_used_vc(ctx) {
        Err(msg) => {
            ctx.error(msg.to_string().as_str());
            return;
        }
        Ok(vc) => vc,
    };

    // 解析参数
    let credential_subject = &vc.credential_subject;
    let account: &String = credential_subject
        .func_params
        .get("account")
        .expect("[account] not in params");

    // 业务逻辑-查询余额
    let erc20_field_name = "erc20";
    let account_amount = ctx
        .get_state(erc20_field_name, didEncode(account).as_str())
        .map(|item| String::from_utf8(item).expect("parse old amount to string fail"))
        .map(|item| {
            if item == "" {
                return 0;
            } else {
                return item.parse().expect("parse old amount to number fail");
            }
        })
        .unwrap_or(0);

    // 返回结果
    ctx.ok(format!("{}", account_amount).as_bytes());
}

#[unsafe(no_mangle)]
pub extern "C" fn erc20GetBalance() {
    // 获取上下文
    let ctx = &mut sim_context::get_sim_context();

    // 获取VC
    let account = ctx.arg_as_utf8_str("account");

    // 业务逻辑-查询余额
    let erc20_field_name = "erc20";
    let account_amount = ctx
        .get_state(erc20_field_name, didEncode(&account).as_str())
        .map(|item| String::from_utf8(item).expect("parse old amount to string fail"))
        .map(|item| {
            if item == "" {
                return 0;
            } else {
                return item.parse().expect("parse old amount to number fail");
            }
        })
        .unwrap_or(0);

    // 返回结果
    ctx.ok(format!("{}", account_amount).as_bytes());
}

// ERC721
#[unsafe(no_mangle)]
pub extern "C" fn erc721MintVcs() {
    // 获取上下文
    let ctx = &mut sim_context::get_sim_context();

    // 获取VC
    let vc = match get_used_vc(ctx) {
        Err(msg) => {
            ctx.error(msg.to_string().as_str());
            return;
        }
        Ok(vc) => vc,
    };

    // 解析参数
    let credential_subject = &vc.credential_subject;
    let to: &String = credential_subject
        .func_params
        .get("to")
        .expect("[account] not in params");

    // 业务逻辑-查询余额
    let erc721_field_name = "erc721";
    let erc721_index_name = "erc721-index";
    let mut erc721_index = ctx
        .get_state(erc721_index_name, erc721_index_name)
        .map(|item| String::from_utf8(item).expect("parse old amount to string fail"))
        .map(|item| {
            if item == "" {
                return 0;
            } else {
                return item.parse().expect("parse old amount to number fail");
            }
        })
        .unwrap_or(0);

    // 分配新的index
    erc721_index = erc721_index + 1;
    ctx.put_state(
        erc721_index_name,
        erc721_index_name,
        format!("{}", erc721_index).as_bytes(),
    );
    ctx.put_state(
        erc721_field_name,
        format!("{}", erc721_index).as_str(),
        to.as_bytes(),
    );

    // 返回结果
    ctx.ok(format!("{}", erc721_index).as_bytes());
}

#[unsafe(no_mangle)]
pub extern "C" fn erc721OwnerOfVcs() {
    // 获取上下文
    let ctx = &mut sim_context::get_sim_context();

    // 获取VC
    let vc = match get_used_vc(ctx) {
        Err(msg) => {
            ctx.error(msg.to_string().as_str());
            return;
        }
        Ok(vc) => vc,
    };

    // 解析参数
    let credential_subject = &vc.credential_subject;
    let account: &String = credential_subject
        .func_params
        .get("account")
        .expect("[account] not in params");
    let token_id: &String = credential_subject
        .func_params
        .get("tokenId")
        .expect("[tokenId] not in params");

    // 业务逻辑-查询余额
    let erc721_field_name = "erc721";
    let old_account = ctx
        .get_state(erc721_field_name, token_id)
        .unwrap_or("".as_bytes().into());
    let old_account = String::from_utf8(old_account).expect("get account fail");

    // 返回结果
    ctx.ok(format!("{}", account == &old_account).as_bytes());
}

fn didEncode(did: &String) -> String {
    hex::encode(did)
}

fn didDecode(didEncoded: &String) -> String {
    String::from_utf8(hex::decode(didEncoded).expect("did decode fail"))
        .expect("parse code to String fail")
}

fn sendCrossChainEventIfNeed(ctx: &mut impl SimContext) {
    // 获取传入参数
    let vcs_str = ctx.arg_as_utf8_str("vcs");
    let ccm: CrossChainMessage =
        serde_json::from_slice(vcs_str.as_bytes()).expect("vc decode fail");

    // 验签
    if ccm.origin.is_some() && ccm.target.is_some() {
        let target_vc = ccm.target.expect("origtargetin not exists");
        let ccm = CrossChainMessage {
            origin: None,
            target: Some(target_vc),
        };

        // todo (优雅排序) 排序输出
        let ccm_str = serde_json::to_string(&ccm).expect("target_vc to string fail");
        let mut ccm_map: Value = serde_json::from_str(&ccm_str).expect("vc parse fail");
        let ccm_map = match ccm_map.as_object_mut() {
            Some(item) => item,
            None => &mut serde_json::Map::new(),
        };
        let ccm_str = serde_json::to_string(&ccm_map).expect("");

        ctx.emit_event("CROSS_CHAIN_VC", &vec![ccm_str]);
    }
}

// save 保存存证数据
#[unsafe(no_mangle)]
pub extern "C" fn save() {
    // 获取上下文
    let ctx = &mut sim_context::get_sim_context();

    // 获取传入参数
    let vcs_str = ctx.arg_as_utf8_str("vcs");
    let ccm: CrossChainMessage =
        serde_json::from_slice(vcs_str.as_bytes()).expect("vc decode fail");
    let mut vc: &mut VerifiableCredential = &mut ccm.origin.expect("origin not exists");
    let credential_subject = &vc.credential_subject;

    // 验签
    let verify_result = valide_vc(&mut vc, ctx);
    // 不使用asset,留下调试空间
    // assert!(verify_result.is_ok());
    if verify_result.is_err() {
        ctx.ok(format!("{:?}", verify_result).as_bytes());
        return;
    }

    // 业务逻辑

    // 事件
    ctx.emit_event("result", &vec![format!("{:?}", verify_result)]);

    // 序列化后存储
    // ctx.put_state(
    //     "fact_ec",
    //     fact.file_hash.as_str(),
    //     fact.marshal().as_slice(),
    // );
}

// 获取需要的VC并验签
fn get_used_vc(ctx: &mut impl SimContext) -> Result<VerifiableCredential, Box<dyn Error>> {
    // 获取传入参数
    let vcs_str = ctx.arg_as_utf8_str("vcs");
    let ccm: Result<CrossChainMessage, serde_json::Error> =
        serde_json::from_slice(vcs_str.as_bytes());
    if let Ok(ccm) = ccm {
        // 验签
        if ccm.origin.is_some() {
            let mut origin_vc = ccm.origin.expect("origin not exists");
            let verify_result = valide_vc(&mut origin_vc, ctx);
            if verify_result.is_err() {
                return Err(format!("origin valid fail, {:?}", verify_result).into());
            }
            return Ok(origin_vc);
        }
        if ccm.target.is_some() {
            let mut target_vc = ccm.target.expect("origtargetin not exists");
            let verify_result = valide_vc(&mut target_vc, ctx);
            if verify_result.is_err() {
                return Err(format!("target valid fail, {:?}", verify_result).into());
            }
            return Ok(target_vc);
        }

        Err("origin and target at least have one".into())
    } else {
        return Err(format!("parse error, {:?}", ccm).into());
    }
}

/**
 * 校验VC
 */
fn valide_vc(
    vc: &mut VerifiableCredential,
    ctx: &mut impl SimContext,
) -> Result<(), Box<dyn Error>> {
    // 1.提取proof
    let proof: crate::did_tool::Proof = vc.proof.take().expect("not have proof");
    let verification_method = &proof
        .verification_method
        .clone()
        .expect("verification_method must need");

    // 2.查询did文档
    let vms: Vec<&str> = verification_method.splitn(2, "#").collect();
    let did = vms.get(0).ok_or(format!(
        "get did fail, verification_method: {}",
        verification_method
    ))?;
    let mut args = EasyCodec::new();
    args.add_string("did", did);
    let result = ctx.call_contract("DID_MANAGER", "getDIDDetails", args);
    let did_doc = result.expect("did doc parse fail");
    let did_doc = String::from_utf8(did_doc)?;

    // 3.解析DID
    let did: DidDocument = serde_json::from_str(&did_doc)
        .map_err(|e| format!("parse did fail, err: {:?}, did_doc: {}", e, did_doc))?;
    let verification_method: &VerificationMethod = did
        .verification_methods
        .iter()
        .filter(|item| &item.id == verification_method)
        .nth(0)
        .ok_or(format!(
            "get verification_method fail, verification_method: {verification_method}, did: {:?}",
            did
        ))?;

    // 4.序列化vc
    // todo (优雅排序) 排序输出
    let vc_str = serde_json::to_string(&vc).expect("");
    let mut vc_map: Value = serde_json::from_str(&vc_str).expect("vc parse fail");
    let vc_map = match vc_map.as_object_mut() {
        Some(item) => item,
        None => &mut serde_json::Map::new(),
    };
    let vc_str = serde_json::to_string(&vc_map).expect("");

    // 5.验证签名(eth+ed25519)
    let verify_result = match verification_method.type_.as_str() {
        "Ed25519VerificationKey2018" => {
            valide_by_ed25519(ctx, &proof, &verification_method, &vc_str)
        }
        "EcdsaSecp256k1VerificationKey2019" => {
            valide_by_eth(ctx, &proof, &verification_method, &vc_str)
        }
        _ => return Err(format!("not support type {}", verification_method.type_).into()),
    };

    return verify_result;
}

fn valide_by_ed25519(
    ctx: &mut impl SimContext,
    proof: &crate::did_tool::Proof,
    verification_method: &VerificationMethod,
    message: &String,
) -> Result<(), Box<dyn Error>> {
    let signature_hex = proof.signature.clone().expect("not have signature");
    ctx.emit_event("signature_hex", &vec![format!("{}", signature_hex)]);
    let signature_bytes = hex::decode(signature_hex.as_str())?;
    let signature = Signature::from_bytes(
        &signature_bytes
            .try_into()
            .map_err(|e| format!("parse signature fail"))?,
    );
    ctx.emit_event(
        "public_key_hex",
        &vec![format!("{:?}", verification_method.public_key_hex)],
    );
    let verify_vec = hex::decode(
        verification_method
            .public_key_hex
            .clone()
            .expect("not have public key"),
    )
    .unwrap();
    let verify_bytes: &[u8; 32] = &verify_vec[..32].try_into().expect("msg");
    let verify_key: VerifyingKey = VerifyingKey::from_bytes(&verify_bytes).unwrap();

    // eth 验签
    let verify_result = verify_key.verify_strict(message.as_bytes(), &signature);

    return Ok(verify_result?);
}

fn valide_by_eth(
    ctx: &mut impl SimContext,
    proof: &crate::did_tool::Proof,
    verification_method: &VerificationMethod,
    message: &String,
) -> Result<(), Box<dyn Error>> {
    let address_in_did = verification_method
        .address
        .clone()
        .ok_or("address not exists")?;
    let r = proof.r.clone().ok_or("'r' not exists")?;
    let s = proof.s.clone().ok_or("'s' not exists")?;
    let v = proof.v.ok_or("'v' not exists")?;

    // 测试
    let mut digest = Keccak256::new();
    digest.update(format!(
        "\x19Ethereum Signed Message:\n{}{}",
        message.len(),
        message
    ));
    let sig = hex::decode(format!("{r}{s}")).expect("gen sig fail");
    let sig = k256::ecdsa::Signature::try_from(sig.as_slice()).unwrap();
    let recid: RecoveryId = RecoveryId::new(v % 2 == 0, false);
    let pk = k256::ecdsa::VerifyingKey::recover_from_digest(digest, &sig, recid).unwrap();
    let pk = pk.to_encoded_point(false);
    let pk = pk.as_bytes();

    // 计算address
    let pubkey = &pk[1..];
    let mut hasher = Keccak256::new();
    hasher.update(pubkey);
    let pubkey_hash = hasher.finalize();
    let address = hex::encode(&pubkey_hash[12..]);

    // 格式化address
    let address_in_did_bytes = hex::decode(&address_in_did)?;
    let address_bytes = hex::decode(&address)?;
    let address_in_did_bytes = address_in_did_bytes[address_in_did_bytes.len() - 20..].to_vec();
    let address_bytes = address_bytes[address_bytes.len() - 20..].to_vec();
    if address_in_did_bytes == address_bytes {
        return Ok(());
    } else {
        return Err(format!(
            "address not match, in did is [{}], in proof is [{}]",
            address_in_did, address
        )
        .into());
    }
}

fn is_address_equal() {

}

// find_by_file_hash 根据file_hash查询存证数据
#[unsafe(no_mangle)]
pub extern "C" fn find_by_file_hash() {
    // 获取上下文
    let ctx = &mut sim_context::get_sim_context();

    // 获取传入参数
    let file_hash = ctx.arg_as_utf8_str("file_hash");

    // 校验参数
    if file_hash.len() == 0 {
        ctx.log("file_hash is null");
        ctx.ok("".as_bytes());
        return;
    }

    // 查询
    let r = ctx.get_state("fact_ec", &file_hash);

    // 校验返回结果
    if r.is_err() {
        ctx.log("get_state fail");
        ctx.error("get_state fail");
        return;
    }
    let fact_vec = r.unwrap();
    if fact_vec.len() == 0 {
        ctx.log("None");
        ctx.ok("".as_bytes());
        return;
    }

    // 查询
    let r = ctx.get_state("fact_ec", &file_hash).unwrap();
    let fact = Fact::unmarshal(&r);
    let json_str = fact.to_json();

    // 返回查询结果
    ctx.ok(json_str.as_bytes());
    ctx.log(&json_str);
}

#[unsafe(no_mangle)]
pub extern "C" fn how_to_use_iterator() {
    let ctx = &mut sim_context::get_sim_context();

    // 构造数据
    ctx.put_state("key1", "field1", "val".as_bytes());
    ctx.put_state("key1", "field2", "val".as_bytes());
    ctx.put_state("key1", "field23", "val".as_bytes());
    ctx.put_state("key1", "field3", "val".as_bytes());
    // 使用迭代器，能查出来  field1，field2，field23 三条数据
    let r = ctx.new_iterator_with_field("key1", "field1", "field3");
    if r.is_ok() {
        let rs = r.unwrap();
        // 遍历
        while rs.has_next() {
            // 获取下一行值
            let row = rs.next_row().unwrap();
            let _key = row.get_string("key").unwrap();
            let _field = row.get_bytes("field");
            let _val = row.get_bytes("value");
            // do something
        }
        // 关闭游标
        rs.close();
    }

    ctx.put_state("key2", "field1", "val".as_bytes());
    ctx.put_state("key3", "field2", "val".as_bytes());
    ctx.put_state("key33", "field2", "val".as_bytes());
    ctx.put_state("key4", "field3", "val".as_bytes());
    // 能查出来 key2，key3，key33 三条数据
    let _ = ctx.new_iterator("key2", "key4");
    // 能查出来 key3，key33 两条数据
    let _ = ctx.new_iterator_prefix_with_key("key3");
    // 能查出来  field2，field23 三条数据
    let _ = ctx.new_iterator_prefix_with_key_field("key1", "field2");

    ctx.put_state_from_key("key5", "val".as_bytes());
    ctx.put_state_from_key("key56", "val".as_bytes());
    ctx.put_state_from_key("key6", "val".as_bytes());
    // 能查出来 key5，key56 两条数据
    let _ = ctx.new_iterator("key5", "key6");
}

// 自定义随机数实现
getrandom::register_custom_getrandom!(my_dummy_getrandom);

fn my_dummy_getrandom(dest: &mut [u8]) -> Result<(), getrandom::Error> {
    // 仅填充 0（非安全，仅用于纯验证场景）
    dest.fill(0);
    Ok(())
}

#[test]
pub fn test_eth_verify() {
    let content = "123456";
    let contentHash = "2912723b3ed60c075b271f075d881d82fa5de112b6c25f7dfa4cab85de25045a";
    let address = "000000000000000000000000bd5d9f6b1b70cfd36c90aa11a8830a43062ab4ea";
    let r = "68c33515e77b9f80404e57512ea112156a1ff058a364f7e96671bf62ca8a63eb";
    let s = "4ccc92170b52985d35be0478786b08a92d4a81bb5186d69ebe4075a30c48b775";
    let v = 27;

    let mut hasher = Keccak256::new();
    hasher.update(format!(
        "\x19Ethereum Signed Message:\n{}{}",
        content.len(),
        content
    ));
    println!("hash: {:?}", hex::encode(hasher.finalize()));

    // 测试
    let mut digest = Keccak256::new();
    digest.update(format!(
        "\x19Ethereum Signed Message:\n{}{}",
        content.len(),
        content
    ));
    let sig = hex::decode(format!("{r}{s}")).expect("gen sig fail");
    let sig = k256::ecdsa::Signature::try_from(sig.as_slice()).unwrap();
    let recid: RecoveryId = RecoveryId::new(v % 2 == 0, false);
    let pk = k256::ecdsa::VerifyingKey::recover_from_digest(digest, &sig, recid).unwrap();
    let pk = pk.to_encoded_point(false);
    let pk = pk.as_bytes();
    println!("pk: {:?}", hex::encode(pk));

    // 计算address
    let pubkey = &pk[1..];
    let mut hasher = Keccak256::new();
    hasher.update(pubkey);
    let pubkey_hash = hasher.finalize();

    println!("address: {address}");
    println!("pubkey_hash: {}", hex::encode(&pubkey_hash[12..]));
}

struct RecoveryTestVector {
    pk: [u8; 33],
    msg: &'static [u8],
    sig: [u8; 64],
    recid: RecoveryId,
}

const RECOVERY_TEST_VECTORS: &[RecoveryTestVector] = &[
    // Recovery ID 0
    RecoveryTestVector {
        pk: hex!("021a7a569e91dbf60581509c7fc946d1003b60c7dee85299538db6353538d59574"),
        msg: b"example message",
        sig: hex!(
            "ce53abb3721bafc561408ce8ff99c909f7f0b18a2f788649d6470162ab1aa032
                     3971edc523a6d6453f3fb6128d318d9db1a5ff3386feb1047d9816e780039d52"
        ),
        recid: RecoveryId::new(false, false),
    },
    // Recovery ID 1
    RecoveryTestVector {
        pk: hex!("036d6caac248af96f6afa7f904f550253a0f3ef3f5aa2fe6838a95b216691468e2"),
        msg: b"example message",
        sig: hex!(
            "46c05b6368a44b8810d79859441d819b8e7cdc8bfd371e35c53196f4bcacdb51
                     35c7facce2a97b95eacba8a586d87b7958aaf8368ab29cee481f76e871dbd9cb"
        ),
        recid: RecoveryId::new(true, false),
    },
];

#[test]
fn public_key_recovery() {
    for vector in RECOVERY_TEST_VECTORS {
        let digest = Sha256::new_with_prefix(vector.msg);
        let sig = k256::ecdsa::Signature::try_from(vector.sig.as_slice()).unwrap();
        let recid = vector.recid;
        let pk = k256::ecdsa::VerifyingKey::recover_from_digest(digest, &sig, recid).unwrap();
        assert_eq!(&vector.pk[..], EncodedPoint::from(&pk).as_bytes());
    }
}

#[test]
pub fn testDecode() {
    let vcs = "{\"origin\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chainmaker001\",\"contract_address\":\"busiWasmCenter6094\",\"contract_func\":\"erc20MintVcs\",\"contract_type\":\"4\",\"func_params\":{\"amount\":\"100\",\"to\":\"did:1\"},\"gateway_id\":\"0\",\"param_name\":\"content\",\"resource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"address\":\"000000000000000000000000bd5d9f6b1b70cfd36c90aa11a8830a43062ab4ea\",\"contentHash\":\"b684b019505fbcf8f7b5980a9332c8cf80bf85f432d85c44907819a1aff13466\",\"r\":\"42cf98b834da98c651fc02438b4cfe7dbf59a56fe459abe444ee101d4a3804f0\",\"s\":\"27c02821638193d7683e586d5db87240ccb007ec5c9e6f1dee4fc1bc48af6f8e\",\"v\":27,\"verificationMethod\":\"did:example:foo#key4\"},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"]},\"target\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain_bsc_01\",\"contract_address\":\"0x812feff56d3241a97235cb8ac83711180a49d675\",\"contract_func\":\"erc20MintVcs\",\"contract_type\":\"2\",\"func_params\":{\"amount\":\"10\",\"to\":\"did:1\"},\"gateway_id\":\"4\",\"param_name\":\"content\",\"resource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"address\":\"000000000000000000000000bd5d9f6b1b70cfd36c90aa11a8830a43062ab4ea\",\"contentHash\":\"8ca5a3fc046ec2f07de9ba3894befe37e281367806f45e78020ebedf8187d8f3\",\"r\":\"5dcdedfad67a7415e9367af818f691c26e75a1227cc07ad797c6442143e584ff\",\"s\":\"7e3649387937334d592dde54f1dd425db2dbd58d8fa2d7fae038248a6cd05e0e\",\"v\":27,\"verificationMethod\":\"did:example:foo#key4\"},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"]}}";
    let ccm: CrossChainMessage = serde_json::from_slice(vcs.as_bytes()).expect("vc decode fail");
}

#[test]
pub fn testDidEncode() {
    let did_raw = "did:rth:11111";
    let did_encoded = didEncode(&did_raw.to_string());
    println!("did_encoded: {}", did_encoded);

    let did = didDecode(&did_encoded);
    println!("did: {}", did);
    assert!(did == did_raw);
}

#[test]
pub fn testVerifyingKey() {
    let vk = "702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a";
    // let vk = "983421564f9133217448450067710d3f9ed0ec740edc8337e30245ba12678953";
    let verify_vec = hex::decode(vk).unwrap();
    let verify_bytes: &[u8; 32] = &verify_vec[..32].try_into().expect("msg");
    let verify_key: VerifyingKey = VerifyingKey::from_bytes(&verify_bytes).unwrap();
}

#[test]
pub fn testVerificationMethod() {
    let vm = "did:1#ed25519";
    let vms: Vec<&str> = vm.splitn(2, "#").collect();
    println!("{:?}", vms);
}

#[test]
pub fn testDidDecode() {
    let did_str = "{\"@context\":[\"https://www.w3.org/ns/did/v1\",\"https://w3c-ccg.github.io/lds-jws2020/contexts/lds-jws2020-v1.json\",{\"Ed25519VerificationKey2018\":\"https://w3id.org/security#Ed25519VerificationKey2018\"}],\"id\":\"did:lsid:1:19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d\",\"verificationMethod\":[{\"id\":\"did:lsid:1:19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d#key0\",\"type\":\"Ed25519VerificationKey2018\",\"controller\":\"did:lsid:1:19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d\",\"publicKeyHex\":\"19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d\"},{\"id\":\"did:lsid:1:19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d#key1\",\"type\":\"EcdsaSecp256k1VerificationKey2019\",\"controller\":\"did:lsid:1:19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d\",\"address\":\"0000000000000000000000007ad158f8244cbaab04d7bac7a1c03d4021f62515\"}],\"assertionMethod\":[\"did:lsid:1:19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d#key0\",\"did:lsid:1:19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d#key1\"],\"authentication\":[\"did:lsid:1:19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d#key0\",\"did:lsid:1:19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d#key1\"],\"capabilityDelegation\":[\"did:lsid:1:19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d#key0\",\"did:lsid:1:19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d#key1\"],\"capabilityInvocation\":[\"did:lsid:1:19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d#key0\",\"did:lsid:1:19d926138e5388c554403ecdf5ac10ab813d162d7c83a8177993c98b03b8e23d#key1\"]}";
    let did: DidDocument = serde_json::from_str(&did_str)
        .map_err(|e| format!("parse did fail, err: {:?}, did_doc: {}", e, did_str))
        .expect("");
    println!("did: {:?}", did);
}

#[test]
fn test_address_equal() {
    let address_in_did = "000000000000000000000000bd5d9f6b1b70cfd36c90aa11a8830a43062ab4ea";
    let address = "bd5d9f6b1b70cfd36c90aa11a8830a43062ab4ea";

    let address_in_did_bytes = hex::decode(&address_in_did).expect("");
    let address_bytes = hex::decode(&address).expect("");
    let address_in_did_bytes = address_in_did_bytes[address_in_did_bytes.len() - 20..].to_vec();
    let address_bytes = address_bytes[address_bytes.len() - 20..].to_vec();

    println!("{:?}", address_in_did_bytes);
    println!("{:?}", address_bytes);

    assert!(address_in_did_bytes == address_bytes)
}
