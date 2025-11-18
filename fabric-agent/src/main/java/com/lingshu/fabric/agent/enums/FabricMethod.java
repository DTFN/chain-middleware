package com.lingshu.fabric.agent.enums;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.ToString;

@Getter
@ToString
@AllArgsConstructor
public enum FabricMethod {

	INVOKE(0, "invoke"),
	QUERY(1, "query"),
	;

	private Integer index;
	private String name;

	public static FabricMethod getByIndex(Integer idx) {
		for (FabricMethod type : FabricMethod.values()) {
			if (type.index.equals(idx)) {
				return type;
			}
		}
		return null;
	}

	public static FabricMethod getByName(String desc) {
		for (FabricMethod type : FabricMethod.values()) {
			if (type.name.equalsIgnoreCase(desc)) {
				return type;
			}
		}
		return null;
	}
}

