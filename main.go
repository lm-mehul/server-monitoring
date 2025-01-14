package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type Metrics struct {
	Timestamp time.Time              `json:"timestamp"`
	CPUUsage  float64                `json:"cpu_usage"`
	Memory    *mem.VirtualMemoryStat `json:"memory"`
	Disk      *disk.IOCountersStat   `json:"disk"`
	Network   []net.IOCountersStat   `json:"network"`
}

var (
	metricsFile = "metrics.json"
	mutex       sync.Mutex
)

func startMonitoring(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var metricsData []Metrics

	for i := 0; i < 60; i++ {
		cpuUsage, _ := cpu.Percent(0, false)
		memory, _ := mem.VirtualMemory()
		diskStats, _ := disk.IOCounters()
		networkStats, _ := net.IOCounters(true)

		metrics := Metrics{
			Timestamp: time.Now(),
			CPUUsage:  cpuUsage[0],
			Memory:    memory,
			Disk:      getFirstDiskStat(diskStats),
			Network:   networkStats,
		}
		metricsData = append(metricsData, metrics)

		time.Sleep(250 * time.Millisecond)
	}

	mutex.Lock()
	defer mutex.Unlock()
	file, err := json.MarshalIndent(metricsData, "", "  ")
	if err != nil {
		http.Error(w, "Error saving metrics", http.StatusInternalServerError)
		return
	}
	_ = ioutil.WriteFile(metricsFile, file, 0644)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Monitoring complete and data saved to file."))
}

func cpuMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	file, err := ioutil.ReadFile(metricsFile)
	if err != nil {
		http.Error(w, "Metrics file not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(file)
}

func getFirstDiskStat(stats map[string]disk.IOCountersStat) *disk.IOCountersStat {
	for _, stat := range stats {
		return &stat
	}
	return nil
}

func main() {
	http.HandleFunc("/start-monitoring", startMonitoring)
	http.HandleFunc("/cpu-metrics", cpuMetrics)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
