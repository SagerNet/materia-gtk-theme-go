package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var pwd string

func shell(dir string, name string, args ...string) {
	println("> cd", dir)
	println(">", name, strings.Join(args, " "))
	proc := exec.Command(name, args...)
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	proc.Dir = dir
	must(proc.Start())
	status, err := proc.Process.Wait()
	must(err)
	if !status.Success() {
		os.Exit(status.ExitCode())
	}
}

func must(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	pwd, _ = filepath.Abs(".")

	materiaDir := fmt.Sprint(pwd, "/build/materia")
	_, err := os.Stat(materiaDir)
	if err != nil {
		shell(pwd, "git", "clone", "--depth", "1", "https://github.com/nana-4/materia-theme", materiaDir)
	}

	distDir := fmt.Sprint(pwd, "/build/dist")
	must(os.RemoveAll(distDir))

	shell(materiaDir, "meson", "_build", fmt.Sprint("-Dprefix=", distDir))
	shell(materiaDir, "meson", "install", "-C", "_build")

	themesDir := fmt.Sprint(distDir, "/share", "/themes")
	themesDirF, err := os.Open(themesDir)
	must(err)
	themeDirs, err := themesDirF.ReadDir(10)
	must(err)
	for _, themeDir := range themeDirs {
		themeDir := fmt.Sprint(themesDir, "/", themeDir.Name())
		if strings.HasSuffix(themeDir, "-compact") {
			println("Delete ", themeDir)
			must(os.RemoveAll(themeDir))
			continue
		}

		themeDirF, err := os.Open(themeDir)
		must(err)

		themeFiles, err := themeDirF.ReadDir(20)
		must(err)

		for _, themeFile := range themeFiles {
			if themeFile.IsDir() && themeFile.Name() != "gtk-3.0" {
				dirToDel := fmt.Sprint(themeDir, "/", themeFile.Name())
				println("Delete ", dirToDel)
				must(os.RemoveAll(dirToDel))
			}
		}
		must(themeDirF.Close())
	}
	must(themesDirF.Close())

	must(os.RemoveAll("themes"))
	shell(pwd, "mv", "-f", "build/dist/share/themes", ".")
	must(os.RemoveAll("build/dist"))
	must(os.RemoveAll("build/materia/_build"))

}
