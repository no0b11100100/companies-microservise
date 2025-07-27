package database

import (
	configparser "companies/cmd/internal/configParser"
	"companies/cmd/internal/consts"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDB struct {
	db *gorm.DB
}

type Database interface {
	CreateRecord(CompanyInfo) (uuid.UUID, error)
	UpdateRecord(CompanyInfo, uuid.UUID) error
	DeleteRecord(uuid.UUID) error
	GetRecord(uuid.UUID) (CompanyInfo, error)
	IsRecordExists(string) bool
}

type Storage interface {
	Database
	io.Closer
}

func NewMySQLDB(config configparser.DB) Storage {
	log.Println(consts.ApplicationPrefix, "Create connection to MySQL DB")

	user := configparser.GetCfgValue("DB_USER", config.User)
	pswd := configparser.GetCfgValue("DB_PASSWORD", config.Password)
	host := configparser.GetCfgValue("DB_HOST", config.Host)
	port := configparser.GetCfgValue("DB_PORT", config.Port)
	dbName := configparser.GetCfgValue("DB_NAME", config.Name)

	initDB(user, pswd, host, port, dbName)

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", user, pswd, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&CompanyInfo{})

	return &MySQLDB{db}
}

func waitForRediness(dsn string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Println(consts.ApplicationPrefix, "Timeout waiting for MySQL: ", ctx.Err())
			return

		default:
			db, err := sql.Open("mysql", dsn)
			if err == nil {
				if pingErr := db.Ping(); pingErr == nil {
					db.Close()
					log.Println(consts.ApplicationPrefix, "MySQL is ready")
					return
				}
				db.Close()
			}

			log.Println(consts.ApplicationPrefix, "Waiting for MySQL...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (msql *MySQLDB) Close() error {
	sqlDB, err := msql.db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Close()
	if err != nil {
		return err
	}

	return nil
}

func initDB(user, psswd, addr, port, dbName string) {
	log.Println(consts.ApplicationPrefix, "InitDB")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/", user, psswd, addr, port)

	waitForRediness(dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
	if err != nil {
		log.Fatal(err)
	}
}

func (msql *MySQLDB) CreateRecord(data CompanyInfo) (uuid.UUID, error) {
	err := msql.db.Create(&data).Error
	if err != nil {
		return uuid.Nil, errors.New("CreateRecord error: " + err.Error())
	}
	return *data.ID, nil
}

func (msql *MySQLDB) UpdateRecord(data CompanyInfo, id uuid.UUID) error {
	if err := msql.db.Where("id = ?", id).First(&data).Error; err != nil {
		return errors.New("UpdateRecord error: " + err.Error())
	}

	if err := msql.db.Save(&data).Error; err != nil {
		return errors.New("UpdateRecord error: " + err.Error())
	}

	return nil
}

func (msql *MySQLDB) DeleteRecord(id uuid.UUID) error {
	if err := msql.db.Where("id = ?", id).Delete(&CompanyInfo{}).Error; err != nil {
		return errors.New("DeleteRecord error: " + err.Error())
	}

	return nil
}

func (msql *MySQLDB) GetRecord(id uuid.UUID) (CompanyInfo, error) {
	record := CompanyInfo{}
	if err := msql.db.Where("id = ?", id).First(&record).Error; err != nil {
		return record, errors.New("GetRecord error: " + err.Error())
	}

	return record, nil
}

func (msql *MySQLDB) IsRecordExists(name string) bool {
	var record CompanyInfo
	err := msql.db.Select("id").Where("name = ?", name).Limit(1).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		return false
	}
	return true
}
