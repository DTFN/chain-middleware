package com.lingshu.bsp.front.event;

import cn.hutool.core.bean.BeanUtil;
import cn.hutool.json.JSONUtil;
import com.lingshu.bsp.front.rpcapi.bsc.contract.BusiCenterBsc;
import com.lingshu.bsp.front.rpcapi.bsc.contract.CrossSaveBsc;
import com.lingshu.bsp.front.rpcapi.bsc.BscApiService;
import io.reactivex.functions.Consumer;
import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Component;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.tx.gas.StaticGasProvider;

import javax.annotation.Resource;
import javax.websocket.OnClose;
import javax.websocket.OnMessage;
import javax.websocket.OnOpen;
import javax.websocket.Session;
import javax.websocket.server.ServerEndpoint;
import java.io.IOException;
import java.math.BigInteger;
import java.util.*;

@Slf4j
@Component
@ServerEndpoint("/bsc")
public class WebSocketServerBscHandler implements CommandLineRunner {

    private static final Set<Session> SESSION_SET = new HashSet<>();

    public static void broadcast(String message) {
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
    private BscApiService bscApiService;

    @Override
    public void run(String... args) throws Exception {

        BusiCenterBsc busiCenterBsc = BusiCenterBsc.load("0xf3dc3b9a436b5d8a1981a07a2d222cce55d95e3e", bscApiService.web3j(),
                bscApiService.transactionManager(), new StaticGasProvider(BigInteger.valueOf(20_000_000_000L), BigInteger.valueOf(2_800_000)));
        try {
            busiCenterBsc.cROSS_CHAIN_VCEventFlowable(DefaultBlockParameterName.EARLIEST, DefaultBlockParameterName.LATEST).subscribe(new Consumer<BusiCenterBsc.CROSS_CHAIN_VCEventResponse>() {
                @Override
                public void accept(BusiCenterBsc.CROSS_CHAIN_VCEventResponse response) throws Exception {
                    List arrList = new ArrayList<>();
                    Map<String, Object> bean = BeanUtil.beanToMap(response.log);
                    arrList.add(bean);
                    String jsonStr = JSONUtil.toJsonStr(arrList);
                    broadcast(jsonStr);
                    log.info("{}", jsonStr);
                }
            });
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}