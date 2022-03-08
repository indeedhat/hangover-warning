package store

import (
	"fmt"
	"time"

	"github.com/indeedhat/hangover-warning/internal/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
	UpdatedAt time.Time `gorm:"index" json:"updated_at"`
}

// Connect to the database
//
// This will create the db if it does not exist
func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.Get(env.DbUser),
		env.Get(env.DbPass),
		env.Get(env.DbHost),
		env.GetFallback(env.DbPort, "3306"),
		env.Get(env.DbName),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		QueryFields:            true,
	})
	if err != nil {
		return db, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDb.SetConnMaxLifetime(time.Hour * 8)

	return db, nil
}

// Migrate schema changes on the database
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Quote{},
	)
}
