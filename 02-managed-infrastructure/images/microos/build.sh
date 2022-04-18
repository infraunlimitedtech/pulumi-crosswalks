#!/usr/bin/env bash

set -ex -o pipefail

COMBUSTION_DIR=init/combustion
KEY_PATH=${COMBUSTION_DIR}/id_edcsa.pub

mkdir -p $COMBUSTION_DIR

user=$(pulumi -s $(pulumi config get --path main.identitystack) stack output --json identity:ssh:server_access:credentials | jq -r .user)
pulumi -s $(pulumi config get --path main.identitystack) stack output --json identity:ssh:server_access:credentials | jq -r .publickey > $KEY_PATH

root_password=$(pulumi -s $(pulumi config get --path main.identitystack) stack output --show-secrets --json identity:local_users:root:password)
USER=${user} \
ROOT_PWD=$(echo -n ${root_password} | tr -d '"' | openssl passwd -6 -stdin ) \
KEY=$(basename ${KEY_PATH}) \
envsubst < script.tmpl > ${COMBUSTION_DIR}/script
chmod +x ${COMBUSTION_DIR}/script


mkisofs -o combustion.iso -V combustion -U -r -v -T -J -joliet-long $(dirname ${COMBUSTION_DIR})
qemu-img convert ./combustion.iso ./combustion.qcow2

rm -rfv $(dirname ${COMBUSTION_DIR})

wget https://download.opensuse.org/tumbleweed/appliances/openSUSE-MicroOS.x86_64-kvm-and-xen.qcow2

pulumi config set --path compute.libvirt.images.base "$(pwd)/openSUSE-MicroOS.x86_64-kvm-and-xen.qcow2"
pulumi config set --path compute.libvirt.images.combustion "$(pwd)/combustion.qcow2"
