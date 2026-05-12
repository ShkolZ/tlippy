package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/ShkolZ/tlippy/internal/config"
	"github.com/ShkolZ/tlippy/internal/download"
)

type CLI struct {
}

func NewCLI() *CLI {
	return &CLI{}
}

func (c *CLI) Run() error {
	initCmds()
	return nil
}

func initCmds() {
	bulkCmd := flag.NewFlagSet("bulk", flag.ContinueOnError)
	categoryPtr := bulkCmd.String("c", "", "Category to download clips from")
	streamerPtr := bulkCmd.String("s", "", "Streamer to download clips from")
	timeRangePtr := bulkCmd.String("t", "24h", "Time range to download clips from")
	countPtr := bulkCmd.String("a", "20", "Number of clips to download")
	pathPtr := bulkCmd.String("o", "./clips", "Path to save clips to")

	singleCmd := flag.NewFlagSet("single", flag.ContinueOnError)
	clipIdPtr := singleCmd.String("id", "", "Clip ID to download")
	spathPtr := singleCmd.String("o", "./clips", "Path to save clips to")

	if len(os.Args) < 2 {
		fmt.Println("Please provide a command.")
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "bulk":
		bulkCmd.Parse(os.Args[2:])
		var input config.UserInput
		if *streamerPtr == "" {
			input = config.UserInput{
				Mode:         config.ModeBulk,
				BulkType:     config.BulkByCategory,
				QueryName:    *categoryPtr,
				TimeRange:    config.TimeRange(*timeRangePtr),
				ClipCount:    *countPtr,
				DownloadPath: *pathPtr,
			}
		} else {
			input = config.UserInput{
				Mode:         config.ModeBulk,
				BulkType:     config.BulkByStreamer,
				QueryName:    *streamerPtr,
				TimeRange:    config.TimeRange(*timeRangePtr),
				ClipCount:    *countPtr,
				DownloadPath: *pathPtr,
			}

		}
		ch := download.StartDownloadChan(&input)
		for p := range ch {
			fmt.Printf("Progress: %v/%v\n", p.Current, p.Total)
		}
		fmt.Println("Done!")
	case "single":
		singleCmd.Parse(os.Args[2:])
		input := config.UserInput{
			Mode:         config.ModeSingle,
			ClipID:       *clipIdPtr,
			DownloadPath: *spathPtr,
		}
		ch := download.StartDownloadChan(&input)
		<-ch
		fmt.Println("Done!")
	default:
		fmt.Println("Unknown command.")
		flag.Usage()
		os.Exit(1)
	}

}
