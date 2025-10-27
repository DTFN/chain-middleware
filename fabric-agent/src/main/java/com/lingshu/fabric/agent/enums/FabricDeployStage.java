package com.lingshu.fabric.agent.enums;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.ToString;

@Getter
@ToString
@AllArgsConstructor
public enum FabricDeployStage {

	INIT_CA(1000, "初始化CA"),
	INIT_CLIENT_ENV(1100, "初始化客户端环境"),
	INIT_ORG(1200, "初始化机构配置"),
	INIT_PEER(1300, "初始化对等节点配置"),
	INIT_ORDERER(1400, "初始化排序节点配置"),
	GEN_CONF_TX(1500, "生成configTx配置"),
	GEN_GENESIS(1600, "生成创世块配置"),

	SCP_CONF(2100, "复制配置文件到节点主机"),
	LAUNCH_NODE(2200, "启动节点"),
	CREATE_CHANNEL(2300, "创建应用通道"),
	JOIN_CHANNEL(2400, "加入应用通道"),
	DEPLOY_SUCCESS(2500, "部署完毕"),
	DEPLOY_FAILED(9999, "部署失败"),
	;

	private Integer index;
	private String name;

	public static FabricDeployStage getByIndex(Integer idx) {
		for (FabricDeployStage type : FabricDeployStage.values()) {
			if (type.index.equals(idx)) {
				return type;
			}
		}
		return null;
	}

	public static FabricDeployStage getByName(String desc) {
		for (FabricDeployStage type : FabricDeployStage.values()) {
			if (type.name.equalsIgnoreCase(desc)) {
				return type;
			}
		}
		return null;
	}
}

