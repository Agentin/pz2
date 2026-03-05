package service

import (
	"fmt"
	"sync"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	Done        bool   `json:"done"`
}

type TaskService struct {
	mu     sync.RWMutex
	tasks  map[string]Task
	nextID int
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks:  make(map[string]Task),
		nextID: 1,
	}
}

func (s *TaskService) Create(title, description, dueDate string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := fmt.Sprintf("t_%03d", s.nextID)
	s.nextID++
	task := Task{
		ID:          id,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Done:        false,
	}
	s.tasks[id] = task
	return task
}

func (s *TaskService) GetAll() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	return result
}

func (s *TaskService) GetByID(id string) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	return task, ok
}

func (s *TaskService) Update(id string, updates map[string]interface{}) (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[id]
	if !ok {
		return Task{}, false
	}
	// Применяем обновления
	if val, ok := updates["title"]; ok {
		task.Title = val.(string)
	}
	if val, ok := updates["description"]; ok {
		task.Description = val.(string)
	}
	if val, ok := updates["due_date"]; ok {
		task.DueDate = val.(string)
	}
	if val, ok := updates["done"]; ok {
		task.Done = val.(bool)
	}
	s.tasks[id] = task
	return task, true
}

func (s *TaskService) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; ok {
		delete(s.tasks, id)
		return true
	}
	return false
}