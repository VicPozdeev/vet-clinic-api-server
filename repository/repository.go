package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"vet-clinic/config"
	"vet-clinic/logging"
)

type Repository interface {
	Model(value interface{}) *gorm.DB
	Select(query interface{}, args ...interface{}) *gorm.DB
	Omit(columns ...string) *gorm.DB
	Find(out interface{}, where ...interface{}) *gorm.DB
	Exec(sql string, values ...interface{}) *gorm.DB
	First(out interface{}, where ...interface{}) *gorm.DB
	Raw(sql string, values ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Updates(value interface{}) *gorm.DB
	Delete(value interface{}, where ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Preload(column string, conditions ...interface{}) *gorm.DB
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB
	ScanRows(rows *sql.Rows, result interface{}) error
	Transaction(fc func(tx Repository) error) (err error)
	Close() error
	DropTableIfExists(value interface{}) error
	AutoMigrate(value interface{}) error
}

type GormRepo struct {
	db *gorm.DB
}

func NewRepository(logger logging.Logger, conf *config.Config) *GormRepo {
	logger.Infof("Try database connection")
	db, err := connectDatabase(logger, conf)
	if err != nil {
		logger.Errorf("Failure database connection")
		os.Exit(config.ErrExitStatus)
	}
	logger.Infof("Success database connection, %s:%s", conf.Database.Host, conf.Database.Port)

	return &GormRepo{db: db}
}

const (
	// SQLITE represents SQLite3
	SQLITE = "sqlite3"
	// POSTGRES represents PostgreSQL
	POSTGRES = "postgres"
	// MYSQL represents MySQL
	MYSQL = "mysql"
)

func connectDatabase(logger logging.Logger, config *config.Config) (*gorm.DB, error) {
	var dsn string
	gormConfig := &gorm.Config{Logger: logger.Gorm()}

	switch config.Database.Dialect {
	case POSTGRES:
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			config.Database.Host,
			config.Database.Port,
			config.Database.Username,
			config.Database.Dbname,
			config.Database.Password)
		return gorm.Open(postgres.Open(dsn), gormConfig)
	case MYSQL:
		dsn = fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.Database.Username,
			config.Database.Password,
			config.Database.Host,
			config.Database.Dbname)
		return gorm.Open(mysql.Open(dsn), gormConfig)
	case SQLITE:
		return gorm.Open(sqlite.Open(config.Database.Host), gormConfig)
	default:
		return nil, errors.New("database dialect configuration not found")
	}
}

// Model specify the models you would like to run db operations
func (rep *GormRepo) Model(value interface{}) *gorm.DB {
	return rep.db.Model(value)
}

// Select specify fields that you want to retrieve from database when querying, by default, will select all fields;
func (rep *GormRepo) Select(query interface{}, args ...interface{}) *gorm.DB {
	return rep.db.Select(query, args...)
}

// Omit specify fields that you want to ignore when creating, updating and querying
func (rep *GormRepo) Omit(columns ...string) *gorm.DB {
	return rep.db.Omit(columns...)
}

// Find find records that match given conditions.
func (rep *GormRepo) Find(out interface{}, where ...interface{}) *gorm.DB {
	return rep.db.Find(out, where...)
}

// Exec exec given SQL using by gorm.DB.
func (rep *GormRepo) Exec(sql string, values ...interface{}) *gorm.DB {
	return rep.db.Exec(sql, values...)
}

// First returns first record that match given conditions, order by primary key.
func (rep *GormRepo) First(out interface{}, where ...interface{}) *gorm.DB {
	return rep.db.First(out, where...)
}

// Raw returns the record that executed the given SQL using gorm.DB.
func (rep *GormRepo) Raw(sql string, values ...interface{}) *gorm.DB {
	return rep.db.Raw(sql, values...)
}

// Create insert the value into database.
func (rep *GormRepo) Create(value interface{}) *gorm.DB {
	return rep.db.Create(value)
}

// Save update value in database, if the value doesn't have primary key, will insert it.
func (rep *GormRepo) Save(value interface{}) *gorm.DB {
	return rep.db.Save(value)
}

// Update update value in database
func (rep *GormRepo) Updates(value interface{}) *gorm.DB {
	return rep.db.Updates(value)
}

// Delete delete value match given conditions.
func (rep *GormRepo) Delete(value interface{}, where ...interface{}) *gorm.DB {
	return rep.db.Delete(value, where...)
}

// Where returns a new relation.
func (rep *GormRepo) Where(query interface{}, args ...interface{}) *gorm.DB {
	return rep.db.Where(query, args...)
}

// Preload preload associations with given conditions.
func (rep *GormRepo) Preload(column string, conditions ...interface{}) *gorm.DB {
	return rep.db.Preload(column, conditions...)
}

// Scopes pass current database connection to arguments `func(*DB) *DB`, which could be used to add conditions dynamically
func (rep *GormRepo) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	return rep.db.Scopes(funcs...)
}

// ScanRows scan `*sql.Rows` to give struct
func (rep *GormRepo) ScanRows(rows *sql.Rows, result interface{}) error {
	return rep.db.ScanRows(rows, result)
}

// Close close current db connection. If database connection is not an io.Closer, returns an error.
func (rep *GormRepo) Close() error {
	sqlDB, _ := rep.db.DB()
	return sqlDB.Close()
}

// DropTableIfExists drop table if it is exist
func (rep *GormRepo) DropTableIfExists(value interface{}) error {
	return rep.db.Migrator().DropTable(value)
}

// AutoMigrate run auto migration for given models, will only add missing fields, won't delete/change current data
func (rep *GormRepo) AutoMigrate(value interface{}) error {
	return rep.db.AutoMigrate(value)
}

// Transaction start a transaction as a block.
// If it is failed, will rollback and return error.
// If it is successed, will commit.
// ref: https://github.com/jinzhu/gorm/blob/master/main.go#L533
func (rep *GormRepo) Transaction(fc func(tx Repository) error) (err error) {
	panicked := true
	tx := rep.db.Begin()
	defer func() {
		if panicked || err != nil {
			tx.Rollback()
		}
	}()

	txrep := &GormRepo{}
	txrep.db = tx
	err = fc(txrep)

	if err == nil {
		err = tx.Commit().Error
	}

	panicked = false
	return
}
