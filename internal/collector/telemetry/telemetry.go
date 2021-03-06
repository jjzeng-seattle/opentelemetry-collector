// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package telemetry controls the telemetry settings to be used in the collector.
package telemetry

import (
	"flag"
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/internal/version"
)

const (
	metricsAddrCfg   = "metrics-addr"
	metricsLevelCfg  = "metrics-level"
	metricsPrefixCfg = "metrics-prefix"

	// Telemetry levels
	//
	// None indicates that no telemetry data should be collected.
	None Level = iota - 1
	// Basic is the default and covers the basics of the service telemetry.
	Basic
	// Normal adds some other indicators on top of basic.
	Normal
	// Detailed adds dimensions and views to the previous levels.
	Detailed
)

var (
	// Command-line flags that control publication of telemetry data.
	metricsLevelPtr  *string
	metricsAddrPtr   *string
	metricsPrefixPtr *string

	useLegacyMetricsPtr *bool
	useNewMetricsPtr    *bool
	addInstanceIDPtr    *bool
)

func Flags(flags *flag.FlagSet) {
	metricsLevelPtr = flags.String(
		metricsLevelCfg,
		"BASIC",
		"Output level of telemetry metrics (NONE, BASIC, NORMAL, DETAILED)")

	// At least until we can use a generic, i.e.: OpenCensus, metrics exporter
	// we default to Prometheus at port 8888, if not otherwise specified.
	metricsAddrPtr = flags.String(
		metricsAddrCfg,
		GetMetricsAddrDefault(),
		"[address]:port for exposing collector telemetry.")

	metricsPrefixPtr = flags.String(
		metricsPrefixCfg,
		"otelcol",
		"Prefix to the metrics generated by the collector.")

	useLegacyMetricsPtr = flags.Bool(
		"legacy-metrics",
		false,
		"Flag to control usage of legacy metrics",
	)

	useNewMetricsPtr = flags.Bool(
		"new-metrics",
		true,
		"Flag to control usage of new metrics",
	)

	addInstanceIDPtr = flags.Bool(
		"add-instance-id",
		true,
		"Flag to control the addition of 'service.instance.id' to the collector metrics.")
}

// GetMetricsAddrDefault returns the default metrics bind address and port depending on
// the current build type.
func GetMetricsAddrDefault() string {
	if version.IsDevBuild() {
		// Listen on localhost by default for dev builds to avoid security prompts.
		return "localhost:8888"
	}
	return ":8888"
}

// Level of telemetry data to be generated.
type Level int8

func GetAddInstanceID() bool {
	return *addInstanceIDPtr
}

// GetLevel returns the Level represented by the string. The parsing is case-insensitive
// and it returns error if the string value is unknown.
func GetLevel() (Level, error) {
	var level Level
	var str string

	if metricsLevelPtr != nil {
		str = strings.ToLower(*metricsLevelPtr)
	}

	switch str {
	case "none":
		level = None
	case "", "basic":
		level = Basic
	case "normal":
		level = Normal
	case "detailed":
		level = Detailed
	default:
		return level, fmt.Errorf("unknown metrics level %q", str)
	}

	return level, nil
}

func GetMetricsAddr() string {
	return *metricsAddrPtr
}

func GetMetricsPrefix() string {
	return *metricsPrefixPtr
}

func UseLegacyMetrics() bool {
	return *useLegacyMetricsPtr
}

func UseNewMetrics() bool {
	return *useNewMetricsPtr
}
