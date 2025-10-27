package com.lingshu.fabric.agent.repo;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.lingshu.fabric.agent.repo.entity.MonitorNodeDo;
import com.lingshu.fabric.agent.repo.mapper.MonitorNodeMapper;
import org.springframework.stereotype.Repository;

/**
 * DeployResultRepo
 *
 * @author XuHang
 * @since 2023/11/28
 **/
@Repository
public class MonitorNodeRepo extends ServiceImpl<MonitorNodeMapper, MonitorNodeDo>{
}
