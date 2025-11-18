package com.lingshu.server.common.metrics;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class CreateAccountMetric {
    private final MeterRegistry meterRegistry;

    public void incrementSuccess(String type) {
        Counter counter = Counter
                .builder("account_count")
                .description("应用启动后创建账户的数量")
                .tag("result", "success")
                .tag("opType", type)
                .register(meterRegistry);
        counter.increment();
    }

    public void incrementFail(String type) {
        Counter counter = Counter
                .builder("account_count")
                .description("应用启动后创建账户的数量")
                .tag("result", "fail")
                .tag("opType", type)
                .register(meterRegistry);
        counter.increment();
    }
}
