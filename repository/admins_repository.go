package repository

import (
    "FoodDecider-TG-Bot/model"
    "gorm.io/gorm"
)

func FindAllActiveAdmins(db *gorm.DB) []model.Admins {
    var admins []model.Admins
    db.Where("status = ?", "A").Find(&admins)

    return admins
}

func FindActiveAdmin(db *gorm.DB, id int64) *model.Admins {
    var admin model.Admins
    result := db.Where("telegram_id = ? AND status = ?", id, "A").First(&admin)
    if result.Error != nil {
        return nil
    }

    return &admin
}

func FindActiveSuperAdmin(db *gorm.DB, id int64) *model.Admins {
    var admin model.Admins
    result := db.Where("telegram_id = ? AND status = ? AND is_superadmin = ?", id, "A", true).First(&admin)
    if result.Error != nil {
        return nil
    }

    return &admin
}

func FindAdmin(db *gorm.DB, id int64) *model.Admins {
    var admin model.Admins
    result := db.Where("telegram_id = ?", id).First(&admin)
    if result.Error != nil {
        return nil
    }

    return &admin
}
