package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func getTemplateDirectoryArg() (string, error) {
	var templateDirectory string
	flag.StringVar(&templateDirectory, "templatedir", "", "directory of template files")
	flag.Parse()

	if templateDirectory == "" {
		return "", errors.New("templatedir unspecified")
	}

	if _, err := os.Stat(templateDirectory); os.IsNotExist(err) {
		return "", errors.New("templatedir does not exist")
	}

	return templateDirectory, nil
}

func getTemplateFiles(templateDir string) []string {
	files := make([]string, 0)

	fs, _ := ioutil.ReadDir(templateDir)
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".template") {
			files = append(files, f.Name())
		}
	}

	return files
}

func main() {
	templateDirectory, err := getTemplateDirectoryArg()
	if err != nil {
		os.Exit(0)
	}

	templateFiles := getTemplateFiles(templateDirectory)
	if len(templateFiles) == 0 {
		os.Exit(0)
	}

	generatedSrc, _ := os.Create("nginx-template-files.go")
	generatedSrc.Write([]byte("package main \n\nconst (\n"))

	for _, filename := range templateFiles {
		varName := fmt.Sprintf("template_%v", strings.Replace(strings.TrimSuffix(filename, ".template"), "-", "_", -1))
		generatedSrc.Write([]byte(varName + " = `"))
		contents, _ := os.Open(path.Join(templateDirectory, filename))
		io.Copy(generatedSrc, contents)
		generatedSrc.Write([]byte("`\n"))

	}
	generatedSrc.Write([]byte(")\n"))
}
