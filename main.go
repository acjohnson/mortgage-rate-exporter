package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

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

func fetchAndParseMortgageData() {
	resp, err := http.Get("https://fred.stlouisfed.org/data/MORTGAGE30US.txt")
	if err != nil {
		log.Printf("Failed to fetch data: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch data: HTTP status %d", resp.StatusCode)
		return
	}

	parseMortgageData(resp.Body)
}

func parseMortgageData(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
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
				value, err := strconv.ParseFloat(parts[1], 64)
				if err != nil {
					log.Printf("Failed to parse value: %v", err)
					continue
				}
				mortgage30USMetric.Set(value)
				log.Printf("Value: %.2f", value)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error while reading data: %v", err)
	}
}

func main() {
	fetchAndParseMortgageData()

	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			fetchAndParseMortgageData()
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Listening on :8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
