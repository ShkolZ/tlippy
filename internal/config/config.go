package config

import "strconv"

type DownloadMode int

const (
	ModeBulk   DownloadMode = iota
	ModeSingle DownloadMode = iota
)

type BulkType int

const (
	BulkByCategory BulkType = iota
	BulkByStreamer  BulkType = iota
)

type TimeRange string

const (
	TimeRange24h TimeRange = "24h"
	TimeRange7d  TimeRange = "7d"
	TimeRangeAll TimeRange = "all"
)

// UserInput holds everything collected from the TUI.
type UserInput struct {
	Mode DownloadMode

	// Bulk fields
	BulkType  BulkType
	QueryName string // category or streamer name
	TimeRange TimeRange
	ClipCount string // number of clips as string

	// Single fields
	ClipID string

	// Shared
	DownloadPath string
}

type Config struct {
	DownloadPath string
	ClipsAmount  int
}

func GetConfig() *Config {
	return nil
}

func SaveConfig() {

}

func SetConfig(path string, amount string) (*Config, error) {
	intAmount, err := strconv.Atoi(amount)
	if err != nil {
		return nil, err
	}
	cfg := Config{
		DownloadPath: path,
		ClipsAmount:  intAmount,
	}
	return &cfg, nil
}
