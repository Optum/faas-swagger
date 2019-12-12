# Swagger for Openfaas

Swagger interface for openfaas functions

### Whats it About

* Document your Openfaas functions in OpenAPI 3.0 Spec and visualize it in a swagger ui.
* Uses the annotations in the function descriptor yml to look for the swagger config.
* Test functions from swagger ui

----------------

### Demo

<a href="http://www.youtube.com/watch?feature=player_embedded&v=p6Vi5vIjO5I
" target="_blank"><img src="http://img.youtube.com/vi/p6Vi5vIjO5I/0.jpg" 
alt="Demo Swagger Faas" width="480" height="360" border="20" /></a>

----------------

### Deploy as Standalone entity 

* Docker image of the utility is murugappans/faas-swagger
* Use the k8s-deploy/k8.yaml to deploy to your kubernetes namespace

```
kubectl apply -f k8s-deploy/k8.yaml -n openfaas
```

-------------

### Deploy as Function

To deploy this utility as a function. Please follow the steps mentioned [here](./swagger-as-function)

-------------

### Using the utility

Add your api spec (json format) in the function descriptor as annotation. Use this [example](./example.yaml).

After deploying your function with this annotation, you should be able see the paths in swagger ui

* We are following open api 3.0 spec
* Please make sure the json is well formatted.
* In swagger 3.0, the spec is defined in yaml. You can define in yaml and convert to json using online editor like [this](https://codebeautify.org/yaml-to-json-xml-csv)
* If annotatation is not provided , by default just the path will be shown.

---------

### Authentication

Currently does not support Openfaas auth plugins. But will work if openfaas is deployed behind a proxy and the proxy is handling the auth needs.

----------
### Contributing to the Project
The team is open to contributions to our project. For more details, see our [Contribution Guide.](./docs/CONTRIBUTING.md)