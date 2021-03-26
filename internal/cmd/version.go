// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"
	"runtime"

	"github.com/caixw/blogit"
	"github.com/issue9/cmdopt"
)

// initVersion 注册 version 子命令
func initVersion(opt *cmdopt.CmdOpt) {
	opt.New("version", "显示版本号\n", func(w io.Writer) error {
		_, err := fmt.Fprintf(w, "blogit %s build with %s\n", blogit.Version(), runtime.Version())
		return err
	})
}
