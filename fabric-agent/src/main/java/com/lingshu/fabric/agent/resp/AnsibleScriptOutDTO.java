package com.lingshu.fabric.agent.resp;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

import java.util.List;

/**
 * AnsibleScriptOutDTO
 *
 * @author XuHang
 * @since 2023/11/21
 **/
@Data
@JsonInclude(JsonInclude.Include.NON_NULL)
public class AnsibleScriptOutDTO {
    private Boolean changed;
    private Integer rc;
    private String stderr;

    @JsonProperty("stderr_lines")
    private List<String> stderrLines;
    private String stdout;

    @JsonProperty("stdout_lines")
    private List<String> stdoutLines;
}
