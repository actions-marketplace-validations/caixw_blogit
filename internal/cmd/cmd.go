// SPDX-License-Identifier: MIT

// Package cmd 提供命令行相关的功能
package cmd

import (
	"flag"
	"os"

	"github.com/issue9/cmdopt"
	"github.com/issue9/term/v3/colors"

	"github.com/caixw/blogit/v2/internal/cmd/console"
	"github.com/caixw/blogit/v2/internal/cmd/create"
	"github.com/caixw/blogit/v2/internal/cmd/preview"
	"github.com/caixw/blogit/v2/internal/cmd/serve"
	"github.com/caixw/blogit/v2/internal/vars"
)

var (
	erro = &console.Logger{
		Prefix:   "[ERRO] ",
		Colorize: colors.New(os.Stderr),
		Color:    colors.Red,
	}

	info = &console.Logger{
		Prefix:   "[INFO] ",
		Colorize: colors.New(os.Stdout),
		Color:    colors.Yellow,
	}

	succ = &console.Logger{
		Prefix:   "[SUCC] ",
		Colorize: colors.New(os.Stdout),
		Color:    colors.Green,
	}
)

// Exec 执行命令行
func Exec(args []string) error {
	p, err := console.NewPrinter()
	if err != nil {
		return err
	}

	opt := &cmdopt.CmdOpt{
		Output:        os.Stdout,
		ErrorHandling: flag.ExitOnError,
		Header:        p.Sprintf("cmd title"),
		Footer:        p.Sprintf("cmd footer", vars.URL),
		CommandsTitle: p.Sprintf("sub command"),
		OptionsTitle:  p.Sprintf("cmd argument"),
		NotFound: func(name string) string {
			return p.Sprintf("sub command not found", name)
		},
	}

	initDrafts(opt, p)
	initBuild(opt, p)
	initVersion(opt, p)
	initStyles(opt, p)
	serve.Init(opt, succ, info, erro, p)
	preview.Init(opt, succ, info, erro, p)
	create.InitInit(opt, erro, p)
	create.InitPost(opt, succ, erro, p)
	opt.Help("help", p.Sprintf("help usage"))

	return opt.Exec(args)
}
