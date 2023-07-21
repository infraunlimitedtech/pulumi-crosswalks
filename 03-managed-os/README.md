
### Useful
```
export KUBECONFIG=/tmp/kube.yml; pulumi stack output "infra:kube:config" --show-secrets > $KUBECONFIG  && k9s
```

## Stacks
### Dev RU stack
Based on hetzner provider.
```
pulumi stack output os:wireguard:config --show-secrets -s hetzner-ru1-dev > ~/wg-dev.conf && wg-quick up ~/wg-dev.conf
```
