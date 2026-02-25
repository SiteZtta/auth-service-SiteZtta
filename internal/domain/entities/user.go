package entities

import "time"

type User struct {
	ID        int64     `json:"id" db:"id`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	PassHash  []byte    `json:"passHash" db:"pass_hash"` // <- BYTEA / BLOB
	Role      int32     `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
