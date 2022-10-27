package main

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/check-markdown/parser"
	"os"
	"path"
	"runtime"
)

func main() {
	md, err := parser.NewMarkdownDocByFilename("/home/wuxu/terraform-provider-azurerm/website/docs/r/aadb2c_directory.html.markdown")
	if err != nil {
		panic(err)
	}
	fmt.Printf("go md: %+v", *md)
	return

	docDir := resourceDocDir()
	fileNames := docFilesNames(docDir)
	for _, f := range fileNames {
		fullPath := path.Join(docDir, f)
		md, err := parser.NewMarkdownDocByFilename(fullPath)
		if err != nil {
			panic(err)
		}
		fmt.Sprintf("go md: %+v", *md)
		return
	}
}

func curDir() string {
	_, file, _, _ := runtime.Caller(0)
	return path.Dir(file)
}

func repoDir() string {
	return path.Dir(path.Dir(path.Dir(curDir())))
}

func resourceDocDir() string {
	return path.Join(repoDir(), "website", "docs", "r")
}

func docFilesNames(path string) (res []string) {
	dir, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, ent := range dir {
		res = append(res, ent.Name())
	}
	return res
}
