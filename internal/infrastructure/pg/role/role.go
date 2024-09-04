package role

import (
	"time"

	"github.com/google/uuid"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/role"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/structure"
)

// type Model struct {
// 	ID        string     `gorm:"column:id;primary_key;size:36;"`
// 	Name      string     `gorm:"column:name;size:100;index;default:'';not null;"`
// 	Sequence  int        `gorm:"column:sequence;index;default:0;not null;"`
// 	Memo      *string    `gorm:"column:memo;size:1024;"`
// 	Status    int        `gorm:"column:status;index;default:0;not null;"`
// 	Creator   string     `gorm:"column:creator;size:36;"`
// 	CreatedAt time.Time  `gorm:"column:created_at;index;"`
// 	UpdatedAt time.Time  `gorm:"column:updated_at;index;"`
// 	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`
// }

type Model struct {
	ID        uuid.UUID
	Name      string
	Sequence  int
	Memo      *string
	Status    int
	Creator   *uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	IDString  *string
}

func (Model) TableName() string {
	return "roles"
}

func (a Model) ToDomain() *role.Role {
	item := new(role.Role)
	structure.CopyWithUUID(a, item)
	return item
}

func toDomainList(ms []*Model) []*role.Role {
	list := make([]*role.Role, len(ms))
	for i, item := range ms {
		list[i] = item.ToDomain()
	}
	return list
}

func domainToModel(r *role.Role) *Model {
	item := new(Model)
	structure.CopyWithUUID(r, item)
	return item
}
