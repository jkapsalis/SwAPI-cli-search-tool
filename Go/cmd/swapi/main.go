package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/jkapsalis/SwAPI-cli-search-tool/Go/internal/api"
	"github.com/jkapsalis/SwAPI-cli-search-tool/Go/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	code := cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr, api.NewClient(api.DefaultBaseURL))
	stop()
	os.Exit(code)
}
