package models

// User represents your database user model used by GORM.
// Note: Password is a []byte to hold bcrypt hashes and is hidden from JSON outputs.
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`
}
