package field

import "time"

type IDAttr struct {
	ID uint `torm:"primary_key" json:"id"`
}

type CreatedAtAttr struct {
	CreatedAt time.Time `torm:"created_at" json:"created_at"`
}

type UpdatedAtAttr struct {
	UpdatedAt time.Time `torm:"updated_at" json:"updated_at"`
}

type DeletedAtAttr struct {
	DeletedAt time.Time `torm:"deleted_at" json:"deleted_at"`
}
