package utils

import (
	"FoodDecider-TG-Bot/repository"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDbConnection() *gorm.DB {
	dsn := GetDbDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	return db

}

func GetDbDSN() string {
	dbUser := GetEnvDefault("DB_USER", "root")
	dbPass := GetEnvDefault("DB_PASS", "root")
	dbName := GetEnvDefault("DB_NAME", "food_decider")
	dbHost := GetEnvDefault("DB_HOST", "localhost")
	dbPort := GetEnvDefault("DB_PORT", "3306")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
}

func CheckIfAdmin(id int64) bool {
	db := GetDbConnection()

	repo := repository.NewAdminsRepository(db)
	admin := repo.FindActiveAdmin(id)
	if admin == nil {
		return false
	}

	return true
}

func CheckIfSuperAdmin(id int64) bool {
	db := GetDbConnection()

	repo := repository.NewAdminsRepository(db)
	admin := repo.FindActiveSuperAdmin(id)
	if admin == nil {
		return false
	}

	return true
}
