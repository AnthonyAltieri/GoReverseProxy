package http

import (
	"strings"
	"fmt"
	"os"
	"strconv"
)

type Request struct {
	ReqType      string
	Path         string
	HttpVersion  string
	Headers      map[string]string
	MessageBody  map[string]string
}

func parseMessageBody(body string) map[string]string {
	messageBody := make(map[string]string)
	split := strings.Split(body, "&")
	for _, keyvalue := range split {
		splitKeyvalue := strings.Split(keyvalue, "=")
		messageBody[splitKeyvalue[0]] = splitKeyvalue[1]
	}
	return messageBody
}

func PrintRequest(request Request) {
	fmt.Println("{")
	fmt.Fprintf(os.Stdout, "\ttype: %s\n", request.ReqType)
	fmt.Fprintf(os.Stdout, "\tpath: %s\n", request.Path)
	fmt.Fprintf(os.Stdout, "\tHttpVersion: %s\n", request.HttpVersion)
	fmt.Fprintf(os.Stdout, "\tHeaders: [\n")
	for key := range request.Headers {
		fmt.Fprintf(os.Stdout, "\t\t %s: %s\n", key, request.Headers[key])
	}
	fmt.Fprintf(os.Stdout, "\t]\n")
	if request.ReqType == "POST" {
		fmt.Fprintf(os.Stdout, "\t{\n")
		for key := range request.MessageBody {
			fmt.Fprintf(os.Stdout, "\t\t%s: %s\n", key, request.MessageBody[key])
		}
		fmt.Fprintf(os.Stdout, "\t}\n")
	}
	fmt.Println("}")
}


func FormatRequest(buffer []byte, bufferLength int) Request {
	var reqType string = ""
	var path string = ""
	var httpVersion string = ""
	keys := []string{}
	values := []string{}
	messagebody := make(map[string]string)

	hasFirstline := false
	handlePost := false
	contentEvaluated := 0
	contentLength := -1

	var accumulator string = ""
	parseResponse:
	for _, byte := range string(buffer[0:]) {
		char := string(byte)
		if char == "\n" {
			if !hasFirstline {
				firstLineSplit := strings.Split(accumulator, " ")
				reqType = firstLineSplit[0]
				path = firstLineSplit[1]
				httpVersion = strings.Split(firstLineSplit[2], "\n")[0]
				hasFirstline = true
			} else {
				split := strings.Split(accumulator, ":")
				if len(split) == 1 {
					if reqType == "GET" {
						break parseResponse
					} else {
						accumulator = ""
						handlePost = true
						for i, key := range keys {
							if key == "content-length" {
								val := values[i]
								parsedVal, _:= strconv.Atoi(val)
								contentLength = parsedVal
							}
						}
						continue
					}
				}
				keys = append(keys, split[0])
				value := split[1]
				values = append(values, (value)[1:len(value) - 1])
			}
			accumulator = ""
		} else if handlePost {
			accumulator += char
			contentEvaluated += 1
			if contentEvaluated == contentLength {
				messagebody = parseMessageBody(accumulator)
			}
		} else {
			accumulator += char
		}
	}

	headers := make(map[string]string)
	for index, key := range keys[0:] {
		headers[key] = values[index]
	}


	return Request { reqType, path, httpVersion, headers, messagebody }
}
