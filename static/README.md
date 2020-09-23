# Swagger UI Static Files

Static files for swagger ui are generated using `github.com/rakyll/statik` .

[index.html](./swaggerui/index.html) has the name of the swagger yaml generator function.

```shell
statik -src=swaggerui/
```