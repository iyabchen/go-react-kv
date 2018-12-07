package web

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/iyabchen/go-react-kv/server/model"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWebController(t *testing.T) {
	// anything not match the router setting returns 405
	ctx := context.Background()
	notfoundErr := fmt.Errorf("not exist")

	testPair, err := model.NewPair("a", "b")
	if err != nil {
		t.Fatal(err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storage := NewMockPairRepo(ctrl)
	// server does not run, so a random port should be OK.
	srv, err := NewWeb(&Options{":12345", storage})
	if err != nil {
		t.Fatal(err)
	}

	// Test API.create
	storage.EXPECT().Create(ctx, gomock.Any()).
		Return(nil)
	req := createTestRequest(http.MethodPost, "/pair", `{"key":"a", "value":"b"}`)
	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	resp := w.Result()
	defer resp.Body.Close()

	req = createTestRequest(http.MethodPost, "/pair", ``)
	w = httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp = w.Result()
	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.Regexp(t, regexp.MustCompile(`\{"error":".+"\}`), string(data))

	req = createTestRequest(http.MethodPost, "/pair", `{}`)
	w = httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp = w.Result()
	data, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.Regexp(t, regexp.MustCompile(`\{"error":".+"\}`), string(data))

	// Test API.getOne
	storage.EXPECT().GetOne(ctx, gomock.Any()).
		Return(testPair, nil)
	req = createTestRequest(http.MethodGet, "/pair/test", "")
	w = httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	resp = w.Result()
	defer resp.Body.Close()

	storage.EXPECT().GetOne(ctx, gomock.Any()).
		Return(nil, notfoundErr)
	req = createTestRequest(http.MethodGet, "/pair/test", "")
	w = httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	resp = w.Result()
	data, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.Regexp(t, regexp.MustCompile(`\{"error":".+"\}`), string(data))

	// Test API.getAll
	storage.EXPECT().GetAll(ctx).
		Return([]*model.Pair{testPair}, nil)
	req = createTestRequest(http.MethodGet, "/pair", "")
	w = httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	resp = w.Result()
	defer resp.Body.Close()

	// Test API.deleteOne
	storage.EXPECT().DeleteOne(ctx, gomock.Any()).
		Return(nil)
	req = createTestRequest(http.MethodDelete, "/pair/test", "")
	w = httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	resp = w.Result()
	data, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	storage.EXPECT().DeleteOne(ctx, gomock.Any()).
		Return(notfoundErr)
	req = createTestRequest(http.MethodDelete, "/pair/test", ``)
	w = httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	resp = w.Result()
	data, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.Regexp(t, regexp.MustCompile(`\{"error":".+"\}`), string(data))

	// Test API.deleteAll
	storage.EXPECT().DeleteAll(ctx).Return(nil)
	req = createTestRequest(http.MethodGet, "/reset", "")
	w = httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	resp = w.Result()
	data, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	// Test API.update
	storage.EXPECT().Update(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	req = createTestRequest(http.MethodPut, "/pair/test", `{"key":"a", "value":"b"}`)
	w = httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	resp = w.Result()
	data, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	storage.EXPECT().Update(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
		Return(notfoundErr)
	req = createTestRequest(http.MethodPut, "/pair/test", `{"key":"a", "value":"b"}`)
	w = httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	resp = w.Result()
	data, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.Regexp(t, regexp.MustCompile(`\{"error":".+"\}`), string(data))

}

func createTestRequest(method string, uri string, body string) *http.Request {
	return httptest.NewRequest(method, uri, strings.NewReader(body))
}
