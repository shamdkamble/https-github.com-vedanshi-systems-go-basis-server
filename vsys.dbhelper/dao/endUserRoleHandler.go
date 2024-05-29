package dao

import (
	"time"

	m "vsys.commons/model"
	u "vsys.commons/utils"
	db "vsys.dbhelper/db"
)

const userRoleTable string = "end_user_role"

// creates end user role in end user table
func CreateEndUserRole(userRole *m.EndUserRole) error {
	now := time.Now()
	userRole.CreatedOn = &now

	if err := db.GetDB().Table(userRoleTable).Create(&userRole).Error; err != nil {
		return err
	}
	return nil
}

// finds end user role associated with provided user id
func FindEndUserRolesUsingEndUserId(userRoles *[]m.EndUserRole, userId uint64) error {
	err := db.GetDB().Table(userRoleTable).Where("end_user_id = ?", userId).Find(&userRoles)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

// return total count of entries in user role table
func GetUserRolesCount() (int64, error) {
	var count int64
	err := db.GetDB().Table(userRoleTable).Count(&count)
	if err != nil {
		return count, err.Error
	}
	return count, nil
}

// return list of all entries from user role table
func GetAllUserRolesList() (userRoles []m.EndUserRole, err error) {
	userRoleQueryErr := db.GetDB().Table(userRoleTable).Find(&userRoles)
	if userRoleQueryErr != nil {
		err = userRoleQueryErr.Error
	}
	return
}

// return end user roles associated with given mobile number
func GetEndUserRolesWithMobileNum(mobileNo string) (userRoles []m.EndUserRole, err error) {
	condition := u.JoinStr("end_user_id in (select (id) from `end_user` where mobile_no=", mobileNo, ")")
	userRoleQueryErr := db.GetDB().Table(userRoleTable).Find(&userRoles, condition)
	if userRoleQueryErr != nil {
		err = userRoleQueryErr.Error
	}
	return
}

// return user roles associated with provided user id
func GetUserRolesUsingUserId(userId string) (userRoles []m.EndUserRole, err error) {
	userRoleQueryErr := db.GetDB().Table(userRoleTable).Where("end_user_id", userId).Find(&userRoles)
	if userRoleQueryErr != nil {
		err = userRoleQueryErr.Error
	}
	return
}
