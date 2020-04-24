package db

import "github.com/jinzhu/gorm"

// NewPostDB returns a new posts database instance
func NewPostDB(ds *gorm.DB) (c *PostDB) {
	c = new(PostDB)
	c.ds = ds
	return
}

// PostDB contains the services to handle posts
type PostDB struct {
	ds *gorm.DB
}
