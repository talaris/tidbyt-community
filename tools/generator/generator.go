package generator

import (
	_ "embed"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"text/template"

	"tidbyt.dev/community/apps/manifest"
)

const (
	appsDir = "apps"
	goExt   = ".go"
)

//go:embed templates/source.star.tmpl
var starSource string

//go:embed templates/source.go.tmpl
var goSource string

//go:embed templates/apps.go.tmpl
var appsSource string

// Generator provides a structure for generating apps.
type Generator struct {
	starTmpl *template.Template
	goTmpl   *template.Template
	appsTmpl *template.Template
}

type appsDef struct {
	Imports  []string
	Packages []string
}

// NewGenerator creates an instantiated generator with the templates parsed.
func NewGenerator() (*Generator, error) {
	starTmpl, err := template.New("star").Parse(starSource)
	if err != nil {
		return nil, err
	}

	goTmpl, err := template.New("go").Parse(goSource)
	if err != nil {
		return nil, err
	}

	appsTmpl, err := template.New("apps").Parse(appsSource)
	if err != nil {
		return nil, err
	}

	return &Generator{
		starTmpl: starTmpl,
		goTmpl:   goTmpl,
		appsTmpl: appsTmpl,
	}, nil
}

// GenerateApp creates the base app starlark, go package, and updates the app
// list.
func (g *Generator) GenerateApp(app *manifest.Manifest) error {
	err := g.createDir(app)
	if err != nil {
		return err
	}

	err = g.generateStarlark(app)
	if err != nil {
		return err
	}

	err = g.generateGo(app)
	if err != nil {
		return err
	}

	return g.updateApps(app)
}

// RemoveApp removes an app from the apps directory.
func (g *Generator) RemoveApp(app *manifest.Manifest) error {
	err := g.removeDir(app)
	if err != nil {
		return err
	}

	return g.updateApps(app)
}

func (g *Generator) createDir(app *manifest.Manifest) error {
	p := path.Join(appsDir, app.PackageName)
	return os.MkdirAll(p, os.ModePerm)
}

func (g *Generator) removeDir(app *manifest.Manifest) error {
	p := path.Join(appsDir, app.PackageName)
	return os.RemoveAll(p)
}

func (g *Generator) updateApps(app *manifest.Manifest) error {
	imports := []string{
		"tidbyt.dev/community/" + appsDir + "/manifest",
	}
	packages := []string{}

	files, err := ioutil.ReadDir(appsDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() && f.Name() != "manifest" {
			imp := "tidbyt.dev/community/" + appsDir + "/" + f.Name()
			imports = append(imports, imp)
			packages = append(packages, f.Name())
		}
	}
	p := path.Join(appsDir, appsDir+goExt)

	file, err := os.Create(p)
	if err != nil {
		return err
	}
	defer file.Close()

	sort.Strings(imports)
	sort.Strings(packages)

	a := &appsDef{
		Imports:  imports,
		Packages: packages,
	}

	return g.appsTmpl.Execute(file, a)
}

func (g *Generator) generateStarlark(app *manifest.Manifest) error {
	p := path.Join(appsDir, app.PackageName, app.FileName)

	file, err := os.Create(p)
	if err != nil {
		return err
	}
	defer file.Close()

	return g.starTmpl.Execute(file, app)
}

func (g *Generator) generateGo(app *manifest.Manifest) error {
	p := path.Join(appsDir, app.PackageName, app.PackageName+goExt)

	file, err := os.Create(p)
	if err != nil {
		return err
	}
	defer file.Close()

	return g.goTmpl.Execute(file, app)
}
