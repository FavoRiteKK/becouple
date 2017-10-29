package appvendor

import (
	"becouple/models/xodb"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/authboss"
	"net/http"
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

func BindAuthbossUser(user *xodb.User, attr authboss.Attributes) error {
	if err := attr.Bind(user, true); err != nil {
		logrus.WithError(err).Errorln("cannot bind attribute to user")
		return err
	}

	return nil

}
