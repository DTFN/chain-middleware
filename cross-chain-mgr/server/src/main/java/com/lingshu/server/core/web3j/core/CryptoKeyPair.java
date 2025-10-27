package com.lingshu.server.core.web3j.core;

import lombok.Data;

/**
 * @author gongrui.wang
 * @since 2025/2/12
 */
@Data
public class CryptoKeyPair {

    private String address;
    private byte[] publicKey;
    private byte[] privateKey;
    private byte[] cert;

    private String hexedPrivateKey;
}
