package com.lingshu.bsp.front.rpcapi.eth;

import cn.hutool.core.bean.BeanUtil;
import com.lingshu.bsp.front.dto.RequestInvoke;
import com.lingshu.bsp.front.dto.TransactionDetail;
import com.lingshu.chain.sdk.client.protocol.model.JsonTransactionResponse;
import com.lingshu.chain.sdk.model.TransactionReceipt;
import com.lingshu.chain.sdk.util.NumberUtil;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiImplicitParam;
import io.swagger.annotations.ApiOperation;
import org.springframework.web.bind.annotation.*;

import javax.annotation.Resource;
import java.math.BigInteger;

/**
 * @author gongrui.wang
 * @since 2025/8/27
 */
@Api(value = "/eth-rpc", tags = "eth-rpc interface")
@RestController
@RequestMapping(value = "/eth/rpc")
public class EthereumController {

    @Resource
    private EthApiService ethApiService;

    @ApiOperation(value = "getBlockNumber", notes = "Get the latest block height of the node")
    @GetMapping("/blockNumber")
    public BigInteger getBlockNumber() {
        return ethApiService.getBlockNumber();
    }


    @ApiOperation(value = "getTransactionByHash",
            notes = "Get transaction information based on transaction hash")
    @ApiImplicitParam(name = "transHash", value = "transHash", required = true, dataType = "String",
            paramType = "path")
    @GetMapping("/transaction/{transHash}")
    public TransactionDetail getTransactionByHash(@PathVariable String transHash) {
        JsonTransactionResponse response = ethApiService.getTransactionByHash(transHash);
        TransactionDetail transactionDetail = BeanUtil.copyProperties(response, TransactionDetail.class);
        transactionDetail.setBlockNumber("0x"+ NumberUtil.toHexStringNoPrefix(response.getBlockNumber()));
        return transactionDetail;
    }


    @ApiOperation(value = "contract invoke", notes = "contract invoke")
    @PostMapping("/invoke")
    public TransactionReceipt invoke(@RequestBody RequestInvoke requestInvoke)
            throws Exception {
        TransactionReceipt transactionReceipt = ethApiService.crossSave(requestInvoke.getContractAddress());
        transactionReceipt.setContractAddress(requestInvoke.getContractAddress());
        return transactionReceipt;
    }


    @ApiOperation(value = "contract deploy", notes = "contract deploy")
    @PostMapping("/deploy")
    public String deploy()
            throws Exception {
        String address = ethApiService.deploy();
        return address;
    }

    @ApiOperation(value = "contract deploy", notes = "contract deploy")
    @PostMapping("busi/deploy")
    public String deployBusiCenter() throws Exception {
        return ethApiService.deployBusiCenter();
    }

    @ApiOperation(value = "contract invoke", notes = "contract invoke")
    @PostMapping("busi/invoke")
    public TransactionReceipt invokeBusiCenter(@RequestBody RequestInvoke requestInvoke)
            throws Exception {
        TransactionReceipt transactionReceipt = ethApiService.invokeBusiCenter(requestInvoke.getContractAddress(), requestInvoke.getArgs());
        transactionReceipt.setContractAddress(requestInvoke.getContractAddress());
        return transactionReceipt;
    }
}
