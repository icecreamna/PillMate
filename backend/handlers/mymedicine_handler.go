package handlers

import (
	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

func AddMyMedicine(db *gorm.DB, mymedicine *models.MyMedicine) (*models.MyMedicine, error) {
	if err := db.Create(mymedicine).Error; err != nil {
        return nil, err
    }
    return mymedicine, nil
}