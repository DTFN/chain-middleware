package com.lingshu.server.common.api;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;

/**
 * REST API 错误码
 *
 * @author Brian
 * @since 2021-03-13
 */
@AllArgsConstructor
@NoArgsConstructor
@Getter
public enum ApiErrorCode {

    FAILED(500, "操作失败"),
    SUCCESS(200, "执行成功"),
    NOT_EXIST_ACCOUNT(600, "链账户地址不存在"),
    BEYOND_MAX_CREATE_ACCOUNT(601, "超出单次最大创建数量"),
    NET_ERROR(602, "网络超时，请稍后再试"),
    SIGN_EXPIRED(603, "签名已过期"),
    DEPLOY_CONTRACT_FAILED(604, "部署合约失败"),
    GRANT_SHARE_FAILED(605, "授予份额失败"),
    TRANSFER_SHARE_FAILED(606, "转让份额失败"),
    ACCOUNT_FROZEN(607, "链账户已被冻结"),
    MISS_PARAM(608, "缺失必要参数"),
    PARAM_TOO_LONG(609, "参数超过指定长度"),
    ILLEGAL_PARAM(610, "非法参数"),
    EXIST_ASSET_ID(611, "已存在的资产编号"),
    NOT_EXIST_ASSET_ID(612, "不存在的资产编号或合约地址"),
    NOT_SUPPORT_ALTER(613, "当前资产不支持修改"),
    NOT_SUPPORT_PUBLISH(614, "当前资产不支持发行"),
    REPEAT_REQUEST_ID(615, "重复的请求ID"),
    NOT_EXIST_SHARE_ID(616, "不存在的份额ID"),
    GRANTED_SHARE_ID(617, "已授予的份额ID"),
    BEYOND_MAX_GRANT_AMOUNT(618, "超出最大授予份额"),
    SHARE_NOT_TRANSFER(619, "当前份额状态不支持转让"),
    NOT_SHARE_OWNER(620, "非当前份额持有人"),
    NOT_OWNER_TOKEN(622, "非当前token持有人不能发起交易"),
    NOT_EXIST_TOKEN(623, "当前tokenId不存在"),
    TRANSFER_TOKEN_ID_FAIL(624, "交易tokenId失败"),
    NOT_USED_TO_TRANSFER_TOKEN_ID(625, "当前tokenId不能用来被交易"),
    TRANSFER_TOKEN_ID_EXISTED(626, "已被交易的tokenId或正在被交易的tokenId"),
    NOT_SUPPORT_MINT_TOKEN_ID(627, "当前状态下的tokenId不支持重铸"),
    NOT_SUPPORT_GRANT_SHARE(628, "当前资产状态不支持授予资产"),
    BEYOND_MAX_THRESHOLD(629, "超出系统处理阈值，请稍后再试"),
    INVALID_APP_ID(630, "无效的appId"),
    SIGN_VERIFICATION_FAILED(631, "验签失败"),
    ACTIVITY_ID_NOT_EXISTED(632, "不存在的活动ID"),
    REPEAT_ACTIVITY_ID(633, "重复的活动ID"),
    ACTIVITY_PUBLISHED_NO_UPDATE(634, "当前活动已发布无法再更改"),
    SHARE_JOIN_ACTIVITY(635, "已参与活动的资产份额"),
    NOT_SUPPORT_JOIN_ACTIVITY(636, "当前资产份额不支持参与活动"),
    HAS_REWARD(637, "当前活动奖励已被兑换"),
    NOT_COLLECT_TOKEN_ID(638, "尚未集齐活动所需的资产份额"),
    NOT_REWARD(639, "当前活动奖励未被兑换"),
    NOT_PUBLISH(640, "当前活动奖励未发布"),
    NOT_OPEN_SHOP_CHARGING(644, "尚未开通调用接口的权限"),
    SHOP_CHARGING_EXPIRED(645, "计费服务已过期"),
    SHOP_CHARGING_NOT_START(647, "未到达计费服务开始日期"),
    NO_TIME_TO_USE(646, "计费接口调用次数已达上限"),
    NOT_EXIST_SHOP(650, "不存在的商户"),
    EXIST_SHOP_CHARGING(651, "已有客户账户，不能新建"),
    NOT_EXIST_SHOP_CHARGING(652, "不存在的客户账户"),
    EXISTED_SHOP(653, "客户名称/机构代码已存在，请勿重复创建备案信息"),
    CHANGE_STATUS_ERROR(654,"变更接入状态失败，请先为该渠道添加客户账户后再试。"),
    INCONSISTENT_PASSWORDS(655, "两次密码不一致"),
    SMS_CODE_ONE_DAY_LIMIT(656, "今日发送验证码数量达到上限"),
    SMS_CODE_INTERVAL_LIMIT(657, "请1分钟后再发送验证码"),
    SMS_CODE_MISSING_MATCH(658, "验证码输入错误"),
    WORKBENCH_USER_EXISTS(659, "该手机号已注册"),
    WORKBENCH_USER_NOT_EXISTS(660, "该用户不存在"),
    SHOP_CHARGING_GAS_NOT_CONSUMED_OUT(661, "当前账户的燃料余额不为0，无法进行账户方案的变更"),
    ALTER_SHOP_CHARGING_MODEL_NOT_CHANGE(662, "计费模式未变化"),
    CAPTCHA_CODE_EXPIRE(663, "验证码已过期"),
    CAPTCHA_CODE_VERIFY_FAILED(664, "验证码输入错误"),
    SHOP_AUDIT_STATUS_INVALID(665, "请完善企业认证信息"),
    TRANSFER_NFT_SIGNATURE_INVALID(666, "验签失败"),
    CREATE_ACCOUNT_REPEAT_DID(667, "重复数字身份"),
    NOT_EXIST_SHOP_NAME_OR_APP_ID(668, "不存在的客户名称或渠道编号"),
    SHOP_USER_LOCKED(706, "当前用户连续10次输入密码错误已经被锁定，请使用忘记密码功能重置登录密码解锁!"),

    NOT_REGISTER_PHONE(707,"当前用户未注册，请先完成注册"),
    ;

    private Integer code;
    private String msg;

    public static ApiErrorCode fromCode(Integer code) {
        ApiErrorCode[] ecs = ApiErrorCode.values();
        for (ApiErrorCode ec : ecs) {
            if (ec.getCode().equals(code)) {
                return ec;
            }
        }
        return FAILED;
    }

}
