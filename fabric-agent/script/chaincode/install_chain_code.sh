#!/bin/bash
set -e

package_file=
chain_path=
peer_msp_id=
peer_endpoint=
peer_name=
peer_org_name=

help()
{
    cat << EOF
Usage:
    -p <package path>            [Required]
    -m <peer msp id>             [Required]
    -f <chain path>              [Required]
    -e <peer endpoint>           [Required]
    -s <peer name>               [Required]
    -g <peer org name>           [Required]
    -h Help
e.g:
    bash $0 -p ~/fabric/config/basic.tar.gz -m Org1MSP -e 192.168.56.102:7051 -s peer0.org1.example.com -g org1.example.com -f ~/fabric
EOF
exit 0
}

parse_params()
{
    while getopts "p:f:m:n:v:e:s:g:h" option;do
        case $option in
        p) package_file="${OPTARG}"
            if [ ! -f "$package_file" ]; then LOG_WARN "$package_file not exist" && exit 1; fi
        ;;
        f) chain_path="${OPTARG}"
            if [ ! -d "$chain_path" ]; then LOG_WARN "$chain_path not exist" && exit 1; fi
        ;;
        m) peer_msp_id="${OPTARG}"
            if [ -z "$peer_msp_id" ]; then LOG_WARN "$peer_msp_id not specified" && exit 1; fi
        ;;
        e) peer_endpoint="${OPTARG}"
            if [ -z "$peer_endpoint" ]; then LOG_WARN "$peer_endpoint not specified" && exit 1; fi
        ;;
        s) peer_name="${OPTARG}"
            if [ -z "$peer_name" ]; then LOG_WARN "$peer_name not specified" && exit 1; fi
        ;;
        g) peer_org_name="${OPTARG}"
            if [ -z "$peer_org_name" ]; then LOG_WARN "$peer_org_name not specified" && exit 1; fi
        ;;
        h) help;;
        *) LOG_WARN "invalid option $option";;
        esac
    done
}

main()
{
    export PATH=${chain_path}/bin:$PATH
    export FABRIC_CA_CLIENT_HOME=${chain_path}/fabric-ca-client/${peer_org_name}
    export FABRIC_CFG_PATH=${chain_path}/config
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="${peer_msp_id}"
    export CORE_PEER_TLS_ROOTCERT_FILE="${chain_path}/organizations/peerOrganizations/${peer_org_name}/peers/${peer_name}/tls/ca.crt"
    export CORE_PEER_MSPCONFIGPATH="${chain_path}/organizations/peerOrganizations/${peer_org_name}/users/Admin@${peer_org_name}/msp"
    export CORE_PEER_ADDRESS="${peer_endpoint}"
    export FABRIC_LOGGING_SPEC=INFO
    # export CORE_PEER_GOSSIP_EXTERNALENDPOINT="${peer_endpoint}"

    log=$(peer lifecycle chaincode install "${package_file}" 2>&1)
    echo "$log"
}

parse_params "$@"
main


