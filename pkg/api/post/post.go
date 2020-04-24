package post

import (
	"go-blog/pkg/util/exception"
	"go-blog/pkg/util/model"
	"net/http"

	"github.com/labstack/echo"
)

// GetBlogPosts returns a list of blog posts, with optional filters
func (p *Post) GetBlogPosts(filters map[string]string, pageSize, page int) (posts []model.Post, pag model.Pagination, err error) {

	postList, pag, errGet := p.database.GetPosts(filters, pageSize, page)
	if errGet != nil {
		p.logger.Error("error loading posts from database", errGet, nil)

		err = echo.NewHTTPError(
			http.StatusInternalServerError,
			exception.GetErrorMap(exception.CodeInternalServerError, errGet.Error()))

		return
	}

	posts = postList
	return
}
