package api

import (
	post "go-blog/pkg/api/post"
	pt "go-blog/pkg/api/post/transport"
	"go-blog/pkg/util/config"
	"go-blog/pkg/util/log"
	"go-blog/pkg/util/model"
	"go-blog/pkg/util/server"
	"go-blog/pkg/util/template"
	"go-blog/pkg/util/watcher"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rwbm/go-tools/files"
)

// Internal consts
const (
	TemplatesExtension = ".tpl"
	DatabaseDriver     = "sqlite3"
)

// Start starts the API service
func Start(cfg *config.Configuration) (err error) {

	logger := log.New() // default logger

	// check if databse exists, so we can recreate it
	recreateDatabase := false
	if !files.Exists(cfg.Database.Filename) {
		recreateDatabase = true
	}

	// create DB connection
	ds, errDB := gorm.Open(DatabaseDriver, cfg.Database.Filename)
	if errDB != nil {
		return errDB
	}

	// create database structure
	if recreateDatabase {
		logger.Info("database NOT found; recreating from scratch", map[string]interface{}{"dbfile": cfg.Database.Filename})
		if errRecreate := ds.AutoMigrate(
			&model.Post{},
			&model.PostCategory{},
			&model.PostTag{}).Error; errRecreate != nil {
			return errRecreate
		}
	}

	// watcher for the templates folder
	templateProcessor := template.NewProcessor(
		ds,
		logger,
		cfg.Template.ProcessedOK,    // location where templates are moved if processed OK
		cfg.Template.ProcessedError) // location where templates are moved if processed with ERROR

	fileWatcher := watcher.NewWatcher(
		cfg.Template.Base,  // location to look for templates
		TemplatesExtension, // templates extension to look for
		time.Duration(cfg.Template.CheckCycle)*time.Second, // interval to check for new templates
		logger,
		templateProcessor.ProcessTemplate)

	go fileWatcher.Start()

	// +++++++++++ SERVICES ++++++++++++

	e := server.New()
	pt.NewHTTP(post.Initialize(ds, nil, logger, cfg.Server.DryRun), e)

	// +++++++++++++++++++++++++++++++++

	// start HTTP server
	server.Start(e,
		&server.Config{
			ServiceName:         cfg.Server.Name,
			Port:                cfg.Server.Port,
			ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
			WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		},
		logger)

	return
}
