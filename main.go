package main

import (
	"embed"

	"github.com/mogenius/punq/cmd"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/utils"
)

//go:embed ui/dist/*
var htmlDirFs embed.FS

//go:embed config/local.yaml
var localConfig string

//go:embed config/dev.yaml
var devConfig string

//go:embed config/prod.yaml
var prodConfig string

func main() {
	utils.DefaultConfigLocalFile = localConfig
	utils.DefaultConfigClusterFileDev = devConfig
	utils.DefaultConfigClusterFileProd = prodConfig

	services.HtmlDirFs = htmlDirFs

	logger.Init()
	cmd.Execute()
}
