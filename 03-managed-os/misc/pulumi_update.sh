#/usr/bin/env bash
set -e

stack=$(pulumi stack --show-name)

function ask() {
  read -p "Do you wish to procced?" yn
    case $yn in
        [Yy]* ) echo Continue;;
        [Nn]* ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
}

for idx in {1..3}; do
  pulumi up -fy -t "urn:pulumi:${stack}::managed-os::k3os:index:Node::k3s-server0${idx}"
  echo "Waiting"
  sleep 60
  if [[ ! $NO_CONFIRMATION ]]; then continue; fi 
  ask
done

for idx in {1..3}; do
  pulumi up -t -fy "urn:pulumi:${stack}::managed-os::k3os:index:Node::k3s-agent0${idx}"
  echo "Waiting"
  sleep 60
  if [[ ! $NO_CONFIRMATION ]]; then continue; fi 
  ask
done

