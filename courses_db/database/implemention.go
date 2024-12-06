package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/antoniofmoliveira/courses/db/configs"
	"github.com/antoniofmoliveira/courses/db/database/mariadb"
	"github.com/antoniofmoliveira/courses/db/database/sqlite"
)

type DBImplementation struct {
	db                 *sql.DB
	CategoryRepository CategoryRepositoryInterface
	CourseRepository   CourseRepositoryInterface
	UserRepository     UserRepositoryInterface
}

var dbi *DBImplementation

func GetDBImplementation() *DBImplementation {
	if dbi != nil {
		return dbi
	}

	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	if cfg.DBDriver == "mysql" {
		conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
		db, err := sql.Open(cfg.DBDriver, conn)
		if err != nil {
			log.Fatalf("failed to open database: %v", err)
		}
		dbi = &DBImplementation{
			db:                 db,
			CategoryRepository: mariadb.NewCategoryRepository(db),
			CourseRepository:   mariadb.NewCourseRepository(db),
			UserRepository:     mariadb.NewUserRepository(db),
		}
		return dbi
	}
	if cfg.DBDriver == "sqlite3" {
		conn := fmt.Sprintf("%s.db", cfg.DBName)
		db, err := sql.Open(cfg.DBDriver, conn)
		if err != nil {
			log.Fatalf("failed to open database: %v", err)
		}
		dbi = &DBImplementation{
			db:                 db,
			CategoryRepository: sqlite.NewCategoryRepository(db),
			CourseRepository:   sqlite.NewCourseRepository(db),
			UserRepository:     sqlite.NewUserRepository(db),
		}
		return dbi
	}
	if cfg.DBDriver == "mongodb" {
		// TODO!
	}

	return nil
}
