package apiTests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
)

func TestSystemStats(t *testing.T) {
	// Reset database before test
	resp, err := resetDatabase()
	require.Equal(t, 200, resp.StatusCode)
	require.NoError(t, err)
	
	// Test with random number of tasks (between 5 and 10)
	numTasks := 15 + rand.Intn(100)
	taskIDs := make([]string, numTasks)
	expectedPriorities := make(map[int]int)
	
	// Create tasks with random priorities and times
	now := time.Now()
	for i := 0; i < numTasks; i++ {
		priority := 1 + rand.Intn(5)
		startOffset := time.Duration(1.0+rand.Float64()*72.0)*time.Hour
		duration := time.Duration(1.0+rand.Float64()*72.0)*time.Hour
		
		startTime := now.Add(-startOffset)
		endTime := now.Add(duration)
		
		taskIDs[i] = createTask(fmt.Sprintf("Task %d", i+1), startTime, endTime, priority)
		require.NotEmpty(t, taskIDs[i], "Failed to create task")
		expectedPriorities[priority]++
	}

	// Randomly complete some tasks
	completedCount := 0
	completionTimes := make([]float64, 0)
	for i, taskID := range taskIDs {
		if rand.Float32() < 0.5 { // 50% chance to complete each task
			resp = completeTask(taskID)
			require.Equal(t, 200, resp.StatusCode, "Failed to complete task %d", i+1)
			completedCount++
			completionTimes = append(completionTimes, 1.0) // Simplified completion time
		}
	}

	// Get and verify system stats
	stats, err := getSystemStats()
	require.NoError(t, err, "Failed to get system stats")

	expectedRate := float64(completedCount) / float64(numTasks) * 100
	assert.Equal(t, expectedRate, stats.CompletionRate, "Unexpected completion rate")
	assert.Equal(t, expectedPriorities, stats.TasksByPriority, "Unexpected tasks by priority")
	
	if completedCount > 0 {
		assert.Greater(t, stats.AvgCompletionHours, 0.0, "Average completion hours should be positive")
	} else {
		assert.Equal(t, 0.0, stats.AvgCompletionHours, "Average completion hours should be zero")
	}
}
