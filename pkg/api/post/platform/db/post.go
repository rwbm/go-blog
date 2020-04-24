package db

import (
	"encoding/base64"
	"fmt"
	"go-blog/pkg/util/model"
	"strings"
)

type filterResult struct {
	Query           string
	Args            []interface{}
	CategoriesFound bool
	TagsFound       bool
}

// GetPosts retusn a list of posts based on the indicated filters
func (p *PostDB) GetPosts(filters map[string]string, pageSize, page int) (posts []model.Post, pag model.Pagination, err error) {

	res := p.buildFilters(filters)
	sb := strings.Builder{}

	// pagination
	offset := 0
	if page > 1 {
		offset = ((page - 1) * pageSize)
	}

	// select
	sb.WriteString("SELECT p.id_post,p.date_created,p.date_updated,p.title,p.author,p.content FROM post p ")

	// joins
	if res.CategoriesFound {
		sb.WriteString("INNER JOIN post_category pc ON pc.id_post = p.id_post")
	}
	if res.TagsFound {
		sb.WriteString("INNER JOIN post_tag pt ON pt.id_post = p.id_post")
	}

	// where
	if len(res.Query) > 0 {
		sb.WriteString(" WHERE " + res.Query)
	}

	// order
	sb.WriteString(" ORDER BY p.id_post ASC")

	q := p.ds.Raw(sb.String(), res.Args...).Offset(offset).Limit(pageSize)
	rows, errQuery := q.Rows()
	defer rows.Close()

	if errQuery != nil {
		err = fmt.Errorf("error loading posts: %s", errQuery)
		return
	}

	// get results
	for rows.Next() {
		post := model.Post{}
		if errScan := p.ds.ScanRows(rows, &post); errScan != nil {
			err = fmt.Errorf("error loading posts: %s", errScan)
			return
		}

		// convert content to base64
		post.Content = p.encodeToBase64(post.Content)

		// load categories
		cats := []model.PostCategory{}
		if errGetCats := p.ds.Where("id_post = ?", post.ID).Find(&cats).Error; errGetCats != nil {
			err = fmt.Errorf("error loading post categories: %s", errGetCats)
			return
		}
		post.Categories = p.categoriesToString(cats)

		// load tags
		tags := []model.PostTag{}
		if errGetTags := p.ds.Where("id_post = ?", post.ID).Find(&tags).Error; errGetTags != nil {
			err = fmt.Errorf("error loading post tags: %s", errGetTags)
			return
		}
		post.Tags = p.tagsToString(tags)

		posts = append(posts, post)
	}

	pag.Page = page
	pag.PageSize = pageSize

	return
}

func (p *PostDB) buildFilters(filters map[string]string) (result filterResult) {

	filterArgs := []interface{}{}
	sbWhere := strings.Builder{}

	for k, v := range filters {

		key := strings.ToLower(k)
		switch key {

		case model.FilterAuthor:
			sbWhere.WriteString(" p.author=? AND ")
			filterArgs = append(filterArgs, v)

		case model.FilterID:
			sbWhere.WriteString(" p.id_post=? AND ")
			filterArgs = append(filterArgs, v)

		case model.FilterDateFrom:
			sbWhere.WriteString(" p.date_created >= ? AND ")
			filterArgs = append(filterArgs, v)

		case model.FilterDateTo:
			sbWhere.WriteString(" p.date_created <= ? AND ")
			filterArgs = append(filterArgs, v)

		case model.FilterCategories:
			if filterValues := p.parseMultipleValuesFilter(v); len(filterValues) > 0 {

				paramStr := strings.Repeat("?,", len(filterValues))
				sbWhere.WriteString("pc.name IN ( ")
				sbWhere.WriteString(paramStr[0 : len(paramStr)-1])
				sbWhere.WriteString(") AND ")

				// add parameter values
				for i := range filterValues {
					filterArgs = append(filterArgs, strings.Trim(filterValues[i], " "))
				}

				result.CategoriesFound = true
			}

		case model.FilterTags:
			if filterValues := p.parseMultipleValuesFilter(v); len(filterValues) > 0 {

				paramStr := strings.Repeat("?,", len(filterValues))
				sbWhere.WriteString("pt.name IN (")
				sbWhere.WriteString(paramStr[0 : len(paramStr)-1])
				sbWhere.WriteString(") AND ")

				// add parameter values
				for i := range filterValues {
					filterArgs = append(filterArgs, strings.Trim(filterValues[i], " "))
				}

				result.TagsFound = true
			}

		}

	}

	result.Query = strings.TrimRight(sbWhere.String(), "AND ")
	result.Args = filterArgs

	return
}

func (p *PostDB) parseMultipleValuesFilter(values string) []string {
	result := strings.Split(values, ",")
	return result
}

func (p *PostDB) encodeToBase64(data string) (encoded string) {
	encoded = base64.StdEncoding.EncodeToString([]byte(data))
	return
}

func (p *PostDB) categoriesToString(cats []model.PostCategory) (s string) {
	sa := []string{}
	for i := range cats {
		sa = append(sa, cats[i].Name)
	}

	return strings.Join(sa, ",")
}

func (p *PostDB) tagsToString(tags []model.PostTag) (s string) {
	sa := []string{}
	for i := range tags {
		sa = append(sa, tags[i].Name)
	}

	return strings.Join(sa, ",")
}
