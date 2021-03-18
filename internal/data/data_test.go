// SPDX-License-Identifier: MIT

package data

import (
	"testing"
	"time"

	"github.com/issue9/assert"
)

func TestLoad(t *testing.T) {
	a := assert.New(t)

	data, err := Load("../../testdata/src")
	a.NotError(err).NotNil(data)

	a.Equal(data.Icon.Type, "image/png").Equal(data.Icon.Sizes, "256x256")
	a.Equal(3, len(data.Index.Posts)).
		Equal(data.Index.Posts[1].Prev, data.Index.Posts[0]).
		Equal(data.Index.Posts[1].Next, data.Index.Posts[2])
	a.NotNil(data.Index.Posts[1].Authors)
	a.NotNil(data.Index.Posts[1].License)
	a.NotNil(data.License)
	a.NotEmpty(data.Authors)

	a.True(data.Builded.After(time.Time{}))
}

func TestBuildURL(t *testing.T) {
	a := assert.New(t)

	base := "https://example.com/"
	a.Equal(buildURL(base, "/p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(buildURL(base, "p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(buildURL(base, ""), "https://example.com/")
	a.Equal(buildURL(base, "/"), "https://example.com/")

	base = "https://example.com"
	a.Equal(buildURL(base, "/p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(buildURL(base, "p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(buildURL(base, ""), "https://example.com/")
	a.Equal(buildURL(base, "/"), "https://example.com/")

	base = ""
	a.Equal(buildURL(base, "/p1/p2.md"), "/p1/p2.md")
	a.Equal(buildURL(base, "p1/p2.md"), "/p1/p2.md")
	a.Equal(buildURL(base, ""), "/")
	a.Equal(buildURL(base, "/"), "/")
}

func TestBuildThemeURL(t *testing.T) {
	a := assert.New(t)

	url := "https://example.com/"
	id := "def"
	a.Equal(buildThemeURL(url, id, "/p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, "p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, ""), "https://example.com/themes/def")
	a.Equal(buildThemeURL(url, id, "/"), "https://example.com/themes/def")

	url = "https://example.com"
	id = "def"
	a.Equal(buildThemeURL(url, id, "/p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, "p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, ""), "https://example.com/themes/def")
	a.Equal(buildThemeURL(url, id, "/"), "https://example.com/themes/def")

	url = "https://example.com/"
	id = "/def"
	a.Equal(buildThemeURL(url, id, "/p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, "p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, ""), "https://example.com/themes/def")
	a.Equal(buildThemeURL(url, id, "/"), "https://example.com/themes/def")

	url = "https://example.com"
	id = "/def"
	a.Equal(buildThemeURL(url, id, "/p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, "p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, ""), "https://example.com/themes/def")
	a.Equal(buildThemeURL(url, id, "/"), "https://example.com/themes/def")

	url = "https://example.com/"
	id = ""
	a.Equal(buildThemeURL(url, id, "/p1/p2.md"), "https://example.com/themes/p1/p2.md")
	a.Equal(buildThemeURL(url, id, "p1/p2.md"), "https://example.com/themes/p1/p2.md")
	a.Equal(buildThemeURL(url, id, ""), "https://example.com/themes")
	a.Equal(buildThemeURL(url, id, "/"), "https://example.com/themes")

	url = "https://example.com"
	id = ""
	a.Equal(buildThemeURL(url, id, "/p1/p2.md"), "https://example.com/themes/p1/p2.md")
	a.Equal(buildThemeURL(url, id, "p1/p2.md"), "https://example.com/themes/p1/p2.md")
	a.Equal(buildThemeURL(url, id, ""), "https://example.com/themes")
	a.Equal(buildThemeURL(url, id, "/"), "https://example.com/themes")
}

func TestBuildPath(t *testing.T) {
	a := assert.New(t)

	a.Panic(func() { buildPath("") })

	a.Equal(buildPath("slug"), "slug.xml")
	a.Equal(buildPath("/slug"), "slug.xml")
	a.Equal(buildPath("/slug.xml"), "slug.xml.xml")
}
