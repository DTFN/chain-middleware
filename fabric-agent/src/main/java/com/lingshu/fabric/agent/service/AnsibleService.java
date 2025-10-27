package com.lingshu.fabric.agent.service;

import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.lingshu.fabric.agent.config.code.ConstantCode;
import com.lingshu.fabric.agent.config.properties.ConstantProperties;
import com.lingshu.fabric.agent.enums.ScpTypeEnum;
import com.lingshu.fabric.agent.exception.AgentException;
import com.lingshu.fabric.agent.util.CleanPathUtil;
import com.lingshu.fabric.agent.util.ExecuteResult;
import com.lingshu.fabric.agent.util.IPUtil;
import com.lingshu.fabric.agent.util.JavaCommandExecutor;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.tuple.Pair;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.nio.file.Files;
import java.nio.file.Paths;
import java.time.Duration;
import java.time.Instant;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * @author: zehao.song
 */
@Slf4j
@Component
public class AnsibleService {

    @Autowired
    private ConstantProperties constant;

    private static final String NOT_FOUND_FLAG = "not found";
    private static final String FREE_MEMORY_FLAG = "free memory";

    /**
     * check ansible installed
     */
    public void checkAnsible() {
        log.info("checkAnsible installed");
//        String command = "ansible --version";
        String command = "ansible --version | grep \"ansible.cfg\"";
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecShellTimeout());
        if (result.failed()) {
            throw new AgentException(ConstantCode.ANSIBLE_NOT_INSTALLED.attach(result.getExecuteOut()));
        }
    }

    /**
     * ansible exec command
     */
    public void exec(String ip, String command) {
        exec(ip, command, true);
    }

    public void exec(String ip, String command, boolean throwException) {
        String ansibleCommand = String.format("ansible %s -m command -a \"%s\"", ip, command);
        ExecuteResult result = JavaCommandExecutor.executeCommand(ansibleCommand, constant.getExecShellTimeout());
        if (result.failed() && throwException) {
            throw new AgentException(ConstantCode.ANSIBLE_COMMON_COMMAND_ERROR.attach(result.getExecuteOut()));
        }
    }

    /**
     * ansible exec command
     */
    public String execShell(String ip, String shellCmd) {
        String ansibleCommand = String.format("ansible %s -m shell -a %s", ip, shellCmd);
        ExecuteResult result = JavaCommandExecutor.executeCommand(ansibleCommand, constant.getExecShellTimeout());
        if (result.failed()) {
            throw new AgentException(ConstantCode.ANSIBLE_COMMON_COMMAND_ERROR.attach(result.getExecuteOut()));
        }
        String[] split = result.getExecuteOut().split("\n");
        if (split.length <= 1) {
            return null;
        } else {
            return split[1];
        }
    }

    /**
     * check ansible ping, code is always 0(success)
     * @case1: ip configured in ansible, output not empty. ex: 127.0.0.1 | SUCCESS => xxxxx
     * @case2: if ip not in ansible's host, output is empty. ex: Exec command success: code:[0], OUTPUT:[]
     */
    public void execPing(String ip) {
        // ansible bsp(ip) -m ping
        String command = String.format("ansible %s -m ping", ip);
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecShellTimeout());
        // if success
        if (result.getExecuteOut().contains(ip)) {
            log.info("execPing success output:{}", result.getExecuteOut());
            return;
        } else {
            throw new AgentException(ConstantCode.ANSIBLE_PING_NOT_REACH.attach(result.getExecuteOut()));
        }
    }


    /**
     * copy, fetch, file(dir
     * scp: copy from local to remote, fetch from remote to local
     */
    public void scp(ScpTypeEnum typeEnum, String ip, String src, String dst) {
        log.info("scp typeEnum:{},ip:{},src:{},dst:{}", typeEnum, ip, src, dst);
        Instant startTime = Instant.now();
        log.info("scp startTime:{}", startTime.toEpochMilli());
        boolean isSrcDirectory = Files.isDirectory(Paths.get(CleanPathUtil.cleanString(src)));
        boolean isSrcFile = Files.isRegularFile(Paths.get(CleanPathUtil.cleanString(src)));
        // exec ansible copy or fetch
        String command;
        if (typeEnum == ScpTypeEnum.UP) {
            // handle file's dir local or remote
            if (isSrcFile) {
                // if src is file, create parent directory of dst on remote
                String parentOnRemote = Paths.get(CleanPathUtil.cleanString(dst)).getParent().toAbsolutePath().toString();
                this.execCreateDir(ip, parentOnRemote);
            }
            if (isSrcDirectory) {
                // if src is directory, create dst on remote
                this.execCreateDir(ip, dst);
            }
            // synchronized cost less time
            command = String.format("ansible %s -m synchronize -a \"src=%s dest=%s\"", ip, src, dst);
            log.info("exec scp copy command: [{}]", command);
            ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecShellTimeout());
            log.info("scp usedTime:{}", Duration.between(startTime, Instant.now()).toMillis());
            if (result.failed()) {
                throw new AgentException(ConstantCode.ANSIBLE_SCP_COPY_ERROR.attach(result.getExecuteOut()));
            }
        } else { // DOWNLOAD
            // fetch file from remote
            if (isSrcDirectory) {
                // fetch not support fetch directory
                log.error("ansible fetch not support fetch directory!");
                throw new AgentException(ConstantCode.ANSIBLE_FETCH_NOT_DIR);
            }
            // use synchronize, mode=pull
            command = String.format("ansible %s -m synchronize -a \"mode=pull src=%s dest=%s\"", ip, src, dst);
            log.info("exec scp copy command: [{}]", command);
            ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecShellTimeout());
            log.info("scp usedTime:{}", Duration.between(startTime, Instant.now()).toMillis());
            if (result.failed()) {
                throw new AgentException(ConstantCode.ANSIBLE_SCP_FETCH_ERROR.attach(result.getExecuteOut()));
            }
        }
    }

    /**
     * host_check shell
     * @param ip
     * @return
     */
    public void execHostCheckShell(String ip, int nodeCount) {
        log.info("execHostCheckShell ip:{},nodeCount:{}", ip, nodeCount);
        String command = String.format("ansible %s -m script -a \"%s -C %d\"", ip, constant.getHostCheckShell(), nodeCount);
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecShellTimeout());
        if (result.failed()) {
            if (result.getExecuteOut().contains(FREE_MEMORY_FLAG)) {
                throw new AgentException(ConstantCode.EXEC_HOST_CHECK_SCRIPT_ERROR_FOR_MEM.attach(result.getExecuteOut()));
            }
            throw new AgentException(ConstantCode.EXEC_CHECK_SCRIPT_FAIL_FOR_PARAM.attach(result.getExecuteOut()));
        }
    }

    /**
     * @param ip
     */
    public void execDockerCheckShell(String ip) {
        log.info("execDockerCheckShell ip:{}", ip);
        String command = String.format("ansible %s -m script -a \"%s\"", ip, constant.getDockerCheckShell());
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecShellTimeout());
        if (result.failed()) {
            throw new AgentException(ConstantCode.EXEC_DOCKER_CHECK_SCRIPT_ERROR.attach(result.getExecuteOut()));
        }
    }

    /**
     * operate include host_init
     * 1. run host_init_shell
     * 2. mkdir -p ${node_root}
     * //&& sudo chown -R ${user} ${node_root} && sudo chgrp -R ${user} ${node_root}
     * param ip        Required.
     * param chainRoot chain root on host, default is /opt/lingshu/{chain_name}.
     */
    public void execHostInit(String ip, String chainRoot) {
        this.execHostInitScript(ip);
        this.execCreateDir(ip, chainRoot);
    }

    public void execHostInitScript(String ip) {
        log.info("execHostInitScript ip:{}", ip);
        String command = String.format("ansible %s -m script -a \"%s\"", ip, constant.getHostInitShell());
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecShellTimeout());
        if (result.failed()) {
            throw new AgentException(ConstantCode.ANSIBLE_PING_NOT_REACH.attach(result.getExecuteOut()));
        }
    }

    /**
     * mkdir directory on target ip
     * @param ip
     * @param dir absolute path
     * @return
     */
    public ExecuteResult execCreateDir(String ip, String dir) {
        log.info("execCreateDir ip:{},dir:{}", ip, dir);
        // not use sudo to make dir, check access
        String mkdirCommand = String.format("mkdir -p %s", dir);
        String command = String.format("ansible %s -m command -a \"%s\"", ip, mkdirCommand);
        return JavaCommandExecutor.executeCommand(command, constant.getExecShellTimeout());
    }

    /* docker operation: checkImageExists, pullImage run stop */

    /**
     * if not found, ansible exit code not same as script's exit code, use String to distinguish "not found"
     * @param ip
     * @param imageFullName
     * @return
     */
    public boolean checkImageExists(String ip, String imageFullName) {
        log.info("checkImageExists ip:{},imageFullName:{}", ip, imageFullName);
        String command = String.format("ansible %s -m script -a \"%s -i %s\"", ip,
                constant.getAnsibleImageCheckShell(), imageFullName);
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecDockerCheckTimeout());
        if (result.failed()) {
            // NOT FOUND IMAGE
            if (result.getExecuteOut().contains(NOT_FOUND_FLAG)) {
                return false;
            }
            // PARAM ERROR
            if (result.getExitCode() == 2) {
                throw new AgentException(ConstantCode.ANSIBLE_CHECK_DOCKER_IMAGE_ERROR.attach(result.getExecuteOut()));
            }
        }
        // found
        return true;
    }

    public boolean checkContainerExists(String ip, String containerName) {
        log.info("checkContainerExists ip:{},containerName:{}", ip, containerName);

        // docker ps | grep "${containerName}"
        String command = String.format("ansible %s -m script -a \"%s -c %s\"", ip, constant.getAnsibleContainerCheckShell(), containerName);
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecDockerCheckTimeout());
        if (result.failed()) {
            // NOT FOUND CONTAINER
            if (result.getExecuteOut().contains(NOT_FOUND_FLAG)) {
                return false;
            }
            // PARAM ERROR
            if (result.getExitCode() == 2) {
                throw new AgentException(ConstantCode.ANSIBLE_CHECK_CONTAINER_ERROR.attach(result.getExecuteOut()));
            }
        }
        // found
        return true;
    }

    /**
     * pull and load image by cdn
     * @param ip
     * @param outputDir
     * @param bspVersion
     * @return
     */
    public void execPullDockerCdnShell(String ip, String outputDir, String imageTag, String bspVersion) {
        log.info("execPullDockerCdnShell ip:{},outputDir:{},imageTag:{},bspVersion:{}", ip, outputDir, imageTag, bspVersion);
        Instant startTime = Instant.now();
        log.info("execPullDockerCdnShell startTime:{}", startTime.toEpochMilli());
        boolean imageExist = this.checkImageExists(ip, imageTag);
        if (imageExist) {
            log.info("image of {} already exist, jump over pull", imageTag);
            return;
        }
        String command = String.format("ansible %s -m script -a \"%s -d %s -v %s\"", ip, constant.getDockerPullCdnShell(), outputDir, bspVersion);
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecShellTimeout());
        log.info("execPullDockerCdnShell usedTime:{}", Duration.between(startTime, Instant.now()).toMillis());
        if (result.failed()) {
            throw new AgentException(ConstantCode.ANSIBLE_PULL_DOCKER_CDN_ERROR.attach(result.getExecuteOut()));
        }
    }


    public ExecuteResult execDocker(String ip, String dockerCommand) {
        log.info("execDocker ip:{},dockerCommand:{}", ip, dockerCommand);
        Instant startTime = Instant.now();
        String command = String.format("ansible %s -m command -a \"%s\"", ip, dockerCommand);
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getDockerRestartPeriodTime());
        log.info("execDocker usedTime:{}", Duration.between(startTime, Instant.now()).toMillis());
        return result;
    }

    /**
     * mv dir on remote
     */
    public void mvDirOnRemote(String ip, String src, String dst){
        if (StringUtils.isNoneBlank(ip, src, dst)) {
            String rmCommand = String.format("mv -fv %s %s", src, dst);
            log.info("Remove config on remote host:[{}], command:[{}].", ip, rmCommand);
            this.exec(ip, rmCommand);
        }
    }

    /**
     * get exec result in host check
     * @param ip
     * @param ports
     * @return
     */
    public ExecuteResult checkPortArrayInUse(String ip, int ... ports) {
        log.info("checkPortArrayInUse ip:{},ports:{}", ip, ports);
        if (ArrayUtils.isEmpty(ports)){
            return new ExecuteResult(0, "ports input is empty");
        }
        StringBuilder portArray = new StringBuilder();
        for (int port : ports) {
            if (portArray.length() == 0) {
                portArray.append(port);
                continue;
            }
            portArray.append(",").append(port);
        }
        String command = String.format("ansible %s -m script -a \"%s -p %s\"", ip, constant.getHostCheckPortShell(), portArray);
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecShellTimeout());
        return result;
    }

    /**
     * check port, if one is in use, break and return false
     * used in restart chain to make sure process is on
     * @param ip
     * @param portArray
     * @return Pair of <true, port> true: not in use, false: in use
     */
    public Pair<Boolean, Integer> checkPorts(String ip, int ... portArray) {
        log.info("checkPorts ip:{},port:{}", ip, portArray);
        if (ArrayUtils.isEmpty(portArray)){
            return Pair.of(true,0);
        }

        for (int port : portArray) {
            boolean notInUse = checkPortInUse(ip, port);
            // if false, in use
            if (!notInUse){
                return Pair.of(false, port);
            }
        }
        return Pair.of(true,0);
    }

    /**
     * check port by ansible script
     * @param ip
     * @param port
     * @return Pair of <true, port> true: not in use, false: in use
     */
    public boolean checkPortInUse(String ip, int port) {
        log.info("checkPortInUse ip:{},port:{}", ip, port);
        String command = String.format("ansible %s -m script -a \"%s -p %s\"", ip, constant.getHostCheckPortShell(), port);
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecShellTimeout());
        return result.success();
    }

    /**
     * exec on 127.0.0.1
     * check 127.0.0.1 if same with other host ip
     * @param ipList ip to check same with local ip 127.0.0.1
     * @return true-success, false-failed
     */
    public boolean checkLocalIp(List<String> ipList) {
        log.info("checkLoopIp ipArray:{}", ipList);
        if (ipList == null || ipList.isEmpty()){
            return true;
        }
        StringBuilder ipStrArray = new StringBuilder();
        for (String ip : ipList) {
            if (ipStrArray.length() == 0) {
                ipStrArray.append(ip);
                continue;
            }
            ipStrArray.append(",").append(ip);
        }
        // ansible 127.0.0.1 -m script hostCheckIpShell
        String command = String.format("ansible %s -m script -a \"%s -p %s\"", IPUtil.LOCAL_IP_127,
                constant.getHostCheckIpShell(), ipStrArray);
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constant.getExecShellTimeout());
        return result.success();
    }

    public int getAvailablePort(String ip, Integer port) {
        if (port > 65535) {
            throw new AgentException(ConstantCode.NODE_PORT_CONFIG_ERROR);
        }
        ExecuteResult executeResult = checkPortArrayInUse(ip, port);
        while (executeResult.failed()) {
            port++;
            executeResult = checkPortArrayInUse(ip, port);
        }
        return port;
    }

    public Map<Integer, Boolean> getPortUseMap(String ip, int ... portArray) {
        if (ArrayUtils.isEmpty(portArray)){
            return Collections.emptyMap();
        }
        for (int port : portArray) {
            if (port > 65535) {
                throw new AgentException(ConstantCode.NODE_PORT_CONFIG_ERROR);
            }
        }
        ExecuteResult executeResult = checkPortArrayInUse(ip, portArray);
        if (executeResult.failed()) {
            //throw new NodeMgrException(ConstantCode.ANSIBLE_COMMON_COMMAND_ERROR.attach(executeResult.getExecuteOut()));
        }
        Map<Integer, Boolean> portUseMap = new HashMap<>();

        String out = executeResult.getExecuteOut();
        if (out == null || !out.contains("=>")) {
            throw new AgentException(ConstantCode.ANSIBLE_COMMON_COMMAND_ERROR.attach(executeResult.getExecuteOut()));
        }
        String jsonStr = out.substring(out.indexOf("=>") + 2);
        JsonObject jsonObject = (JsonObject) JsonParser.parseString(jsonStr);
        if (jsonObject == null) {
            throw new AgentException(ConstantCode.ANSIBLE_COMMON_COMMAND_ERROR.attach(executeResult.getExecuteOut()));
        }
        String content = jsonObject.get("stdout").getAsString();
        String[] lineArr = content.split("\\n");
        String notPassedMsg = "ERROR: port check NOT PASSED!";
        String inUseMsgStart = "ERROR: port";
        String inUseMsgEnd = "is in use!";
        String readyMsgStart = "port";
        String readyMsgEnd = "is ready";
        for (String line : lineArr) {
            if (line.contains(notPassedMsg)) {
                continue;
            }
            if (line.contains(inUseMsgStart)) { //is in use
                int startIndex = line.indexOf(inUseMsgStart);
                int endIndex = line.indexOf(inUseMsgEnd);
                String portStr = line.substring(startIndex + inUseMsgStart.length(), endIndex);
                Integer port = Integer.parseInt(portStr.trim());
                portUseMap.put(port, true);
            } else if (line.contains(readyMsgEnd)) { // is ready
                int startIndex = line.indexOf(readyMsgStart);
                int endIndex = line.indexOf(readyMsgEnd);
                String portStr = line.substring(startIndex + readyMsgStart.length(), endIndex);
                Integer port = Integer.parseInt(portStr.trim());
                portUseMap.put(port, false);
            }
        }
        return portUseMap;
    }
}
