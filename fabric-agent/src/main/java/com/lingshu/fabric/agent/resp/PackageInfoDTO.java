package com.lingshu.fabric.agent.resp;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class PackageInfoDTO {
    private String lang;
    private String label;
    private String filePath;
}
