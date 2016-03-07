package util

type Info struct {
	Ip        string   `json:"ip"`
	Avg1min   float64  `json:"avg1min"`
	Avg5min   float64  `json:"avg5min"`
	Avg15min  float64  `json:"avg15min"`
	Buffers   uint64   `json:"buffers"`
	Cached    uint64   `json:"cached"`
	MemTotal  uint64   `json:"memtotal"`
	MemFree   uint64   `json:"memfree"`
	SwapTotal uint64   `json:"swaptotal"`
	SwapUsed  uint64   `json:"swapused"`
	SwapFree  uint64   `json:"swapfree"`
	DiskInfo  []DfInfo `json:"diskinfo"`
	Timestamp int64    `json:"timestamp"`
}

type DfInfo struct {
	Mount     string  `json:"mount"`
	InodeUsed float64 `json:"inodeuserd"`
	DiskUsed  float64 `json:"diskuserd"`
}
