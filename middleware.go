package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/sirupsen/logrus"
)

var ExceptionURLS = []string{}

type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
	response   []byte
}

func NewResponseWriter(w http.ResponseWriter) *customResponseWriter {
	return &customResponseWriter{w, http.StatusOK, nil}
}

func (rw *customResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *customResponseWriter) Write(content []byte) (int, error) {
	rw.response = content
	return rw.ResponseWriter.Write(content)
}

// ExceptionHandlerMiddleware is a handlers which notfies on slack if any exceptions are encounters in
// http requests such as status code >= 400
func ExceptionHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create custom writer
		rw := NewResponseWriter(w)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusTeapot)
			return
		}

		// setup recovery
		defer func() {
			err := recover()
			if err != nil {
				// capture stacks trace
				stackTrace := string(debug.Stack())
				logrus.Error(err) // May be log this error?
				// print stack trace as well
				fmt.Println(stackTrace)
				w.WriteHeader(http.StatusInternalServerError)
				NotifyError("Exception Handler Middleware Recovery", "Check stack trace above", fmt.Sprintf("%v\n%s", err, stackTrace), "Request", string(body), "URI", r.RequestURI)
			}
		}()

		// clone request and send it ahead
		cloneRequest := r.Clone(r.Context())
		cloneRequest.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(rw, cloneRequest)

		description := fmt.Sprintf("failed request on *%s* with StatusCode: %d", r.RequestURI, rw.statusCode)

		body = func() []byte {
			var request map[string]interface{}
			if err := json.Unmarshal(body, &request); err != nil {
				return body
			}

			request = deleteFields(request)
			modifiedBody, err := json.Marshal(request)
			if err != nil {
				return body
			}

			return modifiedBody
		}()

		if !ArrayContains(ExceptionURLS, r.RequestURI) {
			if rw.statusCode >= 500 {
				NotifyError("Exception Handler Middleware Error", description, "", "Response", string(rw.response), "Request", string(body))
			} else if rw.statusCode >= 400 {
				NotifyWarn("Exception Handler Middleware Warn", description, "", "Response", string(rw.response), "Request", string(body))
			}
		}

		logrus.Debug("Exception Handler Middleware Passed!!!")
	})
}

func deleteFields(content map[string]interface{}) map[string]interface{} {
	for key, value := range content {
		if strings.EqualFold(key, "password") {
			delete(content, key)
			continue
		}

		switch value := value.(type) {
		case map[string]interface{}:
			content[key] = deleteFields(value)
		case []interface{}:
			// array to temporarily store decrypted objects
			tempArr := []interface{}{}
			for _, val := range value {
				// if the value is a map/obj
				if mp, isMap := val.(map[string]interface{}); isMap {
					obj := deleteFields(mp)
					tempArr = append(tempArr, obj)
				} else {
					tempArr = append(tempArr, val)
				}
			}
			// replace original array of objects with decrypted objects
			content[key] = tempArr
		}
	}

	return content
}
