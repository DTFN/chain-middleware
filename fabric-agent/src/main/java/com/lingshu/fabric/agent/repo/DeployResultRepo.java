package com.lingshu.fabric.agent.repo;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.lingshu.fabric.agent.repo.entity.DeployResultDo;
import com.lingshu.fabric.agent.repo.mapper.DeployResultMapper;
import org.springframework.stereotype.Repository;

/**
 * DeployResultRepo
 *
 * @author XuHang
 * @since 2023/11/28
 **/
@Repository
public class DeployResultRepo extends ServiceImpl<DeployResultMapper, DeployResultDo>{
}
