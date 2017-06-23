package telegram

import (
	"strconv"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"log"
)

func DoAsyncInlineQueryAPICall(id int, output *chan SentItem, apiurl string) {
	var sent SentItem
	sent.Id = id

	r, e := http.Get(apiurl)

	if r != nil {
		defer r.Body.Close()
	}
	if e != nil {
		sent.Error = e
		if (output != nil) { *output <- sent }
		return
	}

	sent.Http_code = r.StatusCode
	b, e := ioutil.ReadAll(r.Body)
	if e != nil {
		sent.Error = e
		if (output != nil) { *output <- sent }
		return
	}

	var out TGenericResponse
	e = json.Unmarshal(b, &out)

	if e != nil {
		sent.Error = e
		if (output != nil) { *output <- sent }
		return
	}

	e = HandleSoftError(&out)
	if e != nil {
		sent.Error = e
		if (output != nil) { *output <- sent }
		return
	}

	sent.Success = true
	if (output != nil) { *output <- sent }
	return
}

func DoInlineQueryAPICall(apiurl string) (error) {
	ch := make(chan SentItem, 1)

	DoAsyncInlineQueryAPICall(-1, &ch, apiurl)
	close(ch)
	sent := <- ch

	return sent.Error
}


func BuildAnswerInlineQueryURL(q TInlineQuery, last_offset int) (string) {
	// next_offset should get stuck at -1 forever if pagination breaks somehow, to prevent infinite loops.
	next_offset := ""
	if last_offset == -1 {
		next_offset = "-1"
	} else {
		next_offset = strconv.Itoa(last_offset + 1)
	}

	return apiEndpoint + apiKey + "/answerInlineQuery?" +
		"inline_query_id=" + url.QueryEscape(q.Id) +
		"&next_offset=" + next_offset +
		"&cache_time=30" +
		"&results="
}

func AnswerInlineQuery(q TInlineQuery, out []interface{}, last_offset int) (error) {
	b, e := json.Marshal(out)
	if e != nil {
		return e
	}

	surl := BuildAnswerInlineQueryURL(q, last_offset)
	log.Printf("[telegram] API call: %s\n", surl + "[snip]")
	e = DoInlineQueryAPICall(surl + url.QueryEscape(string(b)))

	if e == nil {
		log.Printf("[telegram] Pushed %d inline query results (id: %s)", len(out), q.Id)
	}

	return e
}

func AnswerInlineQueryAsync(q TInlineQuery, out []interface{}, last_offset int, output *chan SentItem) (int) {
	b, e := json.Marshal(out)
	if e != nil {
		return -1
	}

	surl := BuildAnswerInlineQueryURL(q, last_offset)
	log.Printf("[telegram] Async API call: %s\n", surl + "[snip]")

	current_id := GetNextId()
	go DoAsyncInlineQueryAPICall(current_id, output, surl + url.QueryEscape(string(b)))
	return current_id
}

