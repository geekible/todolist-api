package services

import (
	"errors"
	"time"

	"geekible.todolist/src/config"
	"geekible.todolist/src/domain"
)

type ToDoService struct {
	config *config.ServiceConfig
}

func InitToDoService(config *config.ServiceConfig) *ToDoService {
	return &ToDoService{
		config: config,
	}
}

func (s *ToDoService) isEntityValid(entity domain.ToDoEntity) error {
	if len(entity.Title) <= 0 {
		return errors.New("title cannot be empty")
	}

	if len(entity.Description) <= 0 {
		return errors.New("description cannot be empty")
	}

	if entity.DueDate.Before(time.Now()) {
		return errors.New("due date cannot be in the past")
	}

	return nil
}

func (s *ToDoService) Add(entity domain.ToDoEntity) (domain.ToDoEntity, error) {
	if err := s.isEntityValid(entity); err != nil {
		return entity, err
	}

	if err := s.config.Db.Create(&entity).Error; err != nil {
		s.config.Logger.Errorf("error inserting todo: %v", err)
		return entity, err
	}

	return entity, nil
}
