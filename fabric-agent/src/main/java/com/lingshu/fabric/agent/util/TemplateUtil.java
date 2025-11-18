package com.lingshu.fabric.agent.util;

import com.alibaba.fastjson2.JSON;
import com.alibaba.fastjson2.JSONObject;
import com.lingshu.fabric.agent.req.GenCaServerConfReq;
import com.lingshu.fabric.agent.req.GenCaServerReq;
import com.lingshu.fabric.agent.req.GenConfigTxReq;
import com.lingshu.fabric.agent.req.GenMspConfReq;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.tuple.Pair;
import org.thymeleaf.TemplateEngine;
import org.thymeleaf.context.Context;
import org.thymeleaf.templatemode.TemplateMode;
import org.thymeleaf.templateresolver.ClassLoaderTemplateResolver;

import java.io.File;
import java.io.Writer;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Locale;
import java.util.Map;
import java.util.stream.Collectors;
import java.util.stream.StreamSupport;

/**
 * @author lin
 * @since 2023-11-16
 */
public class TemplateUtil {

    public static final String CONFIG_TX_YAML = "configtx-tpl.yaml";
    public static final String CA_SERVER_CONF_YAML = "ca-server-config-tpl.yaml";
    public static final String MSP_CONF_YAML = "config-tpl.yaml";
    public static final String CA_SERVER_DOCKER_YAML = "ca-server-docker-tpl.yaml";
    public final static TemplateEngine templateEngine = new TemplateEngine();

    static {
        final ClassLoaderTemplateResolver templateResolver = new ClassLoaderTemplateResolver();
        templateResolver.setOrder(1);
        templateResolver.setPrefix("/template/");
        templateResolver.setTemplateMode(TemplateMode.TEXT);
        templateResolver.setCharacterEncoding("UTF-8");
        templateResolver.setCacheable(false);

        templateEngine.addTemplateResolver(templateResolver);
    }

    /**
     * @param tpl
     * @param varMap
     */
    public static String generate(String tpl, Map<String, Object> varMap) {
        final Context ctx = new Context(Locale.CHINA);
        ctx.setVariables(varMap);
        String fileContent =  templateEngine.process(tpl, ctx);
        return fileContent;
    }

    /**
     * @param tpl
     * @param array
     */
    public static String generate(String tpl, Pair<String, Object>... array) {
        return generate(tpl, toMap(array));
    }

    /**
     * @param tpl
     * @param varMap
     */
    public static void write(String tpl, Writer writer, Map<String, Object> varMap) {
        final Context ctx = new Context(Locale.CHINA);
        ctx.setVariables(varMap);
        templateEngine.process(tpl, ctx, writer);
    }

    /**
     * @param tpl
     * @param array
     * @return
     */
    public static void write(String tpl, Writer writer, Pair<String, Object>... array) {
        write(tpl, writer, toMap(array));
    }

    private static <K, V> Map<K, V> toMap(Pair<K, V> ... array) {
        if (ArrayUtils.isNotEmpty(array)) {
            return toMap(Arrays.asList(array));
        }
        return new HashMap<>();
    }

    private static <K, V> Map<K, V> toMap(Iterable<Pair<K, V>> array) {
        if (array != null) {
            return StreamSupport.stream(array.spliterator(), false)
                    .collect(Collectors.toMap(Pair::getKey, Pair::getValue));
        }
        return new HashMap<>();
    }

    public static String genConfigTxYaml(GenConfigTxReq req) {
        String configtxYaml = TemplateUtil.generate(CONFIG_TX_YAML,
                Pair.of("chainPath", req.getChainPath())
                , Pair.of("orderOrg", req.getOrderOrg())
                , Pair.of("peerOrg", req.getPeerOrg())
                , Pair.of("ordererList", req.getOrdererList())
                , Pair.of("peerList", req.getPeerList())
                , Pair.of("channelList", req.getChannelList())
        );
        return configtxYaml;
    }

    public static void newConfigTxYaml(Path chainPath, GenConfigTxReq req) {
        if (chainPath == null ) {
            throw new RuntimeException("chainPath cannot be null.");
        }
        File file = chainPath.toFile();
        if (!file.exists()) {
            file.mkdirs();
        }
        String configtxYaml = genConfigTxYaml(req);

        try {
            Files.write(chainPath.resolve("config").resolve("configtx.yaml"), configtxYaml.getBytes());
        } catch (Exception e) {
            throw new RuntimeException("fail to generate configtx.yaml");
        }

    }

    public static String genCaServerConfigYaml(GenCaServerConfReq req) {
        String caServerConfigYaml = TemplateUtil.generate(CA_SERVER_CONF_YAML,
                Pair.of("caName", req.getCaName())
                ,Pair.of("caIp", req.getCaIp())
                ,Pair.of("caPort", req.getCaPort())
                ,Pair.of("caAdmin", req.getCaAdmin())
                ,Pair.of("caPw", req.getCaPw())
                ,Pair.of("caDomain", req.getCaDomain())
                ,Pair.of("orgDomain", req.getOrgDomain())
        );
        return caServerConfigYaml;
    }

    public static void newCaServerConfigYaml(Path caServerPath, GenCaServerConfReq req) {
        if (caServerPath == null ) {
            throw new RuntimeException("caServerPath cannot be null.");
        }
        File file = caServerPath.toFile();
        if (!file.exists()) {
            file.mkdirs();
        }
        String caServerConfigYaml = genCaServerConfigYaml(req);

        try {
            Files.write(caServerPath.resolve("fabric-ca-server-config.yaml"), caServerConfigYaml.getBytes());
        } catch (Exception e) {
            throw new RuntimeException("fail to generate ca-server-config.yaml");
        }
    }

    public static String genMspConfigYaml(GenMspConfReq req) {
        String configYaml = TemplateUtil.generate(MSP_CONF_YAML,
                Pair.of("caCertName", req.getCaCertName())
        );
        return configYaml;
    }

    public static void newConfigYaml(Path chainPath, GenMspConfReq req) {
        if (chainPath == null ) {
            throw new RuntimeException("chainPath cannot be null.");
        }
        File file = chainPath.toFile();
        if (!file.exists()) {
            file.mkdirs();
        }
        String configYaml = genMspConfigYaml(req);

        try {
            Files.write(chainPath.resolve("config.yaml"), configYaml.getBytes());
        } catch (Exception e) {
            throw new RuntimeException("fail to generate config.yaml");
        }

    }

    public static String genCaServerYaml(GenCaServerReq req) {
        String configYaml = TemplateUtil.generate(CA_SERVER_DOCKER_YAML,
                (JSONObject) JSON.toJSON(req)
        );
        return configYaml;
    }

    public static void newCaServerDockerYaml(Path caServerPath, GenCaServerReq req) {
        if (caServerPath == null ) {
            throw new RuntimeException("caServerPath cannot be null.");
        }
        File file = caServerPath.toFile();
        if (!file.exists()) {
            file.mkdirs();
        }
        String caServerDockerYaml = genCaServerYaml(req);

        try {
            Files.write(caServerPath.resolve("docker-compose.yaml"), caServerDockerYaml.getBytes());
        } catch (Exception e) {
            throw new RuntimeException("fail to generate ca-server docker-compose.yaml");
        }
    }
}
