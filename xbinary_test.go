package xbinary

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXBinaryHeader(t *testing.T) {
	a := assert.New(t)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	a.NoError(err)
	a.NotNil(req)
	res := httptest.NewRecorder()
	xb := &XBinary{Header: "test-x-binary-sum", ErrorHeader: "test-x-binary-err"}
	xb.ServeHTTP(res, req, func(rw http.ResponseWriter, r *http.Request) {})
	a.NotEmpty(res.Header().Get(xb.Header))
	a.Empty(res.Header().Get(xb.ErrorHeader))
}
