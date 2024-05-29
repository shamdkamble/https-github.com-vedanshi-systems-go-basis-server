package model

import "time"

type EndUser struct {
	ID            uint64     `json:"id"`
	FirstName     string     `json:"firstName"`
	LastName      string     `json:"lastName"`
	Email         string     `json:"email"`
	DeviceType    string     `json:"deviceType"`
	ModelNumber   string     `json:"modelNumber"`
	MobileNo      uint64     `json:"mobileNo"`
	Password      string     `json:"password"`
	ProfilePicUrl string     `json:"profilePicUrl"`
	SocialIdent   string     `json:"socialIdent"`
	Status        string     `json:"status"`
	Activated     BitBool    `json:"activated"`
	ResetKey      int        `json:"resetKey"`
	ResetDate     *time.Time `json:"resetDate"`
	CreatedOn     *time.Time `json:"createdOn"`
	CreatedBy     uint64     `json:"createdBy"`
	ModifiedOn    *time.Time `json:"modifiedOn"`
	ModifiedBy    uint64     `json:"modifiedBy"`
}


