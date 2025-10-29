package com.lingshu.bsp.front.dto;

import lombok.Data;

import java.util.List;

@Data
public class RequestInvoke {

    private String contractAddress;
    private String abi;
    private String method;
    private List<Object> args;
}
