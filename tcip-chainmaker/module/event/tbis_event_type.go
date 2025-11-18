/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	errorCrossIDEmpty       = errors.New("the crossID is empty")
	errorEventDataEmpty     = errors.New("the content of event data is empty")
	errorIdentityEmpty      = errors.New("the identity is empty")
	errorContractConfig     = errors.New("the config of contract is error")
	errorChainEmpty         = errors.New("the chain is empty")
	errorParamsFormatConfig = errors.New("the params format of contract is error")
)

const (
	zeroIndex      = 0
	firstIndex     = 1
	keyValueLength = 2
	aBI            = "abi"
	groupSplit     = ";"
	keyValueSplit  = ":"
	segmentSplit   = ","

	crossIDKey   = "CrossID"
	eventNameKey = "EventName"
	businessKey  = "Business"

	commitContractNameKey         = "CommitContractName"
	commitContractVersionKey      = "CommitContractVersion"
	commitContractMethodKey       = "CommitContractMethod"
	commitContractParamsFormatKey = "CommitContractParamsFormat"
	commitContractParamsKey       = "CommitContractParams"
	commitIdentityKey             = "CommitIdentity"
	commitExtraDataKey            = "CommitExtraData"

	chainRIDKey             = "ChainRID"
	identityKey             = "Identity"
	contractNameKey         = "ContractName"
	contractVersionKey      = "ContractVersion"
	contractMethodKey       = "ContractMethod"
	contractParamsFormatKey = "ContractParamsFormat"
	contractParamsKey       = "ContractParams"
	contractExtraDataKey    = "ContractExtraData"

	paramSplit  = "@"
	keyValSplit = "#"

	// TbisFlag tbis协议标记
	TbisFlag = "tbis"

	// SubSuccess tbis协议成功
	SubSuccess = 1
	// SubFailed tbis协议失败
	SubFailed = 2

	// SubResult tbis协议结果字段
	SubResult = "sub_results"
)

// ContractListenData 由合约事件监听到的内容转换后获得
type ContractListenData struct {
	CrossID   string // 跨链ID
	EventName string // 事件名称
	EventData string // 事件数据
	Business  string // 跨链事件对应的业务关键字
	CommitCtt *CommitContract
	SubEvents []*CrossSubEvent
}

// CommitContract 提交的合约信息
type CommitContract struct {
	Identity      string // 提交时的身份
	Ctt           *InnerContract
	ParamsFormat  string   // 参数格式化
	ParamsData    []string // 跨链事件写入的内容
	ParamsContent string   // 转换后的具体参数
}

// CrossSubEvent 子链跨链事件
type CrossSubEvent struct {
	ChainRID  string
	GatewayID string
	ExecIdx   int
	ChainCtt  *CommitContract
}

// InnerContract Inner Contract
type InnerContract struct {
	Name      string            // 合约名称
	Version   string            // 合约版本号信息
	Method    string            // 合约方法
	ExtraData map[string]string // 扩展的一些数据，例如ABI
}

// NewEmptyContractListenData 新建一个空的listen数据
//  @return *ContractListenData
func NewEmptyContractListenData() *ContractListenData {
	return &ContractListenData{}
}

// CheckAndInit 检查并初始化内部参数
//  @receiver d
//  @return error
func (d *ContractListenData) CheckAndInit() error {
	if d.CrossID == "" {
		return errorCrossIDEmpty
	}
	if d.EventData == "" {
		return errorEventDataEmpty
	}
	if err := d.CommitCtt.Check(); err != nil {
		return err
	}
	for _, subEvent := range d.SubEvents {
		if err := subEvent.CheckAndInit(); err != nil {
			return err
		}
	}
	return nil
}

// CheckAvailability 检查链签名身份可用性
//  @receiver d
//  @param chainRID
//  @return error
func (d *ContractListenData) CheckAvailability(chainRID string) error {
	return d.CommitCtt.CheckIdentity(chainRID)
}

// AppendSubEvent 添加SubEvent
//  @receiver d
//  @param crossSubEvent
func (d *ContractListenData) AppendSubEvent(crossSubEvent *CrossSubEvent) {
	d.SubEvents = append(d.SubEvents, crossSubEvent)
}

// GetSubEvents 返回subEvents
//  @receiver d
//  @return []*CrossSubEvent
func (d *ContractListenData) GetSubEvents() []*CrossSubEvent {
	return d.SubEvents
}

