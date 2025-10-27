package com.lingshu.fabric.agent.event;

import org.hyperledger.fabric.gateway.Gateway;
import org.springframework.context.ApplicationEvent;

/**
 * @author: zehao.song
 */
public class GatewayInitedEvent extends ApplicationEvent {

    public GatewayInitedEvent(Object source) {
        super(source);
    }

    public Gateway getGateway() {
        return (Gateway) getSource();
    }
}
