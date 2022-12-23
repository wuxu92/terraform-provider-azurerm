package md

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
)

var (
	docRDir string
	docDDir string
)

var resourceFilePathMap map[string]string
var file2Reosurce = map[string]string{}

var once sync.Once

// MDPathFor return full path of markdown file of resource
func MDPathFor(resourceType string) string {
	// find source
	fullPath := path.Join(ResourceDir(), fmt.Sprintf("%s.html.markdown", strings.TrimPrefix(resourceType, "azurerm_")))
	// check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return getMappingPath(resourceType)
	}
	return fullPath
}

func getMappingPath(resourceName string) (res string) {
	if resourceFilePathMap == nil {
		once.Do(func() {
			start := time.Now()
			tmpMap := map[string]string{}

			dir, err := os.ReadDir(ResourceDir())
			_ = err
			for _, en := range dir {
				if en.IsDir() {
					continue
				}
				fullPath := path.Join(ResourceDir(), en.Name())
				name := fileResource(fullPath)
				tmpMap[name] = fullPath
				if _, ok := file2Reosurce[fullPath]; !ok {
					file2Reosurce[fullPath] = name
				} else {
					panic(fmt.Sprintf("name: %s, path: %s", name, fullPath))
				}
			}
			resourceFilePathMap = tmpMap
			log.Printf("load %d resource, costs: %v", len(resourceFilePathMap), time.Now().Sub(start))
		})
	}
	return resourceFilePathMap[resourceName]
}

var titleReg = regexp.MustCompile(`\npage_title:[^\n]*(azurerm_[a-zA-Z0-9_]+)"?`)

func fileResource(path string) string {
	// match content
	f, _ := os.Open(path)
	rd := bufio.NewReader(f)
	content := make([]byte, 512)
	_, _ = rd.Read(content)
	// if content match pattern
	if subs := titleReg.FindStringSubmatch(string(content)); len(subs) > 1 {
		return subs[1]
	}
	return ""
}

func docDir() string {
	file, _ := util.FuncFileLine(utils.Int32)
	return path.Join(path.Dir(path.Dir(file)), "website", "docs")
}

func ResourceDir() string {
	if docRDir == "" {
		docRDir = path.Join(docDir(), "r")
	}
	return docRDir
}
