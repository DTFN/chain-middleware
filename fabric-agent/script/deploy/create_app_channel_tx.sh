#!/bin/bash
set -e

logfile=${PWD}/build.log
channel_id=
chain_path=

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
    -c <channel id>        [Required]
    -p <chain path>        [Required]
    -h Help
e.g:
    bash $0 -c channel1 -p ./NODES_ROOT/1ef243fd87
EOF
exit 0
}

parse_params()
{
    while getopts "c:p:h" option;do
        case $option in
        c) channel_id="${OPTARG}"
            if [ -z "$channel_id" ]; then LOG_WARN "$channel_id not specified" && exit 1; fi
        ;;
        p) chain_path="${OPTARG}"
            if [ ! -d "$chain_path" ]; then LOG_WARN "$chain_path not exist" && exit 1; fi
        ;;
        h) help;;
        *) LOG_WARN "invalid option $option";;
        esac
    done
}

init_env_param() {
  export PATH="${chain_path}/bin:$PATH"
  export FABRIC_CFG_PATH="${chain_path}/config"
}

main()
{
    init_env_param
    configtxgen -profile Channel1 -outputCreateChannelTx "${chain_path}/config/channel-artifacts/${channel_id}.tx" -channelID "${channel_id}" > "${logfile}" 2>&1
    #rm "${logfile}"
}

print_result()
{
    echo "=============================================================="
    LOG_INFO "Channel Tx Path   : ${chain_path}/config/channel-artifacts/${channel_id}.tx"
    LOG_INFO "All completed."
}

echo "=======$(date)========" >>"${logfile}"
echo "=======start create app channel tx===============" >>"${logfile}"
parse_params "$@"
main
print_result
echo "=======$(date)========" >>"${logfile}"
echo "=======end create app channel tx=================" >>"${logfile}"

