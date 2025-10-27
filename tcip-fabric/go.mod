module chainmaker.org/chainmaker/tcip-fabric/v2

go 1.16

require (
	chainmaker.org/chainmaker/common/v2 v2.3.1
	chainmaker.org/chainmaker/pb-go/v2 v2.3.2
	chainmaker.org/chainmaker/tcip-go/v2 v2.3.1
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible
	github.com/cloudflare/cfssl v1.6.1
	github.com/common-nighthawk/go-figure v0.0.0-20210622060536-734e95fb86be
	github.com/emirpasic/gods v1.18.1
	github.com/gogo/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/hyperledger/fabric-protos-go v0.0.0-20200707132912-fee30f3ccd23
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.11.0
	github.com/stretchr/testify v1.7.1
	github.com/syndtr/goleveldb v1.0.1-0.20210305035536-64b5b1c73954
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.0.0-20221014081412-f15817d10f9b
	//golang.org/x/net v0.0.0-20220517181318-183a9ca12b87 // indirect
	google.golang.org/grpc v1.50.1
	google.golang.org/protobuf v1.28.1
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/go-kit/kit => github.com/go-kit/kit v0.8.0
	google.golang.org/grpc v1.40.0 => google.golang.org/grpc v1.26.0
)
