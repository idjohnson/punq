package main

import (
	"embed"

	"github.com/mogenius/punq/cmd"
	"github.com/mogenius/punq/logger"
)

// //go:embed ui/dist/punq/*
var htmlDirFs embed.FS

//go:embed config/dev.yaml
var defaultEnvFile string

func main() {
	// utils.DefaultEnvFile = defaultEnvFile
	// api.HtmlDirFs = htmlDirFs
	logger.Init()
	cmd.Execute()
}
