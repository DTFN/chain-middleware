#!/bin/bash
set -e

logfile=${PWD}/build.log
channel_id=
chain_path=
orderer_endpoint=
orderer_name=
orderer_org_name=
peer_msp_id=
peer_endpoint=
peer_name=
peer_org_name=

LOG_WARN()
{
    local content=${1}
    echo -e "\033[31m[WARN] ${content}\033[0m"
}

LOG_INFO()
{
    local content=${1}
    echo -e "\033[32m[INFO] ${content}\033[0m"
}

help()
{
    cat << EOF
Usage:
    -c <channel id>              [Required]
    -p <chain path>              [Required]
    -o <orderer endpoint>        [Required]
    -j <orderer name>            [Required]
    -r <orderer org name>        [Required]
    -m <peer msp id>             [Required]
    -n <peer endpoint>           [Required]
    -s <peer name>               [Required]
    -g <peer org name>           [Required]
    -h Help
e.g:
    bash $0 -c channel1 -p ./NODES_ROOT/1ef243fd87 -o 192.168.3.128:7050 -j orderer0.org1.example.com -r org1.example.com -m Org1MSP -n 192.168.3.128:7051 -s peer0.org1.example.com -g org1.example.com
EOF
exit 0
}

parse_params()
{
    while getopts "c:p:o:j:r:m:n:s:g:h" option;do
        case $option in
        c) channel_id="${OPTARG}"
            if [ -z "$channel_id" ]; then LOG_WARN "$channel_id not specified" && exit 1; fi
        ;;
        p) chain_path="${OPTARG}"
            if [ ! -d "$chain_path" ]; then LOG_WARN "$chain_path not exist" && exit 1; fi
        ;;
        o) orderer_endpoint="${OPTARG}"
            if [ -z "$orderer_endpoint" ]; then LOG_WARN "$orderer_endpoint not specified" && exit 1; fi
        ;;
        j) orderer_name="${OPTARG}"
            if [ -z "$orderer_name" ]; then LOG_WARN "$orderer_name not specified" && exit 1; fi
        ;;
        r) orderer_org_name="${OPTARG}"
            if [ -z "$orderer_org_name" ]; then LOG_WARN "$orderer_org_name not specified" && exit 1; fi
        ;;
        m) peer_msp_id="${OPTARG}"
            if [ -z "$peer_msp_id" ]; then LOG_WARN "$peer_msp_id not specified" && exit 1; fi
        ;;
        n) peer_endpoint="${OPTARG}"
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

init_env_param() {
  export PATH="${chain_path}/bin:$PATH"
  export FABRIC_CA_CLIENT_HOME="${chain_path}/fabric-ca-client/${peer_org_name}"
  export FABRIC_CFG_PATH="${chain_path}/config"

  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="${peer_msp_id}"
  export CORE_PEER_TLS_ROOTCERT_FILE="${chain_path}/organizations/peerOrganizations/${peer_org_name}/peers/${peer_name}/tls/ca.crt"
  export CORE_PEER_MSPCONFIGPATH="${chain_path}/organizations/peerOrganizations/${peer_org_name}/users/Admin@${peer_org_name}/msp"
  export CORE_PEER_ADDRESS="${peer_endpoint}"
  # export CORE_PEER_GOSSIP_EXTERNALENDPOINT="${peer_endpoint}"
}

main()
{
    init_env_param
    peer channel update -o "${orderer_endpoint}" -c "${channel_id}" -f "${chain_path}/config/channel-artifacts/${peer_msp_id}anchors.tx" --tls --cafile "${chain_path}/organizations/ordererOrganizations/${orderer_org_name}/orderers/${orderer_name}/msp/tlscacerts/tlsca.${orderer_org_name}-cert.pem"
    #rm "${logfile}"
}

print_result()
{
    echo "=============================================================="
    LOG_INFO "Update Channel ${channel_id} config : ${chain_path}/config/channel-artifacts/${peer_msp_id}anchors.tx"
    LOG_INFO "All completed."
}

echo "=======$(date)========" >>"${logfile}"
echo "=======start update app channel===============" >>"${logfile}"
parse_params "$@"
main
print_result
echo "=======$(date)========" >>"${logfile}"
echo "=======end update app channel=================" >>"${logfile}"


