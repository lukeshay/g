package datadog

import (
	"fmt"
	"strings"

	"github.com/DataDog/datadog-go/statsd"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

type IntializeOptions struct {
	DDAgentHost string
	DDEnv       string
	DDVersion   string
	DDService   string
}

func Initialize(options IntializeOptions) func() error {
	if !strings.Contains(options.DDEnv, "local") && options.DDAgentHost != "" {
		tracer.Start(
			tracer.WithEnv(options.DDEnv),
			tracer.WithService(options.DDService),
			tracer.WithServiceVersion(options.DDVersion),
			tracer.WithLogStartup(false),
			tracer.WithDebugMode(false),
			tracer.WithAgentAddr(fmt.Sprintf("%s:8126", options.DDAgentHost)),
		)
		profiler.Start(
			profiler.WithEnv(options.DDEnv),
			profiler.WithService(options.DDService),
			profiler.WithVersion(options.DDVersion),
			profiler.WithProfileTypes(
				profiler.CPUProfile,
				profiler.HeapProfile,
				profiler.BlockProfile,
				profiler.GoroutineProfile,
				profiler.MetricsProfile,
				profiler.MutexProfile,
			),
			profiler.WithLogStartup(false),
			profiler.WithAgentAddr(fmt.Sprintf("%s:8126", options.DDAgentHost)),
		)
		statsd.New(fmt.Sprintf("http://%s:8125", options.DDAgentHost))
	}

	return func() error {
		tracer.Stop()

		return nil
	}
}
