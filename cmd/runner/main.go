//go:generate go run -mod=vendor

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

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
	mainWG := &sync.WaitGroup{}

	mainWG.Add(1)

	go func(c chan result) {
		defer mainWG.Done()

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
			mainWG.Add(1)
			waitGroup.Add(1)

			go func() {
				defer mainWG.Done()
				defer waitGroup.Done()

				execHook(channel, filepath.Join(hook.FullPath(), file), args)
			}()
		}

		waitGroup.Wait()
	}

	close(channel)
	mainWG.Wait()
}

func getHooks(path string) [][]string {
	files, _ := os.ReadDir(path)
	list := make(map[string][]string)
	groups := make(map[string]string, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		parts := strings.SplitN(file.Name(), "-", 2)

		list[parts[0]] = append(list[parts[0]], file.Name())

		if _, ok := groups[parts[0]]; !ok {
			groups[parts[0]] = parts[0]
		}
	}

	return sortHooks(list, groups)
}

func execHook(c chan result, src string, args []string) {
	out, err := helpers.ScriptExec(src, args)

	c <- result{
		out: string(out),
		err: err,
	}
}

func sortHooks(hooks map[string][]string, groups map[string]string) [][]string {
	list := make([][]string, 0)
	keys := make([]string, 0)

	for key := range hooks {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, group := range keys {
		list = append(list, hooks[groups[group]])
	}

	return list
}
