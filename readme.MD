# arex analysis

## json-schema
### resource
* [json-schema](http://json-schema.org/)
* [cn json-schema](https://json-schema.apifox.cn/)

#### RESTful API
##### Query all json-schemas
```
GET    /schemas                  --> github.com/arextest/arexAnalysis/arex.getSchemas (6 handlers)
DEMO http://{{analysis_url}}/schemas?limit=10
```

##### Query json-schema By KEY
```[GIN-debug] GET    /schema/:key              --> github.com/arextest/arexAnalysis/arex.getSchemaByKey (6 handlers)
DEMO http://{{analysis_url}}/schema/aaabccc
```

#### Save json-schema (json format) to database
```
[GIN-debug] POST   /schema/:key              --> github.com/arextest/arexAnalysis/arex.postSchema (6 handlers)
DEMO
POST http://{{analysis_url}}/schema/aaaaaaaaa
{json-schema}
return success/fail
```

#### Parse JSON and save json-schema to database
```
[GIN-debug] PUT    /schema/:key              --> github.com/arextest/arexAnalysis/arex.putSchema (6 handlers)
DEMO
PUT http://{{analysis_url}}/schema/prometheus
{json}
return {json-schema}
```

#### Parse JSON to json-schema, and merge to base json-schema and save
```
[GIN-debug] PATCH  /schema/:key              --> github.com/arextest/arexAnalysis/arex.patchSchema (6 handlers)
PATCH http://{{analysis_url}}/schema/prometheus
{json}
return {merged json-schema}
```

#### Delete json-schema by key
```
[GIN-debug] DELETE /schema/:key              --> github.com/arextest/arexAnalysis/arex.deleteSchema (6 handlers)
```

### Validate JSON By schema
#### Valid JSON by json-schema GET request
```
[GIN-debug] GET    /validation/:key          --> github.com/arextest/arexAnalysis/arex.getValidation (6 handlers)
DEMO 
GET http://{{analysis_url}}/validation/aaabccc
{json}
return
{
    "message": "success"
}
```

##### 
``` Valid JSON by json-schema POST request
[GIN-debug] POST   /validation               --> github.com/arextest/arexAnalysis/arex.postValidation (6 handlers)
POST http://{{analysis_url}}/validation
{
    "key": "aaabccc",
    "schema": "",
    "input":"",
    "result": ""
}
return 
{
    "message": "validation failed.jsonschema: '/id' does not validate with file:///app/jason-schema#/properties/id/type: expected number, but got string"
}
```


### Compare two json and result differ
```
[GIN-debug] POST   /comparing                --> github.com/arextest/arexAnalysis/arex.postComparing (6 handlers)
DEMO
POST http://{{analysis_url}}/comparing
{
    vx:"",
    vy:"",
    options:""
}
return 
[
    {
        "path": "root[\"panelId\"]",
        "vx": "18",
        "vy": "[18 19]"
    },
    {
        "path": "root[\"state\"]",
        "vx": "ok",
        "vy": "[ok cancel]"
    }
]
```


