// SPDX-License-Identifier: MIT

package builder

import (
	"net/http"
	"os"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"

	"github.com/caixw/blogit/filesystem"
	"github.com/caixw/blogit/internal/vars"
)

func TestIsIgnore(t *testing.T) {
	a := assert.New(t)

	a.True(isIgnore("abc/def/a.md"))
	a.True(isIgnore("abc/def/.git"))
	a.True(isIgnore("themes/d/layout/header.html"))
	a.True(isIgnore("themes/layout/layout/header.html")) // 第一个 layout 为主题名称
	a.False(isIgnore("themes/d/theme.yaml"))
	a.False(isIgnore("themes/layout/header.html")) // 不符合 themes/xx/layout 的格式
}

func TestBuilder_ServeHTTP(t *testing.T) {
	a := assert.New(t)
	a.NotError(os.RemoveAll("../testdata/dest"))
	src := filesystem.Dir("../testdata/src")

	// Memory

	b := New(filesystem.Memory(), nil)
	srv := rest.NewServer(t, b, nil)

	// b 未加载任何数据。返回都是 404
	srv.Get("/robots.txt").Do().Status(http.StatusNotFound)

	a.NotError(b.Rebuild(src, "http://localhost:8080"))
	srv.Get("/robots.txt").Do().Status(http.StatusOK)
	srv.Get("/posts/p1" + vars.Ext).Do().Status(http.StatusOK)
	srv.Get("/posts/not-exists.html").Do().Status(http.StatusNotFound)
	srv.Get("/themes/default/style.css").Do().Status(http.StatusOK)

	// index.html
	srv.Get("/").Do().Status(http.StatusOK)
	srv.Get("/index" + vars.Ext).Do().Status(http.StatusOK)
	srv.Get("/themes/").Do().Status(http.StatusNotFound) // 目录下没有 index.html

	// Dir

	b = New(filesystem.Dir("../testdata/dest"), nil)
	srv = rest.NewServer(t, b, nil)

	// b 未加载任何数据。返回都是 404
	srv.Get("/robots.txt").Do().Status(http.StatusNotFound)

	a.NotError(b.Rebuild(src, "http://localhost:8080"))
	srv.Get("/robots.txt").Do().Status(http.StatusOK)
	srv.Get("/posts/p1" + vars.Ext).Do().Status(http.StatusOK)
	srv.Get("/posts/not-exists.html").Do().Status(http.StatusNotFound)
	srv.Get("/themes/default/style.css").Do().Status(http.StatusOK)

	// index.html
	srv.Get("/").Do().Status(http.StatusOK)
	srv.Get("/index" + vars.Ext).Do().Status(http.StatusOK)
	srv.Get("/themes/").Do().Status(http.StatusNotFound) // 目录下没有 index.html
}