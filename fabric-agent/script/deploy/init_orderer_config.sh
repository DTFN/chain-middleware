#!/bin/bash
set -e

logfile=${PWD}/build.log
chain_path=
org_domain=
ca_name=
ca_ip=
ca_port=
orderer_name=
orderer_ip=
orderer_pw=ordererPw
node_name=

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
    -c <ca name>                 [Required]
    -i <ca ip>                   [Required]
    -t <ca port>                 [Required]
    -n <orderer name>            [Required]
    -e <orderer ip>              [Required]
    -w <orderer password>        [Optional]
    -h Help
e.g:
    bash $0 -p ./fabric/test_chain -o org1.example.com
EOF
exit 0
}

parse_params()
{
    while getopts "v:c:e:i:n:o:p:t:w:h" option;do
        case $option in
        p) chain_path="${OPTARG}"
            if [ ! -d "${chain_path}" ]; then LOG_WARN "${chain_path} not exist" && exit 1; fi
        ;;
        o) org_domain="${OPTARG}"
            if [ -z "${org_domain}" ]; then LOG_WARN "${org_domain} is empty" && exit 1; fi
        ;;
        c) ca_name="${OPTARG}"
            if [ -z "${ca_name}" ]; then LOG_WARN "${ca_name} is empty" && exit 1; fi
        ;;
        i) ca_ip="${OPTARG}"
            if [ -z "${ca_ip}" ]; then LOG_WARN "${ca_ip} is empty" && exit 1; fi
        ;;
        t) ca_port="${OPTARG}"
            if [ -z "${ca_port}" ]; then LOG_WARN "${ca_port} is empty" && exit 1; fi
        ;;
        n) orderer_name="${OPTARG}"
            if [ -z "${orderer_name}" ]; then LOG_WARN "${orderer_name} is empty" && exit 1; fi
        ;;
        e) orderer_ip="${OPTARG}"
            if [ -z "${orderer_ip}" ]; then LOG_WARN "${orderer_ip} is empty" && exit 1; fi
        ;;
        v) node_name="${OPTARG}"
            if [ -z "${node_name}" ]; then LOG_WARN "${node_name} is empty" && exit 1; fi
        ;;
        w) orderer_pw="${OPTARG}"
            if [ -z "${orderer_pw}" ]; then orderer_pw="${orderer_pw}pw"; fi
        ;;
        h) help;;
        *) LOG_WARN "invalid option $option";;
        esac
    done
}

