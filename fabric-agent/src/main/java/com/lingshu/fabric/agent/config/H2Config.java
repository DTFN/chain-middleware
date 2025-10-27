package com.lingshu.fabric.agent.config;

import com.lingshu.fabric.agent.repo.mapper.DbInfoMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.ibatis.session.SqlSession;
import org.apache.ibatis.session.SqlSessionFactory;
import org.springframework.beans.factory.InitializingBean;
import org.springframework.core.io.Resource;
import org.springframework.core.io.ResourceLoader;
import org.springframework.stereotype.Component;
import org.springframework.util.StreamUtils;

import java.nio.charset.StandardCharsets;
import java.sql.Connection;

/**
 * H2Config
 * 第一次运行的时候执行数据库
 *
 * @author XuHang
 * @since 2023/11/28
 **/
@Slf4j
@Component
@RequiredArgsConstructor
public class H2Config implements InitializingBean {
    private final ResourceLoader rl;
    private final SqlSessionFactory sqlSessionFactory;
    private final DbInfoMapper dbInfoMapper;

    @Override
    public void afterPropertiesSet() throws Exception {
        // 判断是否初次执行
        if (dbInfoMapper.showTables().size() == 0) {
            Resource ddlRes = rl.getResource("classpath:sql/ddl.sql");
            Resource dmlRes = rl.getResource("classpath:sql/dml.sql");
            String ddl = StreamUtils.copyToString(ddlRes.getInputStream(), StandardCharsets.UTF_8);
            String dml = StreamUtils.copyToString(dmlRes.getInputStream(), StandardCharsets.UTF_8);

            // 执行初始化脚本
            log.info("init db");
            try (SqlSession sqlSession = sqlSessionFactory.openSession(true)) {
                try (Connection connection = sqlSession.getConnection()) {
                    connection.prepareStatement(ddl).execute();
                    connection.prepareStatement(dml).execute();
                }
            } catch (Exception e) {
                log.error("init db error", e);
                throw e;
            }
        } else {
            log.info("skip init db");
        }
    }
}
