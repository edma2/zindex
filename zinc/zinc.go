package zinc

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/edma2/navigator/index"
	"github.com/edma2/navigator/zinc/parsing"
)

func Watch(paths chan string) chan *index.Update {
	return analysisChanges(analysisFileChanges(paths))
}

func isAnalysisFile(name string) bool {
	return strings.HasPrefix(path.Base(name), "inc_compile_") ||
		strings.HasSuffix(name, ".analysis") &&
			isRegularFile(name)
}

func isRegularFile(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}

func analysisFileChanges(paths chan string) chan string {
	files := make(chan string)
	go func() {
		for path := range paths {
			if isAnalysisFile(path) {
				files <- path
			}
		}
	}()
	return files
}

func analysisChanges(analysisFiles chan string) chan *index.Update {
	updates := make(chan *index.Update)
	go func() {
		for file := range analysisFiles {
			err := parsing.Parse(file, func(class, path string) {
				updates <- &index.Update{Class: class, Path: path}
			})
			if err != nil {
				log.Printf("error reading %s: %s\n", file, err)
			}
		}
	}()
	return updates
}
