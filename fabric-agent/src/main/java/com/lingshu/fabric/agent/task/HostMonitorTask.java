package com.lingshu.fabric.agent.task;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.lingshu.fabric.agent.repo.DeployResultNodeRepo;
import com.lingshu.fabric.agent.repo.MonitorNodeRepo;
import com.lingshu.fabric.agent.repo.MonitorNodeStableRepo;
import com.lingshu.fabric.agent.repo.entity.MonitorNodeDo;
import com.lingshu.fabric.agent.repo.entity.MonitorNodeStableDo;
import com.lingshu.fabric.agent.resp.NodeInfoDTO;
import com.lingshu.fabric.agent.service.impl.ShellMonitorImpl;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import java.util.Date;
import java.util.List;
import java.util.Objects;
import java.util.stream.Collectors;

/**
 * HostCheckTask
 *
 * @author XuHang
 * @since 2023/11/28
 **/
@Slf4j
@Component
@RequiredArgsConstructor
@ConditionalOnProperty(value = "monitor.type", havingValue = "SHELL")
public class HostMonitorTask {
    private final ShellMonitorImpl monitorService;
    private final DeployResultNodeRepo deployResultNodeRepo;
    private final MonitorNodeRepo monitorNodeRepo;
    private final MonitorNodeStableRepo monitorNodeStableRepo;

    @Value("${constant.hostCheck.unStableTimeout:31536000000}")
    private long unStableTimeout;
    @Value("${constant.hostCheck.stableTimeout:604800000}")
    private long stableTimeout;

    /**
     * 容易变化的指标
     */
    @Scheduled(fixedRateString = "${constant.hostCheck.unstableInterval:30000}")
    public void updateUnstable() {
        List<String> hosts = deployResultNodeRepo.allHosts();
        log.info("monitor unstable");
        List<MonitorNodeDo> moitors = hosts.parallelStream()
                .map(hostIp -> {
                    try {
                        NodeInfoDTO allInfo = monitorService.hostInfo(hostIp);
                        MonitorNodeDo monitorNodeDo = new MonitorNodeDo();
                        monitorNodeDo
                                .setIp(hostIp)
                                .setCpuUsage(allInfo.getCpuRatio())
                                .setMemUsage(allInfo.getUsedMem())
                                .setDiskUsage(allInfo.getUsedDisk())
                                .setNetIn(allInfo.getNetIn())
                                .setNetOut(allInfo.getNetOut());
                        return monitorNodeDo;
                    } catch (Exception e) {
                        log.error("get host info fail, hostIp:{}", hostIp, e);
                        return null;
                    }
                })
                .filter(Objects::nonNull)
                .collect(Collectors.toList());
        monitorNodeRepo.saveBatch(moitors);

        removeUnstableTimeout();
    }

    /**
     * 删除过期数据
     */
    private void removeUnstableTimeout() {
        // 计算地板时间
        Date floorTime = new Date(System.currentTimeMillis() - unStableTimeout);
        monitorNodeRepo.remove(
                new LambdaQueryWrapper<MonitorNodeDo>()
                        .lt(MonitorNodeDo::getCreateTime, floorTime)
        );
    }

    /**
     * 不易变化的指标
     */
    @Scheduled(fixedRateString = "${constant.hostCheck.stableInterval:86400000}")
    public void updateStable() {
        List<String> hosts = deployResultNodeRepo.allHosts();
        log.info("monitor stable, hosts:{}", hosts);
        hosts.parallelStream()
                .forEach(hostIp -> {
                    try {
                        NodeInfoDTO allInfo = monitorService.hostInfo(hostIp);
                        log.info("updateStable hostIp:{}, info:{}", hostIp, allInfo);
                        MonitorNodeStableDo monitorNodeStableDo = new MonitorNodeStableDo();
                        monitorNodeStableDo.setIp(hostIp)
                                .setCpuNumber(allInfo.getCpuNum())
                                .setMaxMemory(allInfo.getMaxMem())
                                .setMaxDisk(allInfo.getMaxDisk());
                        monitorNodeStableRepo.saveOrUpdate(monitorNodeStableDo);
                    } catch (Exception e) {
                        log.error("get host info fail, hostIp:{}", hostIp, e);
                    }
                });
    }

    /**
     * 删除已经删除的节点数据
     * 被删除的节点不会被更新
     */
    private void removeStableTimeout() {
        // 计算地板时间
        Date floorTime = new Date(System.currentTimeMillis() - stableTimeout);
        monitorNodeStableRepo.remove(
                new LambdaQueryWrapper<MonitorNodeStableDo>()
                        .lt(MonitorNodeStableDo::getModifyTime, floorTime)
        );
    }
}
