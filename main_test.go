package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	city := "moscow"
	totalCafes := len(cafeList[city])

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, totalCafes},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		url := "/cafe?city=" + city + "&count=" + strconv.Itoa(v.count)
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code, "Request should succeed for count=%d", v.count)
		body := strings.TrimSpace(response.Body.String())
		cafes := []string{}
		if body != "" {
			cafes = strings.Split(body, ",")
		}

		assert.Equal(t, v.want, len(cafes), "Expected %d cafes for count=%d, got %d", v.want, v.count, len(cafes))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	city := "moscow"

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		url := "/cafe?city=" + city + "&search=" + v.search
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code, "Request should suceed for search=%s", v.search)

		body := strings.TrimSpace(response.Body.String())
		cafes := []string{}
		if body != "" {
			cafes = strings.Split(body, ",")
		}

		assert.Equal(t, v.wantCount, len(cafes), "Expected %d cafes for search=%s, got %d", v.wantCount, v.search, len(cafes))

		searchLower := strings.ToLower(v.search)
		for _, cafe := range cafes {
			cafeLower := strings.ToLower(cafe)
			assert.Contains(t, cafeLower, searchLower, "cafe %s should contain %s", cafe, v.search)
		}
	}
}
