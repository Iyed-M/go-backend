package handlers

import (
	"github.com/Iyed-M/go-backend/database"
)

type ApiConfig struct {
	Db                   *database.DB
	FileServerHits       int
	JWTSecret            string
	RefreshTokenLifeDays int
	AccessTokenLifeHours int
}
