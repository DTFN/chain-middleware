package com.lingshu.fabric.agent.repo;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.lingshu.fabric.agent.repo.entity.OrgInfoDo;
import com.lingshu.fabric.agent.repo.mapper.OrgInfoMapper;
import org.springframework.stereotype.Repository;

/**
 * ChainInfoMapper
 *
 * @author XuHang
 * @Date 2023/11/29 下午1:34
 **/
@Repository
public class OrgInfoRepo extends ServiceImpl<OrgInfoMapper, OrgInfoDo> {
}
