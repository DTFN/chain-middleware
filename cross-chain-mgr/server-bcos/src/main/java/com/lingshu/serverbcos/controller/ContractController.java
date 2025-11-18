package com.lingshu.serverbcos.controller;

import com.lingshu.serverbcos.common.api.ApiErrorCode;
import com.lingshu.serverbcos.common.api.OpenAPIResp;
import com.lingshu.serverbcos.common.api.OpenAPIRespBuilder;
import com.lingshu.serverbcos.core.bcos.service.ContractService;
import com.lingshu.serverbcos.dto.BusinessContractRequest;
import com.lingshu.serverbcos.dto.ContractRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import javax.validation.Valid;
import java.util.Map;

/**
 * @author: derrick
 * @since: 2025-08-25
 */
@RestController
@RequestMapping("/v1/openapi/contract/")
@Slf4j
public class ContractController {

    @Autowired
    private ContractService contractService;

    //合约部署
    @PostMapping("deploy")
    public OpenAPIResp deploy(@Valid @RequestBody BusinessContractRequest request) {
        try {
            Map<String, Object> result = contractService.deploy(request);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    //合约调用
//    @PostMapping("call")
//    public OpenAPIResp call(@Valid @RequestBody ContractRequest request) {
//        try {
//            Map<String, Object> result = contractService.call(request.getParam());
//            return OpenAPIRespBuilder.success(result);
//        } catch (Exception e) {
//            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
//        }
//    }
}
