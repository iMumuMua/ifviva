package ifviva

import (
	"encoding/json"
	"net/http"
)

const (
	contentType    = "Content-Type"
	contentText    = "text/plain"
	contentJSON    = "application/json"
	defaultCharset = "UTF-8"
)

type Controller struct {
	Req        *http.Request
	Res        http.ResponseWriter
	statusCode int
	err        error
	Charset    string
}

func (ctrl *Controller) Init(ctx Context) {
	ctrl.Req = ctx.Req
	ctrl.Res = ctx.Res
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

func (ctrl *Controller) InternalError(err error) {
	ctrl.err = err
	ctrl.Res.WriteHeader(500)
	ctrl.Res.Write([]byte("Internal Server Error"))
}

func appendCharset(content string, charset string) string {
	return content + "; charset=" + charset
}
