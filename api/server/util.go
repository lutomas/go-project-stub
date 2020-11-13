package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/lutomas/go-project-stub/pkg/zap_logger"
	"github.com/lutomas/go-project-stub/types"
	"go.uber.org/zap"
)

// errorResponse writes response in a structured error message.
// Leaving http.Request for future implementation of multi-lang error messages (should be passed with the context)
func errorResponse(err error, statusCode int, resp http.ResponseWriter, _ *http.Request) {
	if statusCode == 0 {
		statusCode = 500
	}
	resp.WriteHeader(statusCode)

	// Set up the pipe to write data directly into the Reader.
	pr, pw := io.Pipe()

	var msg string
	if err != nil {
		msg = err.Error()
	} else {
		// setting default text from the code
		msg = http.StatusText(statusCode)
	}

	go func() {
		if err := pw.CloseWithError(json.NewEncoder(pw).Encode(&types.HttpErrorResponse{Msg: msg})); err != nil {
			zap_logger.GetInstance().With(zap.Error(err)).Error("Response encode")
		}
	}()

	_, writeErr := io.Copy(resp, pr)
	if writeErr != nil {
		zap_logger.GetInstance().Error("Http response", zap.Error(writeErr))
	}
}

func response(obj interface{}, statusCode int, err error, resp http.ResponseWriter, _ *http.Request) {
	// Check for an error
	if err != nil {

		code := http.StatusInternalServerError
		errMsg := err.Error()
		if strings.Contains(errMsg, "Permission denied") {
			code = http.StatusForbidden
		}
		resp.WriteHeader(code)
		resp.Write([]byte(err.Error()))
		return
	}

	resp.WriteHeader(statusCode)
	// Write out the JSON object
	if obj != nil {

		resp.Header().Set("Content-Type", "application/json")

		// Set up the pipe to write data directly into the Reader.
		pr, pw := io.Pipe()

		// Write JSON-encoded data to the Writer end of the pipe.
		// Write in a separate concurrent goroutine, and remember
		// to Close the PipeWriter, to signal to the paired PipeReader
		// that we’re done writing.
		go func() {
			pw.CloseWithError(json.NewEncoder(pw).Encode(obj))
		}()

		io.Copy(resp, pr)

		// encoding/json library has a specific bug(feature) to turn empty slices into json null object,
		// let's make an empty array instead
		// resp.Write(buf)
	}
}
