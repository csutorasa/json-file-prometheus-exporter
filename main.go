package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var logger = log.Default()
var metricsRegex = regexp.MustCompile(`[^0-9a-zA-Z_]`)

func main() {
	file, labels, separator, port, path, metricName, err := parseArgs()
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		logger.Fatalf("%v", err)
	}
	go func() {
		reader := NewReader(file, separator)
		reporter := NewReporter(metricName)
		reporter.Init(labels)
		for {
			data, err := reader.Read()
			if err != nil {
				logger.Printf("Reader has exited")
				return
			}
			if data == nil {
				time.Sleep(1 * time.Second)
				continue
			}
			reporter.Report(data)
		}
	}()
	http.Handle(path, promhttp.Handler())
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		logger.Fatalf("%v", err)
	}
}

func parseArgs() (*os.File, []string, byte, int, string, string, error) {
	metricName := flag.String("metric-name", "", "Name of the metric")
	labels := flag.String("labels", "", "Labels of the metric, comma separated list")
	path := flag.String("path", "/metrics", "Url path to expose metrics to")
	port := flag.Int("port", 8080, "Port to expose metrics to")
	separator := flag.String("separator", "\n", "Line separator string")
	stdin := flag.Bool("stdin", false, "If the data should be read from stdin")
	flag.Parse()
	filepath := flag.Arg(0)
	var labelsArray []string = nil
	if *metricName == "" {
		return nil, labelsArray, 0, *port, *path, *metricName, fmt.Errorf("metric-name is required")
	}
	*metricName = metricsRegex.ReplaceAllString(*metricName, "_")
	logger.Printf("Using metric name %s", *metricName)
	if *labels == "" {
		return nil, labelsArray, 0, *port, *path, *metricName, fmt.Errorf("separator is not valid")
	}
	labelsArray = strings.Split(*labels, ",")
	if len(*separator) != 1 {
		return nil, labelsArray, 0, *port, *path, *metricName, fmt.Errorf("separator is not valid")
	}
	if len(*separator) != 1 {
		return nil, labelsArray, 0, *port, *path, *metricName, fmt.Errorf("separator is not valid")
	}
	sep := (*separator)[0]
	if filepath != "" {
		if *stdin {
			return nil, labelsArray, sep, *port, *path, *metricName, fmt.Errorf("filepath and stdin are both present")
		}
		logger.Printf("Opening file %s", filepath)
		file, err := os.OpenFile(filepath, os.O_RDONLY, 0)
		return file, labelsArray, sep, *port, *path, *metricName, err
	}
	if *stdin {
		logger.Print("Using stdin")
		return os.Stdin, labelsArray, sep, *port, *path, *metricName, nil
	}
	return nil, labelsArray, sep, *port, *path, *metricName, fmt.Errorf("filepath and stdin are both missing")
}
