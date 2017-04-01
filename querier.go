package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/common/model"
)

type Querier interface {
	Query(t time.Time, query string) (float64, error)
}

type prometheusQuerier struct {
	ctx  context.Context
	qapi prometheus.QueryAPI
}

func (pq *prometheusQuerier) Query(t time.Time, query string) (float64, error) {
	res, err := pq.qapi.Query(pq.ctx, query, t)
	if err != nil {
		return 0.0, err
	}
	switch res.Type() {
	case model.ValScalar:
		scalar := res.(*model.Scalar)
		return float64(scalar.Value), nil
	case model.ValVector:
		vector := res.(model.Vector)
		s := 0.0
		for _, vi := range vector {
			s += float64(vi.Value)
		}
		return s, nil
	default:
		return 0.0, fmt.Errorf("unknown prometheus model.ValueType: %s", res.Type().String())
	}
}

func NewQuerier(url string) (Querier, error) {
	config := prometheus.Config{
		Address: url,
	}
	client, err := prometheus.New(config)
	if err != nil {
		return nil, err
	}
	return &prometheusQuerier{ctx: context.Background(), qapi: prometheus.NewQueryAPI(client)}, nil
}

// LoadQuerier loads and unmarshals a querier definition from the querier.yaml file.
// The YAML file can be located in the working directory or in /etc/agni. If the
// file is not present then LoadQuerier builds a new Prometheus querier with default
// configuration.
//
// Example plugin definitions:
//   magic:
//     gauge-magic_smoke_level: rate(go_memstats_alloc_bytes[20s])
//     counter-magic_bunnies: go_memstats_frees_total
func LoadQuerier() (Querier, error) {
	url := "http://localhost:9090"

	if fn, err := findYaml("querier.yaml"); err == nil {
		data, err := ioutil.ReadFile(fn)
		if err == nil {
			config := map[string]string{}
			if err := yaml.Unmarshal(data, &config); err == nil && config["url"] != "" {
				url = config["url"]
			}
		}
	}

	return NewQuerier(url)
}
