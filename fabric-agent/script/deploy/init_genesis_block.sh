#!/bin/bash
set -e

logfile=${PWD}/build.log
chain_path=
org_domain=""

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
    -p <chain path>              [Required]
    -o <org domain>              [Required]
    -h Help
e.g:
    bash $0 -p ./fabric/test_chain -o
EOF
exit 0
}

parse_params()
{
    while getopts "o:p:h" option;do
        case $option in
        p) chain_path="${OPTARG}"
            if [ ! -d "$chain_path" ]; then LOG_WARN "$chain_path not exist" && exit 1; fi
        ;;
        o) org_domain="${OPTARG}"
            if [ -z "$org_domain" ]; then LOG_WARN "$org_domain is empty" && exit 1; fi
        ;;
        h) help;;
        *) LOG_WARN "invalid option $option";;
        esac
    done
}


init_env_param() {
  export PATH="${chain_path}/bin:$PATH"
  export FABRIC_CA_CLIENT_HOME="${chain_path}/fabric-ca-client/${org_domain}"
  export FABRIC_CFG_PATH="${chain_path}/config"
}

main()
{
  init_env_param
  configtxgen -profile OneOrgOrdererGenesis -channelID system-channel -outputBlock "${chain_path}/config/system-genesis-block/genesis.block"
}

print_result()
{
    echo "=============================================================="
    LOG_INFO "Init system genesis block conf in : ${chain_path}"
    LOG_INFO "All completed."
}

echo "=======$(date)========" >>"${logfile}"
echo "=======start genesis block env===============" >>"${logfile}"
parse_params "$@"
main
print_result
echo "=======$(date)========" >>"${logfile}"
echo "=======end genesis block env=================" >>"${logfile}"


