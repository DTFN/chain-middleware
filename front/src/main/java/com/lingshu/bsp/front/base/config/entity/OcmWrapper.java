package com.lingshu.bsp.front.base.config.entity;

import cn.hutool.core.lang.Assert;
import com.lingshu.chain.sdk.client.IClient;
import com.lingshu.chain.sdk.ocm.Ocm;

import java.lang.reflect.InvocationHandler;
import java.lang.reflect.Method;
import java.lang.reflect.Proxy;
import java.util.concurrent.locks.ReentrantReadWriteLock;

/**
 * OcmWrapper
 * 防止连接重置，继续使用旧的IClient连接导致jvm崩溃
 *
 * @author XuHang
 * @since 2023/10/31
 **/
public class OcmWrapper implements InvocationHandler {
    private final Ocm ocm;
    private final long modCount;
    private ReentrantReadWriteLock.ReadLock opLock;
    private LingShuChainSDKWrapper lscsw;

    public static Ocm wrap(Ocm ocm, LingShuChainSDKWrapper lscsw) {
        OcmWrapper ocmWrapper = new OcmWrapper(ocm, lscsw);
        return (Ocm) Proxy.newProxyInstance(ocm.getClass().getClassLoader(), new Class[]{IClient.class}, ocmWrapper);
    }

    private OcmWrapper(Ocm ocm, LingShuChainSDKWrapper lscsw) {
        this.ocm = ocm;
        this.modCount = lscsw.getModCount();
        this.lscsw = lscsw;
        this.opLock = lscsw.getOpLock();
    }

    @Override
    public Object invoke(Object o, Method method, Object[] args) throws Throwable {
        opLock.lock();
        try {
            Assert.isTrue(lscsw.getModCount() == modCount, () -> new RuntimeException("connection already reset"));
            return method.invoke(ocm, args);
        } finally {
            opLock.unlock();
        }
    }
}
