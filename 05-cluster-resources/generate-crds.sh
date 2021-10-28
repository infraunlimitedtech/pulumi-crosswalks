#!/usr/bin/env bash
set -xe

function generate_kilo_crds() {
  version=$(pulumi config get --path 'infra.kilo.version')
  local sources_path=$(pulumi config get --path 'infra.kilo.crds.path')

  mkdir -p ${sources_path}

  curl -s -O -L https://raw.githubusercontent.com/squat/kilo/${version}/manifests/crds.yaml --output-dir ${sources_path}
  crd2pulumi --goPath crds/generated/squat ${sources_path}/*.yaml --force
}

function generate_nginxinc_crds() {

  declare -a crds=(
    "k8s.nginx.org_globalconfigurations.yaml"
    #"k8s.nginx.org_virtualservers.yaml"
    "k8s.nginx.org_transportservers.yaml"
  )
  local version=$(pulumi config get --path 'infra.lb.nginxIngress.helm.version')
  local sources_path="crds/sources/nginxinc/nginx-ingress"

  mkdir -p ${sources_path}

  for f in ${crds[@]}; do 
    curl -s -O -L https://raw.githubusercontent.com/nginxinc/kubernetes-ingress/${version}/deployments/helm-chart/crds/${f} --output-dir ${sources_path}
    # name=$(echo $f | cut -f 2 -d '_' | cut -f 1 -d '.')
    #~/crd2pulumi/crd2pulumi --goPath crds/generated/nginxinc/kubernetes-ingress/${name} ${sources_path}/${f} --force
  done
  crd2pulumi --goPath crds/generated/nginxinc/kubernetes-ingress/ ${sources_path}/*.yaml --force

}

generate_kilo_crds
generate_nginxinc_crds


