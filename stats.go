package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func getStats() (CombinedStats, error) {
	repoStats, err := getIpfsRepoStat()
	if err != nil {
		return CombinedStats{}, err
	}

	bwStats, err := getBandwidthStats()
	if err != nil {
		return CombinedStats{}, err
	}

	radicleNodeStats, err := getRadicleStats()
	if err != nil {
		return CombinedStats{}, err
	}

	radicleRepoStats, err := getRadicleRepos()
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
		RadNodeInfo:    radicleNodeStats,
		RadNodeRepos:   radicleRepoStats,
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

func getRadicleStats() (RadNodeInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8888/api/v1/node", nil)
	if err != nil {
		return RadNodeInfo{}, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return RadNodeInfo{}, fmt.Errorf("error fetching Rad Node Stats: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RadNodeInfo{}, fmt.Errorf("error reading response body: %v", err)
	}

	var stats RadNodeInfo
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return RadNodeInfo{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	simplifiedStats := RadNodeInfo{
		ID:    stats.ID,
		Agent: stats.Agent,
		State: stats.State,
		Config: struct {
			SeedingPolicy struct {
				Default string `json:"default"`
			} `json:"seedingPolicy"`
		}{
			SeedingPolicy: stats.Config.SeedingPolicy,
		},
	}

	return simplifiedStats, nil
}

func getRadicleRepos() (RadNodeRepos, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8888/api/v1/stats", nil)
	if err != nil {
		return RadNodeRepos{}, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return RadNodeRepos{}, fmt.Errorf("error fetching Rad Node Stats: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RadNodeRepos{}, fmt.Errorf("error reading response body: %v", err)
	}

	var stats RadNodeRepos
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return RadNodeRepos{}, fmt.Errorf("error unmarshaling JSON: %v", err)
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
