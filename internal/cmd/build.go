// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"time"

	"github.com/issue9/cmdopt"

	"github.com/caixw/blogit"
	"github.com/caixw/blogit/filesystem"
)

var (
	buildSrc  string
	buildDest string
)

// initBuild 注册 build 子命令
func initBuild(opt *cmdopt.CmdOpt) {
	fs := opt.New("build", "编译内容\n", build)
	fs.StringVar(&buildSrc, "src", "./", "指定源码目录")
	fs.StringVar(&buildDest, "dest", "./dest", "指定输出目录")
}

func build(w io.Writer) error {
	start := time.Now()

	info.println("开始编译内容")
	if err := blogit.Build(buildSrc, filesystem.Dir(buildDest)); err != nil {
		erro.println(err.Error())
		return nil
	}

	succ.printf("完成编译，用时：%v\n", time.Since(start))
	return nil
}
