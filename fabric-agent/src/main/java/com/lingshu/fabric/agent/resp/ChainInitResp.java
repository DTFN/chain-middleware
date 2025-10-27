package com.lingshu.fabric.agent.resp;

import com.lingshu.fabric.agent.enums.FabricNodeType;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.experimental.Accessors;

import java.util.List;

/**
 * FabricInitVO
 *
 * @author XuHang
 * @since 2023/11/15
 **/
@Data
@Accessors(chain = true)
@NoArgsConstructor
@AllArgsConstructor
public class ChainInitResp {
    private List<Result> results;

    @Data
    @Accessors(chain = true)
    @AllArgsConstructor
    @NoArgsConstructor
    public static class Result {
        private String ip;
        private FabricNodeType type;
        private int number;
        private Boolean result = Boolean.FALSE;
        private String errorReason;
    }
}
