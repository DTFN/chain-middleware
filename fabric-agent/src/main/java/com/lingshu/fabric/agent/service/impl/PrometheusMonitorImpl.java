package com.lingshu.fabric.agent.service.impl;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.lingshu.fabric.agent.repo.entity.prometheus.DataResultDto;
import com.lingshu.fabric.agent.repo.entity.prometheus.ResultDataDto;
import com.lingshu.fabric.agent.repo.entity.prometheus.ResultDto;
import com.lingshu.fabric.agent.repo.entity.prometheus.ValuesItemDto;
import com.lingshu.fabric.agent.resp.NodeInfoDTO;
import com.lingshu.fabric.agent.resp.UnStableNodeInfosDTO;
import com.lingshu.fabric.agent.service.MonitorService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.math.NumberUtils;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.net.URI;
import java.net.URLEncoder;
import java.text.MessageFormat;
import java.util.*;
import java.util.stream.Collectors;

/**
 * ShellMonitorImpl
 *
 * @author XuHang
 * @since 2023/12/5
 **/
@Slf4j
@Service
@RequiredArgsConstructor
@ConditionalOnProperty(value = "monitor.type", havingValue = "PROMETGEUS")
public class PrometheusMonitorImpl implements MonitorService {
    @Value("${prometheus.uri}")
    private String domain;

    private ObjectMapper om = new ObjectMapper();

    private static final String API_QUERY = "/api/v1/query?query={0}";
    private static final String API_QUERY_RANGE = "/api/v1/query_range?query={0}&start={1,number,#}&end={2,number,#}&step={3,number,#}";
    private static final String CPU_NUMBER = "count(count(node_cpu_seconds_total'{'instance=~\"{0}:?[0-9]*\"}) by (cpu))";
    private static final String CPU_IDLE_RATIO = "avg(rate(node_cpu_seconds_total'{'instance=~\"{0}:?[0-9]*\",mode=\"idle\"}[1m]))";
    private static final String MEMORY_MAX = "node_memory_MemTotal_bytes'{'instance=~\"{0}:?[0-9]*\"}";
    private static final String MEMORY_USE = "node_memory_MemAvailable_bytes'{'instance=~\"{0}:?[0-9]*\"}";
    private static final String DISK_MAX = "sum(node_filesystem_size_bytes'{'instance=~\"{0}:?[0-9]*\"})";
    private static final String DISK_FREE = "sum(node_filesystem_free_bytes'{'instance=~\"{0}:?[0-9]*\"})";
    private static final String NET_IN = "sum(rate(node_network_receive_bytes_total'{'instance=~\"{0}:?[0-9]*\"}[1m]))";
    private static final String NET_OUT = "sum(rate(node_network_transmit_bytes_total'{'instance=~\"{0}:?[0-9]*\"}[1m]))";
    private static final int RANGE_STEP = 60;
    private static final int RANGE_SECOND = RANGE_STEP * (11000 - 1);

    public NodeInfoDTO hostInfo(String ip) {
        // cpu数量
        Integer cpuNum = getValue(MessageFormat.format(CPU_NUMBER, ip), Integer.class, 0);

        // cpu空闲率
        Float cpuIdle = getValue(MessageFormat.format(CPU_IDLE_RATIO, ip), Float.class, 1F);

        // 最大内存
        Long maxMem = getValue(MessageFormat.format(MEMORY_MAX, ip), Long.class, 0L);

        // 已用内存
        Long usedMem = getValue(MessageFormat.format(MEMORY_USE, ip), Long.class, 0L);

        // 最大磁盘
        Long maxDisk = getValue(MessageFormat.format(DISK_MAX, ip), Long.class, 0L);

        // 未使用磁盘
        Long freeDisk = getValue(MessageFormat.format(DISK_FREE, ip), Long.class, 0L);

        // 入网速度
        Long netIn = getValue(MessageFormat.format(NET_IN, ip), Long.class, 0L);

        // 出网速度
        Long netOut = getValue(MessageFormat.format(NET_OUT, ip), Long.class, 0L);

        NodeInfoDTO nodeInfoDTO = new NodeInfoDTO()
                .setCpuNum(cpuNum)
                .setCpuRatio(1 - cpuIdle)
                .setMaxMem(maxMem)
                .setUsedMem(usedMem)
                .setMaxDisk(maxDisk)
                .setUsedDisk(maxDisk - freeDisk)
                .setNetIn(netIn)
                .setNetOut(netOut);
        return nodeInfoDTO;
    }

