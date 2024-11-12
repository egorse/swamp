package main

import (
	"embed"
	"flag"
	"log/slog"
	"os"

	"github.com/cloudcopper/swamp"
	"github.com/cloudcopper/swamp/infra"
	"github.com/cloudcopper/swamp/infra/config"
	"github.com/cloudcopper/swamp/lib"
)

const (
	retNoErrorCode      = 0
	retGenericErrorCode = 1
)

//go:embed templates/**
//go:embed static/**
//go:embed lake_repos.yml
var mainFS embed.FS

func main() {
	// Use config file name from env LAKE_REPO_CONFIG
	// or swamp_repos.yml
	// Note the config file might be embedded!!!
	config.ReposConfigFileName = "lake_repos.yml"
	config.ReposConfigFileName = lib.GetEnvDefault("LAKE_REPO_CONFIG", config.ReposConfigFileName)

	// First filesystem layer location (default is current working dir)
	config.TopRootFileSystemPath = lib.GetEnvDefault("LAKE_ROOT", config.TopRootFileSystemPath)
	// Second layer is this app embed fs - see mainFS
	// Last layer is the swamp embed fs - see appFS

	// Handle command line arguments
	flag.StringVar(&config.Listen, "listen", config.Listen, "web server listen address")
	flag.StringVar(&config.ReposConfigFileName, "repos", config.ReposConfigFileName, "repos config file name")
	flag.StringVar(&config.TopRootFileSystemPath, "root", config.TopRootFileSystemPath, "first layer of filesystem (optional)")
	flag.DurationVar(&config.TimerExpiredStart, "exp-start", config.TimerExpiredStart, "expired start timer")
	flag.DurationVar(&config.TimerExpiredInterval, "exp-interval", config.TimerExpiredInterval, "expired check interval")
	flag.IntVar(&config.TimerExpiredLimit, "exp-limit", config.TimerExpiredLimit, "expired check limit")
	flag.DurationVar(&config.TimerBrokenStart, "broken-start", config.TimerBrokenStart, "broken start timer")
	flag.DurationVar(&config.TimerBrokenInterval, "broken-interval", config.TimerBrokenInterval, "broken check interval")
	flag.IntVar(&config.TimerBrokenLimit, "broken-limit", config.TimerBrokenLimit, "broken check limit")
	flag.Parse()

	//
	// Create logger
	//
	log := slog.Default()
	log.Info("starting")

	topFS, err := infra.NewLayerFileSystem(config.TopRootFileSystemPath, mainFS)
	if err != nil {
		log.Error("unable create topFS", slog.Any("err", err))
	}

	code := retNoErrorCode
	err = swamp.App(log, topFS)
	if err != nil {
		code = retGenericErrorCode
		if i, ok := err.(lib.ErrorCode); ok {
			code = i.Code()
		}
		log.Error("exit", slog.Int("code", code), slog.Any("err", err))
	} else {
		log.Info("exit")
	}

	os.Exit(code)
}
