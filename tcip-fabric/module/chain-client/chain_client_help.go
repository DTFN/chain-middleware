package chain_client

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	"go.uber.org/zap"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/utils"
)

const (
	clientLogLevel = "info"
	keySuffix      = ".key"
	crtSuffix      = ".crt"
	userName       = "user"
)

// FabricSDKTemplate 长安链SDK模板类
type FabricSDKTemplate struct {
	Version        string                   `yaml:"version"`
	Client         *Client                  `yaml:"client"`
	Channels       map[string]*Channel      `yaml:"channels"`
	Organizations  map[string]*Organization `yaml:"organizations"`
	Peers          map[string]*Peer         `yaml:"peers"`
	Orderers       map[string]*Orderer      `yaml:"orderers"`
	EntityMatchers *EntityMatchers          `yaml:"entityMatchers"`
}

// newFabricSDKTemplate 创建 Fabric-SDK 模板
//  @return *FabricSDKTemplate
func newFabricSDKTemplate() *FabricSDKTemplate {
	template := &FabricSDKTemplate{
		Version:        "1.0.0",
		Client:         &Client{},
		Channels:       map[string]*Channel{},
		Organizations:  map[string]*Organization{},
		Peers:          map[string]*Peer{},
		Orderers:       map[string]*Orderer{},
		EntityMatchers: &EntityMatchers{},
	}
	return template
}

// newClient 创建新的client
//  @param organization
//  @param key
//  @param crt
//  @return *Client
func newClient(organization string , key []byte, crt []byte) *Client {
	keyFile, err := saveBytesToFile(key, keySuffix)
	if err != nil {
		// TODO logger
		return nil
	}
	crtFile, err := saveBytesToFile(crt, crtSuffix)
	if err != nil {
		// TODO logger
		return nil
	}
	return &Client{
		Organization: organization,
		Logging: &Logging{
			clientLogLevel,
		},
		TlsCerts: &TlsCerts{
			SystemCertPool: false,
			ClientCert: &ClientCert{
				Key: &FilePath{
					Path: keyFile,
				},
				Cert: &FilePath{
					Path: crtFile,
				},
			},
		},
	}
}

// newChannel 创建新的channel
//  @param orderers
//  @param peers
//  @return *Channel
func newChannel(orderers []string, peers []string) *Channel {
	peerOpts := make(map[string]*PeerOpt)
	for _, pname := range peers {
		opt := &PeerOpt{
			EndorsingPeer:  true,
			ChaincodeQuery: true,
			LedgerQuery:    true,
			EventSource:    true,
		}
		peerOpts[pname] = opt
	}
	return &Channel{
		Orderers: orderers,
		Peers:    peerOpts,
	}
}

// newUser 创建新的user
//  @param name
//  @param keyBytes
//  @param crtBytes
//  @return *User
func newUser(name string, keyBytes, crtBytes []byte) *User {
	keyFile, err := saveBytesToFile(keyBytes, keySuffix)
	if err != nil {
		// TODO logger
		return nil
	}
	crtFile, err := saveBytesToFile(crtBytes, crtSuffix)
	if err != nil {
		// TODO logger
		return nil
	}
	return &User{
		name: name,
		KeyFile: &FilePath{
			Path: keyFile,
		},
		CertFile: &FilePath{
			Path: crtFile,
		},
	}
}

// newOrganization 创建新的组织
//  @param mspID
//  @param peers
//  @param users
//  @return *Organization
func newOrganization(mspID string, peers []string, users map[string]*User) *Organization {
	return &Organization{
		MspID: mspID,
		Peers: peers,
		Users: users,
	}
}

// newPeer 创建新的peer节点
//  @param url
//  @param address
//  @param tlsCrt
//  @return *Peer
func newPeer(url string, address string, tlsCrt []byte) *Peer {
	keyFile, err := saveBytesToFile(tlsCrt, crtSuffix)
	if err != nil {
		// TODO logger
		return nil
	}
	return &Peer{
		URL: url,
		GRPCOptions: &GRPCOptions{
			Override:      address,
			AliveTime:     "10s",
			AliveTimeOut:  "20s",
			AlivePermit:   false,
			FailFast:      false,
			AllowInsecure: false,
		},
		TlsCACerts: &FilePath{
			Path: keyFile,
		},
	}
}

