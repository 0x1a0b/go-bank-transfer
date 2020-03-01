package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//Logger armazena a estrutura de logger para entrada e saídas da API
type Logger struct {
	logger *logrus.Logger
}

//NewLogger constrói um logger com suas dependências
func NewLogger(log *logrus.Logger) Logger {
	return Logger{logger: log}
}

//Logging cria logs de entrada e saída da API
func (l Logger) Logging(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	const logKey = "logger_middleware"

	body, err := getRequestPayload(r)
	if err != nil {
		l.logger.WithFields(logrus.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Error("error when getting payload")

		return
	}

	l.logger.WithFields(logrus.Fields{
		"key":         "api_request",
		"payload":     body,
		"url":         r.URL.Path,
		"http_method": r.Method,
	}).Info("request received by the API")

	next.ServeHTTP(w, r)

	l.logger.WithFields(logrus.Fields{
		"key":         "api_response",
		"url":         r.URL.Path,
		"http_method": r.Method,
	}).Info("response returned from the API")
}

func getRequestPayload(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", errors.New("body not defined")
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", errors.Wrap(err, "error read body")
	}

	// re-adiciona o payload ao buffer da request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	return strings.TrimSpace(string(payload)), nil
}
