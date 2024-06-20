// onchange: react on changes
// Copyright 2024 Martin Holst Swende @holiman
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/holiman/onchange/watch"
	"github.com/urfave/cli/v2"
	"os/signal"
	"strings"
	"time"
)

var (
	stdoutFlag = &cli.StringFlag{
		Name:    "stdout",
		Aliases: []string{"1"},
		Usage:   "Where to direct stdout-output from the executed command",
	}
	stderrFlag = &cli.StringFlag{
		Name:    "stderr",
		Aliases: []string{"2"},
		Usage:   "Where to direct stderr-output from the executed command",
	}
	outFlag = &cli.StringFlag{
		Name:  "output",
		Usage: "Where to direct both stdout- and stderr-output from the executed command",
	}
	intervalFlag = &cli.DurationFlag{
		Name:  "polling",
		Usage: "How long to wait between rechecks (for polling checks: tcp/web), default is 5m",
		Value: 5 * time.Minute,
	}
)

func initApp() *cli.App {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "onchange"
	app.Usage = "Do things when files are modified"
	app.Copyright = "Copyright 2024 @holiman"
	app.Flags = []cli.Flag{
		outFlag,
		stdoutFlag,
		stderrFlag,
		intervalFlag,
	}
	app.Action = onchange
	app.Commands = []*cli.Command{}
	return app
}

var app = initApp()

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func onchange(c *cli.Context) error {
	if c.NArg() < 2 {
		return fmt.Errorf("need 2 args, got %d", c.NArg())
	}
	var (
		subject = c.Args().First()
		stderr  = os.Stderr
		stdout  = os.Stdout
		args    = c.Args().Tail()
	)
	if c.IsSet(outFlag.Name) {
		f, err := os.Create(c.String(outFlag.Name))
		if err != nil {
			return fmt.Errorf("failed to open output: %v", err)
		}
		defer f.Close()
		stdout = f
		stderr = f
	}
	if c.IsSet(stdoutFlag.Name) {
		f, err := os.Create(c.String(stdoutFlag.Name))
		if err != nil {
			return fmt.Errorf("failed to open output: %v", err)
		}
		defer f.Close()
		stdout = f
	}
	if c.IsSet(stderrFlag.Name) {
		f, err := os.Create(c.String(stderrFlag.Name))
		if err != nil {
			return fmt.Errorf("failed to open output: %v", err)
		}
		defer f.Close()
		stderr = f
	}
	callback := func(status string) {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		slog.Info("Running", "cmd", cmd.String())
		if err := cmd.Run(); err != nil {
			slog.Info("Command errored", "err", err)
		}
	}

	var w watch.Watcher
	// Is it an HTTP URL?
	if strings.HasPrefix(subject, "http://") || strings.HasPrefix(subject, "https://") {
		if ww, err := watch.NewWebWatcher(subject, c.Duration(intervalFlag.Name), callback); err != nil {
			return err
		} else {
			w = ww
		}
	} else if strings.HasPrefix(subject, "tcp://") {
		if ww, err := watch.NewTcpWatcher(strings.TrimPrefix(subject, "tcp://"),
			c.Duration(intervalFlag.Name),
			callback); err != nil {
			return err
		} else {
			w = ww
		}
	} else {
		if ww, err := watch.NewFileWatcher(subject, callback); err != nil {
			return err
		} else {
			w = ww
		}
	}
	slog.Info("Watching", "subject", subject,
		"cmd", c.Args().Tail()[0],
		"args", c.Args().Tail()[1:])
	w.Start()
	defer w.Stop()
	abortChan := make(chan os.Signal, 1)
	signal.Notify(abortChan, os.Interrupt)
	<-abortChan
	fmt.Fprintf(os.Stderr, "Shutting down\n")
	return nil
}
