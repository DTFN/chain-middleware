package com.lingshu.fabric.agent.enums;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.ToString;

@Getter
@ToString
@AllArgsConstructor
public enum FabricDeployNodeStage {

	INIT_PEER(1100, "初始化对等节点配置"),
	INIT_ORDERER(1200, "初始化排序节点配置"),

	SCP_CONF(2100, "复制配置文件到节点主机"),
	LAUNCH_NODE(2200, "启动节点"),
	DEPLOY_SUCCESS(2300, "部署完毕"),
	DEPLOY_FAILED(9999, "部署失败"),
	;

	private Integer index;
	private String name;

	public static FabricDeployNodeStage getByIndex(Integer idx) {
		for (FabricDeployNodeStage type : FabricDeployNodeStage.values()) {
			if (type.index.equals(idx)) {
				return type;
			}
		}
		return null;
	}

	public static FabricDeployNodeStage getByName(String desc) {
		for (FabricDeployNodeStage type : FabricDeployNodeStage.values()) {
			if (type.name.equalsIgnoreCase(desc)) {
				return type;
			}
		}
		return null;
	}
}

