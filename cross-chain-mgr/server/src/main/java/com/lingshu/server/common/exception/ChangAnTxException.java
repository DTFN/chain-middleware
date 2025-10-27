package com.lingshu.server.common.exception;


import cn.hutool.json.JSONUtil;
import lombok.Data;
import org.chainmaker.pb.common.ResultOuterClass;

import java.util.Arrays;
import java.util.HashSet;
import java.util.Optional;
import java.util.Set;

/**
 * <p>
 * 全局异常
 * </p>
 *
 * @author xiaoxiang.zhang
 * @date Created in 2021-08-20 17:24
 */
@Data
public class ChangAnTxException extends RuntimeException {
    private ResultOuterClass.TxResponse txResponse;
    private String message;

    public ChangAnTxException(ResultOuterClass.TxResponse txResponse, String message) {
        super(message);
        this.txResponse = txResponse;
        this.message = message;
    }

    public ChangAnTxException(ResultOuterClass.TxResponse txResponse) {
        super("Send Invoker Contract Request Failed, result = " + Optional.ofNullable(txResponse).map(ResultOuterClass.TxResponse::getMessage).orElse("not known"));
        this.txResponse = txResponse;
        this.message = "Send Invoker Contract Request Failed, result = " + Optional.ofNullable(txResponse).map(ResultOuterClass.TxResponse::getMessage).orElse("not known");
    }

    /**
     * 判断是否在合约执行抛出的异常
     * 文档：https://docs.chainmaker.org.cn/tech/%E4%B8%AD%E7%BB%A7%E8%B7%A8%E9%93%BE%E5%8D%8F%E8%AE%AE.html#code
     *
     * @return 是否是合约异常
     */
    public boolean isContractExecFail() {
        // ResultOuterClass.TxStatusCode.SUCCESS 和 txResponse.getCode()
        ResultOuterClass.TxStatusCode txCode = txResponse.getCode();
        if (ResultOuterClass.TxStatusCode.SUCCESS.equals(txCode)) {
            return false;
        }

        // 合约执行出错的状态
        Set<ResultOuterClass.TxStatusCode> contractFailCode = new HashSet<>(Arrays.asList(ResultOuterClass.TxStatusCode.CONTRACT_FAIL));
        if (contractFailCode.contains(txCode)) {
            return true;
        } else {
            return false;
        }
    }

    @Override
    public String toString() {
        try {
            return String.format("txResponse, code: %s, message: %s, transactionHash: %s, blockheight: %s", txResponse.getCode(), txResponse.getMessage(), txResponse.getTxId(), txResponse.getTxBlockHeight());
        } catch (Exception e) {
            if (txResponse == null) {
                return "txResponse is null";
            } else {
                return JSONUtil.toJsonStr(txResponse);
            }
        }
    }
}
