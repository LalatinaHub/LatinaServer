package model

import "time"

// User represents a VPN user entity.
type User struct {
	ID         int64
	Token      string
	Password   string
	Expired    time.Time
	ServerCode string
	Quota      int64 // Quota in bytes
	Relay      string
	Adblock    bool
	VPN        string
}

// IsActive returns true if the user has remaining quota, has not expired,
// and has assigned ServerCode and VPN protocol.
func (u *User) IsActive(now time.Time) bool {
	return u.Quota > 0 && !now.After(u.Expired) && u.ServerCode != "" && u.VPN != ""
}
