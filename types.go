package main

type RadNodeInfo struct {
	ID     string `json:"id"`
	Agent  string `json:"agent"`
	State  string `json:"state"`
	Config struct {
		SeedingPolicy struct {
			Default string `json:"default"`
		} `json:"seedingPolicy"`
	} `json:"config"`
}

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
	RadNodeInfo
}
