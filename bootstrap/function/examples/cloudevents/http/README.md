```shell
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/cloudevents+json" \
  -d '{
    "data": {
      "message": "hello world!!!"
    },
    "id": "214b0346-7386-11ea-b6ae-acde48001122",
    "source": "changeme",
    "specversion": "1.0",
    "subject": "changeme",
    "type": "changeme",
    "time": "0001-01-01T00:00:00Z"
    }'
```