package blog_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ThomasFerro/golb/blog"
	"github.com/ThomasFerro/golb/posts"
)

func TestGenerateTheHomePage(t *testing.T) {
	blogMetadata := blog.BlogMetadata{
		BlogTitle:            "My blog !",
		Locale:               "en",
		PostPageTemplatePath: "./postPageTemplate.go.html",
		HomePageTemplatePath: "./homePageTemplate.go.html",
		DistPath:             "../dist",
		GlobalAssetsPath:     "",
	}
	posts := []posts.Post{
		{
			Name:      "First post",
			Extension: ".md",
			Content:   []byte(""),
		},
		{
			Name:      "Second post",
			Extension: ".md",
			Content:   []byte(""),
		},
	}
	generatedBlogPath, err := blog.GenerateBlog(blogMetadata, posts)

	if err != nil {
		t.Fatalf("cannot generate the blog: %v", err)
	}

	expectedHomePage := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My blog !</title>
</head>
<body>
    <header><a href="/">My blog !</a></header>
    <ul class="posts">
        
            <li><a href="posts/second-post/">Second post</a></li>
        
            <li><a href="posts/first-post/">First post</a></li>
        
    </ul>
</body>
</html>`

	homePagePath := fmt.Sprintf("%v/index.html", generatedBlogPath)
	homePage, err := ioutil.ReadFile(homePagePath)
	if err != nil {
		t.Fatalf("cannot open the generated home page: %v", err)
	}

	if string(homePage) != expectedHomePage {
		t.Fatalf("the generated home page is not as expected, got: %v\nexpected: %v", string(homePage), expectedHomePage)
	}
}

func TestGenerateThePostsPages(t *testing.T) {
	blogMetadata := blog.BlogMetadata{
		BlogTitle:            "My blog !",
		Locale:               "en",
		PostPageTemplatePath: "./postPageTemplate.go.html",
		HomePageTemplatePath: "./homePageTemplate.go.html",
		DistPath:             "../dist",
		GlobalAssetsPath:     "",
	}
	posts := []posts.Post{
		{
			Name:      "First post",
			Extension: ".md",
			Content:   []byte("<h1>First post !</h1>"),
		},
		{
			Name:      "Second post",
			Extension: ".md",
			Content: []byte(`
<h1>Second post !</h1>

Content...`),
		},
	}
	generatedBlogPath, err := blog.GenerateBlog(blogMetadata, posts)

	if err != nil {
		t.Fatalf("cannot generate the blog: %v", err)
	}

	expectedFirstPost := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>First post - My blog !</title>
</head>
<body>
    <header><a href="/">My blog !</a></header>
    <main><h1>First post !</h1></main>
</body>
</html>`

	firstPostPagePath := fmt.Sprintf("%v/posts/first-post/index.html", generatedBlogPath)
	firstPostPage, err := ioutil.ReadFile(firstPostPagePath)
	if err != nil {
		t.Fatalf("cannot open the generated first post page: %v", err)
	}

	if string(firstPostPage) != expectedFirstPost {
		t.Fatalf("the generated first post page is not as expected, got: %v\nexpected: %v", string(firstPostPage), expectedFirstPost)
	}

	expectedSecondPost := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Second post - My blog !</title>
</head>
<body>
    <header><a href="/">My blog !</a></header>
    <main>
<h1>Second post !</h1>

Content...</main>
</body>
</html>`

	secondPostPagePath := fmt.Sprintf("%v/posts/second-post/index.html", generatedBlogPath)
	secondPostPage, err := ioutil.ReadFile(secondPostPagePath)
	if err != nil {
		t.Fatalf("cannot open the generated second post page: %v", err)
	}

	if string(secondPostPage) != expectedSecondPost {
		t.Fatalf("the generated second post page is not as expected, got: %v\nexpected: %v", string(secondPostPage), expectedSecondPost)
	}
}

func TestAssets(t *testing.T) {
	err := os.RemoveAll("../assets")
	if err != nil {
		t.Fatalf("unable to remove assets directory: %v", err)
	}
	err = os.MkdirAll("../assets/my/path", 0755)
	if err != nil {
		t.Fatalf("unable to create dir: %v", err)
	}

	firstAssetToCopy := `
	.title {
		font-size: 12px;
	}
	`
	err = ioutil.WriteFile("../assets/my/path/design.css", []byte(firstAssetToCopy), 0755)
	if err != nil {
		t.Fatalf("unable to write file: %v", err)
	}

	blogMetadata := blog.BlogMetadata{
		BlogTitle:            "My blog !",
		Locale:               "en",
		PostPageTemplatePath: "./postPageTemplate.go.html",
		HomePageTemplatePath: "./homePageTemplate.go.html",
		DistPath:             "../dist",
		GlobalAssetsPath:     "../assets",
	}

	generatedBlogPath, err := blog.GenerateBlog(blogMetadata, []posts.Post{})

	if err != nil {
		t.Fatalf("cannot generate the blog: %v", err)
	}

	firstAssetPath := fmt.Sprintf("%v/assets/my/path/design.css", generatedBlogPath)
	copiedFirstAsset, err := ioutil.ReadFile(firstAssetPath)
	if err != nil {
		t.Fatalf("cannot open the first asset: %v", err)
	}

	if string(copiedFirstAsset) != firstAssetToCopy {
		t.Fatalf("the first asset is not as expected, got: %v\nexpected: %v", string(copiedFirstAsset), firstAssetToCopy)
	}
}
