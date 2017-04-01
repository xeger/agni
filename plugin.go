package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type Series struct {
	Type, Name, Query string
}

type Plugin struct {
	Name   string
	Series []Series
}

// pluginSeries is the YAML representation of all the series defined "under"
// a given plugin. It is itself a map of series name to Prometheus query.
type pluginSeries map[string]string

// plugins is the YAML representation of several plugins. It is a map, where
// each key is a plugin name and the corresponding value is a map of series.
type plugins map[string]pluginSeries

// LoadPlugins loads and unmarshals plugin definitions from the plugins.yaml file.
// The YAML file can be located in the working directory or in /etc/agni.
//
// Example plugin definitions:
//   magic:
//     gauge-magic_smoke_level: rate(go_memstats_alloc_bytes[20s])
//     counter-magic_bunnies: go_memstats_frees_total
func LoadPlugins() ([]Plugin, error) {
	fn, err := findYaml("plugins.yaml")
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	return unmarshalPlugins(data)
}

func findYaml(fn string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	etcYaml := filepath.Join("/etc/agni", fn)
	pwdYaml := filepath.Join(pwd, fn)
	candidates := []string{etcYaml, pwdYaml}

	for _, c := range candidates {
		s, err := os.Stat(c)
		if err == nil && !s.IsDir() {
			return c, nil
		}
	}

	return "", fmt.Errorf("%s not found in any of: %v", fn, candidates)
}

func unmarshalPlugins(data []byte) ([]Plugin, error) {
	pm := plugins{}

	if err := yaml.Unmarshal(data, &pm); err != nil {
		return nil, err
	}

	out := make([]Plugin, len(pm))
	i := 0
	for pName, sm := range pm {
		series := make([]Series, len(sm))
		j := 0
		for k, v := range sm {
			bits := strings.SplitN(k, "-", 2)
			if len(bits) < 2 {
				return nil, fmt.Errorf("Malformed series name '%s'; must begin with 'gauge-' or 'counter-'", k)
			}
			switch bits[0] {
			case Counter, Gauge:
				series[j] = Series{Type: bits[0], Name: bits[1], Query: v}
			default:
				return nil, fmt.Errorf("Malformed series name '%s'; must begin with 'gauge-' or 'counter-'", k)
			}
			j++
		}
		out[i] = Plugin{Name: pName, Series: series}
		i++
	}

	return out, nil
}
