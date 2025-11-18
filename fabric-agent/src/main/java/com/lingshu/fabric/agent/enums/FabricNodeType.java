package com.lingshu.fabric.agent.enums;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.ToString;

@Getter
@ToString
@AllArgsConstructor
public enum FabricNodeType {

	ORDERER(1, "ORDERER"),
	PEER(2, "PEER"),
	;

	private Integer index;
	private String name;

	public static FabricNodeType getByIndex(Integer idx) {
		for (FabricNodeType type : FabricNodeType.values()) {
			if (type.index.equals(idx)) {
				return type;
			}
		}
		return null;
	}

	public static FabricNodeType getByName(String desc) {
		for (FabricNodeType type : FabricNodeType.values()) {
			if (type.name.equalsIgnoreCase(desc)) {
				return type;
			}
		}
		return null;
	}
}

