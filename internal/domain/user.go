package domain

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	Name      string    `bun:"name,notnull" json:"name"`
	Email     string    `bun:"email,unique,notnull" json:"email"`
	Password  string    `bun:"password,notnull" json:"-"`
	Role      string    `bun:"role,default:'user'" json:"role"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}
