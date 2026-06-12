package project

import (
	"os"
	"path/filepath"
	"strings"
)

/**
 * project structure (running):
 * audio/
 * effects/
 * project.json
 */

func CreateProject(name string, dir string) error {
	path := filepath.Join(dir, strings.TrimSuffix(name, filepath.Ext(name)))
	allDirs := []string{
		path,
		filepath.Join(path, "audio"),
		filepath.Join(path, "effects"),
	}

	allFiles := []string{
		filepath.Join(path, "project.json"),
	}

	for _, dir := range allDirs {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return err
		}
	}

	for _, file := range allFiles {
		file, err := os.OpenFile(file, os.O_RDWR | os.O_CREATE, os.FileMode(0644))
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}
