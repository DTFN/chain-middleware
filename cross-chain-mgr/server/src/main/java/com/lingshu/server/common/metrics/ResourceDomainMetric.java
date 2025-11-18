package com.lingshu.server.common.metrics;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class ResourceDomainMetric {
    private final MeterRegistry meterRegistry;

    public void incrementSuccess(String gatewayId, String chainRid, String contractType) {
        Counter counter = Counter
                .builder("resource_domain")
                .description("应用启动后创建账户的数量")
                .tag("result", "success")
                .tag("gatewayId", gatewayId)
                .tag("chainRid", chainRid)
                .tag("contractType", contractType)
                .register(meterRegistry);
        counter.increment();
    }

    public void incrementFail(String gatewayId, String chainRid, String contractType) {
        Counter counter = Counter
                .builder("resource_domain")
                .description("应用启动后创建账户的数量")
                .tag("result", "fail")
                .tag("gatewayId", gatewayId)
                .tag("chainRid", chainRid)
                .tag("contractType", contractType)
                .register(meterRegistry);
        counter.increment();
    }
}
