module github.com/xgodev/boost

go 1.26.5

require (
	cloud.google.com/go/bigquery v1.77.0
	cloud.google.com/go/firestore v1.22.0

	cloud.google.com/go/pubsub v1.50.3
	github.com/DataDog/dd-trace-go/contrib/aws/aws-sdk-go-v2/v2 v2.9.1
	github.com/DataDog/dd-trace-go/contrib/database/sql/v2 v2.9.1
	github.com/DataDog/dd-trace-go/contrib/go.mongodb.org/mongo-driver.v2/v2 v2.9.1
	github.com/DataDog/dd-trace-go/contrib/go.mongodb.org/mongo-driver/v2 v2.9.1
	github.com/DataDog/dd-trace-go/contrib/google.golang.org/grpc/v2 v2.9.1
	github.com/DataDog/dd-trace-go/contrib/labstack/echo.v4/v2 v2.9.1
	github.com/DataDog/dd-trace-go/contrib/redis/go-redis.v9/v2 v2.9.1
	github.com/DataDog/dd-trace-go/v2 v2.9.1
	github.com/XSAM/otelsql v0.42.0
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/allegro/bigcache/v3 v3.1.0
	github.com/aws/aws-sdk-go-v2 v1.42.1
	github.com/aws/aws-sdk-go-v2/config v1.32.28
	github.com/aws/aws-sdk-go-v2/credentials v1.19.27
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.45.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.105.0
	github.com/aws/aws-sdk-go-v2/service/sns v1.41.0
	github.com/aws/aws-sdk-go-v2/service/sqs v1.45.0

	github.com/coocood/freecache v1.2.7
	github.com/dubonzi/otelresty v1.6.0
	github.com/elastic/go-elasticsearch/v8 v8.19.6
	github.com/globocom/echo-prometheus v0.1.2
	github.com/go-logr/logr v1.4.3
	github.com/go-playground/validator/v10 v10.30.3
	github.com/go-redis/redis/v7 v7.4.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-resty/resty/v2 v2.17.2
	github.com/gobeam/stringy v0.0.7
	github.com/goccy/go-json v0.10.6
	github.com/gocql/gocql v1.7.0
	github.com/godror/godror v0.51.0
	github.com/google/uuid v1.6.0
	github.com/graphql-go/graphql v0.8.1
	github.com/graphql-go/handler v0.2.4
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/go-memdb v1.3.5
	github.com/hiko1129/echo-pprof v1.0.1

	github.com/jedib0t/go-pretty/v6 v6.8.1
	github.com/jlaffaye/ftp v0.2.1
	github.com/knadh/koanf v1.5.0
	github.com/labstack/echo/v4 v4.15.2

	github.com/nats-io/nats.go v1.52.0
	github.com/panjf2000/ants/v2 v2.12.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.23.3-0.20251103151724-a5ae20370e5e

	go.mongodb.org/mongo-driver/v2 v2.6.0
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.69.0
	go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo v0.69.0
	go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/v2/mongo/otelmongo v0.0.0-20260611090623-6090c504d87e
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.69.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.69.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.69.0
	go.opentelemetry.io/otel v1.44.0

	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.44.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.44.0
	go.opentelemetry.io/otel/metric v1.44.1-0.20260625150014-c84013202f01
	go.opentelemetry.io/otel/sdk v1.44.1-0.20260625150014-c84013202f01
	go.opentelemetry.io/otel/sdk/metric v1.44.1-0.20260625150014-c84013202f01
	go.opentelemetry.io/otel/trace v1.44.1-0.20260625150014-c84013202f01
	go.uber.org/fx v1.24.0
	go.uber.org/zap v1.28.0
	gocloud.dev v0.46.0
	gocloud.dev/pubsub/kafkapubsub v0.46.0

	golang.org/x/sync v0.21.0
	golang.org/x/text v0.38.0
	google.golang.org/api v0.284.0
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.12-0.20260120151049-f2248ac996af
	gopkg.in/matryer/try.v1 v1.0.0-20150601225556-312d2599e12e
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	k8s.io/client-go v0.36.1

	golang.org/x/sync v0.22.0
	golang.org/x/text v0.39.0
	google.golang.org/api v0.287.1
	google.golang.org/grpc v1.82.0
	google.golang.org/protobuf v1.36.12-0.20260120151049-f2248ac996af
	gopkg.in/matryer/try.v1 v1.0.0-20150601225556-312d2599e12e
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	k8s.io/client-go v0.36.2
)
