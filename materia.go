package materia_gtk_theme_go

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
)

//go:embed themes
var Themes embed.FS

func InitTheme(cacheDir string) error {
	themesDir := fmt.Sprint(cacheDir, "/themes")
	_, notExists := os.Stat(fmt.Sprint(themesDir, "/Materia/index.theme"))
	if notExists != nil {
		themes, err := Themes.ReadDir("themes")
		if err != nil {
			return err
		}
		for _, theme := range themes {
			err := extract("themes", theme, themesDir)
			if err != nil {
				return err
			}
		}
	}
	return os.Setenv("GTK_DATA_PREFIX", cacheDir)
}

func extract(path string, entry fs.DirEntry, dst string) error {
	if entry.IsDir() {
		path = fmt.Sprint(path, "/", entry.Name())
		dst = fmt.Sprint(dst, "/", entry.Name())
		entries, err := Themes.ReadDir(path)
		if err != nil {
			return err
		}
		if _, notExists := os.Stat(dst); notExists != nil {
			err := os.MkdirAll(dst, 0777)
			if err != nil {
				return err
			}
		}
		for _, entry := range entries {
			err := extract(path, entry, dst)
			if err != nil {
				return err
			}
		}
	} else {
		path = fmt.Sprint(path, "/", entry.Name())
		dst = fmt.Sprint(dst, "/", entry.Name())
		file, err := Themes.ReadFile(path)
		if err != nil {
			return err
		}
		var output *os.File
		if _, notExists := os.Stat(dst); notExists != nil {
			output, err = os.Create(dst)
		} else {
			output, err = os.Open(dst)
		}
		if err != nil {
			return err
		}
		_, err = io.Copy(output, bytes.NewReader(file))
		if err != nil {
			return err
		}
		err = output.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
