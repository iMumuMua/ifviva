package ifviva

import (
	// "errors"
	"log"
	"net/http"
	// "strings"
)

type Application struct {
	Router
}

func (app *Application) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func(rw http.ResponseWriter) {
		if r := recover(); r != nil {
			rw.WriteHeader(500)
			rw.Write([]byte("Internal Server Error."))
		}
	}(rw)

	isMatch, handler, params := app.Match(r.Method, r.URL.Path)
	if isMatch == true {
		handler(Context{r, rw, params})
	} else {
		rw.WriteHeader(404)
		rw.Write([]byte("Not Found."))
	}
}

func (app *Application) Run(port string) {
	log.Println("Ifviva Application listen on", port)
	log.Fatal(http.ListenAndServe(":"+port, app))
}
