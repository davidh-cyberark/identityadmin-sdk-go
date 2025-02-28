package identity

import (
	"fmt"
	"net/http"
	"strings"
)

func headersToString(header http.Header) string {
	var sb strings.Builder
	for key, values := range header {
		sb.WriteString(fmt.Sprintf("%s: %s\n", key, strings.Join(values, ",")))
	}
	return sb.String()
}
