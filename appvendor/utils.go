package appvendor

import (
	"becouple/models/xodb"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/volatiletech/authboss"
)

func InternalServerError(w http.ResponseWriter, err string) {
	logrus.WithField("error", err).Errorln("Internal Server Error")
	http.Error(w, err, http.StatusInternalServerError)
}

func ConcateErrorWith(errs []error, delim string) string {
	var s string
	for i, e := range errs {
		s += e.Error()
		if i < len(errs)-1 {
			s += delim
		}
	}

	return s
}

func BindAuthbossUser(user *xodb.User, attr authboss.Attributes) {
	if err := attr.Bind(user, true); err != nil {
		// if there is error, just warning, no error returned
		logrus.WithError(err).Warnln("cannot bind attribute to user")
	}
}
