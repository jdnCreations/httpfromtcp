package headers

import (
	"errors"
	"regexp"
	"strings"
)

var validHeaderKeyPattern = regexp.MustCompile("^[A-Za-z0-9!#$%&'*+\\-.^_`|~]+$")


type Headers map[string]string

func (headers Headers) Get(key string) string {
  if val, ok := headers[strings.ToLower(key)]; ok {
    return val
  }
  return "" 
}

func NewHeaders() (headers Headers) {
	return make(Headers)
}

func (headers Headers) Override(key, value string) {
	key = strings.ToLower(key)
	headers[key] = value
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	dataToParse := string(data)
	if len(dataToParse) >= 2 && dataToParse[0] == '\r' && dataToParse[1] == '\n' {
		return 2, true, nil
	}

	crlfIndex := strings.Index(dataToParse, "\r\n")
	if crlfIndex == -1 {
		return 0, false, nil
	}

	line := dataToParse[:crlfIndex]

	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return 0, false, errors.New("invalid header format: missing colon")
	}

	key := line[:colonIndex]

	if len(key) <= 0 {
		return 0, false, errors.New("invalid header format: key must be min of 1 character")
	}

	if !validHeaderKeyPattern.MatchString(key) {
		return 0, false, errors.New("invalid header format: invalid character")
	}

	if strings.Contains(key, " ") {
		return 0, false, errors.New("invalid header format: space in key")
	}

	key = strings.ToLower(strings.TrimSpace(key))
	value := strings.TrimSpace(line[colonIndex+1:])

	val, ok := h[key]
	if ok {
		h[key] = val + ", " + value
	} else {
		h[key] = value
	}

	return crlfIndex + 2, false, nil
}