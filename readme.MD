# arex analysis

## json-schema
### resource
* [json-schema](http://json-schema.org/)
* [cn json-schema](https://json-schema.apifox.cn/)

#### RESTful API
* GET    /schemas                  --> local/arex-reporter/arex.getSchemas (6 handlers)
* GET    /schema/:key              --> local/arex-reporter/arex.getSchemaByKey (6 handlers)
* POST   /schema/:key              --> local/arex-reporter/arex.postSchema (6 handlers)
* PUT    /schema/:key              --> local/arex-reporter/arex.putSchema (6 handlers)
* DELETE /schema/:key              --> local/arex-reporter/arex.deleteSchema (6 handlers)
* GET    /validation/:key          --> local/arex-reporter/arex.getValidation (6 handlers)
* POST   /validation               --> local/arex-reporter/arex.postValidation (6 handlers)
* POST   /comparing                --> local/arex-reporter/arex.postComparing (6 handlers)
