package http

import (
	"io"
	"strconv"
)


type Client struct {
}

type dummyReader struct {
	content string
	lastIndex int
}

type Response struct {
	StatusCode int
	Body dummyReader
	Status string
}


var DefaultClient Client
var newRequestUrlReturn string
var newRequestErrReturn error
var doResponseReturn Response
var doErrorReturn error

const StatusOK = 200
const StatusInternalServerError = 500

func SetNewRequestReturn(s string, e error) {
	newRequestUrlReturn = s
	newRequestErrReturn = e
}

func SetDoReturn(body string, statusCode int, e error) {
	doResponseReturn.Body.content = body
	doResponseReturn.Body.lastIndex = 0
	doResponseReturn.StatusCode = statusCode
	doResponseReturn.Status  = "CODE: " + strconv.Itoa(statusCode)
	doErrorReturn = e
}


func NewRequest(_, _, _ any) (url string, err error) {
	return newRequestUrlReturn, newRequestErrReturn
}

func (c Client) Do (req string) (Response, error) {
	return doResponseReturn, doErrorReturn
}

func (d dummyReader) Read(p []byte) (n int, err error) {
	if d.lastIndex >= len(d.content) {
		return 0, io.EOF
	}
	n = 0
	for i := 0; i < len(p) && d.lastIndex < len(d.content); i++ {
		p[i] = d.content[d.lastIndex]
		d.lastIndex++
		n++
	}
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}