package main

import (
	"context"
	"fmt"
	"os"
	"snips/internal"
)

func main() {
	if !(len(os.Args) == 4 || len(os.Args) == 5) {
		fmt.Print(`usage:
  snips {user} {githubAuthToken} {dateStart} [{dayCount}] 
e.g.
  go run . monopole deadbeef0000deadbeef 2020-04-06 
`)
		os.Exit(1)
	}
	user := os.Args[1]
	token := os.Args[2]
	dayStart := internal.ParseDate(os.Args[3])
	dayCount := 6
	if len(os.Args) == 5 {
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
