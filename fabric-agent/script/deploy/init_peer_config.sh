#!/bin/bash
set -e


logfile=${PWD}/build.log
chain_path=
org_domain=
ca_name=
ca_ip=
ca_port=
peer_name=
peer_ip=
peer_pw=peerPw
node_id=
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
    -g <node name>               [Optional]
    -o <org domain>              [Required]
    -c <ca name>                 [Required]
    -i <ca ip>                   [Required]
    -t <ca port>                 [Required]
    -n <peer name>               [Required]
    -e <peer ip>                 [Required]
    -w <peer password>           [Optional]
    -h Help
e.g:
    bash $0 -p ./fabric/test_chain -o org1.example.com
EOF
exit 0
}

parse_params()
{
    while getopts "v:g:c:e:i:n:o:p:t:w:h" option;do
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
        n) peer_name="${OPTARG}"
            if [ -z "${peer_name}" ]; then LOG_WARN "${peer_name} is empty" && exit 1; fi
        ;;
        e) peer_ip="${OPTARG}"
            if [ -z "${peer_ip}" ]; then LOG_WARN "${peer_ip} is empty" && exit 1; fi
        ;;
        g) node_id="${OPTARG}"
            if [ -z "${node_id}" ]; then node_id="node_id is empty"; fi
        ;;
        v) node_name="${OPTARG}"
            if [ -z "${node_name}" ]; then LOG_WARN "${node_name} is empty" && exit 1; fi
        ;;
        w) peer_pw="${OPTARG}"
            if [ -z "${peer_pw}" ]; then peer_pw="${peer_name}pw"; fi
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

config_peer() {
  local check_result=$(check_user_exist ${peer_name})
  if [[ ${check_result} =~ 'already exists' ]] ;then
    echo "peer:[${peer_name}] already exists." >>"${logfile}"
    return;
  fi
  # 注册peer0
  fabric-ca-client register --caname "${ca_name}" --id.name "${peer_name}" --id.secret "${peer_pw}" --id.type peer --tls.certfiles "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem"
  # 登记peer0
  fabric-ca-client enroll -u "https://${peer_name}:${peer_pw}@${ca_ip}:${ca_port}" --caname "${ca_name}" -M "${chain_path}/organizations/peerOrganizations/${org_domain}/peers/${node_name}/msp" --tls.certfiles "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem"
  cp "${chain_path}/config.yaml" "${chain_path}/organizations/peerOrganizations/${org_domain}/msp/config.yaml"
  cp "${chain_path}/organizations/peerOrganizations/${org_domain}/msp/config.yaml" "${chain_path}/organizations/peerOrganizations/${org_domain}/peers/${node_name}/msp/config.yaml"

  #登记peer0的tls
  fabric-ca-client enroll -u "https://${peer_name}:${peer_pw}@${ca_ip}:${ca_port}" --caname "${ca_name}" -M "${chain_path}/organizations/peerOrganizations/${org_domain}/peers/${node_name}/tls" --enrollment.profile tls --csr.hosts "${peer_ip}" --tls.certfiles "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem"

  #五、将peer0的msp目录下的cacert复制到机构msp下
  cp -r "${chain_path}/organizations/peerOrganizations/${org_domain}/peers/${node_name}/msp/cacerts" "${chain_path}/organizations/peerOrganizations/${org_domain}/msp"
}

config_peer_user() {
    local username="user_${node_id}";
    local userPassword="${username}pw";
    local check_result=$(check_user_exist ${username})
    if [[ ${check_result} =~ 'already exists' ]] ;then
      echo "user:[${username}] already exists." >>"${logfile}"
      return;
    fi
    # 注册peer组织org1的user
    fabric-ca-client register --caname "${ca_name}" --id.name "${username}" --id.secret "${userPassword}" --id.type client --tls.certfiles "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem"
    #登记peer组织org1的user
    cp "${chain_path}/organizations/peerOrganizations/${org_domain}/peers/${node_name}/tls/tlscacerts/"* "${chain_path}/organizations/peerOrganizations/${org_domain}/peers/${node_name}/tls/ca.crt"
    cp "${chain_path}/organizations/peerOrganizations/${org_domain}/peers/${node_name}/tls/signcerts/"* "${chain_path}/organizations/peerOrganizations/${org_domain}/peers/${node_name}/tls/server.crt"
    cp "${chain_path}/organizations/peerOrganizations/${org_domain}/peers/${node_name}/tls/keystore/"* "${chain_path}/organizations/peerOrganizations/${org_domain}/peers/${node_name}/tls/server.key"

    fabric-ca-client enroll -u "https://${username}:${userPassword}@${ca_ip}:${ca_port}" --caname "${ca_name}" -M "${chain_path}/organizations/peerOrganizations/${org_domain}/users/user1_${node_id}@${org_domain}/msp" --tls.certfiles "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem"

    cp "${chain_path}/organizations/peerOrganizations/${org_domain}/msp/config.yaml" "${chain_path}/organizations/peerOrganizations/${org_domain}/users/user1_${node_id}@${org_domain}/msp/config.yaml"
}

config_peer_admin() {
     local username="admin_${org_domain}";
     local userPassword="${username}pw";
     local check_result=$(check_user_exist ${username})
     if [[ ${check_result} =~ 'already exists' ]] ;then
       echo "admin:[${username}] already exists." >>"${logfile}"
       return;
     fi
    #注册peer组织org1的admin
    fabric-ca-client register --caname "${ca_name}" --id.name "${username}" --id.secret "${userPassword}" --id.type admin --tls.certfiles "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem"
    #登记peer组织org1的admin
    fabric-ca-client enroll -u "https://${username}:${userPassword}@${ca_ip}:${ca_port}" --caname "${ca_name}" -M "${chain_path}/organizations/peerOrganizations/${org_domain}/users/Admin@${org_domain}/msp" --tls.certfiles "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem"
    mv "${chain_path}/organizations/peerOrganizations/${org_domain}/users/Admin@${org_domain}/msp/keystore/"* "${chain_path}/organizations/peerOrganizations/${org_domain}/users/Admin@${org_domain}/msp/keystore/priv_sk"
    cp "${chain_path}/organizations/peerOrganizations/${org_domain}/msp/config.yaml" "${chain_path}/organizations/peerOrganizations/${org_domain}/users/Admin@${org_domain}/msp/config.yaml"
}

init_peer_dir() {
    mkdir -p "${chain_path}/organizations/peerOrganizations/${org_domain}/peers/${node_name}"
}

init_env_param() {
  export PATH="${chain_path}/bin:$PATH"
  export FABRIC_CA_CLIENT_HOME="${chain_path}/fabric-ca-client/${org_domain}"
  export FABRIC_CFG_PATH="${chain_path}/config"
}

main()
{
  init_env_param
  init_peer_dir

  config_peer
  config_peer_user
  config_peer_admin
}

print_result()
{
    echo "=============================================================="
    LOG_INFO "Init peer[${peer_name}] conf in : ${chain_path}"
    LOG_INFO "All completed."
}

echo "=======$(date)========" >>"${logfile}"
echo "=======start init client env===============" >>"${logfile}"
parse_params "$@"
main
print_result
echo "=======$(date)========" >>"${logfile}"
echo "=======end init client env=================" >>"${logfile}"


