package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	url           = "http://srv.msk01.gigacorp.local/_stats"
	checkInterval = 5 * time.Second
)

func main() {
	consecutiveErrors := 0

	for {
		resp, err := http.Get(url)

		if err != nil || resp.StatusCode != 200 {
			consecutiveErrors++
			if consecutiveErrors >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
			time.Sleep(checkInterval)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			consecutiveErrors++
			if consecutiveErrors >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(checkInterval)
			continue
		}

		dataStr := strings.TrimSpace(string(body))
		parts := strings.Split(dataStr, ",")

		if len(parts) != 7 {
			consecutiveErrors++
			if consecutiveErrors >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(checkInterval)
			continue
		}

		values := make([]float64, 7)
		parseError := false
		for i, v := range parts {
			val, err := strconv.ParseFloat(v, 64)
			if err != nil {
				parseError = true
				break
			}
			values[i] = val
		}

		if parseError {
			consecutiveErrors++
			if consecutiveErrors >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(checkInterval)
			continue
		}

		consecutiveErrors = 0

		loadAvg := values[0]
		memTotal := values[1]
		memUsed := values[2]
		diskTotal := values[3]
		diskUsed := values[4]
		netTotal := values[5]
		netUsed := values[6]

		if loadAvg > 30 {
			fmt.Printf("Load Average is too high: %v\n", loadAvg)
		}

		if memTotal > 0 {
			memUsagePercent := (memUsed / memTotal) * 100
			if memUsagePercent > 80 {
				fmt.Printf("Memory usage too high: %d%%\n", int(memUsagePercent))
			}
		}

		if diskTotal > 0 {
			diskUsagePercent := (diskUsed / diskTotal) * 100
			if diskUsagePercent > 90 {
				freeBytes := diskTotal - diskUsed
				freeMb := int(freeBytes / 1024 / 1024)
				fmt.Printf("Free disk space is too low: %d Mb left\n", freeMb)
			}
		}

		if netTotal > 0 {
			netUsagePercent := (netUsed / netTotal) * 100
			if netUsagePercent > 90 {
				freeBytes := netTotal - netUsed
				freeBits := freeBytes * 8
				freeMbit := int(freeBits / 1000 / 1000)
				fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeMbit)
			}
		}

		time.Sleep(checkInterval)
	}
}
