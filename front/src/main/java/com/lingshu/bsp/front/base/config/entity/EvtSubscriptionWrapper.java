package com.lingshu.bsp.front.base.config.entity;

import cn.hutool.core.lang.Assert;
import com.lingshu.chain.sdk.evtsub.EvtSubCallback;
import com.lingshu.chain.sdk.evtsub.IEvtSubscription;
import com.lingshu.chain.sdk.evtsub.model.EvtSubParams;

import java.util.concurrent.locks.ReentrantReadWriteLock;

/**
 * EvtSubscriptionWrapper
 * 防止连接重置，继续使用旧的IClient连接导致jvm崩溃
 *
 * @author XuHang
 * @since 2023/10/31
 **/
public class EvtSubscriptionWrapper implements IEvtSubscription {
    private final IEvtSubscription evt;
    private final long modCount;
    private ReentrantReadWriteLock.ReadLock opLock;
    private LingShuChainSDKWrapper lscsw;

    public EvtSubscriptionWrapper(IEvtSubscription evt, LingShuChainSDKWrapper lscsw) {
        this.evt = evt;
        this.modCount = lscsw.getModCount();
        this.lscsw = lscsw;
        this.opLock = lscsw.getOpLock();
    }

    @Override
    public String subscribe(EvtSubParams evtSubParams, EvtSubCallback evtSubCallback) {
        opLock.lock();
        try {
            Assert.isTrue(lscsw.getModCount() == modCount, () -> new RuntimeException("connection already reset"));
            return evt.subscribe(evtSubParams, evtSubCallback);
        } finally {
            opLock.unlock();
        }
    }

    @Override
    public void unsubscribe(String s) throws Exception {
        opLock.lock();
        try {
            Assert.isTrue(lscsw.getModCount() == modCount, () -> new RuntimeException("connection already reset"));
            evt.unsubscribe(s);
        } finally {
            opLock.unlock();
        }
    }
}
