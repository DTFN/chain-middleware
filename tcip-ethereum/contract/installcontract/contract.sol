pragma solidity ^0.4.25;

contract CrossSave {
    // kv保存
    mapping (string => string) saver;
    mapping (string => string) saverFlag;

    // 事件，用来通知跨链事件
    event Test(string key, string value);

    // 这里是构造函数, 实例创建时候执行
    function CrossChainSave(string key, string value) public {
        saver[key] = value;
        saverFlag[key] = "false";

        emit Test(key, value);
    }

    function CrossChainTry(string key, string value) public {
        saver[key] = value;
        saverFlag[key] = "false";
    }

    function CrossChainConfirm(string key) public {
        saverFlag[key] = "true";
    }

    function CrossChainCancel(string key) public {
        saverFlag[key] = "failed";
    }

    function query(string key) public returns (string value, string flag) {
        return (saver[key], saverFlag[key]);
    }
}