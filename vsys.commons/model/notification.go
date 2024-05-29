package model

import "time"

type Notification struct {
	Id               uint64     `json:"id"`
	Text             string     `json:"text"`
	NotificationType string     `json:"notificationType"`
	PostedOn         *time.Time `json:"postedOn"`
	PostedFor        uint64     `json:"postedFor"`
	PostedBy         uint64     `json:"postedBy"`
	RoleType         string     `json:"roleType"`
	NotificationHash string     `json:"notificationHash"`
}
