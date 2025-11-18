/*
 * Copyright 2014-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package com.lingshu.bsp.front.transaction;

import com.lingshu.bsp.front.base.controller.BaseController;
import com.lingshu.bsp.front.transaction.entity.DIDUpdate;
import com.lingshu.bsp.front.transaction.entity.VcCalReq;
import com.lingshu.bsp.front.util.JsonUtils;
import com.lingshu.chain.sdk.model.TransactionReceipt;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiImplicitParam;
import io.swagger.annotations.ApiOperation;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;

/**
 * TransController.
 * handle transactions with sign to deploy/call contract
 */
@Api(value = "/trans", tags = "transaction interface")
@Slf4j
@RestController
@RequestMapping(value = "/trans")
public class TransController extends BaseController {

    @Autowired
    TransService transServiceImpl;

    /**
     * transHandle through signer
     * @return
     */
    @GetMapping("/addCrossChainMsgListener/eth/{address}")
    public void addEthCrossChainMsgListener(@PathVariable("address") String address) throws Exception {
        transServiceImpl.addEthCrossChainMsgListener(address);
    }

    /**
     * transHandle through signer
     * @return
     */
    @GetMapping("/addCrossChainMsgListener/{address}")
    public void addCrossChainMsgListener(@PathVariable("address") String address) throws Exception {
        transServiceImpl.addCrossChainMsgListener(address);
    }

    /**
     * transHandle through signer
     * @return
     */
    @ApiOperation(value = "transaction handling", notes = "transaction handling")
    @ApiImplicitParam(name = "reqTransHandle", value = "transaction info", required = true, dataType = "ReqTransHandleWithSign")
    @PostMapping("/eth/vc-cal")
    public TransactionReceipt ethVcCal(@Valid @RequestBody VcCalReq req) throws Exception {
        log.info("vcCal, {}", JsonUtils.objToString(req));
        TransactionReceipt obj =  transServiceImpl.ethVcCal(req);
        return obj;
    }

    /**
     * transHandle through signer
     * @return
     */
    @ApiOperation(value = "transaction handling", notes = "transaction handling")
    @ApiImplicitParam(name = "reqTransHandle", value = "transaction info", required = true, dataType = "ReqTransHandleWithSign")
    @PostMapping("/eth/did-update")
    public Object ethDidUpdate(@Valid @RequestBody DIDUpdate didUpdate) throws Exception {
        log.info("eth didUpdate, {}", JsonUtils.objToString(didUpdate));
        Object obj =  transServiceImpl.ethDidUpdate(didUpdate);
        return obj;
    }

    /**
     * transHandle through signer
     * @return
     */
    @ApiOperation(value = "transaction handling", notes = "transaction handling")
    @ApiImplicitParam(name = "reqTransHandle", value = "transaction info", required = true, dataType = "ReqTransHandleWithSign")
    @PostMapping("/vc-cal")
    public TransactionReceipt vcCal(@Valid @RequestBody VcCalReq req) throws Exception {
        log.info("vcCal, {}", JsonUtils.objToString(req));
        TransactionReceipt obj =  transServiceImpl.vcCal(req);
        return obj;
    }

    /**
     * transHandle through signer
     * @return
     */
    @ApiOperation(value = "transaction handling", notes = "transaction handling")
    @ApiImplicitParam(name = "reqTransHandle", value = "transaction info", required = true, dataType = "ReqTransHandleWithSign")
    @PostMapping("/did-update")
    public Object didUpdate(@Valid @RequestBody DIDUpdate didUpdate) throws Exception {
        log.info("lingshu didUpdate, {}", JsonUtils.objToString(didUpdate));
        Object obj =  transServiceImpl.didUpdate(didUpdate);
        return obj;
    }

}
