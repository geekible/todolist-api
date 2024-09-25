package config

import (
	"geekible.todolist/src/domain"
	"go.uber.org/zap/zapcore"
)

type Migration struct {
	serviceConfig *ServiceConfig
}

func InitMigration(serviceConfig *ServiceConfig) *Migration {
	return &Migration{
		serviceConfig: serviceConfig,
	}
}

func (m *Migration) DoMigration() {
	m.serviceConfig.Logger.Log(zapcore.InfoLevel, "migrating ToDoEntity")
	if err := m.serviceConfig.Db.AutoMigrate(&domain.ToDoEntity{}); err != nil {
		m.serviceConfig.Logger.Errorf("error migrating table ToDoEntity: %v", err)
	}
	m.serviceConfig.Logger.Log(zapcore.InfoLevel, "migrating ToDoEntity completed")
}
