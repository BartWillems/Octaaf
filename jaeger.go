package main

import (
	"fmt"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	jaegerClient "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// Tracer is the opentracing instance
var Tracer opentracing.Tracer

func initJaeger() io.Closer {
	cfg := &jaegercfg.Configuration{
		ServiceName: settings.Jaeger.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			// Only log the spans in development mode
			LogSpans:           settings.Environment == development,
			LocalAgentHostPort: fmt.Sprintf("%v:%v", settings.Jaeger.AgentHost, settings.Jaeger.AgentPort),
		},
	}

	var err error
	var closer io.Closer
	Tracer, closer, err = cfg.New(settings.Jaeger.ServiceName, jaegercfg.Logger(jaegerClient.StdLogger))

	if err != nil {
		log.Panicf("Jaeger init error: %v", err)
	}

	return closer
}
