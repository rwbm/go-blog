package transport

import (
	"fmt"
	post "go-blog/pkg/api/post"
	"go-blog/pkg/util/exception"
	"go-blog/pkg/util/model"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Default values
const (
	DefaultPageSize = 25
)

// HTTP represents auth http service
type HTTP struct {
	svc              post.Service
	jwtSigningKey    string
	jwtSigningMethod *jwt.SigningMethodHMAC
}

// NewHTTP creates new http service to handle request to /posts
func NewHTTP(svc post.Service, e *echo.Echo) (h HTTP) {
	h = HTTP{
		svc: svc,
	}

	e.GET("/posts", h.getPostsHandler)

	return
}

//
// --- GET BLOG POSTS ---
//
func (h *HTTP) getPostsHandler(c echo.Context) error {

	// get pagination data
	var page, pageSize int
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("page-size")

	if pageStr == "" {
		pageStr = "1"
	}
	if pageSizeStr == "" {
		pageSize = DefaultPageSize
	}

	page, errConv := strconv.Atoi(pageStr)
	if errConv != nil {
		return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(exception.CodeInvalidPage, ""))
	}

	if pageSize == 0 {
		pageSize, errConv = strconv.Atoi(pageSizeStr)
		if errConv != nil {
			return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(exception.CodeInvalidPageSize, ""))
		}
	}

	filters, errFilters := h.buildFilterMap(c)
	if errFilters != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			exception.GetErrorMap(exception.CodeBadRequest, errFilters.Error()))
	}

	// get posts
	posts, pageInfo, errPosts := h.svc.GetBlogPosts(filters, pageSize, page)
	if errPosts != nil {
		return errPosts
	}

	// if we got no records, return an empty array
	if posts == nil {
		posts = []model.Post{}
	}

	payload := make(map[string]interface{})
	payload["posts"] = posts
	payload["pagination"] = pageInfo

	return c.JSON(http.StatusOK, payload)
}

func (h *HTTP) buildFilterMap(c echo.Context) (filters map[string]string, err error) {
	filters = make(map[string]string)

	for k := range c.QueryParams() {

		switch k {

		case "author", "tags", "categories", "id_post":
			filters[k] = c.QueryParam(k)

		case "date-from", "date-to":
			v := c.QueryParam(k)
			if len(v) > 10 {
				v = v[0:10]
			}

			_, errParse := time.Parse("2006-01-02", v)
			if errParse != nil {
				err = fmt.Errorf("error parsing date value from '%s'", k)
				return
			}

			// add extra time section to fitler in the database
			timePart := "00:00:00"
			if k == "date-to" {
				timePart = "23:95:59"
			}

			v += " " + timePart
			filters[k] = v
		}

	}

	return
}
