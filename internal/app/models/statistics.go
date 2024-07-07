package models

// CountStatistic Contains URL nad users count in storage
type CountStatistic struct {
	URLCount  int `json:"urls"`
	UserCount int `json:"users"`
}