// NewEmptyCommitContract 空的提交对象
//  @return *CommitContract
func NewEmptyCommitContract() *CommitContract {
	return &CommitContract{
		Ctt: NewEmptyInnerContract(),
	}
}

// NewCommitContract 新建一个提交对象
//  @param identity
//  @param cttName
//  @param cttVersion
//  @param cttMethod
//  @param cttExtraConf
//  @param paramsFormat
//  @param paramsData
//  @return *CommitContract
func NewCommitContract(identity string, cttName, cttVersion, cttMethod string,
	cttExtraConf map[string]string, paramsFormat string, paramsData []string) *CommitContract {
	innerContact := NewInnerContract(cttName, cttVersion, cttMethod)
	innerContact.Puts(cttExtraConf)
	return &CommitContract{
		Identity:     identity,
		Ctt:          innerContact,
		ParamsFormat: paramsFormat,
		ParamsData:   paramsData,
	}
}

// GetParamsData 获取参数
//  @receiver c
//  @return string
func (c *CommitContract) GetParamsData() string {
	if len(c.ParamsData) == 0 {
		return ""
	}
	bytes, err := json.Marshal(c.ParamsData)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// Check 检查内部参数，该类允许内部为空
//  @receiver c
//  @return error
func (c *CommitContract) Check() error {
	if c.Identity == "" {
		return errorIdentityEmpty
	}
	if !c.Ctt.isValid() {
		return errorContractConfig
	}
	return nil
}

// CheckIdentity 检查链身份
//  @receiver c
//  @param ChainRID
//  @return error
func (c *CommitContract) CheckIdentity(chainRID string) error {
	if c.Identity == "" {
		return errorIdentityEmpty
	}
	//chain, err := resource.GetResourceManager().GetChain(ChainRID)
	//if err != nil {
	//	return fmt.Errorf(util.CommitIdentityNotExist, c.Identity, ChainRID)
	//}
	//if chain.SignIdentity(c.Identity) == nil {
	//	return fmt.Errorf(util.CommitIdentityNotExist, c.Identity, ChainRID)
	//}
	return nil
}

// NewEmptyCrossSubEvent 新建空的子事件
//  @return *CrossSubEvent
func NewEmptyCrossSubEvent() *CrossSubEvent {
	return &CrossSubEvent{
		ChainCtt: NewEmptyCommitContract(),
	}
}

// NewCrossSubEvent 新建子事件
//  @param chainRID
//  @param dstGatewayID
//  @param chainCtt
//  @return *CrossSubEvent
func NewCrossSubEvent(chainRID string, dstGatewayID string, chainCtt *CommitContract) *CrossSubEvent {
	return &CrossSubEvent{
		ChainRID:  chainRID,
		GatewayID: dstGatewayID,
		ChainCtt:  chainCtt,
	}
}

// CheckAndInit 检查并初始化内部参数
//  @receiver se
//  @return error
func (se *CrossSubEvent) CheckAndInit() error {
	if se.ChainRID == "" {
		return errorChainEmpty
	}
	if err := se.ChainCtt.Check(); err != nil {
		return err
	}
	if se.ChainCtt.ParamsFormat == "" {
		return errorParamsFormatConfig
	}
	// 都正确的情况下，返回nil
	return nil
}

// NewEmptyInnerContract 新建新的内部合约对象
//  @return *InnerContract
func NewEmptyInnerContract() *InnerContract {
	return &InnerContract{
		ExtraData: make(map[string]string),
	}
}

// NewInnerContract 新建内部合约
//  @param name
//  @param version
//  @param method
//  @return *InnerContract
func NewInnerContract(name, version, method string) *InnerContract {
	return &InnerContract{
		Name:      name,
		Version:   version,
		Method:    method,
		ExtraData: make(map[string]string),
	}
}

// Put 设置key值
//  @receiver c
//  @param key
//  @param val
func (c *InnerContract) Put(key, val string) {
	c.ExtraData[key] = val
}

// Puts 批量设置key值
//  @receiver c
//  @param kvs
func (c *InnerContract) Puts(kvs map[string]string) {
	if len(kvs) > 0 {
		for k, v := range kvs {
			c.ExtraData[k] = v
		}
	}
}

// PutBatch 批量设置key值
//  @receiver c
//  @param kvs
func (c *InnerContract) PutBatch(kvs map[string]interface{}) {
	for k, v := range kvs {
		c.ExtraData[k] = fmt.Sprintf("%v", v)
	}
}

// Abi 获取abi
//  @receiver c
//  @return string
//  @return bool
func (c *InnerContract) Abi() (string, bool) {
	val, ok := c.ExtraData[aBI]
	return val, ok
}

// isValid 是否合法
//  @receiver c
//  @return bool
func (c *InnerContract) isValid() bool {
	// 要求合约名称和方法不可以为空，允许版本号为空
	if c.Name != "" && c.Method != "" {
		return true
	}
	return false
}

// GetExtraData 获取扩展字段
//  @receiver c
//  @return string
func (c *InnerContract) GetExtraData() string {
	if len(c.ExtraData) == 0 {
		return ""
	}
	bytes, err := json.Marshal(c.ExtraData)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// KeyValues kv对象
type KeyValues struct {
	kvs map[string]*KeyValue
}

// NewKeyValues 新建kv对象
//  @return *KeyValues
func NewKeyValues() *KeyValues {
	return &KeyValues{
		kvs: make(map[string]*KeyValue, 8),
	}
}

// NewNewKeyValuesByString 通过字符串创建kev对象
//  @param data
//  @return *KeyValues
func NewNewKeyValuesByString(data string) *KeyValues {
	kvs := NewKeyValues()
	kvStrs := strings.Split(data, segmentSplit)
	for i := 0; i < len(kvStrs); i++ {
		kvStr := kvStrs[i]
		if len(kvStr) == 0 {
			continue
		}
		kv := NewKeyValueByString(kvStr)
		if kv != nil {
			kvs.Append(kv)
		}
	}
	return kvs
}

// Append 追加kev对象
//  @receiver kvs
//  @param kv
func (kvs *KeyValues) Append(kv *KeyValue) {
	kvs.kvs[kv.Key()] = kv
}

// Length 获取kv的长度
//  @receiver kvs
//  @return int
func (kvs *KeyValues) Length() int {
	return len(kvs.kvs)
}

// hasKey 检查是否存在某个key
//  @receiver kvs
//  @param key
//  @return bool
func (kvs *KeyValues) hasKey(key string) bool {
	_, ok := kvs.kvs[key]
	return ok
}

// getValue 获取key对应的value
//  @receiver kvs
//  @param key
//  @return string
//  @return bool
func (kvs *KeyValues) getValue(key string) (string, bool) {
	if kv, ok := kvs.kvs[key]; ok {
		return kv.Value(), true
	}
	return "", false
}

// KeyValue kv对象
type KeyValue struct {
	key   string
	value string
}

// NewKeyValueByString 通过字符串新建kv对象
//  @param data
//  @return *KeyValue
func NewKeyValueByString(data string) *KeyValue {
	kv := strings.Split(data, keyValueSplit)
	if len(kv) == keyValueLength {
		return NewKeyValue(kv[zeroIndex], kv[firstIndex])
	}
	return nil
}

// NewKeyValue 新建kv对象
//  @param key
//  @param value
//  @return *KeyValue
func NewKeyValue(key, value string) *KeyValue {
	return &KeyValue{
		key:   key,
		value: value,
	}
}

// Key 获取key值
//  @receiver kv
//  @return string
func (kv *KeyValue) Key() string {
	return kv.key
}

// Value 获取value值
//  @receiver kv
//  @return string
func (kv *KeyValue) Value() string {
	return kv.value
}

// isEmpty 检查是否为空
//  @param data
//  @return bool
func isEmpty(data string) bool {
	if len(data) == 0 {
		return true
	}
	// 判断其内容是否全部为空
	dataBytes := []byte(data)
	for _, bz := range dataBytes {
		if bz != 0x0 {
			return false
		}
	}
	return true
}

// selfProtocolToMap 自定义协议字符串转为Map
// 自定义协议格式如下
// 格式：base64(${Key})#base64(${Value})@base64(${Key})#base64(${Value})......
// 说明：key和value都需要使用base64进行转换，并且多个k-v间使用@符号进行分割
//  @param data
//  @return map[string]string
func selfProtocolToMap(data string) map[string]string {
	if len(data) > 0 {
		dataMap := make(map[string]string)
		datas := strings.Split(data, paramSplit)
		if len(datas) > 0 {
			for i, str := range datas {
				// 最后一项需要特殊处理
				if i == len(datas)-1 {
					// 首先判断最后一项是否全部为0
					if isEmpty(str) {
						break
					}
				}
				strs := strings.Split(str, keyValSplit)
				if len(strs) == 2 {
					keyBytes, err := base64.URLEncoding.DecodeString(strs[0])
					if err != nil || len(keyBytes) == 0 {
						// 解码报错，则不处理该key-value
						continue
					}
					valBytes, err := base64.URLEncoding.DecodeString(strs[1])
					if err != nil || len(valBytes) == 0 {
						// 解码报错，则不处理该key-value
						continue
					}
					dataMap[string(keyBytes)] = string(valBytes)
				}
			}
		}
		return dataMap
	}
	return nil
}

// convertToParamArray 将原始数据转换为数组，可选是否处理base64
//  @param paramsData
//  @param saveToBase64
//  @return []string
//  @return error
func convertToParamArray(paramsData string, saveToBase64 bool) ([]string, error) {
	newParams := strings.Trim(paramsData, paramSplit)
	params := strings.Split(newParams, paramSplit)
	datas := make([]string, 0)
	// 遍历
	for _, param := range params {
		if isEmpty(param) {
			continue
		}
		var (
			data []byte
			err  error
		)
		if !saveToBase64 {
			data, err = base64.URLEncoding.DecodeString(param)
			if err != nil {
				return nil, err
			}
		} else {
			data = []byte(param)
		}
		datas = append(datas, string(data))
	}
	return datas, nil
}

// ChainTxContext 某条链交易的处理上下文，含链资源ID/证明情况及合约结果
type ChainTxContext struct {
	ChainRID       string
	ProveStatus    int
	ContractStatus int
	ContractResult string
}

// NewChainTxContext 新建某条链交易的处理上下文
//  @param chainRID
//  @param proveStatus
//  @param contractStatus
//  @param contractResult
//  @return *ChainTxContext
func NewChainTxContext(chainRID string, proveStatus, contractStatus int, contractResult string) *ChainTxContext {
	return &ChainTxContext{
		ChainRID:       chainRID,
		ProveStatus:    proveStatus,
		ContractStatus: contractStatus,
		ContractResult: contractResult,
	}
}

// ToContractParam 序列化为tbis协议参数
//  @receiver c
//  @return string
func (c *ChainTxContext) ToContractParam() string {
	// 防止合约部支持json处理，采用字符串拼接的方式
	base64Result := base64.URLEncoding.EncodeToString([]byte(c.ContractResult))
	// ${ChainRID},${ProveStatus},${ContractStatus},${base64Result};
	// ${ChainRID},${ProveStatus},${ContractStatus},${base64Result}
	return c.ChainRID + "," + fmt.Sprintf("%v", c.ProveStatus) + "," +
		fmt.Sprintf("%v", c.ContractStatus) + "," + base64Result
}

// ChainTxContexts 交易的处理上下文
type ChainTxContexts struct {
	Contexts []*ChainTxContext
}

// NewChainTxContexts 新建交易的处理上下文
//  @return *ChainTxContexts
func NewChainTxContexts() *ChainTxContexts {
	return &ChainTxContexts{
		Contexts: make([]*ChainTxContext, 0),
	}
}

// ToContractParam 序列化为tbis协议参数
//  @receiver c
//  @return string
func (c *ChainTxContexts) ToContractParam() string {
	if len(c.Contexts) == 0 {
		return ""
	}
	var buffer bytes.Buffer //Buffer是一个实现了读写方法的可变大小的字节缓冲
	for _, txContext := range c.Contexts {
		txParam := txContext.ToContractParam()
		if buffer.Len() > 0 {
			buffer.WriteString(";")
		}
		buffer.WriteString(txParam)
	}
	return buffer.String()
}

// Append 追加
//  @receiver c
//  @param chainTxCtx
func (c *ChainTxContexts) Append(chainTxCtx *ChainTxContext) {
	c.Contexts = append(c.Contexts, chainTxCtx)
}

// Len 获取长度
//  @receiver c
//  @return int
func (c *ChainTxContexts) Len() int {
	return len(c.Contexts)
}
