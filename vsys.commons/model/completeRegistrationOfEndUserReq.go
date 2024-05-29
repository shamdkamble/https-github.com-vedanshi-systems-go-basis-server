package model

type CompleteRegistrationOfUser struct {
	EndUser  EndUser `json:"endUser"`
	RoleType int     `json:"roleType"`
}
