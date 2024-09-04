package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/structure"
)

// type Model struct {
// 	ID        string     `gorm:"column:id;primary_key;size:36;"`
// 	UserName  string     `gorm:"column:user_name;size:64;index;default:'';not null;"`
// 	RealName  string     `gorm:"column:real_name;size:64;index;default:'';not null;"`
// 	Password  string     `gorm:"column:password;size:40;default:'';not null;"`
// 	Email     *string    `gorm:"column:email;size:255;index;"`
// 	Phone     *string    `gorm:"column:phone;size:20;index;"`
// 	Status    int        `gorm:"column:status;index;default:0;not null;"`
// 	Creator   string     `gorm:"column:creator;size:36;"`
// 	CreatedAt time.Time  `gorm:"column:created_at;index;"`
// 	UpdatedAt time.Time  `gorm:"column:updated_at;index;"`
// 	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`
// }

type Model struct {
	ID        uuid.UUID
	UserName  string
	RealName  string
	Password  string
	Email     *string
	Phone     *string
	Status    int
	Creator   *uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	IDString  *string
}

func (Model) TableName() string {
	return "user"
}

func (a Model) ToDomain() *user.User {
	item := new(user.User)
	structure.CopyWithUUID(a, item)
	return item
}

func toDomainList(ms []*Model) []*user.User {
	list := make([]*user.User, len(ms))
	for i, item := range ms {
		list[i] = item.ToDomain()
	}
	return list
}

func domainToModel(u *user.User) *Model {
	item := new(Model)
	structure.CopyWithUUID(u, item)
	return item
}
