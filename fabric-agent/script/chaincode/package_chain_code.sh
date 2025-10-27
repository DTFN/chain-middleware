#!/bin/bash
set -e

logfile=${PWD}/deploy.log
package_file=
chain_code_path=
lang=
chain_code_name=
chain_code_version=
chain_path=
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
    -p <package file>            [Required]
    -c <chain code path>         [Required]
    -l <lang>                    [Required]
    -n <chain code name>         [Required]
    -v <chain code version>      [Required]
    -m <peer msp id>             [Required]
    -f <chain path>              [Required]
    -e <peer endpoint>           [Required]
    -s <peer name>               [Required]
    -g <peer org name>           [Required]
    -h Help
e.g:
    bash $0 -p ~/fabric/config -l java -c ~/fabric/chain-code/basic -n basic -v 1.0 -m Org1MSP -e 192.168.56.102:7051 -s peer0.org1.example.com -g org1.example.com -f ~/fabric
EOF
exit 0
}

parse_params()
{
    while getopts "c:p:l:n:v:f:m:e:s:g:h" option;do
        case $option in
        c) chain_code_path="${OPTARG}"
            if [ ! -d "$chain_code_path" ]; then LOG_WARN "$chain_code_path not exist" && exit 1; fi
        ;;
        p) package_file="${OPTARG}"
            if [ -z "$package_file" ]; then LOG_WARN "$package_file not exist" && exit 1; fi
        ;;
        l) lang="${OPTARG}"
            if [ -z "$lang" ]; then LOG_WARN "$lang not specified" && exit 1; fi
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

    peer lifecycle chaincode package "${package_file}" --path "${chain_code_path}" --lang "${lang}" --label "${chain_code_name}_${chain_code_version}"
    #rm "${logfile}"
}

print_result()
{
    echo "=============================================================="
    LOG_INFO "chain_code_version : ${chain_code_name}_${chain_code_version}"
    LOG_INFO "All completed."
}

echo "=======$(date)========" >>"${logfile}"
echo "=======start package chain code===============" >>"${logfile}"
parse_params "$@"
main
print_result
echo "=======$(date)========" >>"${logfile}"
echo "=======end package chain code=================" >>"${logfile}"


