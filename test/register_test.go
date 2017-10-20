package test

import (
	"cloud.google.com/go/profiler/mocks"
	"gopkg.in/authboss.v1"
	"gopkg.in/authboss.v1/register"
	"html/template"
	"net/http"
)

func setup() *register.Register {
	ab := authboss.New()
	ab.RegisterOKPath = "/regsuccess"
	ab.Layout = template.Must(template.New("").Parse(`{{template "authboss" .}}`))
	ab.XSRFName = "xsrf"
	ab.XSRFMaker = func(_ http.ResponseWriter, _ *http.Request) string {
		return "xsrfvalue"
	}
	ab.ConfirmFields = []string{"password", "confirm_password"}
	ab.Storer = mocks.NewMockStorer()

	reg := register.Register{}
	if err := reg.Initialize(ab); err != nil {
		panic(err)
	}

}
