#!/bin/bash
set -e

logfile=${PWD}/build.log
chain_path=
org_domain=
ca_cert=

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
    -c <ca cert file>            [Required]
    -h Help
e.g:
    bash $0 -p ./fabric/test_chain -o org1.example.com
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
        c) ca_cert="${OPTARG}"
            if [ -z "$ca_cert" ]; then LOG_WARN "$ca_cert is empty" && exit 1; fi
        ;;
        h) help;;
        *) LOG_WARN "invalid option $option";;
        esac
    done
}

init_org_dir() {
  mkdir -p "${chain_path}/organizations/ordererOrganizations"
  mkdir -p "${chain_path}/organizations/peerOrganizations"
  mkdir -p "${chain_path}/log"

  mkdir -p "${chain_path}/organizations/peerOrganizations/${org_domain}/msp"
  mkdir -p "${chain_path}/organizations/ordererOrganizations/${org_domain}/msp"
}

gen_org_msp() {
  mkdir -p "${chain_path}/organizations/peerOrganizations/${org_domain}/msp/tlscacerts"
  cp "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem" "${chain_path}/organizations/peerOrganizations/${org_domain}/msp/tlscacerts/ca.crt"
  
  mkdir -p "${chain_path}/organizations/peerOrganizations/${org_domain}/tlsca"
  cp "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem" "${chain_path}/organizations/peerOrganizations/${org_domain}/tlsca/tlsca.${org_domain}-cert.pem"
  
  mkdir -p "${chain_path}/organizations/peerOrganizations/${org_domain}/ca"
  cp "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem" "${chain_path}/organizations/peerOrganizations/${org_domain}/ca/ca.${org_domain}-cert.pem"
}



main()
{
  init_org_dir
  gen_org_msp
}

print_result()
{
    echo "=============================================================="
    LOG_INFO "Init org[${org_domain}] conf in : ${chain_path}"
    LOG_INFO "All completed."
}

echo "=======$(date)========" >>"${logfile}"
echo "=======start init org config===============" >>"${logfile}"
parse_params "$@"
main
print_result
echo "=======$(date)========" >>"${logfile}"
echo "=======end init org config=================" >>"${logfile}"


