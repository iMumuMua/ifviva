package ifviva

import (
	// "errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func expect(t *testing.T, actual interface{}, expect interface{}) {
	if actual != expect {
		t.Errorf("expect: %v (type: %v) but got: %v (type: %v)\n", expect, reflect.TypeOf(expect), actual, reflect.TypeOf(actual))
	}
}

func createReqRes(method string, url string) (req *http.Request, res *httptest.ResponseRecorder) {
	req, _ = http.NewRequest(method, url, nil)
	res = httptest.NewRecorder()
	return
}

type HomeCtrl struct {
	Controller
}

func (ctrl *HomeCtrl) Home() {
	ctrl.Status(200)
	ctrl.Text("ok")
}

func (ctrl *HomeCtrl) GetArticle(id string) {
	ctrl.Text(id)
}

func createApp() *Application {
	app := Application{}

	app.All("/", func(ctx Context) {
		homeCtrl := HomeCtrl{}
		homeCtrl.Init(ctx)
		homeCtrl.Home()
	})

	app.Get("/articles/:id", func(ctx Context) {
		HomeCtrl := HomeCtrl{}
		HomeCtrl.Init(ctx)
		HomeCtrl.GetArticle(ctx.Params["id"])
	})

	return &app
}

func Test_App_Run(t *testing.T) {
	app := Application{}
	go app.Run("3000")
}

func Test_App_Base(t *testing.T) {
	app := createApp()

	req, res := createReqRes("GET", "/")
	app.ServeHTTP(res, req)
	expect(t, res.Code, 200)
	expect(t, res.Body.String(), "ok")
}

func Test_App_Params(t *testing.T) {
	app := createApp()

	req, res := createReqRes("GET", "/articles/123")
	app.ServeHTTP(res, req)
	expect(t, res.Code, 200)
	expect(t, res.Body.String(), "123")
}
