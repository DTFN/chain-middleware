package com.lingshu.server.dto.resp.busi;

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

/**
 * @author: derrick
 * @since: 2025-09-08
 */


@Data
public class MultiDidDocResp {
    @ApiModelProperty("长安链")
    private String chainMaker;

    @ApiModelProperty("零数链")
    private String lingshu;

    @ApiModelProperty("以太坊")
    private String ethereum;

    @ApiModelProperty("BCOS")
    private String bcos;

    @ApiModelProperty("Fabric")
    private String fabric;
}
