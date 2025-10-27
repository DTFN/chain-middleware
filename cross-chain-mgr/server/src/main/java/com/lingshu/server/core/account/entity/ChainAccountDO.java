package com.lingshu.server.core.account.entity;

//import com.baomidou.mybatisplus.annotation.*;
import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

import java.util.Date;

/**
 * @Author wang jian
 * @Date 2021/12/31 10:37
 */
//@TableName(value = "chain_account")
@Data
public class ChainAccountDO {
//    @TableId(value = "id", type = IdType.AUTO)
    private Integer id;
    private Integer shopId;
    private String requestId;
    private String operationId;
    private String did;
    private String didPubkey;
    private String didAddress;
    private String publicKey;
    private String privateKey;
    private String cert;
    // 账户地址
    private String address;
    private Integer source;
    @ApiModelProperty(value = "钱包信息")
    private String wallet;
    @ApiModelProperty(value = "钱包密码")
    private String password;
    //账户状态，参考AccountStatusEnum枚举值
    private Integer accountStatus;
//    @TableField(fill = FieldFill.INSERT)
    private Date createdTime;
//    @TableField(fill = FieldFill.INSERT_UPDATE)
    private Date updatedTime;


}
