package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/monopole/snips/internal"
	"os"
	"time"
)

//go:embed README.md
var readMeMd string

func main() {
	if len(os.Args) < 3 || len(os.Args) > 6 {
		fmt.Print(readMeMd)
		os.Exit(1)
	}
	token := os.Args[1]
	user := os.Args[2]
	dayStart := time.Now().Round(24 * time.Hour)
	if len(os.Args) > 3 {
		dayStart = internal.ParseDate(os.Args[3])
	}
	dayCount := 6
	if len(os.Args) > 4 {
		dayCount = internal.ParseDayCount(os.Args[4])
	}
	ctx := context.Background()
	questioner{
		user:      user,
		dateStart: dayStart,
		dateEnd:   dayStart.AddDate(0, 0, dayCount),
		ctx:       ctx,
		client:    internal.MakeClient(ctx, token),
	}.doIt()
}
