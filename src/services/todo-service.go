package services

import (
	"errors"
	"time"

	"geekible.todolist/src/config"
	"geekible.todolist/src/domain"
	"gorm.io/gorm"
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

func (s *ToDoService) Update(entity domain.ToDoEntity) error {
	if err := s.isEntityValid(entity); err != nil {
		return err
	}

	if err := s.config.Db.Save(&entity).Error; err != nil {
		s.config.Logger.Errorf("error updating todo with id: %d with error: %v", entity.ID, err)
		return err
	}

	return nil
}

func (s *ToDoService) Delete(entity domain.ToDoEntity) error {
	if err := s.config.Db.Delete(&entity).Error; err != nil {
		s.config.Logger.Errorf("error deleting todo with id: %d with error: %v", entity.ID, err)
		return err
	}

	return nil
}

func (s *ToDoService) GetById(id, userId uint) (domain.ToDoEntity, error) {
	var entity domain.ToDoEntity

	if err := s.config.Db.First(&entity, "id = ? and user_id = ?", id, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ToDoEntity{}, errNoMatchingRecords
		}

		s.config.Logger.Errorf("error attempting to find todo with id: %d with error: %v", id, err)
		return domain.ToDoEntity{}, err
	}

	return entity, nil
}

func (s *ToDoService) GetByUserId(userId uint, startAt, pageSize int) ([]domain.ToDoEntity, error) {
	var entities []domain.ToDoEntity

	if err := s.config.Db.
		Where("user_id = ?", userId).
		Order("due_date desc").
		Find(&entities).
		Limit(pageSize).
		Offset(startAt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []domain.ToDoEntity{}, errNoMatchingRecords
		}

		s.config.Logger.Errorf("error attempting to find todos for user id: %d with error: %v", userId, err)
		return []domain.ToDoEntity{}, err
	}

	return entities, nil
}
