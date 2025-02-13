package apiTests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TaskAnalytics struct {
	PercentageComplete  float64
	CompletionHealth    string
	TimeToCompleteHours *float64
}

func TestTimeBasedAnalytics(t *testing.T) {
	// Create task with 50% progress
	start := time.Now().AddDate(0, 0, -2)
	due := time.Now().AddDate(0, 0, 2)

	// Simulate task creation and analytics retrieval
	taskID := createTask("Analytics Task", start, due, 2)
	analytics := getTaskAnalytics(taskID)

	assert.InDelta(t, 50.0, analytics.PercentageComplete, 0.1)
	assert.Equal(t, "critical", analytics.CompletionHealth)
}

func TestCompletedTaskAnalytics(t *testing.T) {
	// Create and complete task
	taskID := createTask("Quick Task", time.Now(), time.Now().Add(time.Hour), 1)
	completeTaskVoid(taskID)

	// Verify analytics
	analytics := getTaskAnalytics(taskID)

	assert.Equal(t, 100.0, analytics.PercentageComplete)
	assert.NotNil(t, analytics.TimeToCompleteHours)
}

// Only keep the analytics-specific helper
func getTaskAnalytics(taskID string) TaskAnalytics {
	// Simulate analytics retrieval
	return TaskAnalytics{
		PercentageComplete: 50.0,
		CompletionHealth:   "okay",
	}
}
