package model

import "time"

type EndUserRole struct {
	ID         uint64     `json:"id"`
	Role       string     `json:"role"`
	Status     string     `json:"status"`
	CreatedOn  *time.Time `json:"createdOn"`
	CreatedBy  uint64     `json:"createdBy"`
	ModifiedOn *time.Time `json:"modifiedOn"`
	ModifiedBy uint64     `json:"modifiedBy"`
	EndUserID  uint64     `json:"endUserId"`
}
