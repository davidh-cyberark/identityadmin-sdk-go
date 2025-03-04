package identity

import (
	"context"
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
func AddRequestHeaders(ctx context.Context, req *http.Request) error {
	headers := ctx.Value(HeadersKey).(map[string]string)
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	return nil
}
