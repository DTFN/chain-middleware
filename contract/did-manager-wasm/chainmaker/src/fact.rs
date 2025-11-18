use contract_sdk_rust::sim_context;
use contract_sdk_rust::sim_context::SimContext;

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

// 创建did文档
#[unsafe(no_mangle)]
pub extern "C" fn createDID() {
    let ctx = &mut sim_context::get_sim_context();
    let did = ctx.arg_as_utf8_str("did");
    let did_doc = ctx.arg_as_utf8_str("didDoc");
    ctx.emit_event("ADD_DID", &vec![format!("{}", did)]);

    // todo 检查是否已经存在

    // 保存
    let did_manager_prefix = "did_manager_prefix";
    ctx.put_state(
        did_manager_prefix,
        hexEncode(&did).as_str(),
        hexEncode(&did_doc).as_bytes(),
    );

    ctx.ok(format!("add did {}", did).as_bytes());
}

// 更新did文档
#[unsafe(no_mangle)]
pub extern "C" fn updateDID() {
    let ctx = &mut sim_context::get_sim_context();
    let did = ctx.arg_as_utf8_str("did");
    let did_doc = ctx.arg_as_utf8_str("didDoc");
    ctx.emit_event("ADD_DID", &vec![format!("{}", did)]);

    // 保存
    let did_manager_prefix = "did_manager_prefix";
    ctx.put_state(
        did_manager_prefix,
        hexEncode(&did).as_str(),
        hexEncode(&did_doc).as_bytes(),
    );

    ctx.ok(format!("add did {}", did).as_bytes());
}

// did文档是否存在
#[unsafe(no_mangle)]
pub extern "C" fn doesDIDExist() {
    let ctx = &mut sim_context::get_sim_context();
    let did = ctx.arg_as_utf8_str("did");

    // 保存
    let did_manager_prefix = "did_manager_prefix";
    let r = ctx.get_state(did_manager_prefix, hexEncode(&did).as_str());

    if r.is_ok() {
        ctx.ok("true".as_bytes());
    } else {
        ctx.ok("false".as_bytes());
    }
}

// 获取did文档详情
#[unsafe(no_mangle)]
pub extern "C" fn getDIDDetails() {
    let ctx = &mut sim_context::get_sim_context();
    let did = ctx.arg_as_utf8_str("did");

    // 保存
    let did_manager_prefix = "did_manager_prefix";
    let r = ctx.get_state(did_manager_prefix, hexEncode(&did).as_str());

    if let Ok(did_doc_bytes) = r {
        let did_doc = hexDecode(&String::from_utf8(did_doc_bytes).expect("did content fail"));
        ctx.ok(did_doc.as_bytes());
    } else {
        ctx.error(format!("did: {} not exists", did).as_str());
    }
}

fn hexEncode(did: &String) -> String {
    hex::encode(did)
}

fn hexDecode(didEncoded: &String) -> String {
    String::from_utf8(hex::decode(didEncoded).expect("did decode fail"))
        .expect("parse code to String fail")
}
