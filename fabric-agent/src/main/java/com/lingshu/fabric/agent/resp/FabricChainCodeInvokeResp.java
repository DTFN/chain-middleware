package com.lingshu.fabric.agent.resp;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class FabricChainCodeInvokeResp {

    private int status;
    private Object payload;
}
