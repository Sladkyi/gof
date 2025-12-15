package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	cnt := 0
	for {
		resp, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
		if err != nil || resp.StatusCode != 200 {
			cnt++
			if cnt >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
			time.Sleep(5 * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		arr := strings.Split(strings.TrimSpace(string(body)), ",")

		if len(arr) != 7 {
			cnt++
			if cnt >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(5 * time.Second)
			continue
		}

		data := make([]float64, 7)
		badData := false
		for i, v := range arr {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				badData = true
				break
			}
			data[i] = f
		}

		if badData {
			cnt++
			if cnt >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(5 * time.Second)
			continue
		}
		cnt = 0

		if data[0] > 30 {
			fmt.Printf("Load average is too high: %v\n", data[0])
		}

		if data[1] > 0 {
			usage := (data[2] / data[1]) * 100
			if usage > 80 {
				fmt.Printf("Memory usage too high: %d%%\n", int(usage))
			}
		}

		if data[3] > 0 {
			usage := (data[4] / data[3]) * 100
			if usage > 90 {

				fmt.Printf("Free disk space is too low: %d Mb left\n", int((data[3]-data[4])/1024/1024))
			}
		}

		if data[5] > 0 {
			usage := (data[6] / data[5]) * 100
			if usage > 90 {
				fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", int((data[5]-data[6])/1000/1000))
			}
		}

		time.Sleep(5 * time.Second)
	}
}
