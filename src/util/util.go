package util

type Info struct {
	Ip       string  `json:"ip"`
	Avg1min  float64 `json:"avg1min"`
	Avg5min  float64 `json:"avg5min"`
	Avg15min float64 `json:"avg15min"`
}
