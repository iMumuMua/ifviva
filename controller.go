package ifviva

import (
	"encoding/json"
	"html/template"
	"net/http"
)

const (
	contentType    = "Content-Type"
	contentText    = "text/plain"
	contentJSON    = "application/json"
	contentHTML    = "text/html"
	defaultCharset = "UTF-8"
)

type Controller struct {
	Req        *http.Request
	Res        http.ResponseWriter
	Params     map[string]string
	statusCode int
	Err        error
	Charset    string
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

func (ctrl *Controller) View(views map[string]string, data interface{}) {
	viewPaths := []string{}
	for _, viewPath := range views {
		viewPaths = append(viewPaths, viewPath)
	}
	t, err := template.ParseFiles(viewPaths...)
	if err != nil {
		ctrl.InternalError(err)
		return
	}

	ctrl.Res.Header().Set(contentType, appendCharset(contentHTML, ctrl.Charset))
	ctrl.Res.WriteHeader(ctrl.statusCode)
	err = t.ExecuteTemplate(ctrl.Res, "main", data)
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
