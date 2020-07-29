package httpflow

import "fmt"

//HTTPHeaders deals with http.header and values as a map
type HTTPHeaders map[string][]string

func (h HTTPHeaders) String() (s string) {
	s = ""
	for key, values := range h {
		for _, value := range values {
			s = fmt.Sprintf("%s\n%s: %s", s, key, value)
		}
	}
	return s
}
