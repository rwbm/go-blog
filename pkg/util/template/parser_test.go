package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTemplate(t *testing.T) {

	htmlExample := `
<head>
	<meta name="title" content="My First Blog Post"/>
	<meta name="author" content="John Doe"/>
	<meta name="post-date" content="2020-04-15 12:09:57"/>
	<meta name="edit-date" content="2020-04-15 12:19:05"/>
	<meta name="categories" content="Go Programming"/>
	<meta name="tags" content="go, programming, web"/>
<head>
<body>
	<h1>My First Blog Post</h1>
	<p>This is my first Blog post, just to try if templates are working OK.</p>
</body>
	`

	post, err := ParseTemplate(htmlExample)
	assert.NoError(t, err)
	assert.Equal(t, "My First Blog Post", post.Title)
	assert.Equal(t, "John Doe", post.Author)
	assert.Equal(t, "Go Programming", post.Categories)
	assert.Equal(t, "go, programming, web", post.Tags)
	assert.NotEmpty(t, post.Content)

	t.Logf("%+v", post)

}
