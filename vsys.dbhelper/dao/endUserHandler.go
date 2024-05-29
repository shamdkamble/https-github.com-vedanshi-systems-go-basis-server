package dao

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"gorm.io/gorm"
	m "vsys.commons/model"
	u "vsys.commons/utils"
	db "vsys.dbhelper/db"
)

const (
	USER_TABLE           string = "end_user"
	PASSWORD_PLACEHOLDER string = "********"
)

// return count of entries in users table
func GetUsersCount() (int64, error) {
	var count int64
	err := db.GetDB().Table(USER_TABLE).Count(&count)
	if err != nil {
		return count, err.Error
	}
	return count, nil
}

// return all users from users list
func GetAllUsersList() (users []m.EndUser, err error) {
	userQueryErr := db.GetDB().Table(USER_TABLE).Find(&users)
	if userQueryErr != nil {
		err = userQueryErr.Error
	}
	for user := range users {
		users[user].Password = PASSWORD_PLACEHOLDER
	}
	return
}

func CheckIfEndUserExists(endUser *m.EndUser, mobileNoInt uint64) error {
	// Attempt to find the first record matching the mobile number.
	result := db.GetDB().Table(USER_TABLE).Where("mobile_no = ?", mobileNoInt).First(&endUser)

	// Check if a record was found
	if result.Error == nil {
		// A record was found, return an error indicating the user exists.
		return errors.New("user already exists")
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// No record was found, which is the expected outcome.
		return nil
	}

	// An error occurred that wasn't due to the record not being found (e.g., DB connection issue).
	return result.Error
}

func CreateNewOrUpdateExistingEndUser(endUser *m.EndUser) error {

	now := time.Now()
	if endUser.ID != 0 {
		endUser.ModifiedOn = &now

		if err := db.GetDB().Table(USER_TABLE).Save(&endUser).Error; err != nil {
			return err
		}
	} else {
		endUser.CreatedOn = &now

		if err := db.GetDB().Table(USER_TABLE).Create(&endUser).Error; err != nil {
			return err
		}
	}
	return nil
}

// return user associated with given id
func GetUserWithID(Id string) (user m.EndUser, err error) {
	condtition := u.JoinStr("id=", Id)
	userQueryErr := db.GetDB().Table(USER_TABLE).Find(&user, condtition)
	if userQueryErr != nil {
		err = userQueryErr.Error
	}
	user.Password = PASSWORD_PLACEHOLDER
	return
}

// return user associated with given mobile number
func GetUserWithMobileNo(mobileNo string) (user m.EndUser, err error) {

	mobileNoInt, _ := strconv.ParseUint(mobileNo, 10, 64)

	userQueryErr := db.GetDB().Table(USER_TABLE).Where("mobile_no = ?", mobileNoInt).First(&user)
	if userQueryErr != nil {
		err = userQueryErr.Error
	}
	user.Password = PASSWORD_PLACEHOLDER
	return
}

// return user associated with given email
func GetUserWithEmail(email string) (user m.EndUser, err error) {
	result := db.GetDB().Table(USER_TABLE).Where("email = ?", email).First(&user)
	if result.Error != nil {
		err = result.Error
	}
	user.Password = PASSWORD_PLACEHOLDER
	return
}

// return accumulated result from end user table matching provided string
func SearchUsersForString(str string) (users []m.EndUser, err error) {
	condtition := u.JoinStr(
		"first_name LIKE '%", str, "%' OR last_name LIKE '%", str, "%' OR ",
		"email LIKE '%", str, "%' OR mobile_no LIKE '%", str, "%' OR ", "status LIKE '%", str, "%'")
	userQueryErr := db.GetDB().Table(USER_TABLE).Find(&users, condtition)
	if userQueryErr != nil {
		err = userQueryErr.Error
	}
	for user := range users {
		users[user].Password = PASSWORD_PLACEHOLDER
	}
	return
}

// return accumulated result from users table matching provided string to user role
func SearchUsersByUserRoleString(str string) (users []m.EndUser, err error) {
	condtition := u.JoinStr("id in (select (end_user_id) from `end_user_role` where role like '%", str, "%')")
	userQueryErr := db.GetDB().Table(USER_TABLE).Find(&users, condtition)
	if userQueryErr != nil {
		err = userQueryErr.Error
	}
	for user := range users {
		users[user].Password = PASSWORD_PLACEHOLDER
	}
	return
}

// return classified result from end user table matching provided string
func GetPaginatedUsersDataDB(str, limitstr, pagestr string) (users []m.EndUser, err error) {
	var count int
	var limit int
	var page int

	limit, _ = strconv.Atoi(limitstr)
	page, _ = strconv.Atoi(pagestr)
	offset := (page - 1) * limit

	condition := u.JoinStr(
		"first_name LIKE '%", str, "%' OR last_name LIKE '%", str, "%' OR ",
		"email LIKE '%", str, "%' OR mobile_no LIKE '%", str, "%' OR status LIKE '%", str, "%'")

	countErr := db.GetDB().Table(USER_TABLE).Select("COUNT(*)").Find(&count, condition)
	if countErr != nil {
		err = countErr.Error
	}

	res := float64(count) / float64(limit)
	totalPages := int(math.Round(res))
	if page > totalPages {
		offset = 0
	}

	userQueryErr := db.GetDB().Table(USER_TABLE).Limit(limit).Offset(offset).Find(&users, condition)
	if userQueryErr != nil {
		err = userQueryErr.Error
	}
	for user := range users {
		users[user].Password = PASSWORD_PLACEHOLDER
	}

	return
}

// GenerateToken creates a JWT token for the given end user and logs any errors
func GenerateToken(endUser m.EndUser) (string, time.Time, error) {

	token, expire, err := u.CreateJwtToken(endUser.MobileNo)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", time.Now(), fmt.Errorf("failed to generate token: %v", err)
	}

	return token, expire, nil
}

// Check User Credentials i.e. username and password are correct or not
// On success, return token and user data i.e. user and userRole
// On failure, return error
func CheckUserCredentials(username string, password string) (m.EndUser, []m.EndUserRole, error) {
	// log.Printf("Inside CheckUserCredentials, checkpoint: %d\n", 1)
	var userFound m.EndUser
	var userRoles []m.EndUserRole

	// log.Printf("Inside CheckUserCredentials, checkpoint: %d\n", 2)

	mobileNoInt, err := strconv.ParseUint(username, 10, 64)
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"errorMsg": "bad request(102)" + err.Error()})
		return userFound, userRoles, err
	}

	findEndUserErr := CheckIfEndUserExists(&userFound, mobileNoInt)
	// if record is not found, proceed with registration
	if gorm.ErrRecordNotFound == findEndUserErr {
		return userFound, userRoles, findEndUserErr
	}

	// log.Printf("Inside CheckUserCredentials, checkpoint: %d\n", 2)
	if password == userFound.Password {
		// password matching successful, get userRoles information

		queryEndUserRolesErr := FindEndUserRolesUsingEndUserId(&userRoles, userFound.ID)

		// log.Printf("Inside CheckUserCredentials, checkpoint: %d\n", 3)
		if queryEndUserRolesErr != nil {
			log.Printf("Inside CheckUserCredentials, checkpoint: %f\n", 3.5)
			return userFound, userRoles, queryEndUserRolesErr
		}
		// log.Printf("Inside CheckUserCredentials, checkpoint: %d\n", 4)
		return userFound, userRoles, nil
	} else {
		// authentication failure
		return userFound, userRoles, errors.New("authentication failed")
	}
}
