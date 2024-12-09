package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config interface {
	DSN() string
}

func New(config Config) gorm.Dialector {
	return postgres.Open(config.DSN())
}
