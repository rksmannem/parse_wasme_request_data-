package main

import (
	"errors"
	"fmt"
	"strings"
)

type HeaderName string

const (
	EmptyString = ""
	InHeader    = "header"
	InBody      = "body"

	ParamName HeaderName = "paramName"
)

var (
	ErrorHeaderNotFound        = errors.New("header not found")
	ErrorParamLocationNotFound = errors.New("`paramLocation` not found in resolve specs")
	ErrorInvalidParamLocation  = errors.New("invalid paramLocation value")
)

type ResolveSpec struct {
	ParamName     string `json:"paramName"`
	ParamLocation string `json:"ParamLocation"`
}

type ResolveSpecs []ResolveSpec

// ParseRequestHeaders ...
func ParseRequestHeaders(headers [][2]string) (map[string]string, error) {
	headersMap := make(map[string]string)
	for _, hdrEntry := range headers {
		fmt.Printf("header_name:'%s', header_value:'%s'\n", hdrEntry[0], hdrEntry[1])
		if hdrName := hdrEntry[0]; hdrName != "" {
			headersMap[strings.ToLower(hdrName)] = hdrEntry[1]
		}
	}

	return headersMap, nil
}

func GetHeader(hdrName string, headers map[string]string) (string, error) {
	hdrValue, ok := headers[strings.ToLower(hdrName)]
	if ok {
		return hdrValue, nil
	}
	return EmptyString, ErrorHeaderNotFound
}

func ProcessRequestDataForKey(headers [][2]string, resolveSpecs ResolveSpecs) (string, error) {
	headersMap, err := ParseRequestHeaders(headers)
	if err != nil {
		return EmptyString, err
	}

	// 1. get value of `paramName` from headers
	paramNameHdrVal, err := GetHeader(string(ParamName), headersMap)
	if err != nil {
		return EmptyString, err
	}

	// 2. check for the `paramNameHdrVal` in resolve specs List
	paramLocVal, err := resolveSpecs.GetParamLocation(paramNameHdrVal)
	if err != nil {
		return EmptyString, err
	}
	fmt.Printf("paramLocation: %s\n", paramLocVal)

	switch {
	case paramLocVal == InHeader:
		//use the value of `paramName` from header and check for it
		//in headers to fetch value to be used for discovery API key
		keyValue, err := GetHeader(paramNameHdrVal, headersMap)
		if err != nil {
			return EmptyString, err
		}
		return keyValue, nil
	}

	fmt.Printf("ParamLocation value not found in headers\n")
	return EmptyString, ErrorInvalidParamLocation
}

func (rsps ResolveSpecs) GetParamLocation(paramNameVal string) (string, error) {
	for _, spec := range rsps {
		if strings.EqualFold(spec.ParamName, paramNameVal) {
			return spec.ParamLocation, nil
		}
	}
	return EmptyString, ErrorParamLocationNotFound
}

func main() {
	fmt.Println("Hello World!")
}
