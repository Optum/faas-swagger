# faas-swagger as Openfaas function

You can deploy this utility as a openfaas function

### Requirements

This function requires 

1. Environment variable "openfaas_gateway"
2. 2 files to be created as secrets.

```
kubectl create secret generic swagger --from-file=swagger.yaml --from-file=sample.yaml
```

### Deployment

This function has been built and pushed to murugappans/faas-swagger:latest

You can use faas-swagger.yml to deploy this function to your platform. (Change the gateway in the file)

### Swagger UI

This function will serve the swagger.yaml, for the UI, please use the below image.

```
docker run -it -e openfaas_gateway=http://gateway:8080 -p 8080:8080 murugappans/swaggerui-openfaas
```