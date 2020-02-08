package models

import (
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

//ResetPassword type
type ResetPassword struct {
	gorm.Model
	Email string `gorm:"size:100;not null;" json:"email"`
	Token string `gorm:"size:255;not null;" json:"token"`
}

//Prepare prepare resetPassword struct
func (resetPassword *ResetPassword) Prepare() {
	resetPassword.Token = html.EscapeString(strings.TrimSpace(resetPassword.Token))
	resetPassword.Email = html.EscapeString(strings.TrimSpace(resetPassword.Email))
}

//SaveDatails save resetPassword
func (resetPassword *ResetPassword) SaveDatails(db *gorm.DB) (*ResetPassword, error) {
	var err error
	err = db.Debug().Create(&resetPassword).Error
	if err != nil {
		return &ResetPassword{}, err
	}
	return resetPassword, nil
}

//DeleteDatails delete resetpassword
func (resetPassword *ResetPassword) DeleteDatails(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&ResetPassword{}).Where("id = ?", resetPassword.ID).Take(&ResetPassword{}).Delete(&ResetPassword{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
