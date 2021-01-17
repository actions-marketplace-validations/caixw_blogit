// SPDX-License-Identifier: MIT

package cmd

import (
	"io"

	"github.com/issue9/cmdopt"

	"github.com/caixw/blogit"
)

var (
	serveSrc  string
	serveAddr string
	servePath string
)

// initServe 注册 serve 子命令
func initServe(opt *cmdopt.CmdOpt) {
	fs := opt.New("serve", "以 HTTP 服务运行", serve)
	fs.StringVar(&serveSrc, "src", "./", "指定源码目录")
	fs.StringVar(&serveAddr, "addr", ":8080", "服务端口")
	fs.StringVar(&servePath, "path", "/", "根路径")
}

func serve(w io.Writer) error {
	return blogit.Serve(serveSrc, serveAddr, servePath, info.asLogger())
}