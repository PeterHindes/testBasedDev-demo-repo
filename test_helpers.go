package apiTests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8000"

// Response represents an HTTP response structure
type Response struct {
	StatusCode int
	Status     string
	Body       []byte
}

// TaskRequest represents the request body for creating a task
type TaskRequest struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	DueDate   time.Time `json:"due_date"`
	Priority  int       `json:"priority"`
}

// TaskResponse represents the API response for a task
type TaskResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	StartDate   time.Time  `json:"start_date"`
	DueDate     time.Time  `json:"due_date"`
	Priority    int        `json:"priority"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// TaskAnalytics represents task analytics data
type TaskAnalytics struct {
	DaysUntilDue        int      `json:"days_until_due"`
	CompletionHealth    string   `json:"completion_health"`
	PercentageComplete  float64  `json:"percentage_complete"`
	TimeToCompleteHours float64 `json:"time_to_complete_hours,omitempty"`
}

// SystemStats represents system-wide statistics
type SystemStats struct {
	AvgCompletionHours float64    `json:"avg_completion_hours,omitempty"`
	CompletionRate     float64     `json:"completion_rate"`
	TasksByPriority    map[int]int `json:"tasks_by_priority"`
}

func makeRequest(method, path string, body interface{}) (Response, error) {
	var jsonBody []byte
	var err error

	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return Response{}, err
		}
	}

	req, err := http.NewRequest(method, baseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return Response{}, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	return Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       respBody,
	}, nil
}

func createTaskWithResponse(name string, start, due time.Time, priority int) Response {
	task := TaskRequest{
		Name:      name,
		StartDate: start,
		DueDate:   due,
		Priority:  priority,
	}

	resp, err := makeRequest("POST", "/tasks", task)
	if err != nil {
		return Response{StatusCode: 500}
	}
	return resp
}

func createTask(name string, start, due time.Time, priority int) string {
	resp := createTaskWithResponse(name, start, due, priority)
	if resp.StatusCode != 201 {
		return ""
	}
	var task TaskResponse
	if err := json.Unmarshal(resp.Body, &task); err != nil {
		return ""
	}
	return task.ID
}

func completeTask(taskID string) Response {
	resp, err := makeRequest("PUT", fmt.Sprintf("/tasks/%s/complete", taskID), nil)
	if err != nil {
		return Response{StatusCode: 500}
	}
	return resp
}

func uncompleteTask(taskID string) Response {
	resp, err := makeRequest("PUT", fmt.Sprintf("/tasks/%s/uncomplete", taskID), nil)
	if err != nil {
		return Response{StatusCode: 500}
	}
	return resp
}

func getTaskAnalytics(taskID string) (TaskAnalytics, error) {
	resp, err := makeRequest("GET", fmt.Sprintf("/tasks/%s/analytics", taskID), nil)
	if err != nil {
		return TaskAnalytics{}, err
	}

	var analytics TaskAnalytics
	if err := json.Unmarshal(resp.Body, &analytics); err != nil {
		return TaskAnalytics{}, err
	}
	return analytics, nil
}

func getSystemStats() (SystemStats, error) {
	resp, err := makeRequest("GET", "/analytics", nil)
	if err != nil {
		return SystemStats{}, err
	}

	var response struct {
		SystemStats SystemStats `json:"system_stats"`
	}
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return SystemStats{}, err
	}
	return response.SystemStats, nil
}

// resetDatabase clears all data from the tasks table
func resetDatabase() (Response, error) {
	return makeRequest("DELETE", "/reset?confirm=true", nil)
}
