package com.lingshu.fabric.agent.repo;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.service.additional.query.impl.LambdaQueryChainWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.lingshu.fabric.agent.repo.entity.ChainCodePeerDo;
import com.lingshu.fabric.agent.repo.entity.MonitorNodeDo;
import com.lingshu.fabric.agent.repo.mapper.ChainCodePeerMapper;
import com.lingshu.fabric.agent.repo.mapper.MonitorNodeMapper;
import org.springframework.stereotype.Repository;

import java.util.List;

/**
 * DeployResultRepo
 *
 * @author XuHang
 * @since 2023/11/28
 **/
@Repository
public class ChainCodePeerRepo extends ServiceImpl<ChainCodePeerMapper, ChainCodePeerDo>{

    public void savePeer(ChainCodePeerDo chainCodePeerDo){
        LambdaQueryWrapper<ChainCodePeerDo> query = new LambdaQueryWrapper<>();
        query.eq(ChainCodePeerDo::getChainCodeId, chainCodePeerDo.getChainCodeId());
        query.eq(ChainCodePeerDo::getPeerName, chainCodePeerDo.getPeerName());
        ChainCodePeerDo one = this.getOne(query);
        if (one == null){
            this.save(chainCodePeerDo);
        }
    }

    public List<ChainCodePeerDo> queryChainCodePeers(long chainCodeId){
        LambdaQueryWrapper<ChainCodePeerDo> query = new LambdaQueryWrapper<>();
        query.eq(ChainCodePeerDo::getChainCodeId, chainCodeId);
        return this.list(query);
    }
}
