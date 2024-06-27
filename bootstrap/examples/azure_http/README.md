```shell
curl -X POST -H "Content-Type: application/json" -d '{
    "data": {
        "message": "hello world!!!"
    },
    "id": "214b0346-7386-11ea-b6ae-acde48001122",
    "source": "changeme",
    "specversion": "1.0",
    "subject": "changeme",
    "type": "changeme",
    "time": "0001-01-01T00:00:00Z"
}' 'http://localhost:7071/api/handler'
```


```shell
curl -X POST -H "Content-Type: application/json" -d '{
  "Data": {
    "req": {
      "Url": "http://localhost:7071/api/order",
      "Method": "POST",
      "Query": "{}",
      "Headers": {
        "Content-Type": [
          "application/json"
        ]
      },
      "Params": {},
      "Body": "{\"id\":1005,\"quantity\":2,\"color\":\"black\"}"
    }
  },
  "Metadata": {
  }
}' 'http://localhost:7071/api/handler'
```


