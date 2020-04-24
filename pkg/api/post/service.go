package post

import (
	"go-blog/pkg/api/post/platform/db"
	"go-blog/pkg/util/log"
	"go-blog/pkg/util/model"

	"github.com/jinzhu/gorm"
)

// Service holds the functions delcared in the service interface
type Service interface {
	GetBlogPosts(filters map[string]string, pageSize, page int) (posts []model.Post, pag model.Pagination, err error)
}

// DB holds the functions for database access
type DB interface {
	GetPosts(filters map[string]string, pageSize, page int) (posts []model.Post, pag model.Pagination, err error)
}

// Post defines the module for posts related operations
type Post struct {
	database   DB
	logger     *log.Log
	dryRunMode bool
}

// creates new post service
func new(database DB, l *log.Log, dryRunMode bool) *Post {
	return &Post{
		database:   database,
		logger:     l,
		dryRunMode: dryRunMode,
	}
}

// Initialize initializes Post application service
func Initialize(ds *gorm.DB, dbService DB, l *log.Log, dryRunMode bool) *Post {
	if dbService == nil {
		dbService = db.NewPostDB(ds)
	}
	return new(dbService, l, dryRunMode)
}
