package com.lingshu.fabric.agent.repo;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.lingshu.fabric.agent.repo.entity.DeployResultNodeDo;
import com.lingshu.fabric.agent.repo.mapper.DeployResultNodeMapper;
import org.springframework.stereotype.Repository;

import java.util.List;

/**
 * DeployResultRepo
 *
 * @author XuHang
 * @since 2023/11/28
 **/
@Repository
public class DeployResultNodeRepo extends ServiceImpl<DeployResultNodeMapper, DeployResultNodeDo>{
    public List<String> allHosts() {
        return this.baseMapper.allHosts();
    }
}
