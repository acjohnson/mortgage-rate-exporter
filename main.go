package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mortgage30USMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mortgage30us",
		Help: "30-Year Fixed Rate Mortgage Average in the United States",
	})
	whitespaceRegex = regexp.MustCompile(`\s+`)
)

func parseMortgageData(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inDataSection := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "DATE        VALUE" {
			inDataSection = true
			continue
		}

		if inDataSection {
			parts := whitespaceRegex.Split(line, -1)
			if len(parts) == 2 {
				date := parts[0]
				value, err := strconv.ParseFloat(parts[1], 64)
				if err != nil {
					log.Printf("Failed to parse value: %v", err)
					continue
				}
				mortgage30USMetric.Set(value)
				log.Printf("Date: %s, Value: %.2f", date, value)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %v", err)
	}
}

func main() {
	filePath := "sample_mortgage_data.txt"
	parseMortgageData(filePath)

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Listening on :8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
