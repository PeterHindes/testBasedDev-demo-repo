package apiTests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemStats(t *testing.T) {
	// Create mix of tasks
	taskID1 := createTask("Task 1", time.Now(), time.Now().AddDate(0, 0, 1), 1)
	require.NotEmpty(t, taskID1, "Failed to create first task")

	taskID2 := createTask("Task 2", time.Now(), time.Now().Add(2*time.Hour), 2)
	require.NotEmpty(t, taskID2, "Failed to create second task")

	resp := completeTask(taskID2)
	require.Equal(t, 200, resp.StatusCode, "Failed to complete task")

	// Get system stats
	stats, err := getSystemStats()
	require.NoError(t, err, "Failed to get system stats")

	assert.Equal(t, 50.0, stats.CompletionRate)
	assert.NotNil(t, stats.AvgCompletionHours)
	assert.Equal(t, map[int]int{1: 1, 2: 1}, stats.TasksByPriority)
}
