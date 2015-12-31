package ifviva

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

const (
	contentType    = "Content-Type"
	contentText    = "text/plain"
	contentJSON    = "application/json"
	contentHTML    = "text/html"
	defaultCharset = "UTF-8"
)

var (
	cacheTemplate *template.Template
)

type Controller struct {
	Req        *http.Request
	Res        http.ResponseWriter
	Params     map[string]string
	statusCode int
	Err        error
	Charset    string
}

func SetViewPath(dir string) {
	viewPaths := []string{}
	scanDir(dir, func(viewPath string) {
		viewPaths = append(viewPaths, viewPath)
	})
	var err error
	cacheTemplate, err = template.ParseFiles(viewPaths...)
	if err != nil {
		log.Println("[ifviva]Set view path error: ", err)
	}
}

func scanDir(dir string, fn func(string)) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("[ifviva]Set view path error: ", err)
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

func (ctrl *Controller) Init(ctx Context) {
	ctrl.Req = ctx.Req
	ctrl.Res = ctx.Res
	ctrl.Params = ctx.Params
	ctrl.statusCode = 200
	ctrl.Charset = defaultCharset
}

func (ctrl *Controller) Status(status int) {
	ctrl.statusCode = status
}

func (ctrl *Controller) Text(text string) {
	ctrl.Res.Header().Set(contentType, appendCharset(contentText, ctrl.Charset))
	ctrl.Res.WriteHeader(ctrl.statusCode)
	ctrl.Res.Write([]byte(text))
}

func (ctrl *Controller) Json(v interface{}) {
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		ctrl.InternalError(err)
		return
	}
	ctrl.Res.Header().Set(contentType, appendCharset(contentJSON, ctrl.Charset))
	ctrl.Res.WriteHeader(ctrl.statusCode)
	ctrl.Res.Write(result)
}

func (ctrl *Controller) View(name string, data interface{}) {
	ctrl.Res.Header().Set(contentType, appendCharset(contentHTML, ctrl.Charset))
	ctrl.Res.WriteHeader(ctrl.statusCode)
	err := cacheTemplate.ExecuteTemplate(ctrl.Res, name, data)
	if err != nil {
		ctrl.InternalError(err)
		return
	}
}

func (ctrl *Controller) InternalError(err error) {
	ctrl.Err = err
	ctrl.Res.WriteHeader(500)
	ctrl.Res.Write([]byte("Internal Server Error"))
}

func appendCharset(content string, charset string) string {
	return content + "; charset=" + charset
}
