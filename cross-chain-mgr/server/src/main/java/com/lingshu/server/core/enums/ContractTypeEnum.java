package com.lingshu.server.core.enums;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;

import java.util.Arrays;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;

/**
 * @author: derrick
 * @since: 2025-09-05
 */
@NoArgsConstructor
@AllArgsConstructor
@Getter
public enum ContractTypeEnum {
    //合约类型( 1-长安链solidity,2-长安链wasm,3-零数链solidity,4-以太坊solidity,5-bcos的solidity,6-fabric的wasm )
//    CHAINMAKER_SOLIDITY("1", "长安链solidity"),
    CHAINMAKER_WASM("2", "长安链wasm"),
    LINGSHU_SOLIDITY("3", "零数链solidity"),
    ETHEREUM_SOLIDITY("4", "以太坊solidity"),
    BCOS_SOLIDITY("5", "bcos的solidity"),
    BSC_SOLIDITY("6", "bsc的solidity"),
    FABRIC_WASM("7", "fabric的wasm"),
    ;
    private String code;
    private String desc;

    public static String getDescByCode(String code) {
        for (ContractTypeEnum value : ContractTypeEnum.values()) {
            if (value.getCode().equals(code)) {
                return value.getDesc();
            }
        }
        return null;
    }

    public static ContractTypeEnum getByCode(String code) {
        for (ContractTypeEnum value : ContractTypeEnum.values()) {
            if (value.getCode().equals(code)) {
                return value;
            }
        }
        return null;
    }

    public static boolean isEd25519Sign(ContractTypeEnum type) {
        if (CHAINMAKER_WASM.equals(type) || FABRIC_WASM.equals(type)) {
            return true;
        }
        return false;
    }

    public static boolean isEthSign(ContractTypeEnum type) {
        List<ContractTypeEnum> list = Arrays.asList(
                LINGSHU_SOLIDITY, ETHEREUM_SOLIDITY, BCOS_SOLIDITY, BSC_SOLIDITY
        );
        return new HashSet<ContractTypeEnum>(list).contains(type);
    }

}
