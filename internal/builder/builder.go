// SPDX-License-Identifier: MIT

// Package builder 提供编译成 HTML 的相关功能
package builder

import (
	"bytes"
	"encoding/xml"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/issue9/errwrap"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/filesystem"
	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

// Builder 提供了一个可重复生成 HTML 内容的对象
type Builder struct {
	log  *log.Logger
	fs   filesystem.WritableFS
	site *site
	tpl  *template.Template
}

// Build 编译内容
func Build(src, dest string) error {
	return New(filesystem.Dir(dest), log.Default()).Build(src, "")
}

// New 声明 Builder 实例
//
// fs 表示用于保存编译后的 HTML 文件的系统。可以是内存或是文件系统，
// 以及任何实现了 filesystem.WritableFS 接口都可以；
// l 表示的是在把 Builder 当作 http.Handler 处理时，在出错时的日志输出通道。
// 如果为空，则会采用 log.Default() 作为默认值。
func New(fs filesystem.WritableFS, l *log.Logger) *Builder {
	if l == nil {
		l = log.Default()
	}
	return &Builder{fs: fs, log: l}
}

// Build 重新生成数据
func (b *Builder) Build(src, base string) error {
	paths := make([]string, 0, 100)

	err := filepath.Walk(src, func(path string, d os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && !isIgnore(path) {
			paths = append(paths, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	for _, p := range paths {
		stat, err := os.Stat(p)
		if err != nil {
			return err
		}

		data, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		if err = b.appendFile(loader.Slug(src, p), stat.ModTime(), data); err != nil {
			return err
		}
	}

	return b.buildData(src, base)
}

func (b *Builder) buildData(src, base string) (err error) {
	d, err := data.Load(src, base)
	if err != nil {
		return err
	}

	b.tpl, err = newTemplate(d, src)
	if err != nil {
		return err
	}

	b.site = newSite(d)

	call := func(f func(*data.Data) error) {
		if err == nil {
			err = f(d)
		}
	}

	call(b.buildTags)
	call(b.buildPosts)
	call(b.buildSitemap)
	call(b.buildArchive)
	call(b.buildAtom)
	call(b.buildRSS)
	call(b.buildRobots)
	call(b.buildProfile)

	return
}

func isIgnore(src string) bool {
	// themes/**/layout/file 这种格式将忽略
	layout := path.Dir(src)
	if path.Base(layout) == vars.LayoutDir &&
		path.Base(path.Dir(path.Dir(layout))) == vars.ThemesDir {
		return true
	}

	ext := strings.ToLower(path.Ext(src))
	return ext == vars.MarkdownExt ||
		ext == ".yaml" ||
		ext == ".yml" ||
		ext == ".gitignore" ||
		ext == ".git"
}

// path 表示输出的文件路径，相对于源目录；
func (b *Builder) appendTemplateFile(path string, p *page) error {
	buf := &bytes.Buffer{}

	if err := b.tpl.ExecuteTemplate(buf, p.Type, p); err != nil {
		return err
	}

	return b.appendFile(path, time.Now(), buf.Bytes())
}

// path 表示输出的文件路径，相对于源目录；
// xsl 表示关联的 xsl，如果不需要则可能为空；
func (b *Builder) appendXMLFile(d *data.Data, path, xsl string, v interface{}) error {
	data, err := xml.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	buf := &errwrap.Buffer{}
	buf.WString(xml.Header)
	if xsl != "" {
		buf.Printf(`<?xml-stylesheet type="text/xsl" href="%s"?>`, xsl).WByte('\n')
	}
	buf.WBytes(data)

	if buf.Err != nil {
		return buf.Err
	}

	return b.appendFile(path, time.Now(), buf.Bytes())
}

// 如果 path 以 / 开头，则会自动去除 /
func (b *Builder) appendFile(p string, mod time.Time, data []byte) error {
	return b.fs.WriteFile(p, data, os.ModePerm)
}

// ServeHTTP 作为 HTTP 服务接口使用
func (b *Builder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const index = "index" + vars.Ext

	p := r.URL.Path
	if p != "" && p[0] == '/' {
		p = p[1:]
	}
	if p == "" || p[len(p)-1] == '/' {
		p += index
	}

	f, err := b.fs.Open(p)
	if errors.Is(err, os.ErrNotExist) {
		http.NotFound(w, r)
		return
	}
	if errors.Is(err, os.ErrPermission) {
		errStatus(w, http.StatusForbidden)
		return
	}
	if err != nil {
		b.log.Println(err)
		errStatus(w, http.StatusInternalServerError)
		return
	}

	stat, err := f.Stat()
	if err != nil {
		b.log.Println(err)
		errStatus(w, http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		b.log.Println(err)
		errStatus(w, http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, p, stat.ModTime(), bytes.NewReader(data))
}

func errStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
