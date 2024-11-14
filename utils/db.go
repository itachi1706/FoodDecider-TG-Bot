package utils

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "strconv"
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
    dbPort := GetEnvDefaultInt("DB_PORT", 3306)

    return dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + strconv.Itoa(dbPort) + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
}
