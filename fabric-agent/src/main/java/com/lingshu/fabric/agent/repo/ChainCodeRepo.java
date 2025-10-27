package com.lingshu.fabric.agent.repo;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.lingshu.fabric.agent.repo.entity.ChainCodeDo;
import com.lingshu.fabric.agent.repo.entity.ChainCodePeerDo;
import com.lingshu.fabric.agent.repo.entity.MonitorNodeDo;
import com.lingshu.fabric.agent.repo.mapper.ChainCodeMapper;
import com.lingshu.fabric.agent.repo.mapper.MonitorNodeMapper;
import org.hyperledger.fabric.sdk.Peer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Repository;

import java.util.Collection;
import java.util.List;

/**
 * DeployResultRepo
 *
 * @author XuHang
 * @since 2023/11/28
 **/
@Repository
public class ChainCodeRepo extends ServiceImpl<ChainCodeMapper, ChainCodeDo>{

    @Autowired
    private ChainCodePeerRepo chainCodePeerRepo;

    public void save(String channelId, String chainCodeName, String lang, Collection<Peer> peers){
        LambdaQueryWrapper<ChainCodeDo> queryWrapper = new LambdaQueryWrapper<>();
        queryWrapper.eq(ChainCodeDo::getChainCodeName, chainCodeName);
        queryWrapper.eq(ChainCodeDo::getChannelId, channelId);
        ChainCodeDo chainCodeDo = this.getOne(queryWrapper);
        if (chainCodeDo == null){
            chainCodeDo = new ChainCodeDo();
            chainCodeDo.setChainCodeName(chainCodeName);
            chainCodeDo.setLang(lang);
            chainCodeDo.setChannelId(channelId);
            this.save(chainCodeDo);
            chainCodeDo = this.getOne(queryWrapper);
        }

        for (Peer peer : peers){
            ChainCodePeerDo chainCodePeerDo = new ChainCodePeerDo();
            chainCodePeerDo.setChainCodeId(chainCodeDo.getId());
            chainCodePeerDo.setPeerName(peer.getName());
            chainCodePeerDo.setPeerUrl(peer.getUrl());
            chainCodePeerRepo.savePeer(chainCodePeerDo);
        }
    }

    public List<ChainCodePeerDo> getChainCodePeers(String channelId, String chainCodeName){
        LambdaQueryWrapper<ChainCodeDo> queryWrapper = new LambdaQueryWrapper<>();
        queryWrapper.eq(ChainCodeDo::getChainCodeName, chainCodeName);
        queryWrapper.eq(ChainCodeDo::getChannelId, channelId);
        ChainCodeDo chainCodeDo = this.getOne(queryWrapper);
        if (chainCodeDo == null){
            return null;
        }
        return chainCodePeerRepo.queryChainCodePeers(chainCodeDo.getId());
    }
}
