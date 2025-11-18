package com.lingshu.fabric.agent.service;

import com.lingshu.fabric.agent.resp.NodeInfoDTO;
import com.lingshu.fabric.agent.resp.UnStableNodeInfosDTO;

import java.util.Date;

/**
 * MonitorService
 *
 * @author XuHang
 * @Date 2023/12/5 下午1:35
 **/
public interface MonitorService {
    /**
     * 实时数据
     * @param ip
     * @return
     */
    NodeInfoDTO hostInfo(String ip);

    /**
     * 历史数据
     * @param ip
     * @param start
     * @param end
     * @return
     */
    UnStableNodeInfosDTO hostInfoHistory(String ip, Date start, Date end);

    /**
     * 最近记录的数据
     * @param ip
     * @return
     */
    NodeInfoDTO hostLatestInfo(String ip);
}
