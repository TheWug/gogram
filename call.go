package gogram

import (
	"errors"
	"encoding/json"
	"io/ioutil"
	"io"
	"github.com/thewug/reqtify"
)

var CallResponseChannel chan HandlerBox = make(chan HandlerBox, 10)

type HandlerBox struct {
	Success   bool
	Error     error
	Http_code int
	Handler   ResponseHandler
	Output   *json.RawMessage
	Bytes   []byte
	Reader   io.ReadCloser
}

// call this in a goroutine.
func DoAsyncGetReader(request *reqtify.Request, handler ResponseHandler, output *chan HandlerBox) {
	var hbox HandlerBox
	hbox.Handler = handler

	r, e := request.Do()
	if e != nil {
		hbox.Error = e
		if (output != nil) { *output <- hbox }
		return
	}

	hbox.Http_code = r.StatusCode
	hbox.Reader = r.Body

	hbox.Success = true
	if (output != nil) { *output <- hbox }
	return
}

// call this in a goroutine.
func DoAsyncFetch(request *reqtify.Request, handler ResponseHandler, output *chan HandlerBox) {
	temp := make(chan HandlerBox, 1)
	DoAsyncGetReader(request, handler, &temp)
	close(temp)
	hbox := <- temp

	if hbox.Reader != nil { defer hbox.Reader.Close() }

	if !hbox.Success {
		if (output != nil) { *output <- hbox }
		return
	}

	hbox.Success = false

	hbox.Bytes, hbox.Error = ioutil.ReadAll(hbox.Reader)
	if hbox.Error != nil {
		if (output != nil) { *output <- hbox }
		return
	}

	hbox.Success = true
	if (output != nil) { *output <- hbox }
	return
}

// call this in a goroutine.
func DoAsyncCall(request *reqtify.Request, handler ResponseHandler, output *chan HandlerBox) {
	temp := make(chan HandlerBox, 1)
	DoAsyncFetch(request, handler, &temp)
	close(temp)
	hbox := <- temp

	if !hbox.Success {
		if (output != nil) { *output <- hbox }
		return
	}

	hbox.Success = false

	var out TGenericResponse
	e := json.Unmarshal(hbox.Bytes, &out)

	if e != nil {
		hbox.Error = e
		if (output != nil) { *output <- hbox }
		return
	}

	e = HandleSoftError(&out)
	if e != nil {
		hbox.Error = e
		if (output != nil) { *output <- hbox }
		return
	}

	hbox.Output = out.Result
	hbox.Bytes = nil
	hbox.Success = true
	if (output != nil) { *output <- hbox }
	return
}

func DoGetReader(request *reqtify.Request) (io.ReadCloser, error) {
	ch := make(chan HandlerBox, 1)

	DoAsyncGetReader(request, nil, &ch)
	close(ch)
	output := <- ch

	return output.Reader, output.Error
}

func DoFetch(request *reqtify.Request) ([]byte, error) {
	ch := make(chan HandlerBox, 1)

	DoAsyncFetch(request, nil, &ch)
	close(ch)
	output := <- ch

	return output.Bytes, output.Error
}

func DoCall(request *reqtify.Request) (*json.RawMessage, error) {
	ch := make(chan HandlerBox, 1)

	DoAsyncCall(request, nil, &ch)
	close(ch)
	output := <- ch

	return output.Output, output.Error
}

// Type Helpers

func OutputToMessage(raw *json.RawMessage, err error) (*TMessage, error) {
	if err != nil { return nil, err }

	var msg TMessage
	err = json.Unmarshal(*raw, &msg)

	if err != nil { return nil, err }
	if msg.Message_id == 0 { return nil, errors.New("Missing message") }
	return &msg, nil
}

func OutputToStickerSet(raw *json.RawMessage, err error) (*TStickerSet, error) {
	if err != nil { return nil, err }

	var set TStickerSet
	err = json.Unmarshal(*raw, &set)

	if err != nil { return nil, err }
	if set.Stickers == nil { return nil, errors.New("Missing sticker set") }
	return &set, nil
}

func OutputToChatMember(raw *json.RawMessage, err error) (*TChatMember, error) {
	if err != nil { return nil, err }

	var cm TChatMember
	err = json.Unmarshal(*raw, &cm)

	if err != nil { return nil, err }
	if cm.Status == "" { return nil, errors.New("Missing chat member") }
	return &cm, nil
}

func OutputToFile(raw *json.RawMessage, err error) (*TFile, error) {
	if err != nil { return nil, err }

	var out TFile
	err = json.Unmarshal(*raw, &out)

	if err != nil { return nil, err }
	if out.File_id == "" { return nil, errors.New("Missing file metadata") }
	return &out, nil
}
