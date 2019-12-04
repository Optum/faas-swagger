# Swagger for Openfaas

Utility for generating open api spec documentation and enabling swagger ui for openfaas functions.

Uses the openfaas annotations to pulls the swagger config

### About

* Swagger UI to serve the swagger yaml
* On loading the swagger UI, the yml is generated on the fly by looking at the annotations on openfaas functions
* Just a page refresh would load any new changes
* Test functions from swagger

----------------

### Build and Deploy the Utility

#### Static content

Use [statik](https://github.com/rakyll/statik) to generate the static files related to swagger ui

`Before generating the static binary, please change the following swagger.yaml path in swaggerui/index.html. 
This is the path where this utility is serving the swagger yaml. http://host:port/swagger.yaml`

```
statik -src=swaggerui/
```
This should generate a folder called statik

#### Build the go binary, docker and push
```
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o release/faas-swagger .
```

#### Deploy

Use the k8.yaml to deploy to your kubernetes environment.

Change the following before deployment

##### 1. Environment variables

Provide the below environment variables while deploying

* `openfaas_gateway` - gateway url
* `swagger_file_path` - location of swagger files (mentioned in dockerfile)

##### 2. Configuration Map

CM named Swagger has the skeleton swagger yaml.
Change the server url in the k8.yaml to gateway endpoint

##### 3. Docker Image

Change the docker image to point to the docker repo where this utility image was pushed.

-------------

### Using the utility

Add your api spec (json format) in the function descriptor as annotation. Use this [example](./example.yaml).

After deploying your function with this annotation, you should be able see the paths in swagger ui

* We are following open api 3.0 spec
* Please make sure the json is well formatted.
* In swagger 3.0 the spec is defined in yaml. You can define in yaml and convert to json using online editor like [this](https://codebeautify.org/yaml-to-json-xml-csv)
* If annotatation is not provided , by default just the path will be shown.

---------

### Authentication

Authentication in faas-swagger is designed to work with openfaas deployed behind a proxy. 
Currently does not support Openfaas auth plugins.

----------
### Contributing to the Project
The team is open to contributions to our project. For more details, see our [Contribution Guide.](./docs/CONTRIBUTING.md)