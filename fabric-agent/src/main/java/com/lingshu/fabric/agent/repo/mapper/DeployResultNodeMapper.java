package com.lingshu.fabric.agent.repo.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.lingshu.fabric.agent.repo.entity.DeployResultNodeDo;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Select;

import java.util.List;

/**
 * DeployResultMapper
 *
 * @author XuHang
 * @Date 2023/11/28 上午11:06
 **/
@Mapper
public interface DeployResultNodeMapper extends BaseMapper<DeployResultNodeDo> {
    @Select("select distinct node_ip from `deploy_result_node`")
    List<String> allHosts();
}
