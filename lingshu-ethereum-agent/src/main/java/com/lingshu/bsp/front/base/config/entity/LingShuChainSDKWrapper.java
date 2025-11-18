package com.lingshu.bsp.front.base.config.entity;

import cn.hutool.core.lang.Assert;
import com.lingshu.bsp.front.util.JavaCmdExecutor;
import com.lingshu.chain.sdk.LingShuChainSDK;
import com.lingshu.chain.sdk.LingShuChainSDKException;
import com.lingshu.chain.sdk.client.IClient;
import com.lingshu.chain.sdk.config.ConfigOption;
import com.lingshu.chain.sdk.evtsub.IEvtSubscription;
import com.lingshu.chain.sdk.jni.BlockNotificationHandler;
import com.lingshu.chain.sdk.ocm.Ocm;
import lombok.extern.slf4j.Slf4j;

import java.util.Map;
import java.util.concurrent.atomic.AtomicLong;
import java.util.concurrent.locks.ReentrantReadWriteLock;
import java.util.stream.Collectors;

/**
 * LingShuChainSDKChild
 *
 * @author XuHang
 * @since 2023/10/31
 **/
@Slf4j
public class LingShuChainSDKWrapper extends LingShuChainSDK {
    private LingShuChainSDK chainSDK;

    private final ConfigOption configOption;

    // 连接区块链的次数
    private final AtomicLong modCount = new AtomicLong(0);

    private final ReentrantReadWriteLock lock = new ReentrantReadWriteLock();

    public static LingShuChainSDKWrapper wrap(ConfigOption configOption) {
        LingShuChainSDKWrapper lingShuChainSDKWrapper = new LingShuChainSDKWrapper(configOption);
        lingShuChainSDKWrapper.reNew();
        return lingShuChainSDKWrapper;
    }

    private LingShuChainSDKWrapper(ConfigOption configOption) throws LingShuChainSDKException {
        // 删除连接
        super(configOption);
        super.stop();
        super.destroy();

        // 保存配置文件
        this.configOption = configOption;
    }

    public long getModCount() {
        return modCount.longValue();
    }

    public ReentrantReadWriteLock.ReadLock getOpLock() {
        return lock.readLock();
    }

    public synchronized void reNew() {
        log.info("reNew LingShuChainSDK");
        close();
        chainSDK = new LingShuChainSDK(configOption);
    }

    public synchronized void close() {
        if (chainSDK != null){
            log.info("LingShuChainSDK destroy");
            chainSDK.stop();
            chainSDK.destroy();
            chainSDK = null;
            // 连接次数增长
            modCount.incrementAndGet();
        }
    }

    @Override
    public synchronized void stop() {
        if (chainSDK != null){
            chainSDK.stop();
        }
    }

    @Override
    public synchronized void destroy() {
        if (chainSDK != null){
            chainSDK.destroy();
        }
    }

    @Override
    public synchronized Map<Integer, IClient> getClientMap() {
        checkChainSdk();
        Map<Integer, IClient> clientMap = chainSDK.getClientMap();

        // 包装每个IClient
        return clientMap.entrySet().stream()
                .collect(
                        Collectors.toMap(
                                Map.Entry::getKey,
                                ent -> IClientWrapper.wrap(ent.getValue(), this)
                        )
                );
    }

    @Override
    public ConfigOption getConfig() {
        checkChainSdk();
        return chainSDK.getConfig();
    }

    @Override
    public synchronized void registerBlockNotifier(Integer ledgerId, BlockNotificationHandler notificationHandler) {
        checkChainSdk();
        chainSDK.registerBlockNotifier(ledgerId, notificationHandler);
    }

    @Override
    public synchronized IClient getClient() {
        checkChainSdk();
        IClient client1 = chainSDK.getClient();
        return IClientWrapper.wrap(client1, this);
    }

    @Override
    public synchronized IClient getClient(Integer ledgerId) {
        checkChainSdk();
        IClient client1 = chainSDK.getClient(ledgerId);
        return IClientWrapper.wrap(client1, this);
    }

    @Override
    public synchronized Ocm getOcm() {
        checkChainSdk();
        Ocm ocm = chainSDK.getOcm();
        return OcmWrapper.wrap(ocm, this);
    }

    @Override
    public synchronized IEvtSubscription getEventSubscribe(Integer ledgerId) {
        checkChainSdk();
        IEvtSubscription eventSubscribe = chainSDK.getEventSubscribe(ledgerId);
        return new EvtSubscriptionWrapper(eventSubscribe, this);
    }

    private void checkChainSdk(){
        if (chainSDK == null){
            log.warn("chainSdk not connect, try to reconnect...");
            JavaCmdExecutor.ExecResult execResult = JavaCmdExecutor.executeCommand("bash /status.sh", 0);
            if (execResult.failed()){
                throw new RuntimeException("chain node not running");
            }
            reNew();
        }
    }
}

