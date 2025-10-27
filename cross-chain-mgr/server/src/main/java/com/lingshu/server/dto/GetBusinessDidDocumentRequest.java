package com.lingshu.server.dto;

import lombok.Data;
import lombok.EqualsAndHashCode;

import java.io.Serializable;

/**
 * @author: derrick
 * @since: 2025-09-11
 */
@Data
@EqualsAndHashCode(callSuper = false)
public class GetBusinessDidDocumentRequest implements Serializable {
    private String privateKeyHex;
    private String resourceName;
    private String did;
}
