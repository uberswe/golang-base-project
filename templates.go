// Package baseproject is the main package of the golang-base-project which defines all routes and database connections and settings, the glue to the entire application
package baseproject

import (
	"html/template"
	"io/fs"
	"io/ioutil"
	"strings"
)

func loadTemplates() (*template.Template, error) {
	var err4 error
	t := template.New("")
	err := fs.WalkDir(staticFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		f, err2 := staticFS.Open(path)
		if err2 != nil {
			return err2
		}
		h, err3 := ioutil.ReadAll(f)
		if err3 != nil {
			return err3
		}
		parts := strings.Split(path, "/")
		if len(parts) > 0 && strings.HasSuffix(parts[len(parts)-1], ".html") {
			t, err4 = t.New(parts[len(parts)-1]).Parse(string(h))
			if err4 != nil {
				return err4
			}
		}
		return nil
	})
	return t, err
}
