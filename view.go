package ifviva

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"path"
)

var (
	viewTemplate  *template.Template
	cacheTemplate *template.Template
	viewConfig    ViewConfig
	viewPaths     []string
)

type ViewConfig struct {
	Name    string
	IsCache bool
	ViewDir string
	CSSDir  string
	JSDir   string
	Funcs   template.FuncMap
}

func InitViewRenderer(config ViewConfig) {
	viewConfig = config

	if viewConfig.ViewDir == "" {
		panic("[ifviva.view]ViewDir should not be empty.")
	}
	if viewConfig.Name == "" {
		viewConfig.Name = "ifviva"
	}
	if viewConfig.CSSDir == "" {
		viewConfig.CSSDir = viewConfig.ViewDir
	}
	if viewConfig.JSDir == "" {
		viewConfig.JSDir = viewConfig.ViewDir
	}

	viewTemplate = template.New(viewConfig.Name)
	setViewPaths()
	setViewFuncs()

	var err error
	cacheTemplate, err = parseFiles()
	if err != nil {
		log.Println("[ifviva.view]Parse file error: ", err)
	}
}

func render(wr io.Writer, name string, data interface{}) error {
	var tpl *template.Template
	var err error
	if viewConfig.IsCache == false {
		tpl, err = parseFiles()
		if err != nil {
			return err
		}
	} else {
		tpl = cacheTemplate
	}
	return tpl.ExecuteTemplate(wr, name, data)
}

func setViewPaths() {
	scanDir(viewConfig.ViewDir, func(viewPath string) {
		viewPaths = append(viewPaths, viewPath)
	})
}

func setViewFuncs() {
	viewFuncs := template.FuncMap{
		"css": css,
	}
	for name, method := range viewConfig.Funcs {
		viewFuncs[name] = method
	}
	viewTemplate.Funcs(viewFuncs)
}

func parseFiles() (*template.Template, error) {
	cloneTemplate, err := viewTemplate.Clone()
	if err != nil {
		return nil, err
	}
	return cloneTemplate.ParseFiles(viewPaths...)
}

func scanDir(dir string, fn func(string)) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("[ifviva.view]Read view dir error: ", err)
		return
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			scanDir(path.Join(dir, fileInfo.Name()), fn)
		} else {
			fn(path.Join(dir, fileInfo.Name()))
		}
	}
}

// view methods
func css(filePath string) template.CSS {
	abPath := path.Join(viewConfig.CSSDir, filePath)
	content, err := ioutil.ReadFile(abPath)
	if err != nil {
		log.Println("[ifviva.view]View method {css} error: ", err)
		return template.CSS("")
	} else {
		return template.CSS(content)
	}
}
