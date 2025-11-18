package com.lingshu.fabric.agent.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.lingshu.fabric.agent.config.properties.ConstantProperties;
import com.lingshu.fabric.agent.repo.MonitorNodeRepo;
import com.lingshu.fabric.agent.repo.MonitorNodeStableRepo;
import com.lingshu.fabric.agent.repo.entity.MonitorNodeDo;
import com.lingshu.fabric.agent.repo.entity.MonitorNodeStableDo;
import com.lingshu.fabric.agent.repo.entity.prometheus.ValuesItemDto;
import com.lingshu.fabric.agent.resp.AnsibleScriptOutDTO;
import com.lingshu.fabric.agent.resp.NodeInfoDTO;
import com.lingshu.fabric.agent.resp.UnStableNodeInfosDTO;
import com.lingshu.fabric.agent.service.MonitorService;
import com.lingshu.fabric.agent.util.ExecuteResult;
import com.lingshu.fabric.agent.util.JavaCommandExecutor;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.math.NumberUtils;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.LinkedList;
import java.util.List;

/**
 * ShellMonitorImpl
 *
 * @author XuHang
 * @since 2023/12/5
 **/
@Slf4j
@Service
@RequiredArgsConstructor
@ConditionalOnProperty(value = "monitor.type", havingValue = "SHELL")
public class ShellMonitorImpl implements MonitorService {
    private final ConstantProperties constant;
    private final MonitorNodeStableRepo monitorNodeStableRepo;
    private final MonitorNodeRepo monitorNodeRepo;

    private ObjectMapper om = new ObjectMapper();

    public NodeInfoDTO hostInfo(String ip) {
        NodeInfoDTO nodeInfoDTO = new NodeInfoDTO();

        String command = String.format("ansible %s -m script -a \"%s\"", ip, constant.getAllInfo());
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecShellTimeout());

