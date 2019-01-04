package jaeger

import (
	"fmt"
	"io"
	"log"

	opentracing "github.com/opentracing/opentracing-go"
	jaegerClient "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// Config contains the required runtime configuration for connecting to the jaeger host
type Config struct {
	ServiceName string `toml:"service_name" env:"JAEGER_SERVICE_NAME"`
	AgentHost   string `toml:"agent_host" env:"JAEGER_AGENT_HOST"`
	AgentPort   int    `toml:"agent_port" env:"JAEGER_AGENT_PORT"`
}

// Tracer is the opentracing instance
var Tracer opentracing.Tracer

// Init creates a connection to the jaeger agent and returns a closer
func Init(serviceName string, agentHost string, agentPort int, shouldLog bool) (io.Closer, error) {
	cfg := &jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			// Only log the spans in development mode
			LogSpans:           shouldLog,
			LocalAgentHostPort: fmt.Sprintf("%v:%v", agentHost, agentPort),
		},
	}

	var err error
	var closer io.Closer
	Tracer, closer, err = cfg.New(serviceName, jaegercfg.Logger(jaegerClient.StdLogger))

	if err != nil {
		log.Fatalf("Jaeger init error: %v", err)
		return nil, err
	}

	return closer, nil
}
