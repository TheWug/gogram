package gogram

import (
	"github.com/thewug/reqtify"

	"github.com/thewug/gogram/data"

	"encoding/json"
	"io/ioutil"
	"io"
	"log"
)

var call_response_channel chan HandlerBox = make(chan HandlerBox, 10)

type HandlerBox struct {
	Success   bool
	Error     error
	Http_code int
	Handler   data.ResponseHandler
	Output   *json.RawMessage
	Bytes   []byte
	Reader   io.ReadCloser
}

// call this in a goroutine.
func DoAsyncGetReader(logger *log.Logger, request *reqtify.Request, handler data.ResponseHandler) {
	doAsyncGetReader(logger, request, handler, call_response_channel)
}

func doAsyncGetReader(logger *log.Logger, request *reqtify.Request, handler data.ResponseHandler, output chan HandlerBox) {
	var hbox HandlerBox
	hbox.Handler = handler

	r, e := request.Do()

	status := "Request Failure"
	if r != nil {
		status = r.Status
		hbox.Http_code = r.StatusCode
		hbox.Reader = r.Body
	}

	if logger != nil {
		logger.Printf("[telegram] API call: %s (%s)\n", request.Path, status)
	}

	if e != nil {
		hbox.Error = e
		if (output != nil) { output <- hbox }
		return
	}

	hbox.Success = true
	if (output != nil) { output <- hbox }
	return
}

// call this in a goroutine.
func DoAsyncFetch(logger *log.Logger, request *reqtify.Request, handler data.ResponseHandler) {
	doAsyncFetch(logger, request, handler, call_response_channel)
}

func doAsyncFetch(logger *log.Logger, request *reqtify.Request, handler data.ResponseHandler, output chan HandlerBox) {
	temp := make(chan HandlerBox, 1)
	doAsyncGetReader(logger, request, handler, temp)
	close(temp)
	hbox := <- temp

	if hbox.Reader != nil { defer hbox.Reader.Close() }

	if !hbox.Success {
		if (output != nil) { output <- hbox }
		return
	}

	hbox.Success = false

	hbox.Bytes, hbox.Error = ioutil.ReadAll(hbox.Reader)
	if hbox.Error != nil {
		if (output != nil) { output <- hbox }
		return
	}

	hbox.Success = true
	if (output != nil) { output <- hbox }
	return
}

// call this in a goroutine.
func DoAsyncCall(logger *log.Logger, request *reqtify.Request, handler data.ResponseHandler) {
	doAsyncCall(logger, request, handler, call_response_channel)
}

func doAsyncCall(logger *log.Logger, request *reqtify.Request, handler data.ResponseHandler, output chan HandlerBox) {
	temp := make(chan HandlerBox, 1)
	doAsyncFetch(logger, request, handler, temp)
	close(temp)
	hbox := <- temp

	if !hbox.Success {
		if (output != nil) { output <- hbox }
		return
	}

	hbox.Success = false

	var out data.TGenericResponse
	e := json.Unmarshal(hbox.Bytes, &out)

	if e != nil {
		hbox.Error = e
		if (output != nil) { output <- hbox }
		return
	}

	e = HandleSoftError(&out)
	if e != nil {
		hbox.Error = e
		if (output != nil) { output <- hbox }
		return
	}

	hbox.Output = out.Result
	hbox.Bytes = nil
	hbox.Success = true
	if (output != nil) { output <- hbox }
	return
}

func DoGetReader(logger *log.Logger, request *reqtify.Request) (io.ReadCloser, error) {
	ch := make(chan HandlerBox, 1)

	doAsyncGetReader(logger, request, nil, ch)
	close(ch)
	output := <- ch

	return output.Reader, output.Error
}

func DoFetch(logger *log.Logger, request *reqtify.Request) ([]byte, error) {
	ch := make(chan HandlerBox, 1)

	doAsyncFetch(logger, request, nil, ch)
	close(ch)
	output := <- ch

	return output.Bytes, output.Error
}

func DoCall(logger *log.Logger, request *reqtify.Request) (*json.RawMessage, error) {
	ch := make(chan HandlerBox, 1)

	doAsyncCall(logger, request, nil, ch)
	close(ch)
	output := <- ch

	return output.Output, output.Error
}

// Type Helpers

func OutputToObject(raw *json.RawMessage, err error, output interface{}) (error) {
	if err != nil { return err }
	return json.Unmarshal(*raw, output)
}
