package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServiceConfig struct {
	Port   int
	Logger *zap.SugaredLogger
	Db     *gorm.DB
}

func InitServiceConfig() *ServiceConfig {
	return &ServiceConfig{}
}

func (s *ServiceConfig) BuildConfig() (*ServiceConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return &ServiceConfig{}, err
	}

	logfile, err := os.Create("todo.log")
	if err != nil {
		return &ServiceConfig{}, err
	}
	logger := s.buildLogger(logfile)

	db, err := s.buildDatbaseConnection(
		viper.GetString("service.environment"),
		viper.GetString("database.host"),
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetInt("database.port"))

	if err != nil {
		return &ServiceConfig{}, err
	}

	return &ServiceConfig{
		Port:   viper.GetInt("service.port"),
		Logger: logger,
		Db:     db,
	}, nil
}

func (s *ServiceConfig) BuilderMux() *chi.Mux {
	requestLogger := httplog.NewLogger("todo-service", httplog.Options{
		JSON:     true,
		Concise:  true,
		LogLevel: "debug",
	})

	mux := chi.NewMux()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	mux.Use(httplog.RequestLogger(requestLogger))
	mux.Use(middleware.Compress(5, "application/json"))
	mux.Use(middleware.AllowContentType("application/json", "text/xml"))
	mux.Use(middleware.NoCache)
	mux.Use(middleware.StripSlashes)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	return mux

}

func (s *ServiceConfig) buildLogger(f *os.File) *zap.SugaredLogger {
	pe := zap.NewProductionEncoderConfig()
	fileEncoder := zapcore.NewJSONEncoder(pe)
	pe.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	level := zap.InfoLevel

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(f), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	logger := zap.New(core)
	return logger.Sugar()

}

func (s *ServiceConfig) buildDatbaseConnection(env, host, username, password, dbName string, port int) (*gorm.DB, error) {

	if env == "dev" {
		host = "localhost"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		host,
		username,
		password,
		dbName,
		port)

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger:      dbLogger,
		NowFunc:     time.Now,
		PrepareStmt: true,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