        try {
            if (result.success()) {
                String out = result.getExecuteOut();
                int start = out.indexOf("=>");
                AnsibleScriptOutDTO ansibleScriptOutDTO = om.readValue(out.substring(start + 2), AnsibleScriptOutDTO.class);

                if (ansibleScriptOutDTO.getStdoutLines().size() == 8) {
                    nodeInfoDTO.setCpuNum(NumberUtils.toInt(ansibleScriptOutDTO.getStdoutLines().get(0), 1));
                    nodeInfoDTO.setCpuRatio(NumberUtils.toFloat(ansibleScriptOutDTO.getStdoutLines().get(1), 0));
                    nodeInfoDTO.setMaxMem(NumberUtils.toLong(ansibleScriptOutDTO.getStdoutLines().get(2), 0));
                    nodeInfoDTO.setUsedMem(NumberUtils.toLong(ansibleScriptOutDTO.getStdoutLines().get(3), 0));
                    nodeInfoDTO.setMaxDisk(NumberUtils.toLong(ansibleScriptOutDTO.getStdoutLines().get(4), 0));
                    nodeInfoDTO.setUsedDisk(NumberUtils.toLong(ansibleScriptOutDTO.getStdoutLines().get(5), 0));
                    nodeInfoDTO.setNetIn(NumberUtils.toLong(ansibleScriptOutDTO.getStdoutLines().get(6), 0));
                    nodeInfoDTO.setNetOut(NumberUtils.toLong(ansibleScriptOutDTO.getStdoutLines().get(7), 0));
                }

                return nodeInfoDTO;
            }
        } catch (Exception e) {
            log.error("getNetInfo error, ip:{}", ip, e);
        }
        return nodeInfoDTO;
    }

    @Override
    public UnStableNodeInfosDTO hostInfoHistory(String ip, Date start, Date end) {
        UnStableNodeInfosDTO r = new UnStableNodeInfosDTO();
        r
                .setIp(ip)
                .setCpuUsage(new LinkedList<>())
                .setMemUsage(new LinkedList<>())
                .setDiskUsage(new LinkedList<>())
                .setNetIn(new LinkedList<>())
                .setNetOut(new LinkedList<>());

        List<MonitorNodeDo> ms = monitorNodeRepo.list(
                new LambdaQueryWrapper<MonitorNodeDo>().eq(MonitorNodeDo::getIp, ip)
                        .ge(start != null, MonitorNodeDo::getCreateTime, start)
                        .le(end != null, MonitorNodeDo::getCreateTime, end)
                        .orderByAsc(MonitorNodeDo::getCreateTime)
        );
        for (MonitorNodeDo m : ms) {
            r.getCpuUsage().add(new ValuesItemDto<Float>(m.getCpuUsage(), m.getCreateTime()));
            r.getMemUsage().add(new ValuesItemDto<Long>(m.getMemUsage(), m.getCreateTime()));
            r.getDiskUsage().add(new ValuesItemDto<Long>(m.getDiskUsage(), m.getCreateTime()));
            r.getNetIn().add(new ValuesItemDto<Long>(m.getNetIn(), m.getCreateTime()));
            r.getNetOut().add(new ValuesItemDto<Long>(m.getNetOut(), m.getCreateTime()));
        }
        return r;
    }

    @Override
    public NodeInfoDTO hostLatestInfo(String ip) {
        NodeInfoDTO nodeInfoDTO = new NodeInfoDTO();

        // 稳定数据
        MonitorNodeStableDo stableInfo = monitorNodeStableRepo.getOne(
                new LambdaQueryWrapper<MonitorNodeStableDo>()
                        .eq(MonitorNodeStableDo::getIp, ip)
        );
        if (stableInfo != null) {
            nodeInfoDTO
                    .setCpuNum(stableInfo.getCpuNumber())
                    .setMaxMem(stableInfo.getMaxMemory())
                    .setMaxDisk(stableInfo.getMaxDisk());
        }

        MonitorNodeDo unstableInfo = monitorNodeRepo.getOne(
                new LambdaQueryWrapper<MonitorNodeDo>().eq(MonitorNodeDo::getIp, ip)
                        .orderByDesc(MonitorNodeDo::getCreateTime)
                        .last("limit 1")
        );
        if (unstableInfo != null) {
            nodeInfoDTO
                    .setCpuRatio(unstableInfo.getCpuUsage())
                    .setUsedMem(unstableInfo.getMemUsage())
                    .setUsedDisk(unstableInfo.getDiskUsage())
                    .setNetIn(unstableInfo.getNetIn())
                    .setNetOut(unstableInfo.getNetOut());
        }

        return nodeInfoDTO;
    }



    private MonitorNodeStableDo getStableInfoOrUpdate(String ip) {
        MonitorNodeStableDo stableInfo = monitorNodeStableRepo.getOne(
                new LambdaQueryWrapper<MonitorNodeStableDo>()
                        .eq(MonitorNodeStableDo::getIp, ip)
        );
        if (stableInfo != null) {
            return stableInfo;
        }

        return monitorNodeStableRepo.getOne(
                new LambdaQueryWrapper<MonitorNodeStableDo>()
                        .eq(MonitorNodeStableDo::getIp, ip)
        );
    }

    private MonitorNodeDo getUnstableInfoOrUpdate(String ip) {
        MonitorNodeDo unstableInfo = monitorNodeRepo.getOne(
                new LambdaQueryWrapper<MonitorNodeDo>().eq(MonitorNodeDo::getIp, ip)
                        .orderByDesc(MonitorNodeDo::getCreateTime)
                        .last("limit 1")
        );
        if (unstableInfo != null) {
            return unstableInfo;
        }

        return monitorNodeRepo.getOne(
                new LambdaQueryWrapper<MonitorNodeDo>().eq(MonitorNodeDo::getIp, ip)
                        .orderByDesc(MonitorNodeDo::getCreateTime)
                        .last("limit 1")
        );
    }
}
