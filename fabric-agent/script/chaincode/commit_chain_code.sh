#!/bin/bash
set -e

logfile=${PWD}/deploy.log
package_id=
channel_id=
orderer_name=
orderer_endpoint=
orderer_org_name=
chain_code_name=
chain_code_version=
chain_path=
peer_msp_id=
peer_endpoint=
peer_name=
peer_org_name=
sequence=

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
    -p <package id>              [Required]
    -c <channel id>              [Required]
    -a <order name>              [Required]
    -o <order endpoint>          [Required]
    -r <order org name>          [Required]
    -n <chain code name>         [Required]
    -v <chain code version>      [Required]
    -m <peer msp id>             [Required]
    -f <chain path>              [Required]
    -e <peer endpoint>           [Required]
    -s <peer name>               [Required]
    -g <peer org name>           [Required]
    -q <sequence>               [Required]
    -h Help
e.g:
    bash $0 -p basic_1.0:9a3b699872735315aa3032cd2a9e863b69434f7197f1324b55bf573e9a6884e6 -c channel1 -m Org1MSP -e 192.168.56.102:7051 -s peer0.org1.example.com -g org1.example.com -f ~/fabric -a orderer0.org1.example.com -r org1.example.com -o 192.168.56.102:7050 -n basic -v 1.0 -q 1
EOF
exit 0
}

parse_params()
{
    while getopts "c:p:a:r:o:n:v:f:m:e:s:g:q:h" option;do
        case $option in
        c) channel_id="${OPTARG}"
            if [ -z "$channel_id" ]; then LOG_WARN "$channel_id not specified" && exit 1; fi
        ;;
        p) package_id="${OPTARG}"
            if [ -z "$package_id" ]; then LOG_WARN "$package_id not specified" && exit 1; fi
        ;;
        a) orderer_name="${OPTARG}"
            if [ -z "$orderer_name" ]; then LOG_WARN "$orderer_name not exist" && exit 1; fi
        ;;
        r) orderer_org_name="${OPTARG}"
            if [ -z "$orderer_org_name" ]; then LOG_WARN "$orderer_org_name not specified" && exit 1; fi
        ;;
        o) orderer_endpoint="${OPTARG}"
            if [ -z "$orderer_endpoint" ]; then LOG_WARN "$orderer_endpoint not specified" && exit 1; fi
        ;;
        n) chain_code_name="${OPTARG}"
            if [ -z "$chain_code_name" ]; then LOG_WARN "$chain_code_name not specified" && exit 1; fi
        ;;
        v) chain_code_version="${OPTARG}"
            if [ -z "$chain_code_version" ]; then LOG_WARN "$chain_code_version not specified" && exit 1; fi
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
        q) sequence="${OPTARG}"
            if [ -z "$sequence" ]; then LOG_WARN "$sequence not specified" && exit 1; fi
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
    # export CORE_PEER_GOSSIP_EXTERNALENDPOINT="${peer_endpoint}"

    peer lifecycle chaincode commit -o "${orderer_endpoint}" --channelID "${channel_id}" --name "${chain_code_name}" --peerAddresses "${peer_endpoint}" --tlsRootCertFiles "${chain_path}/organizations/peerOrganizations/${peer_org_name}/peers/${peer_name}/tls/ca.crt" --version "${chain_code_version}" --sequence "${sequence}" --tls --cafile "${chain_path}/organizations/ordererOrganizations/${orderer_org_name}/orderers/${orderer_name}/msp/tlscacerts/tlsca.${orderer_org_name}-cert.pem"

}

print_result()
{
    echo "=============================================================="
    LOG_INFO "package id  : ${package_id}"
    LOG_INFO "All completed."
}

echo "=======$(date)========" >>"${logfile}"
echo "=======start approve and commit chain code===============" >>"${logfile}"
parse_params "$@"
main
print_result
echo "=======$(date)========" >>"${logfile}"
echo "=======end approve and commit chain code=================" >>"${logfile}"


