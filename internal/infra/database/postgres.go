package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/raulaguila/go-template/pkg/helpers"
	"github.com/raulaguila/go-template/pkg/pgerror"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func pgConnect(dbName string) *gorm.DB {
	uri := fmt.Sprintf("host=%s user=%s password=%s dbname=%v port=%s sslmode=disable TimeZone=%v", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"), dbName, os.Getenv("POSTGRES_PORT"), time.Local.String())
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NowFunc: func() time.Time {
			return time.Now()
		},
		PrepareStmt: true,
	})
	helpers.PanicIfErr(err)

	return db
}

func createDataBase() {
	db := pgConnect("postgres")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	con, err := db.WithContext(ctx).DB()
	helpers.PanicIfErr(err)
	defer con.Close()

	if err := db.Exec(fmt.Sprintf("CREATE DATABASE %v;", os.Getenv("POSTGRES_BASE"))).Error; err != nil {
		switch pgerror.HandlerError(err) {
		case pgerror.ErrDatabaseAlreadyExists:
		default:
			helpers.PanicIfErr(err)
		}
	}
}

func ConnectPostgresDB() (*gorm.DB, error) {
	createDataBase()

	db := pgConnect(os.Getenv("POSTGRES_BASE"))
	postgresdb := db.WithContext(context.Background())
	if err := postgresdb.Exec("CREATE EXTENSION IF NOT EXISTS unaccent;").Error; err != nil {
		return nil, err
	}

	if err := postgresdb.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return nil, err
	}

	autoMigrate(postgresdb)
	go createDefaults(postgresdb)
	return postgresdb, nil
}