    @Override
    public UnStableNodeInfosDTO hostInfoHistory(String ip, Date start, Date end) {
        // 时间范围不超过11000
        long startM, endM;
        if (start == null) {
            endM = Optional.ofNullable(end).map(Date::getTime).orElse(System.currentTimeMillis()) / 1000;
            startM = endM - RANGE_SECOND;
        } else {
            startM = start.getTime() / 1000;
            endM = startM + RANGE_SECOND;
        }

        // cpu使用
        List<ValuesItemDto<Float>> cpu = getValues(MessageFormat.format(CPU_IDLE_RATIO, ip), startM, endM, Float.class);
        for (ValuesItemDto<Float> v : cpu) {
            v.setValue(1 - v.getValue());
        }

        // 内存使用
        List<ValuesItemDto<Long>> mem = getValues(MessageFormat.format(MEMORY_USE, ip), startM, endM, Long.class);

        // 磁盘使用
        Long diskMax = getValue(MessageFormat.format(MEMORY_MAX, ip), Long.class, 0l);
        List<ValuesItemDto<Long>> disk = getValues(MessageFormat.format(DISK_FREE, ip), startM, endM, Long.class);
        for (ValuesItemDto<Long> v : disk) {
            v.setValue(diskMax - v.getValue());
        }

        // 入网网速
        List<ValuesItemDto<Long>> netIn = getValues(MessageFormat.format(NET_IN, ip), startM, endM, Long.class);

        // 出网网速
        List<ValuesItemDto<Long>> netOut = getValues(MessageFormat.format(NET_OUT, ip), startM, endM, Long.class);

        UnStableNodeInfosDTO r =
                new UnStableNodeInfosDTO()
                        .setIp(ip)
                        .setCpuUsage(cpu)
                        .setMemUsage(mem)
                        .setDiskUsage(disk)
                        .setNetIn(netIn)
                        .setNetOut(netOut);
        return r;
    }

    @Override
    public NodeInfoDTO hostLatestInfo(String ip) {
        return hostInfo(ip);
    }

    private <T> List<ValuesItemDto<T>> getValues(String pql, long start, long end, Class<T> clazz) {
        // prometheus 查询
        List<List<String>> result = Optional.ofNullable(promQlRange(pql, start, end))
                .map(ResultDto::getData)
                .map(ResultDataDto::getResult)
                .map(rs -> rs.size() > 0 ? rs.get(0) : null)  // 统计结果只有一条
                .map(DataResultDto::getValues)
                .orElse(Collections.emptyList());

        return result.stream()
                .filter(r -> r.size() == 2)
                .map(r -> {
                    // 时间
                    float second = NumberUtils.toFloat(r.get(0), 0);
                    Date time = new Date((long) (second * 1000));

                    // 返回值
                    try {
                        T o = om.readValue(r.get(1), clazz);
                        return new ValuesItemDto<T>(o, time);
                    } catch (Exception e) {
                        log.warn("convert fail, value:{}, class:{}, e:{}", r.get(1), clazz, e.getMessage());
                        return null;
                    }
                })
                .filter(Objects::nonNull)
                .collect(Collectors.toList());
    }

    private <T> T getValue(String pql, Class<T> clazz, T defaultValue) {
        return Optional.ofNullable(promQl(pql))
                .map(ResultDto::getData)
                .map(ResultDataDto::getResult)
                .map(rs -> !rs.isEmpty() ? rs.get(0) : null)  // 统计结果只有一条
                .map(DataResultDto::getValue)
                .map(ls -> ls.size() > 1 ? ls.get(1) : null)   // 第二个值为结果
                .map(v -> {
                    // 字符串转换成对应数据
                    try {
                        return om.readValue(v, clazz);
                    } catch (Exception e) {
                        log.warn("convert fail, value:{}, class:{}, e:{}", v, clazz, e.getMessage());
                        return null;
                    }
                })
                .orElse(defaultValue);
    }

    private ResultDto promQlRange(String pql, long start, long end) {
        String uri = null, url = null;
        try {
            uri = MessageFormat.format(API_QUERY_RANGE, URLEncoder.encode(pql, "utf-8"), start, end, RANGE_STEP);
            url = MessageFormat.format("http://{0}{1}", domain, uri);
            ResponseEntity<ResultDto> forEntity = new RestTemplate().getForEntity(new URI(url), ResultDto.class);
            return forEntity.getBody();
        } catch (Exception e) {
            log.warn("promQl fail, url:{}, error:{}", url, e.getMessage());
            return null;
        }
    }

    private ResultDto promQl(String pql) {
        String uri = null, url = null;
        try {
            uri = MessageFormat.format(API_QUERY, URLEncoder.encode(pql, "utf-8"));
            url = MessageFormat.format("http://{0}{1}", domain, uri);
            ResponseEntity<ResultDto> forEntity = new RestTemplate().getForEntity(new URI(url), ResultDto.class);
            return forEntity.getBody();
        } catch (Exception e) {
            log.warn("promQl fail, url:{}, error:{}", url, e.getMessage());
            return null;
        }
    }
}
