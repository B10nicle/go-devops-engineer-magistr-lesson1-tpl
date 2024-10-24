package main

import (
	"fmt"
	"github.com/B10nicle/go-devops-engineer-magistr-lesson1-tpl/config"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	errorCount := 0

	for {
		statsData, err := getStatistics()
		if err != nil {
			fmt.Println("Error fetching server statistics:", err)
			errorCount++
		} else {
			stats, err := parseStatistics(statsData)

			if err != nil {
				fmt.Println("Error parsing server statistics:", err)
				errorCount++
			} else {
				checkThresholds(stats)
				errorCount = 0
			}
		}

		if errorCount >= 3 {
			fmt.Println("Unable to fetch server statistics")
		}
	}
}

func getStatistics() (string, error) {
	response, err := http.Get(config.ServerURL)

	if err != nil {
		return "", err
	}

	defer func() {
		if bodyError := response.Body.Close(); bodyError != nil {
			fmt.Printf("Error closing response body: %v", bodyError)
			fmt.Println()
		}
	}()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func parseStatistics(data string) ([]float64, error) {
	parts := strings.Split(data, ",")
	stats := make([]float64, len(parts))

	for i, part := range parts {
		part = strings.TrimSpace(part)
		value, err := strconv.ParseFloat(part, 64)

		if err != nil {
			return nil, err
		}

		stats[i] = value
	}

	return stats, nil
}

func checkThresholds(stats []float64) {
	loadAverage := stats[0]
	totalMemory := stats[1]
	usedMemory := stats[2]
	totalDisk := stats[3]
	usedDisk := stats[4]
	totalNet := stats[5]
	usedNet := stats[6]

	currentMemoryUsage := usedMemory / totalMemory * 100
	currentDiskUsage := usedDisk / totalDisk * 100
	currentNetworkBandwidthUsage := usedNet / totalNet * 100
	freeDiskMb := (totalDisk - usedDisk) / 1048576
	freeNet := (totalNet - usedNet) / 1000000

	if loadAverage > config.LoadAverageThreshold {
		fmt.Printf("Load Average is too high: %d\n", int(loadAverage))
	}

	if currentMemoryUsage > config.MemoryUsageLimit {
		fmt.Printf("Memory usage too high: %d%%\n", int(currentMemoryUsage))
	}

	if currentDiskUsage > config.DiskUsageLimit {
		fmt.Printf("Free disk space is too low: %d Mb left\n", int(freeDiskMb))
	}

	if currentNetworkBandwidthUsage > config.NetworkBandwidthLimit {
		fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", int(freeNet))
	}
}
