package utils

import (
	"strings"
)

func FormatResponse(buffer []byte, bufferLength int) {
	var reqType string = ""
	var path string = ""
	var httpVersion string = ""
	keys := []string{}
	values := []string{}

	var firstLine string = ""
	hasFirstline := false
	var accumulator string = ""
	for _, char := range buffer[0:] {
		if char == '\n' {
			if !hasFirstline {
				firstLineSplit := strings.Split(accumulator, " ")
				reqType = firstLineSplit[0]
				path = firstLineSplit[1]
				httpVersion = firstLineSplit[2]
				hasFirstline = true
			}


			split := strings.Split(accumulator, ":")
			keys = append(keys, split[0])
			values = append(values, (split[1])[1:])
		}
	}

}
