#!/bin/bash
set -e

logfile=${PWD}/build.log
chain_path=
ca_name='ca-org1'
org_domain='ca.org1.example.com'
ca_admin='admin'
ca_admin_pw='adminpw'
ca_ip=
ca_port=

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
    -c <ca name>                 [Required]
    -i <ca ip>                   [Required]
    -t <ca port>                 [Required]
    -o <org domain>              [Required]
    -a <ca admin>                [Optional]
    -w <ca admin password>       [Optional]
    -h Help
e.g:
    bash $0 -p ./fabric/test_chain -c Org1 -i 192.168.1.1 -t 7054 -o org1.example.com
EOF
exit 0
}

parse_params()
{
    while getopts "a:c:i:o:p:t:w:h" option;do
        case $option in
        p) chain_path="${OPTARG}"
            if [ ! -d "$chain_path" ]; then mkdir -p "${chain_path}"; fi
        ;;
        c) ca_name="${OPTARG}"
            if [ -z "$ca_name" ]; then LOG_WARN "$ca_name is empty" && exit 1; fi
        ;;
        i) ca_ip="${OPTARG}"
            if [ -z "$ca_ip" ]; then LOG_WARN "$ca_ip is empty" && exit 1; fi
        ;;
        o) org_domain="${OPTARG}"
            if [ -z "$org_domain" ]; then LOG_WARN "$org_domain is empty" && exit 1; fi
        ;;
        t) ca_port="${OPTARG}"
            if [ -z "$ca_port" ]; then LOG_WARN "$ca_port is empty" && exit 1; fi
        ;;
        a) ca_admin="${OPTARG}"
            if [ -z "$ca_admin" ]; then ca_admin='admin'; fi
        ;;
        w) ca_admin_pw="${OPTARG}"
            if [ -z "$ca_admin_pw" ]; then ca_admin_pw='adminpw'; fi
        ;;
        h) help;;
        *) LOG_WARN "invalid option $option";;
        esac
    done
}

init_chain() {
  mkdir -p "${chain_path}/organizations/ordererOrganizations"
  mkdir -p "${chain_path}/organizations/peerOrganizations"
  mkdir -p "${chain_path}/fabric-ca-client/${org_domain}"
  mkdir -p "${chain_path}/log"

  prepare_bin
}

prepare_bin() {
    local scriptPath=$(dirname $0)
    cp "${scriptPath}/../tar/hyperledger-fabric-ca-linux-amd64-1.4.7.tar.gz" "${chain_path}"
    cp "${scriptPath}/../tar/hyperledger-fabric-linux-amd64-2.2.0.tar.gz" "${chain_path}"

    cd "${chain_path}"
    tar -xvf "hyperledger-fabric-ca-linux-amd64-1.4.7.tar.gz"
    tar -xvf "hyperledger-fabric-linux-amd64-2.2.0.tar.gz"
}

init_env_param() {
  export PATH="${chain_path}/bin:$PATH"
  export FABRIC_CA_CLIENT_HOME="${chain_path}/fabric-ca-client/${org_domain}"
  export FABRIC_CFG_PATH="${chain_path}/config"
}

init_ca() {
  fabric-ca-server init -b "${ca_admin}:${ca_admin_pw}"
  nohup fabric-ca-server start -b "${ca_admin}:${ca_admin_pw}" >> "${chain_path}/log/fabric-ca-server.log" 2>&1 &
  echo $! > "${chain_path}/fabric-ca-server.pid"
}

enroll_ca_admin() {
  fabric-ca-client enroll -u "https://${ca_admin}:${ca_admin_pw}@${ca_ip}:${ca_port}" --caname "${ca_name}" --tls.certfiles "${chain_path}/fabric-ca-server/${org_domain}/ca-cert.pem"
}

wait_ca_server_start() {
  local num=60
  local i=0
  while [ $i -le $num ]; do
      local port_status=$(sudo netstat -pan | grep -E "${ca_port}.*LISTEN")
      echo "ca server started.port_status:${port_status}" >>"${logfile}"
      if [[ $port_status =~ "LISTEN" ]]; then
        echo "ca server started." >>"${logfile}"
        break;
      else
        echo "Checking ca server.i=${i}" >>"${logfile}"
        i=$((i+1))
        sleep 1
      fi
  done

  if [[ $i -gt $num ]]; then
    echo "ca server has not started yet." >>"${logfile}"
  fi
}

main()
{
  init_chain
  init_env_param
  #init_ca
  #wait_ca_server_start
  enroll_ca_admin
}

print_result()
{
    echo "=============================================================="
    LOG_INFO "Init client env in : ${chain_path}"
    LOG_INFO "All completed."
}

echo "=======$(date)========" >>"${logfile}"
echo "=======start init client env===============" >>"${logfile}"
parse_params "$@"
main
print_result
echo "=======$(date)========" >>"${logfile}"
echo "=======end init client env=================" >>"${logfile}"


