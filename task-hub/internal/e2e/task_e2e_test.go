package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vagonaizer/workmate/task-hub/internal/app"
	"github.com/vagonaizer/workmate/task-hub/internal/config"
)

type createTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Deadline    time.Time `json:"deadline"`
}

type taskResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Deadline    string `json:"deadline"`
}

func TestTaskE2E(t *testing.T) {
	cfg := config.LoadConfig()
	application := app.NewApp(cfg)
	ts := httptest.NewServer(application.Engine)
	defer ts.Close()

	// 1. Создать задачу
	createReq := createTaskRequest{
		Title:       "E2E Test",
		Description: "desc",
		Priority:    "high",
		Deadline:    time.Now().Add(24 * time.Hour),
	}
	body, _ := json.Marshal(createReq)
	resp, err := http.Post(ts.URL+"/api/tasks", "application/json", bytes.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	respBody, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var created taskResponse
	_ = json.Unmarshal(respBody, &created)
	assert.Equal(t, createReq.Title, created.Title)
	assert.NotEmpty(t, created.ID)

	// 2. Получить задачу по id
	resp, err = http.Get(ts.URL + "/api/tasks/" + created.ID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	respBody, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	var got taskResponse
	_ = json.Unmarshal(respBody, &got)
	assert.Equal(t, created.ID, got.ID)

	// 3. Удалить задачу
	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/tasks/"+created.ID, nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// 4. Проверить, что задача удалена
	resp, err = http.Get(ts.URL + "/api/tasks/" + created.ID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
