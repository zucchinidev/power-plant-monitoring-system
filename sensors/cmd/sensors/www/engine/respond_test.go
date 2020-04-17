package engine_test

import (
	"encoding/json"
	"errors"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/cmd/sensors/www/engine"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeReq struct {
	r *http.Request
	w *httptest.ResponseRecorder
}

func TestRespond_WithErrorResponse(t *testing.T) {
	f := fakeRequest()
	code := http.StatusInternalServerError
	res := new(struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	})
	e := errors.New("invalid operation")
	engine.Respond(f.w, f.r, code, e)
	shouldCheckTheCorrectStatusCode(t, f.w, code)
	shouldUnmarshalTheResponse(t, f.w, &res)

	if res.Status != http.StatusText(code) {
		t.Errorf("error in status property when marshaling the response body: got %s - want %s", res.Status, http.StatusText(code))
	}

	if res.Error != e.Error() {
		t.Errorf("error in error property when marshaling the response body: got %s - want %s", res.Error, e.Error())
	}
}

func TestRespond_WithCorrectResponse(t *testing.T) {
	f := fakeRequest()
	code := http.StatusOK
	res := new(struct {
		Res string `json:"response"`
	})
	fakeRes := "fake response"
	res.Res = fakeRes
	engine.Respond(f.w, f.r, code, res)
	shouldCheckTheCorrectStatusCode(t, f.w, code)
	shouldUnmarshalTheResponse(t, f.w, &res)

	if res.Res != fakeRes {
		t.Errorf("error in Res property when marshaling the response body: got %s - want %s", res.Res, fakeRes)
	}
}

func fakeRequest() fakeReq {
	return fakeReq{
		r: httptest.NewRequest("GET", "/foo", nil),
		w: httptest.NewRecorder(),
	}
}

func shouldUnmarshalTheResponse(t *testing.T, w *httptest.ResponseRecorder, res interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Errorf("error body unmarshaling %s", err)
	}
}

func shouldCheckTheCorrectStatusCode(t *testing.T, w *httptest.ResponseRecorder, code int) {
	if w.Code != code {
		t.Errorf("unexpected respond operation, status code not correct: got %d - want %d", w.Code, code)
	}
}
