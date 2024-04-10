package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/lukeshay/g/datadog"
	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type DatadogContextHandler struct {
	slog.Handler
	options datadog.IntializeOptions
}

type ctxKey string

const (
	slogFields            ctxKey = "slog_fields"
	SlogCustomAttrsCtxKey ctxKey = "slog_custom_attrs"
)

// Handle adds contextual attributes to the Record before calling the underlying
// handler
func (h DatadogContextHandler) Handle(ctx context.Context, r slog.Record) error {
	ddArgs := []any{
		"service", h.options.DDService,
		"version", h.options.DDVersion,
		"env", h.options.DDEnv,
	}

	if span, found := tracer.SpanFromContext(ctx); found {
		ddArgs = append(
			ddArgs,
			"span_id", span.Context().SpanID(),
			"trace_id", span.Context().TraceID(),
		)
	}
	r.Add(
		slog.Group(
			"dd",
			ddArgs...,
		),
	)

	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	if customAttrs, ok := ctx.Value(SlogCustomAttrsCtxKey).([]slog.Attr); ok {
		for _, v := range customAttrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

// AppendCtx adds an slog attribute to the provided context so that it will be
// included in any Record created with such context
func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}

	v := []slog.Attr{}
	v = append(v, attr)

	return context.WithValue(parent, slogFields, v)
}

func Initialize(options datadog.IntializeOptions) {
	zerologL := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
	h := &DatadogContextHandler{
		Handler: slogzerolog.Option{Logger: &zerologL}.NewZerologHandler(),
		options: options,
	}

	logger := slog.New(h)

	slog.SetDefault(logger)
}
