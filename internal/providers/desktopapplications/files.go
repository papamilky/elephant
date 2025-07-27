package main

import (
	"github.com/adrg/xdg"
	"github.com/charlievieth/fastwalk"
	"github.com/fsnotify/fsnotify"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	files        map[string]*DesktopFile
	filesMu      sync.RWMutex
	watcher      *fsnotify.Watcher
	regionLocale = ""
	langLocale   = ""
)

func loadFiles() {
	start := time.Now()

	filesMu.Lock()
	files = make(map[string]*DesktopFile)
	filesMu.Unlock()

	getLocale()

	dirs := xdg.ApplicationDirs

	conf := fastwalk.Config{
		Follow: true,
	}

	for _, root := range dirs {
		if _, err := os.Stat(root); err != nil {
			continue
		}

		walkFn := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				slog.Error(Name, "walk", err)
				os.Exit(1)
			}

			filesMu.RLock()
			_, exists := files[path]
			filesMu.RUnlock()

			if exists {
				return nil
			}

			if !d.IsDir() && filepath.Ext(path) == ".desktop" {
				filesMu.Lock()
				files[path] = parseFile(path, langLocale, regionLocale)
				filesMu.Unlock()
			}

			return err
		}

		if err := fastwalk.Walk(&conf, root, walkFn); err != nil {
			slog.Error(Name, "walk", err)
			os.Exit(1)
		}

	}

	filesMu.RLock()
	fileCount := len(files)
	filesMu.RUnlock()

	slog.Info(Name, "files", fileCount, "time", time.Since(start))
}

func startWatcher() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		slog.Error(Name, "watcher_init", err)
		return
	}

	dirs := xdg.ApplicationDirs
	watchedDirs := make(map[string]bool)

	for _, root := range dirs {
		if _, err := os.Stat(root); err != nil {
			continue
		}

		addDirToWatcher(root, watchedDirs)

		conf := fastwalk.Config{Follow: true}
		walkFn := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if d.IsDir() {
				addDirToWatcher(path, watchedDirs)
			}

			return nil
		}

		fastwalk.Walk(&conf, root, walkFn)
	}

	go watchFiles()

	slog.Info(Name, "watcher_dirs", len(watchedDirs))
}

func addDirToWatcher(dir string, watchedDirs map[string]bool) {
	if watchedDirs[dir] {
		return
	}

	if err := watcher.Add(dir); err != nil {
		slog.Warn(Name, "watcher_add", err, "dir", dir)
		return
	}

	watchedDirs[dir] = true
}

func watchFiles() {
	defer watcher.Close()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			handleFileEvent(event)

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			slog.Error(Name, "watcher", err)
		}
	}
}

func handleFileEvent(event fsnotify.Event) {
	if filepath.Ext(event.Name) != ".desktop" {
		// Handle directory creation to watch new subdirectories
		if event.Op&fsnotify.Create == fsnotify.Create {
			if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
				if err := watcher.Add(event.Name); err != nil {
					slog.Warn(Name, "watcher_add_new", err, "dir", event.Name)
				}
			}
		}
		return
	}

	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		handleFileCreate(event.Name)
	case event.Op&fsnotify.Write == fsnotify.Write:
		handleFileUpdate(event.Name)
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		handleFileRemove(event.Name)
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		handleFileRemove(event.Name)
	}
}

func handleFileCreate(path string) {
	desktopFile := parseFile(path, langLocale, regionLocale)
	if desktopFile != nil {
		filesMu.Lock()
		files[path] = desktopFile
		filesMu.Unlock()

		slog.Debug(Name, "file_created", path)
	}
}

func handleFileUpdate(path string) {
	desktopFile := parseFile(path, langLocale, regionLocale)
	if desktopFile != nil {
		filesMu.Lock()
		files[path] = desktopFile
		filesMu.Unlock()

		slog.Debug(Name, "file_updated", path)
	}
}

func handleFileRemove(path string) {
	filesMu.Lock()
	delete(files, path)
	filesMu.Unlock()

	slog.Debug(Name, "file_removed", path)
}

func getLocale() {
	regionLocale = config.Locale

	if regionLocale == "" {
		regionLocale = os.Getenv("LANG")

		langMessages := os.Getenv("LC_MESSAGES")
		if langMessages != "" {
			regionLocale = langMessages
		}

		langAll := os.Getenv("LC_ALL")
		if langAll != "" {
			regionLocale = langAll
		}

		regionLocale = strings.Split(regionLocale, ".")[0]
	}

	langLocale = strings.Split(regionLocale, "_")[0]
}
