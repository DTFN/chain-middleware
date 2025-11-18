package com.lingshu.server.dto.resp.account;

import lombok.Data;

/**
 * @author: derrick
 * @since: 2025-09-08
 */


@Data
public class CreateAccountResp {
    private String did;
    private String privateKey;
    private String didDocument;
    private String mnemonic;
}
