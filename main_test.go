package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {

	db, _ := redismock.NewClientMock()

	handler := &handler{Client: db}
	router := handler.setupRouter()

	w1 := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w1, req)
	assert.Equal(t, 200, w1.Code)
	assert.Equal(t, "pong", w1.Body.String())

	w2 := httptest.NewRecorder()
	data := `{"name": "rams", "age":35, "books": ["solid design","galaxy", "avathar"] }`
	req2, _ := http.NewRequest("POST", "/admin", strings.NewReader(data))
	req2.Header.Add("Authorization", "Basic Zm9vOmJhcg==")
	req2.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)
	assert.Equal(t, 200, w2.Code)
	assert.Regexp(t, `[a-z]+`, w2.Body.String())

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/user/rams", nil)
	router.ServeHTTP(w3, req3)
	assert.Equal(t, 200, w3.Code)
	assert.Regexp(t, `[a-z]+`, w3.Body.String())

	//db, mock := redismock.NewClientMock()
	// key := "rams"
	// mock.ExpectGet(key).RedisNil()
	// mock.Regexp().ExpectSet(key, `[a-z]+`, 30*time.Minute).SetErr(errors.New("FAIL"))
	// mock.Regexp().ExpectGet(key).SetErr(errors.New("FAIL"))
	// if err := mock.ExpectationsWereMet(); err != nil {
	// 	t.Error(err)
	// }

}
