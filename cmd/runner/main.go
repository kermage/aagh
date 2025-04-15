//go:generate go run -mod=vendor

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"aagh/internal/helpers"

	"github.com/kermage/GO-Mods/pathinfo"
)

type result struct {
	out string
	err error
}

func main() {
	script := filepath.Base(os.Args[0])
	args := os.Args[1:]
	hook := pathinfo.Get(filepath.Join(helpers.DIR, script))

	if !hook.Exists() {
		return
	}

	if !hook.Stats().IsDir() {
		out, err := helpers.ScriptExec(hook.FullPath(), args)
		msg := strings.TrimSpace(string(out))

		if msg != "" {
			fmt.Println(msg)
		}

		os.Exit(helpers.GetExitCode(err))

		return
	}

	channel := make(chan result)

	go func(c chan result) {
		for r := range c {
			msg := strings.TrimSpace(r.out)

			if msg != "" {
				fmt.Println(msg)
			}

			code := helpers.GetExitCode(r.err)

			if code != 0 {
				os.Exit(code)
			}
		}
	}(channel)

	for _, files := range getHooks(hook.FullPath()) {
		waitGroup := &sync.WaitGroup{}

		for _, file := range files {
			waitGroup.Add(1)

			go execHook(channel, waitGroup, filepath.Join(hook.FullPath(), file), args)
		}

		waitGroup.Wait()
	}

	time.Sleep(time.Second)
}

func getHooks(path string) map[string][]string {
	files, _ := os.ReadDir(path)
	list := make(map[string][]string)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		parts := strings.SplitN(file.Name(), "-", 2)

		list[parts[0]] = append(list[parts[0]], file.Name())
	}

	return list
}

func execHook(c chan result, wg *sync.WaitGroup, src string, args []string) {
	defer wg.Done()

	out, err := helpers.ScriptExec(src, args)

	c <- result{
		out: string(out),
		err: err,
	}
}
