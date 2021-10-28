
### Usefull
```
export KUBECONFIG=/tmp/kube.yml; pulumi stack output "infra:kube:config" --show-secrets > $KUBECONFIG  && k9s
```
