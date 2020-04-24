package model

import (
	"time"
)

// Filter names
const (
	FilterID         = "id"
	FilterDateFrom   = "date-from"
	FilterDateTo     = "date-to"
	FilterAuthor     = "author"
	FilterCategories = "categories"
	FilterTags       = "tags"
)

// Post represents a blog post
type Post struct {
	ID               int       `gorm:"column:id_post;primary_key;AUTO_INCREMENT" json:"id_post"`
	DateCreated      time.Time `gorm:"column:date_created;NOT NULL" json:"date_created"`
	DateUpdated      time.Time `gorm:"column:date_updated;NOT NULL" json:"date_updated"`
	Title            string    `gorm:"column:title;NOT NULL;type:varchar(128);NOT NULL" json:"title"`
	Author           string    `gorm:"column:author;NOT NULL;type:varchar(128);NOT NULL" json:"author"`
	Content          string    `gorm:"column:content;NOT NULL;type:text;NOT NULL" json:"content"`
	Categories       string    `gorm:"-" json:"categories"`
	Tags             string    `gorm:"-" json:"tags"`
	OriginalFileName string    `gorm:"column:original_filename;type:varchar(128);NOT NULL" json:"-"`
}

// TableName returns the table name for the model
func (Post) TableName() string {
	return "post"
}

// PostCategory represents a post category
type PostCategory struct {
	IDPost int    `gorm:"column:id_post;NOT NULL;type:integer" json:"id_post"`
	Name   string `gorm:"column:name;NOT NULL;type:varchar(128);NOT NULL" json:"name"`
}

// TableName returns the table name for the model
func (PostCategory) TableName() string {
	return "post_category"
}

// PostTag represents a post tag
type PostTag struct {
	IDPost int    `gorm:"column:id_post;NOT NULL;type:integer" json:"id_post"`
	Name   string `gorm:"column:name;NOT NULL;type:varchar(128);NOT NULL" json:"name"`
}

// TableName returns the table name for the model
func (PostTag) TableName() string {
	return "post_tag"
}
