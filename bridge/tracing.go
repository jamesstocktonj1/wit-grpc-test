package main

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	wrpcnats "wrpc.io/go/nats"
)

var _ propagation.TextMapCarrier = &wrpcProp{}

type wrpcProp struct {
	h nats.Header
}

func (p *wrpcProp) Get(key string) string {
	return p.h.Get(key)
}

func (p *wrpcProp) Set(key string, value string) {
	p.h.Set(key, value)
}

func (p *wrpcProp) Keys() []string {
	keys := []string{}
	for k := range p.h {
		keys = append(keys, k)
	}
	return keys
}

func injectTraceHeader(_ctx context.Context) context.Context {
	carrier := &wrpcProp{
		h: nats.Header{},
	}
	otel.GetTextMapPropagator().Inject(_ctx, carrier)
	return wrpcnats.ContextWithHeader(_ctx, carrier.h)
}
