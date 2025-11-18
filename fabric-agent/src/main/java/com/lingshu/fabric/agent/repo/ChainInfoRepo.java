package com.lingshu.fabric.agent.repo;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.lingshu.fabric.agent.repo.entity.ChainInfoDo;
import com.lingshu.fabric.agent.repo.mapper.ChainInfoMapper;
import org.springframework.stereotype.Repository;

/**
 * ChainInfoMapper
 *
 * @author XuHang
 * @Date 2023/11/29 下午1:34
 **/
@Repository
public class ChainInfoRepo extends ServiceImpl<ChainInfoMapper, ChainInfoDo> {

    public void removeByChainUid(String chainUid) {
        LambdaQueryWrapper<ChainInfoDo> queryWrapper = new LambdaQueryWrapper<>();
        queryWrapper.eq(ChainInfoDo::getChainUid, chainUid);
        this.remove(queryWrapper);
    }
}
