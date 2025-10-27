#!/bin/bash
set -e

logfile=${PWD}/build.log
channel_id=
chain_path=
peer_msp_id=
anchor_peers=

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
    -m <peer msp id>             [Required]
    -x <peers endpoint str>      [Required]
    -h Help
e.g:
    bash $0 -c channel1 -p ./NODES_ROOT/1ef243fd87 -m Org1MSP -x {\"host\": \"192.168.3.128\",\"port\": 7051}
EOF
exit 0
}

parse_params()
{
    while getopts "c:p:m:x:h" option;do
        case $option in
        c) channel_id="${OPTARG}"
            if [ -z "$channel_id" ]; then LOG_WARN "$channel_id not specified" && exit 1; fi
        ;;
        p) chain_path="${OPTARG}"
            if [ ! -d "$chain_path" ]; then LOG_WARN "$chain_path not exist" && exit 1; fi
        ;;
        m) peer_msp_id="${OPTARG}"
            if [ -z "$peer_msp_id" ]; then LOG_WARN "$peer_msp_id not specified" && exit 1; fi
        ;;
        x) anchor_peers="${OPTARG}"
            if [ -z "$anchor_peers" ]; then LOG_WARN "$anchor_peers not specified" && exit 1; fi
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

    configtxlator proto_decode --input "${chain_path}/config/channel-artifacts/config_block.pb" --type common.Block --output "${chain_path}/config/channel-artifacts/config_block.json"

    jq .data.data[0].payload.data.config "${chain_path}/config/channel-artifacts/config_block.json" > "${chain_path}/config/channel-artifacts/${peer_msp_id}config.json"

    jq ".channel_group.groups.Application.groups.${peer_msp_id}.values += {\"AnchorPeers\":{\"mod_policy\": \"Admins\",\"value\":{\"anchor_peers\": [${anchor_peers}]},\"version\": \"0\"}}" "${chain_path}/config/channel-artifacts/${peer_msp_id}config.json" > "${chain_path}/config/channel-artifacts/${peer_msp_id}modified_config.json"

    configtxlator proto_encode --input "${chain_path}/config/channel-artifacts/${peer_msp_id}config.json" --type common.Config --output "${chain_path}/config/channel-artifacts/original_config.pb"

    configtxlator proto_encode --input "${chain_path}/config/channel-artifacts/${peer_msp_id}modified_config.json" --type common.Config --output "${chain_path}/config/channel-artifacts/modified_config.pb"

    configtxlator compute_update --channel_id "${channel_id}" --original "${chain_path}/config/channel-artifacts/original_config.pb" --updated "${chain_path}/config/channel-artifacts/modified_config.pb" --output "${chain_path}/config/channel-artifacts/config_update.pb"

    configtxlator proto_decode --input "${chain_path}/config/channel-artifacts/config_update.pb" --type common.ConfigUpdate --output "${chain_path}/config/channel-artifacts/config_update.json"

    echo "{\"payload\":{\"header\":{\"channel_header\":{\"channel_id\":\"${channel_id}\", \"type\":2}},\"data\":{\"config_update\":$(cat ${chain_path}/config/channel-artifacts/config_update.json)}}}" | jq . > "${chain_path}/config/channel-artifacts/config_update_in_envelope.json"

    configtxlator proto_encode --input "${chain_path}/config/channel-artifacts/config_update_in_envelope.json" --type common.Envelope --output "${chain_path}/config/channel-artifacts/${peer_msp_id}anchors.tx" > "${logfile}" 2>&1
    #rm "${logfile}"
}

print_result()
{
    echo "=============================================================="
    LOG_INFO "Config Channel ${channel_id} config : ${chain_path}/config/channel-artifacts/${peer_msp_id}anchors.tx"
    LOG_INFO "All completed."
}

echo "=======$(date)========" >>"${logfile}"
echo "=======start config app channel===============" >>"${logfile}"
parse_params "$@"
main
print_result
echo "=======$(date)========" >>"${logfile}"
echo "=======end config app channel=================" >>"${logfile}"


