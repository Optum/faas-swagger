# Swagger YAML Generator Function

This function generates a swagger yaml for the functions deployed in a openfaas 
gateway.

### PreReq

This function requires couple of secrets to be created

1. Base [swagger.yaml](./swagger.yaml)

```bash
$ kubectl create secret generic swagger --from-file=swagger.yaml -n openfaas-fn
```

2. Create the basic auth secret in `openfaas-fn` namespace. 
This is the same secret that is created in `openfaas` namespace.

```bash
$ BAP=( $(kubectl get secret basic-auth -n openfaas -ojsonpath='{.data.basic-auth-password}') )
$ BAU=( $(k get secret basic-auth -n openfaas -ojsonpath='{.data.basic-auth-user}') )
$ kubectl create secret generic basic-auth --from-literal basic-auth-password=$BAP  --from-literal basic-auth-user=$BAU -n openfaas-fn
```

### Deployment

```
$ faas up -f faas-swagger.yaml
```