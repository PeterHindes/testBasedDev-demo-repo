package apiTests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateValidTask(t *testing.T) {
	// Reset database before test
	resp, err := resetDatabase()
	require.Equal(t, 200, resp.StatusCode)
	require.NoError(t, err)

	response := createTaskWithResponse("Valid Task", time.Now(), time.Now().AddDate(0, 0, 5), 3)
	require.Equal(t, 201, response.StatusCode)

	var task TaskResponse
	err = json.Unmarshal(response.Body, &task)
	require.NoError(t, err)
	assert.NotEmpty(t, task.ID)
}

func TestCreateInvalidTasks(t *testing.T) {
	// Reset database before test
	resp, err := resetDatabase()
	require.Equal(t, 200, resp.StatusCode)
	require.NoError(t, err)

	startTime := time.Now()
	dueTime := startTime.AddDate(0, 0, 1)
	duePastTime := startTime.AddDate(0, 0, -1)

	invalidTasks := []struct {
		name     string
		start    time.Time
		due      time.Time
		priority int
	}{
		{"AB", startTime, dueTime, 3},               // Name too short
		{"Invalid Date", startTime, duePastTime, 2}, // Due date before start
		{"Bad Priority", startTime, dueTime, 0},     // Invalid priority
	}

	for _, task := range invalidTasks {
		response := createTaskWithResponse(task.name, task.start, task.due, task.priority)
		assert.Equal(t, 400, response.StatusCode, "Task should be invalid: %s", task.name)
	}
}

func TestCompleteTaskWorkflow(t *testing.T) {
	// Reset database before test
	resp, err := resetDatabase()
	require.Equal(t, 200, resp.StatusCode)
	require.NoError(t, err)

	// Create task
	taskID := createTask("Completable Task", time.Now(), time.Now().Add(time.Hour), 1)
	require.NotEmpty(t, taskID, "Failed to create task")

	// Complete task
	response := completeTask(taskID)
	require.Equal(t, 200, response.StatusCode)

	var task TaskResponse
	err = json.Unmarshal(response.Body, &task)
	require.NoError(t, err)
	assert.Equal(t, "completed", task.Status)

	// Uncomplete task
	response = uncompleteTask(taskID)
	require.Equal(t, 200, response.StatusCode)

	err = json.Unmarshal(response.Body, &task)
	require.NoError(t, err)
	assert.Equal(t, "pending", task.Status)
}
