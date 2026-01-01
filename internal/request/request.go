package request

import (
	"errors"
	"io"
	"log"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ErrInvalidInput = errors.New("invalid request format")

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
		return &Request{}, err
	}
	reqInfo := strings.Split(string(data), "\r\n")
	if len(reqInfo) == 1 {
		log.Println("invalid request format")
		return &Request{}, ErrInvalidInput
	}

	reqMethod, reqTarget, httpVersion, err := parseRequestLine(reqInfo[0])
	if err != nil {
		log.Println(err)
		return &Request{}, err
	}

	reqLine := RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: reqTarget,
		Method:        reqMethod,
	}

	return &Request{
		RequestLine: reqLine,
	}, nil
}

func parseRequestLine(reqLine string) (string, string, string, error) {
	reqLineParts := strings.Split(reqLine, " ")
	if len(reqLineParts) != 3 {
		return "", "", "", errors.New("invalid request line")
	}

	for _, c := range reqLineParts[0] {
		if !unicode.IsUpper(c) {
			return "", "", "", errors.New("invalid request method")
		}
	}

	versionInfo := strings.Split(reqLineParts[2], "/")
	if len(versionInfo[1]) != 3 {
		return "", "", "", errors.New("invalid http version, should be HTTP/1.1")
	}
	
	if !(string(versionInfo[1][0]) == "1" && string(versionInfo[1][1]) == "." && string(versionInfo[1][2]) == "1") {
		return "", "", "", errors.New("invalid http version, should be HTTP/1.1")
	}

	return reqLineParts[0], reqLineParts[1], versionInfo[1], nil
}
