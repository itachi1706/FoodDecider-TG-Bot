package repository

import (
    "FoodDecider-TG-Bot/model"
    "gorm.io/gorm"
)

type AdminsRepository struct {
    db *gorm.DB
}

func NewAdminsRepository(db *gorm.DB) AdminsRepository {
    return AdminsRepository{db: db}
}

func (a AdminsRepository) FindAllActiveAdmins() []model.Admins {
    var admins []model.Admins
    a.db.Where("status = ?", "A").Find(&admins)

    return admins
}

func (a AdminsRepository) FindActiveAdmin(id int64) *model.Admins {
    var admin model.Admins
    result := a.db.Where("telegram_id = ? AND status = ?", id, "A").First(&admin)
    if result.Error != nil {
        return nil
    }

    return &admin
}

func (a AdminsRepository) FindActiveSuperAdmin(id int64) *model.Admins {
    var admin model.Admins
    result := a.db.Where("telegram_id = ? AND status = ? AND is_superadmin = ?", id, "A", true).First(&admin)
    if result.Error != nil {
        return nil
    }

    return &admin
}

func (a AdminsRepository) FindAdmin(id int64) *model.Admins {
    var admin model.Admins
    result := a.db.Where("telegram_id = ?", id).First(&admin)
    if result.Error != nil {
        return nil
    }

    return &admin
}
