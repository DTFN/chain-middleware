package com.lingshu.server.dto;

import lombok.Data;
import lombok.EqualsAndHashCode;

import java.io.Serializable;

/**
 * @author: derrick
 * @since: 2025-08-22
 */
@Data
@EqualsAndHashCode(callSuper = false)
public class UpdateAccountRequest implements Serializable {

    private String did;
    private String didDocument;
}
