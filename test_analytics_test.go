package apiTests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeBasedAnalytics(t *testing.T) {
	// Reset database before test
	resp, err := resetDatabase()
	require.Equal(t, 200, resp.StatusCode)
	require.NoError(t, err)

	// Create task with 50% progress
	start := time.Now().AddDate(0, 0, -2)
	due := time.Now().AddDate(0, 0, 2)

	// Simulate task creation and analytics retrieval
	taskID := createTask("Analytics Task", start, due, 2)
	analytics, err := getTaskAnalytics(taskID)
	if err != nil {
		t.Fatal(err)
	}

	assert.InDelta(t, 50.0, analytics.PercentageComplete, 0.1)
	assert.Equal(t, "critical", analytics.CompletionHealth)
}

func TestCompletedTaskAnalytics(t *testing.T) {
	// Reset database before test
	resp, err := resetDatabase()
	require.Equal(t, 200, resp.StatusCode)
	require.NoError(t, err)

	// Create and complete task
	taskID := createTask("Quick Task", time.Now(), time.Now().Add(time.Hour), 1)
	completeTask(taskID)

	// Verify analytics
	analytics, err := getTaskAnalytics(taskID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 100.0, analytics.PercentageComplete)
	assert.NotNil(t, analytics.TimeToCompleteHours)
}