// newOrderer 创建新的orderer节点
//  @param url
//  @param address
//  @param tlsCrt
//  @return *Orderer
func newOrderer(url string, address string, tlsCrt []byte) *Orderer {
	keyFile, err := saveBytesToFile(tlsCrt, crtSuffix)
	if err != nil {
		// TODO logger
		return nil
	}
	return &Orderer{
		URL: url,
		GRPCOptions: &GRPCOptions{
			Override:      address,
			AliveTime:     "10s",
			AliveTimeOut:  "20s",
			AlivePermit:   false,
			FailFast:      false,
			AllowInsecure: false,
		},
		TlsCACerts: &FilePath{
			Path: keyFile,
		},
	}
}

// newEntityMatcher 创建新的EntityMatcher
//  @param peerMappers
//  @param ordererMappers
//  @return *EntityMatchers
func newEntityMatcher(peerMappers, ordererMappers []string) *EntityMatchers {
	pattern := "(\\w*)%s(\\w*)"

	peerMatchers := make([]*Matcher, 0)
	ordererMatchers := make([]*Matcher, 0)

	for _, peerMapper := range peerMappers {
		matcher := &Matcher{
			Pattern:    fmt.Sprintf(pattern, peerMapper),
			MappedHost: peerMapper,
		}
		peerMatchers = append(peerMatchers, matcher)
	}

	for _, ordererMapper := range ordererMappers {
		matcher := &Matcher{
			Pattern:    fmt.Sprintf(pattern, ordererMapper),
			MappedHost: ordererMapper,
		}
		ordererMatchers = append(ordererMatchers, matcher)
	}

	return &EntityMatchers{
		PeerMatcher:    peerMatchers,
		OrdererMatcher: ordererMatchers,
	}
}

// SetClient 设置client
//  @receiver t
//  @param client
func (t *FabricSDKTemplate) SetClient(client *Client) {
	t.Client = client
}

// SetChannel 设置channel
//  @receiver t
//  @param name
//  @param channel
func (t *FabricSDKTemplate) SetChannel(name string, channel *Channel) {
	t.Channels[name] = channel
}

// SetOrganization 设置组织
//  @receiver t
//  @param name
//  @param organization
func (t *FabricSDKTemplate) SetOrganization(name string, organization *Organization) {
	t.Organizations[name] = organization
}

// SetPeers 设置peer节点
//  @receiver t
//  @param peers
func (t *FabricSDKTemplate) SetPeers(peers map[string]*Peer) {
	if len(peers) > 0 {
		for name, peer := range peers {
			t.SetPeer(name, peer)
		}
	}
}

// SetPeer 设置peer节点
//  @receiver t
//  @param name
//  @param peer
func (t *FabricSDKTemplate) SetPeer(name string, peer *Peer) {
	t.Peers[name] = peer
}

// SetOrderers 设置orderer节点
//  @receiver t
//  @param orderers
func (t *FabricSDKTemplate) SetOrderers(orderers map[string]*Orderer) {
	if len(orderers) > 0 {
		for name, orderer := range orderers {
			t.SetOrderer(name, orderer)
		}
	}
}

// SetOrderer 设置orderer节点
//  @receiver t
//  @param name
//  @param orderer
func (t *FabricSDKTemplate) SetOrderer(name string, orderer *Orderer) {
	t.Orderers[name] = orderer
}

// SetEntityMatchers 设置EntityMatchers
//  @receiver t
//  @param entityMatchers
func (t *FabricSDKTemplate) SetEntityMatchers(entityMatchers *EntityMatchers) {
	t.EntityMatchers = entityMatchers
}

// fabricSDKTemplateToYml 生成yml文件
//  @param t
//  @return string
//  @return error
func fabricSDKTemplateToYml(t *FabricSDKTemplate) (string, error) {
	ymlFile, err := utils.WriteTempYmlFile(t)
	if err != nil {
		return "", err
	}
	return ymlFile.Name(), nil
}

// Client client结构
type Client struct {
	Organization string    `yaml:"organization"`
	Logging      *Logging  `yaml:"logging"`
	TlsCerts     *TlsCerts `yaml:"tlsCerts"`
}

// Logging 日志结构
type Logging struct {
	Level string `yaml:"level"`
}

