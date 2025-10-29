package com.lingshu.bsp.front.base.config.entity;

import cn.hutool.core.lang.Assert;
import com.lingshu.chain.sdk.client.IClient;

import java.lang.reflect.InvocationHandler;
import java.lang.reflect.Method;
import java.lang.reflect.Proxy;
import java.util.concurrent.locks.ReentrantReadWriteLock;

/**
 * IClient包装类
 * 防止连接重置，继续使用旧的IClient连接导致jvm崩溃
 *
 * @author XuHang
 * @since 2023/10/31
 **/
public class IClientWrapper implements InvocationHandler {
    private final IClient client;
    private final long modCount;
    private ReentrantReadWriteLock.ReadLock opLock;
    private LingShuChainSDKWrapper lscsw;

    public static IClient wrap(IClient client, LingShuChainSDKWrapper lscsw) {
        IClientWrapper iClientWrapper = new IClientWrapper(client, lscsw);
        return (IClient) Proxy.newProxyInstance(client.getClass().getClassLoader(), new Class[]{IClient.class}, iClientWrapper);
    }

    private IClientWrapper(IClient client, LingShuChainSDKWrapper lscsw) {
        this.client = client;
        this.modCount = lscsw.getModCount();
        this.lscsw = lscsw;
        this.opLock = lscsw.getOpLock();
    }

    @Override
    public Object invoke(Object o, Method method, Object[] args) throws Throwable {
        opLock.lock();
        try {
            Assert.isTrue(lscsw.getModCount() == modCount, () -> new RuntimeException("connection already reset"));
            return method.invoke(client, args);
        } finally {
            opLock.unlock();
        }
    }
}
