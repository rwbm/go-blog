package main

import (
	"flag"
	"go-blog/pkg/api"
	"go-blog/pkg/util/config"
	"path"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rwbm/go-tools/files"
)

func main() {

	defaultConfigFile := path.Join(files.GetAppPath(), "config.yml")

	cfgPath := flag.String("config", defaultConfigFile, "path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	// start api server
	checkErr(api.Start(cfg))
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