// TlsCerts tls证书结构
type TlsCerts struct {
	SystemCertPool bool        `yaml:"systemCertPool"`
	ClientCert     *ClientCert `yaml:"client"`
}

// ClientCert client证书结构
type ClientCert struct {
	Key  *FilePath `yaml:"key"`
	Cert *FilePath `yaml:"cert"`
}

// Channel channel结构
type Channel struct {
	Orderers []string            `yaml:"orderers"`
	Peers    map[string]*PeerOpt `yaml:"peers"`
}

// PeerOpt peer选项
type PeerOpt struct {
	EndorsingPeer  bool `yaml:"endorsingPeer"`
	ChaincodeQuery bool `yaml:"chaincodeQuery"`
	LedgerQuery    bool `yaml:"ledgerQuery"`
	EventSource    bool `yaml:"eventSource"`
}

// Organization 组织结构
type Organization struct {
	MspID string           `yaml:"mspid"`
	Peers []string         `yaml:"peers"`
	Users map[string]*User `yaml:"users"`
}

// User user结构
type User struct {
	name     string
	KeyFile  *FilePath `yaml:"key"`
	CertFile *FilePath `yaml:"cert"`
}

// FilePath 文件路径
type FilePath struct {
	Path string `yaml:"path"`
}

// Peer peer节点
type Peer struct {
	URL         string       `yaml:"url"`
	GRPCOptions *GRPCOptions `yaml:"grpcOptions"`
	TlsCACerts  *FilePath    `yaml:"tlsCACerts"`
}

// Orderer 配置内容与Peer一致
type Orderer struct {
	URL         string       `yaml:"url"`
	GRPCOptions *GRPCOptions `yaml:"grpcOptions"`
	TlsCACerts  *FilePath    `yaml:"tlsCACerts"`
}

// GRPCOptions grpc配置
type GRPCOptions struct {
	Override      string `yaml:"ssl-target-name-override"`
	AliveTime     string `yaml:"keep-alive-time"`
	AliveTimeOut  string `yaml:"keep-alive-timeout"`
	AlivePermit   bool   `yaml:"keep-alive-permit"`
	FailFast      bool   `yaml:"fail-fast"`
	AllowInsecure bool   `yaml:"allow-insecure"`
}

// EntityMatchers EntityMatchers
type EntityMatchers struct {
	PeerMatcher    []*Matcher `yaml:"peer"`
	OrdererMatcher []*Matcher `yaml:"orderer"`
}

// Matcher Matcher
type Matcher struct {
	Pattern    string `yaml:"pattern"`
	MappedHost string `yaml:"mappedHost"`
}

// saveBytesToFile 将指定内容存到文件
//  @param data
//  @param suffix
//  @return string
//  @return error
func saveBytesToFile(data []byte, suffix string) (string, error) {
	// 首先将内容写入文件
	tempFile, err := utils.WriteTempFile(data, suffix)
	if err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}

