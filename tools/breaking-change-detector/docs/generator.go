package docs

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/GoogleCloudPlatform/magic-modules/.ci/breaking-change-detector/constants"
	"github.com/GoogleCloudPlatform/magic-modules/.ci/breaking-change-detector/rules"
	"github.com/golang/glog"
)

func Generate(outputPath string) {
	rs := rules.GetRules()
	tmpl, err := template.New("breaking-changes.md.tmpl").Parse(markdownTemplate)
	if err != nil {
		glog.Exit(err)
	}

	contents := bytes.Buffer{}
	if err = tmpl.ExecuteTemplate(&contents, "breaking-changes.md.tmpl", rs); err != nil {
		glog.Exit(err)
	}

	if outputPath == "" {
		fmt.Printf("%v", contents.String())
	} else {
		outName := constants.BreakingChangeFileName
		err := os.WriteFile(path.Join(outputPath, outName), contents.Bytes(), 0644)
		if err != nil {
			glog.Exit(err)
		}
	}

}
