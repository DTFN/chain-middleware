package com.lingshu.serverbcos.core.enums;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;

/**
 * @author: derrick
 * @since: 2025-09-05
 */
@NoArgsConstructor
@AllArgsConstructor
@Getter
public enum ContractTypeEnum {
    //合约类型( 1-长安链solidity,2-长安链wasm,3-零数链solidity,4-以太坊solidity,5-bcos的solidity,6-fabric的wasm )
    CHAINMAKER_SOLIDITY("1", "长安链solidity"),
    CHAINMAKER_WASM("2", "长安链wasm"),
    LINGSHU_SOLIDITY("3", "零数链solidity"),
    ETHEREUM_SOLIDITY("4", "以太坊solidity"),
    BCOS_SOLIDITY("5", "bcos的solidity"),
    BSC_SOLIDITY("6", "bsc的solidity");
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

}