check_user_exist() {
    local username=$1
    local res=$(fabric-ca-client identity list  --tls.certfiles ${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem)
    if [[ ${res} =~ ${username} ]] ;then
      echo "${username} already exists."
    else
      echo "${username} is new."
    fi
}

config_orderer() {
  local check_result=$(check_user_exist ${node_name})
  if [[ ${check_result} =~ 'already exists' ]] ;then
    echo "orderer:[${node_name}] already exists." >>"${logfile}"
    return;
  fi
  #注册orderer0
  fabric-ca-client register --caname "${ca_name}" --id.name "${node_name}" --id.secret "${orderer_pw}" --id.type orderer --tls.certfiles "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem"

  #登记orderer0
  fabric-ca-client enroll -u "https://${node_name}:${orderer_pw}@${ca_ip}:${ca_port}" --caname "${ca_name}" -M "${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}/msp" --tls.certfiles "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem"
  cp "${chain_path}/config.yaml" "${chain_path}/organizations/ordererOrganizations/${org_domain}/msp/config.yaml"
  cp "${chain_path}/organizations/ordererOrganizations/${org_domain}/msp/config.yaml" "${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}/msp/config.yaml"

  #登记orderer0的tls
  fabric-ca-client enroll -u "https://${node_name}:${orderer_pw}@${ca_ip}:${ca_port}" --caname "${ca_name}" -M "${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}/tls" --enrollment.profile tls --csr.hosts "${orderer_ip}" --tls.certfiles "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem"

  cp "${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}/tls/tlscacerts/"* "${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}/tls/ca.crt"
  cp "${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}/tls/signcerts/"* "${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}/tls/server.crt"
  cp "${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}/tls/keystore/"* "${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}/tls/server.key"

  mkdir -p "${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}/msp/tlscacerts"
  cp "${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}/tls/tlscacerts/"* "${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}/msp/tlscacerts/tlsca.${org_domain}-cert.pem"

  #五、将orderer0的msp目录下的cacert复制机构下
  cp -r "${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}/msp/cacerts" "${chain_path}/organizations/ordererOrganizations/${org_domain}/msp"
}

config_admin(){
  local username="adminOrderer_${org_domain}";
  local userPassword="${username}pw";
  local check_result=$(check_user_exist ${username})
  if [[ ${check_result} =~ 'already exists' ]] ;then
    echo "orderer admin:[${username}] already exists." >>"${logfile}"
    return;
  fi
  #注册orderer组织org1的admin
  fabric-ca-client register --caname "${ca_name}" --id.name "${username}" --id.secret "${userPassword}" --id.type admin --tls.certfiles "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem"

  #登记orderer组织org1的admin
  fabric-ca-client enroll -u "https://${username}:${userPassword}@${ca_ip}:${ca_port}" --caname "${ca_name}" -M "${chain_path}/organizations/ordererOrganizations/${org_domain}/users/Admin@${org_domain}/msp" --tls.certfiles "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem"
  mv "${chain_path}/organizations/ordererOrganizations/${org_domain}/users/Admin@${org_domain}/msp/keystore/"* "${chain_path}/organizations/ordererOrganizations/${org_domain}/users/Admin@${org_domain}/msp/keystore/priv_sk"
  cp "${chain_path}/organizations/ordererOrganizations/${org_domain}/msp/config.yaml" "${chain_path}/organizations/ordererOrganizations/${org_domain}/users/Admin@${org_domain}/msp/config.yaml"
}

init_orderer_dir() {

  local ordererDir="${chain_path}/organizations/ordererOrganizations/${org_domain}/orderers/${node_name}"
  if [ ! -f "${ordererDir}" ]; then mkdir -p "${ordererDir}"; fi
  local orderTlsCertDir="${chain_path}/organizations/ordererOrganizations/${org_domain}/msp/tlscacerts"
  if [ ! -f "${orderTlsCertDir}" ]; then mkdir -p "${orderTlsCertDir}"; fi
  cp "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem" "${chain_path}/organizations/ordererOrganizations/${org_domain}/msp/tlscacerts/tlsca.${org_domain}-cert.pem"
  local ordererTlsCaDir="${chain_path}/organizations/ordererOrganizations/${org_domain}/tlsca"
  if [ ! -f "${ordererTlsCaDir}" ]; then mkdir -p "${ordererTlsCaDir}"; fi
  cp "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem" "${chain_path}/organizations/ordererOrganizations/${org_domain}/tlsca/tlsca.${org_domain}-cert.pem"
}


init_env_param() {
  export PATH="${chain_path}/bin:$PATH"
  export FABRIC_CA_CLIENT_HOME="${chain_path}/fabric-ca-client/${org_domain}"
  export FABRIC_CFG_PATH="${chain_path}/config"
}

main()
{
  init_env_param
  init_orderer_dir

  config_orderer
  config_admin
}

print_result()
{
    echo "=============================================================="
    LOG_INFO "Init orderer[${orderer_name}] conf in : ${chain_path}"
    LOG_INFO "All completed."
}

echo "=======$(date)========" >>"${logfile}"
echo "=======start orderer config ===============" >>"${logfile}"
parse_params "$@"
main
print_result
echo "=======$(date)========" >>"${logfile}"
echo "=======end init orderer config=================" >>"${logfile}"


