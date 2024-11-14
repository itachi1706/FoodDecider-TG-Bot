package utils

import (
    "FoodDecider-TG-Bot/model"
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

func CheckIfAdmin(id int64) bool {
    db := GetDbConnection()

    // Check if exists in admin list and status is A
    var admin model.Admins
    result := db.Where("telegram_id = ? AND status = ?", id, "A").First(&admin)
    if result.Error != nil {
        return false
    }

    return true
}

func CheckIfSuperAdmin(id int64) bool {
    db := GetDbConnection()

    // Check if exists in admin list and status is A
    var admin model.Admins
    result := db.Where("telegram_id = ? AND status = ? AND is_superadmin = ?", id, "A", true).First(&admin)
    if result.Error != nil {
        return false
    }

    return true
}
