package opentelemetry

import (
	"time"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/extra/middleware"
	xotel "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	meter                    metric.Meter
	messagesProcessed        metric.Int64Counter
	messageProcessingLatency metric.Float64Histogram
)

func initMeter() {
	// Inicializando um Meter (usando noop como um exemplo)
	meter = xotel.MeterProvider.Meter("boost_function")

	// Configurando o contador de mensagens processadas
	var err error
	messagesProcessed, err = meter.Int64Counter(
		"boost_function_messages_processed_total",
		metric.WithDescription("Number of messages processed"),
	)
	if err != nil {
		panic("Failed to create counter: " + err.Error())
	}

	// Configurando o histograma para latência de processamento
	messageProcessingLatency, err = meter.Float64Histogram(
		"boost_function_message_processing_latency_seconds",
		metric.WithDescription("Time taken to process message"),
	)
	if err != nil {
		panic("Failed to create histogram: " + err.Error())
	}
}

func init() {
	initMeter() // Inicializando o meter
}

type OpenTelemetry[T any] struct {
}

func (c *OpenTelemetry[T]) Exec(ctx *middleware.AnyErrorContext[T], exec middleware.AnyErrorExecFunc[T], fallbackFunc middleware.AnyErrorReturnFunc[T]) (T, error) {

	// Iniciando um novo span para tracing
	tracer := xotel.NewTracer("boost_function_tracer")

	ctxTrace, span := tracer.Start(ctx.GetContext(), "ProcessMessage")
	defer span.End()

	// Medindo a latência manualmente
	startTime := time.Now()

	// Processando a mensagem
	e, err := ctx.Next(exec, fallbackFunc)

	// Calculando a duração
	duration := time.Since(startTime).Milliseconds()

	// Determinando status
	status := "success"
	if err != nil {
		status = "error"
		span.RecordError(err) // Gravando erro no span
	}

	// Adicionando atributos ao span
	span.SetAttributes(
		attribute.String("function_name", boost.ApplicationName()),
		attribute.String("status", status),
		attribute.Int64("processing_duration_ms", duration),
	)

	// Gravando métricas
	messagesProcessed.Add(ctxTrace, 1, metric.WithAttributes(
		attribute.String("status", status),
		attribute.String("function_name", boost.ApplicationName()),
	))

	messageProcessingLatency.Record(ctxTrace, float64(duration), metric.WithAttributes(
		attribute.String("status", status),
		attribute.String("function_name", boost.ApplicationName()),
	))

	return e, err
}

func NewAnyErrorMiddleware[T any]() middleware.AnyErrorMiddleware[T] {
	return NewOpenTelemetry[T]()
}

func NewOpenTelemetry[T any]() *OpenTelemetry[T] {
	return &OpenTelemetry[T]{}
}
