package com.lingshu.server.dto.resp.busi;

import lombok.Data;

/**
 * @author: derrick
 * @since: 2025-09-08
 */


@Data
public class VcVerifyResp {
    private Boolean originVerifyResult;
    private Boolean targetVerifyResult;
}
