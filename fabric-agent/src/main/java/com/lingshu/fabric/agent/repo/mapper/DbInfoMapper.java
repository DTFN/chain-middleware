package com.lingshu.fabric.agent.repo.mapper;

import com.lingshu.fabric.agent.repo.entity.ShowTableDo;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Select;

import java.util.List;

/**
 * DbInfoMapper
 * 数据库相关
 *
 * @author XuHang
 * @Date 2023/11/28 下午2:44
 **/
@Mapper
public interface DbInfoMapper {
    @Select("show tables")
    List<ShowTableDo> showTables();
}
