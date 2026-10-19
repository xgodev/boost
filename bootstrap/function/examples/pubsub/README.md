# pubsub

App `worker` runs a function on the messages of the Google Cloud Pub/Sub
subscription `messages-worker`: it replies with a greeting, which the
publisher middleware publishes with the pubsub driver on topic `replies`,
named after the subject of the reply.

```shell
gcloud beta emulators pubsub start --project=example --host-port=localhost:8085
```

Create the topics and the subscription, then run the worker:

```shell
export PUBSUB_EMULATOR_HOST=localhost:8085
gcloud pubsub topics create messages replies --project=example
gcloud pubsub subscriptions create messages-worker --topic=messages --project=example
CONF=configs/config.yaml go run ./cmd/worker
```

```shell
gcloud pubsub topics publish messages --message='{"message":"hi"}' --project=example
```
