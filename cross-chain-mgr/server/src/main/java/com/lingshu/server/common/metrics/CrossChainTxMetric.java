package com.lingshu.server.common.metrics;

import cn.hutool.core.date.DatePattern;
import cn.hutool.core.date.DateUtil;
import cn.hutool.extra.mail.MailUtil;
import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class CrossChainTxMetric {
    @Value("${prometheusUrl:http://192.168.1.177:9090}")
    private String prometheusUrl;

    private final MeterRegistry meterRegistry;

    public void incrementSuccess(String did, String originResourceName, String targetResourceName) {
        Counter counter = Counter
                .builder("cross_chain_transaction")
                .description("跨链交易统计")
                .tag("result", "success")
                .tag("failType", "")
                .tag("did", did)
                .tag("originResourceName", originResourceName)
                .tag("targetResourceName", targetResourceName)
                .register(meterRegistry);
        counter.increment();
    }

    public void incrementFail(String did, String originResourceName, String targetResourceName) {
        Counter counter = Counter
                .builder("cross_chain_transaction")
                .description("跨链交易统计")
                .tag("result", "fail")
                .tag("failType", "")
                .tag("did", did)
                .tag("originResourceName", originResourceName)
                .tag("targetResourceName", targetResourceName)
                .register(meterRegistry);
        counter.increment();
    }

    public void incrementSignFail(String did, String originResourceName, String targetResourceName) {
        Counter counter = Counter
                .builder("cross_chain_transaction")
                .description("签名验签失败统计")
                .tag("result", "fail")
                .tag("failType", "signFail")
                .tag("did", did)
                .tag("originResourceName", originResourceName)
                .tag("targetResourceName", targetResourceName)
                .register(meterRegistry);
        counter.increment();

        // 邮件通知
        MailUtil.send(
                "admin@foxmail.com",
                "【异常】系统签名验签流程执行失败告警",
                String.format(
                        "系统监控中心于 %s 监测到一笔业务请求在执行签名验签流程时出现失败异常，相关核心信息及影响说明如下:\n" +
                                "异常账户DID\t%s\n" +
                                "原始资源域名\t%s\n" +
                                "目标资源域名\t%s\n" +
                                "异常账户统计\t%s\n",
                        DateUtil.date().toString(DatePattern.NORM_DATETIME_PATTERN),
                        did,
                        originResourceName,
                        targetResourceName,
                        prometheusUrl + "/graph?g0.expr=cross_chain_transaction_total%7BfailType%3D%22signFail%22%2Cdid%3D%22"
                                + did
                                + "%22%7D&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=1h"
                ),
                false
        );
    }
}
