package com.lingshu.fabric.agent.service;

import com.lingshu.fabric.agent.req.chaincode.FabricCCApproveReq;
import com.lingshu.fabric.agent.req.chaincode.FabricCCCommitReq;
import com.lingshu.fabric.agent.req.chaincode.FabricCCInstallReq;
import com.lingshu.fabric.agent.req.chaincode.FabricCCInvokeReq;
import com.lingshu.fabric.agent.resp.FabricChainCodeInvokeResp;
import com.lingshu.fabric.agent.resp.PackageInfoDTO;
import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;
import org.hyperledger.fabric.sdk.exception.ProposalException;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeoutException;

public interface ChainCodeService {

    PackageInfoDTO resolveChaincodePackage(MultipartFile chainCodeFile);

    String installChainCode(FabricCCInstallReq req) throws InvalidArgumentException, ProposalException, IOException;

    long approveChainCode(FabricCCApproveReq req) throws InvalidArgumentException, ProposalException, ExecutionException, InterruptedException, TimeoutException;

    void commitChainCode(FabricCCCommitReq req) throws InvalidArgumentException, ProposalException, ExecutionException, InterruptedException, TimeoutException;

    FabricChainCodeInvokeResp invokeChainCode(FabricCCInvokeReq req);

    Object queryChainCode(FabricCCInvokeReq req);
}
