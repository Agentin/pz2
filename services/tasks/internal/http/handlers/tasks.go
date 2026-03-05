package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/student/tech-ip-sem2/services/tasks/internal/service"
)

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type taskResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	Done        bool   `json:"done"`
}

type taskListItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func CreateTaskHandler(svc *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if req.Title == "" {
			http.Error(w, "title is required", http.StatusBadRequest)
			return
		}
		task := svc.Create(req.Title, req.Description, req.DueDate)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(taskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Done:        task.Done,
		})
	}
}

func GetTasksHandler(svc *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks := svc.GetAll()
		list := make([]taskListItem, 0, len(tasks))
		for _, t := range tasks {
			list = append(list, taskListItem{ID: t.ID, Title: t.Title, Done: t.Done})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	}
}

func GetTaskHandler(svc *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}
		task, ok := svc.GetByID(id)
		if !ok {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(taskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Done:        task.Done,
		})
	}
}

func UpdateTaskHandler(svc *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}
		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		// Проверяем, что обновляемые поля имеют корректный тип (упрощённо)
		updatedTask, ok := svc.Update(id, updates)
		if !ok {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(taskResponse{
			ID:          updatedTask.ID,
			Title:       updatedTask.Title,
			Description: updatedTask.Description,
			DueDate:     updatedTask.DueDate,
			Done:        updatedTask.Done,
		})
	}
}

func DeleteTaskHandler(svc *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}
		if !svc.Delete(id) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