// createSdk 创建fabricsdk client
//  @param fabricConfig
//  @param log
//  @return *channel.Client
//  @return *ledger.Client
//  @return *event.Client
//  @return []string
//  @return error
func createSdk(fabricConfig *common.FabricConfig,
	log *zap.SugaredLogger) (*channel.Client, *ledger.Client, *event.Client, []string, error) {
	peers, orderers := make(map[string]*Peer), make(map[string]*Orderer)
	sdkTemplate := newFabricSDKTemplate()
	client := newClient(fabricConfig.Org[0].OrgId , []byte(fabricConfig.TlsKey), []byte(fabricConfig.TlsCert))
	sdkTemplate.SetClient(client)
	peerAddresses := make([]string, 0)
	peerNames := make([]string, 0)
	userAndOrgs := make([]fabsdk.ContextOption, 0)

	for _, org := range fabricConfig.Org {
		peerOrgAddresses := make([]string, 0)
		for _, peerNode := range org.Peers {
			cas := make([][]byte, len(peerNode.TrustRoot))
			for i, trustRoot := range peerNode.TrustRoot {
				cas[i] = []byte(trustRoot)
			}
			peer := newPeer(peerNode.NodeAddr, peerNode.TlsHostName, []byte(peerNode.TrustRoot[0]))
			peers[peerNode.TlsHostName] = peer
			peerOrgAddresses = append(peerOrgAddresses, peerNode.TlsHostName)
			peerNames = append(peerNames, peerNode.TlsHostName)
		}
		peerAddresses = append(peerAddresses, peerOrgAddresses...)

		user := newUser(userName, []byte(org.SignKey), []byte(org.SignCert))
		orgConfig := newOrganization(org.MspId, peerAddresses, map[string]*User{user.name: user}) // TODO
		sdkTemplate.SetOrganization(org.OrgId, orgConfig)

		userAndOrgs = append(userAndOrgs, fabsdk.WithOrg(org.OrgId))
		userAndOrgs = append(userAndOrgs, fabsdk.WithUser(userName))
	}

	orderer := newOrderer(fabricConfig.Orderers[0].NodeAddr,
		fabricConfig.Orderers[0].TlsHostName,
		[]byte(fabricConfig.Orderers[0].TrustRoot[0]))
	orderers[fabricConfig.Orderers[0].TlsHostName] = orderer
	ordererAddresses := []string{fabricConfig.Orderers[0].TlsHostName}

	//sdkTemplate := newFabricSDKTemplate()
	//client := newClient(fabricConfig.Org[0].OrgId /*, []byte(fabricConfig.TlsKey), []byte(fabricConfig.TlsCert)*/)
	//sdkTemplate.SetClient(client)

	ch := newChannel(ordererAddresses, peerAddresses)
	sdkTemplate.SetChannel(fabricConfig.ChainId, ch)

	//user := newUser(userName, []byte(fabricConfig.SignKey), []byte(fabricConfig.SignCert))
	//org := newOrganization(fabricConfig.MspId, peerAddresses, map[string]*User{user.name: user}) // TODO
	//sdkTemplate.SetOrganization(fabricConfig.OrgId, org)

	sdkTemplate.SetPeers(peers)
	sdkTemplate.SetOrderers(orderers)

	entity := newEntityMatcher(peerAddresses, ordererAddresses)
	sdkTemplate.SetEntityMatchers(entity)

	yml, err := fabricSDKTemplateToYml(sdkTemplate)
	if err != nil {
		log.Errorf("[createSdk] fabricSDKTemplateToYml: fabricConfig %v, error %v", fabricConfig, err)
		return nil, nil, nil, nil, fmt.Errorf("[createSdk] fabricSDKTemplateToYml: fabricConfig %v, error %v",
			fabricConfig, err)
	}
	log.Infof("[createSdk] fabric config yml file path: %s\n", yml)

	sdk, err := fabsdk.New(config.FromFile(yml))
	if err != nil {
		log.Errorf("[createSdk] New: fabricConfig %v, error %v", fabricConfig, err)
		return nil, nil, nil, nil, fmt.Errorf("[createSdk] New: fabricConfig %v, error %v", fabricConfig, err)
	}

	//userAndOrgs := make([]fabsdk.ContextOption, 0)
	//userAndOrgs = append(userAndOrgs, fabsdk.WithOrg(fabricConfig.OrgId))
	//userAndOrgs = append(userAndOrgs, fabsdk.WithUser(userName))
	ccp := sdk.ChannelContext(fabricConfig.ChainId, userAndOrgs...)

	channelClient, err := channel.New(ccp)
	if err != nil {
		log.Errorf("[InitChainClient] create chain [%s] client error: %v", fabricConfig.ChainId, err)
		return nil, nil, nil, nil, fmt.Errorf("[InitChainClient] create chain [%s] client error: %v",
			fabricConfig.ChainId, err)
	}

	ledgerClinet, err := ledger.New(ccp)
	if err != nil {
		log.Errorf("[InitChainClient] create chain [%s] ledger client error: %v", fabricConfig.ChainId, err)
		return nil, nil, nil, nil, fmt.Errorf("[InitChainClient] create chain [%s] ledger client error: %v",
			fabricConfig.ChainId, err)
	}

	eventClient, err := event.New(ccp, event.WithBlockEvents())
	if err != nil {
		log.Errorf("[InitChainClient] create chain [%s] event client error: %v", fabricConfig.ChainId, err)
		return nil, nil, nil, nil, fmt.Errorf("[InitChainClient] create chain [%s] event client error: %v",
			fabricConfig.ChainId, err)
	}
	return channelClient, ledgerClinet, eventClient, peerNames, nil
}
