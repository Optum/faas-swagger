# Swagger YAML Generator Function

This function generates a swagger yaml for the functions deployed in a openfaas 
gateway.

### PreReq

This function requires couple of secrets to be created

1. Base [swagger.yaml](./swagger.yaml)

```
$ kubectl create secret generic swagger --from-file=swagger.yaml -n openfaas-fn
```

1. Create the basic auth secret in `openfaas-fn` namespace. 
This is the same secret that is created in `openfaas` namespace.

### Deployment

```
$ faas up -f faas-swagger.yaml
```