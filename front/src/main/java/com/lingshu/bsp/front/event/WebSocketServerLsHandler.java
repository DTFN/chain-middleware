package com.lingshu.bsp.front.event;

import cn.hutool.core.collection.CollectionUtil;
import cn.hutool.json.JSONUtil;
import com.lingshu.bsp.front.rpcapi.ls.contract.BusiNormalLs;
import com.lingshu.bsp.front.rpcapi.ls.contract.CrossSaveLs;
import com.lingshu.bsp.front.rpcapi.ls.RpcApiService;
import com.lingshu.chain.sdk.client.IClient;
import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Component;

import javax.annotation.Resource;
import javax.websocket.OnClose;
import javax.websocket.OnMessage;
import javax.websocket.OnOpen;
import javax.websocket.Session;
import javax.websocket.server.ServerEndpoint;
import java.io.IOException;
import java.util.HashSet;
import java.util.Set;

@Slf4j
@Component
@ServerEndpoint("/event")
public class WebSocketServerLsHandler implements CommandLineRunner {

    private static final Set<Session> SESSION_SET = new HashSet<>();

    public static void broadcast(String message) {
        log.info("SESSION_SET size {}", SESSION_SET.size());
        for (Session webSocketSession : SESSION_SET) {
            try {
                webSocketSession.getBasicRemote().sendText(message);
            } catch (IOException e) {
                log.error("", e);
            }
        }
    }

    @OnOpen
    public void onOpen(Session session) {
        log.info("Server: WebSocket connection established");
        SESSION_SET.add(session);
    }

    @OnClose
    public void onClose(Session session) {
        log.info("Server: WebSocket connection closed");
        SESSION_SET.remove(session);
    }

    @OnMessage
    public void onMessage(String message, Session session) {
        //消息直接广播出去
        if ("ping".equals(message)) {
            try {
                session.getBasicRemote().sendText("pong");
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }
    }

    @Resource
    private RpcApiService rpcApiService;
    @Override
    public void run(String... args) throws Exception {
        // listen event
        IClient sdkClient = rpcApiService.getSdkClient(1);

        BusiNormalLs busiNormalLs = BusiNormalLs.load("0x2b684e59f9444373a7ae17a497fe296c24c40a9b", sdkClient, sdkClient.getCryptoSuite().getKeyPair());
        busiNormalLs.subscribeCROSS_CHAIN_VCEvent((s, i, list) -> {
            try {
                if (CollectionUtil.isEmpty(list)) {
                    return;
                }
                log.info("event size: {}", list.size());
                broadcast(JSONUtil.toJsonStr(list));
            } catch (Exception e) {
                log.error("", e);
            }
        });
    }
}