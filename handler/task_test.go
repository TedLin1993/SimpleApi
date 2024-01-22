package handler

import (
	"Gogolook_test/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	r := gin.Default()
	r.POST("/tasks", CreateTask)

	newTask := model.Task{
		Name:   "Test Task",
		Status: new(model.Status),
	}

	jsonTask, err := json.Marshal(newTask)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonTask))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdTask model.Task
	err = json.Unmarshal(w.Body.Bytes(), &createdTask)
	assert.NoError(t, err)

	assert.Equal(t, newTask.Name, createdTask.Name)
	assert.Equal(t, newTask.Status, createdTask.Status)
}

func TestGetTasks(t *testing.T) {
	existingTask := model.Task{
		ID:     "1",
		Name:   "Existing Task",
		Status: new(model.Status),
	}
	cache.Store(existingTask.ID, existingTask)

	r := gin.Default()
	r.GET("/tasks", GetTasks)

	req, err := http.NewRequest("GET", "/tasks", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var tasks []model.Task
	err = json.Unmarshal(w.Body.Bytes(), &tasks)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(tasks))

}
func TestUpdateTask(t *testing.T) {
	existingTask := model.Task{
		ID:     "1",
		Name:   "Existing Task",
		Status: new(model.Status),
	}
	*existingTask.Status = model.Incomplete
	cache.Store(existingTask.ID, existingTask)

	r := gin.Default()
	r.PUT("/tasks/:id", UpdateTask)

	updatedTask := model.Task{
		ID:     existingTask.ID,
		Name:   existingTask.Name,
		Status: new(model.Status),
	}
	*updatedTask.Status = model.Completed
	jsonTask, err := json.Marshal(updatedTask)
	assert.NoError(t, err)
	req, err := http.NewRequest("PUT", fmt.Sprintf("/tasks/%s", updatedTask.ID), bytes.NewBuffer(jsonTask))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &updatedTask)
	assert.NoError(t, err)

	assert.Equal(t, existingTask.ID, updatedTask.ID)
	assert.Equal(t, existingTask.Name, updatedTask.Name)
	assert.NotEqual(t, existingTask.Status, updatedTask.Status)
}

func TestDeleteTask(t *testing.T) {
	existingTask := model.Task{
		ID:     "1",
		Name:   "Existing Task",
		Status: new(model.Status),
	}
	*existingTask.Status = model.Incomplete
	cache.Store(existingTask.ID, existingTask)

	r := gin.Default()
	r.DELETE("/tasks/:id", DeleteTask)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/tasks/%s", existingTask.ID), nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	assert.Contains(t, w.Body.String(), "Task deleted successfully")
}

func TestCreateTaskInvalidInput(t *testing.T) {
	r := gin.Default()
	r.POST("/tasks", CreateTask)

	// Task with missing name
	invalidTask := model.Task{
		Status: new(model.Status),
	}

	jsonTask, err := json.Marshal(invalidTask)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonTask))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestUpdateTaskNotFound(t *testing.T) {
	r := gin.Default()
	r.PUT("/tasks/:id", UpdateTask)

	// Task not found
	nonExistingID := "100"
	updatedTask := model.Task{
		ID:     nonExistingID,
		Name:   "Updated Task",
		Status: new(model.Status),
	}
	*updatedTask.Status = model.Completed

	jsonTask, err := json.Marshal(updatedTask)
	assert.NoError(t, err)
	req, err := http.NewRequest("PUT", fmt.Sprintf("/tasks/%s", nonExistingID), bytes.NewBuffer(jsonTask))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}