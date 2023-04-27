package producer

// TODO: can't get a raw endpoint test to work w/ "encore test"
//import (
//	"encoding/json"
//	"github.com/stretchr/testify/assert"
//	"io"
//	"net/http/httptest"
//	"strings"
//	"testing"
//)
//
//func handleTestRequest(body io.Reader) (map[string]string, error) {
//	w := httptest.NewRecorder()
//	req := httptest.NewRequest("GET", "/oi", body)
//
//	Oi(w, req)
//
//	b, _ := io.ReadAll(w.Body)
//	resp := map[string]string{}
//	err := json.Unmarshal(b, &resp)
//	if err != nil {
//		return nil, err
//	}
//
//	return resp, nil
//}
//
//func TestOi(t *testing.T) {
//	var resp any
//	var err error
//	var body io.Reader
//
//	resp, err = handleTestRequest(nil)
//	assert.Nil(t, err)
//	assert.Equal(
//		t,
//		map[string]string{
//			"response_type": "in_channel",
//			"text":          "Oi what mate!?",
//		},
//		resp,
//	)
//
//	body = strings.NewReader("text=hello")
//	resp, err = handleTestRequest(body)
//	assert.Equal(
//		t,
//		map[string]string{
//			"response_type": "in_channel",
//			"text":          "Oi what mate!?",
//		},
//		resp,
//	)
//}
