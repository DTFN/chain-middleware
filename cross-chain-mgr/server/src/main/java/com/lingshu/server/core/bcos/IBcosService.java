package com.lingshu.server.core.bcos;

import com.lingshu.server.common.api.OpenAPIResp;
import com.lingshu.server.dto.BcosContractRequest;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

import javax.validation.Valid;

/**
 * @author: derrick
 * @since: 2025-09-05
 */
@FeignClient(url = "${bcos.service.url:http://127.0.0.1:13001}", name = "bcos-service", configuration = BcosFeignConfig.class)
public interface IBcosService {

//    @PostMapping("/v1/openapi/bcos-business/deploy")
//    OpenAPIResp deployContract();
//
//    @PostMapping("/v1/openapi/bcos-business/load")
//    OpenAPIResp load(@Valid @RequestBody LingshuBusinessContractRequest request);
//
//    @PostMapping("/v1/openapi/bcos-business/save")
//    OpenAPIResp save(@Valid @RequestBody BusinessContractSaveRequest request);

    @PostMapping("/v1/openapi/bcos-business/get")
    OpenAPIResp get(@Valid @RequestBody BcosContractRequest request);

//    @PostMapping("/v1/openapi/contract/deploy")
//    OpenAPIResp deploy(@Valid @RequestBody BusinessContractRequest request);

    @PostMapping("/v1/openapi/bcos-business/call")
    OpenAPIResp call(@Valid @RequestBody BcosContractRequest request);
}
