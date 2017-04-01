package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err.Error())
		os.Exit(1)
	}
}

func poll(q Querier, hostname, plugin string, series Series, interval int) {
	period := time.Duration(interval) * time.Second

	for {
		t0 := time.Now()

		value, err := q.Query(t0, series.Query)
		if err == nil {
			fmt.Printf("PUTVAL %s/%s/%s-%s interval=%d %d:%g\n", hostname, plugin, series.Type, series.Name, interval, t0.Unix(), value)
		} else {
			fmt.Fprintf(os.Stderr, "%s/%s: %s\n", plugin, series.Name, err.Error())
		}
		dt := time.Now().Sub(t0)
		if sleep := period - dt; sleep > 0 {
			time.Sleep(sleep)
		}
	}
}

func main() {
	plugins, err := LoadPlugins()
	fatal(err)

	hostname := os.Getenv("COLLECTD_HOSTNAME")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	if hostname == "" {
		hostname = "localhost"
	}

	interval := 20
	if env := os.Getenv("COLLECTD_INTERVAL"); env != "" {
		envi, erri := strconv.Atoi(env)
		fatal(erri)
		interval = envi
	}
	sleep := time.Duration(interval) * time.Second

	q, err := LoadQuerier()
	fatal(err)

	for _, p := range plugins {
		for _, s := range p.Series {
			go poll(q, hostname, p.Name, s, interval)
		}
	}

	for {
		time.Sleep(sleep)
	}
}
