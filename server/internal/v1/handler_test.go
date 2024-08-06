package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maximilianhagelstam/speek/internal/mock"
)

func TestGetPostsHandler(t *testing.T) {
	r := mock.NewRepository()
	h := NewHandler(r)

	server := httptest.NewServer(http.HandlerFunc(h.GetPostsHandler))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %v", resp.Status)
	}
}

func TestCreatePostHandler(t *testing.T) {
	r := mock.NewRepository()
	h := NewHandler(r)

	server := httptest.NewServer(http.HandlerFunc(h.CreatePostHandler))
	defer server.Close()

	postData := map[string]string{"caption": "Test Post"}
	jsonData, _ := json.Marshal(postData)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %v", resp.Status)
	}
}
