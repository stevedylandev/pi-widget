package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SystemStats struct {
	CPUUsage    float64 `json:"cpuUsage"`
	MemoryUsed  uint64  `json:"memoryUsed"`
	MemoryTotal uint64  `json:"memoryTotal"`
	OS          string  `json:"os"`
	Kernel      string  `json:"kernel"`
	Uptime      string  `json:"uptime"`
	CPUModel    string  `json:"cpuModel"`
}

type IPFSRepoStats struct {
	RepoSize   int64  `json:"RepoSize"`
	StorageMax int64  `json:"StorageMax"`
	NumObjects int    `json:"NumObjects"`
	RepoPath   string `json:"RepoPath"`
	Version    string `json:"Version"`
}

type BandwidthStats struct {
	RateIn   float64 `json:"RateIn"`
	RateOut  float64 `json:"RateOut"`
	TotalIn  int64   `json:"TotalIn"`
	TotalOut int64   `json:"TotalOut"`
}

type CombinedStats struct {
	IPFSRepoStats
	BandwidthStats
	SystemStats
}

func main() {
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/events", handleSSE)
	fmt.Println("Server is running on http://localhost:4321")
	log.Fatal(http.ListenAndServe(":4321", nil))
}

//go:embed index.html
var indexHTML string

func serveHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(indexHTML))
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	log.Println("SSE connection established")

	for {
		combinedStats, err := getStats()
		if err != nil {
			log.Printf("Error getting IPFS Stats: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		data, err := json.Marshal(combinedStats)
		if err != nil {
			log.Printf("Error marshaling IPFS stats: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		_, err = fmt.Fprintf(w, "data: %s\n\n", data)
		if err != nil {
			log.Printf("Error writing to response: %v", err)
			return
		}
		w.(http.Flusher).Flush()

		time.Sleep(1 * time.Second)
	}
}

func getStats() (CombinedStats, error) {
	repoStats, err := getIpfsRepoStat()
	if err != nil {
		return CombinedStats{}, err
	}

	bwStats, err := getBandwidthStats()
	if err != nil {
		return CombinedStats{}, err
	}

	sysStats, err := getSystemStats()
	if err != nil {
		return CombinedStats{}, err
	}

	return CombinedStats{
		IPFSRepoStats:  repoStats,
		BandwidthStats: bwStats,
		SystemStats:    sysStats,
	}, nil
}

func getIpfsRepoStat() (IPFSRepoStats, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:5001/api/v0/repo/stat", nil)
	if err != nil {
		return IPFSRepoStats{}, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return IPFSRepoStats{}, fmt.Errorf("error fetching IPFS Repo Stats: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return IPFSRepoStats{}, fmt.Errorf("error reading response body: %v", err)
	}

	var stats IPFSRepoStats
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return IPFSRepoStats{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return stats, nil
}

func getBandwidthStats() (BandwidthStats, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:5001/api/v0/stats/bw", nil)
	if err != nil {
		return BandwidthStats{}, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return BandwidthStats{}, fmt.Errorf("error fetching Bandwidth Stats: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return BandwidthStats{}, fmt.Errorf("error reading response body: %v", err)
	}

	var stats BandwidthStats
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return BandwidthStats{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return stats, nil
}

func getSystemStats() (SystemStats, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return SystemStats{}, err
	}

	c, err := cpu.Percent(0, false)
	if err != nil {
		return SystemStats{}, err
	}

	hostInfo, err := host.Info()
	if err != nil {
		return SystemStats{}, err
	}

	uptime := time.Duration(hostInfo.Uptime) * time.Second

	return SystemStats{
		CPUUsage:    c[0],
		MemoryUsed:  v.Used,
		MemoryTotal: v.Total,
		OS:          "Debian GNU/Linux 12 (bookworm) aarch64",
		Kernel:      hostInfo.KernelVersion,
		Uptime:      FormatUptime(uptime),
		CPUModel:    "BCM2835 (4) @ 1.800GHz",
	}, nil
}
