package model

import (
	"context"
	"github.com/jinzhu/gorm"
)

var err error
var InsertOmitFields = []string{"xxx_unrecognized", "xxx_sizecache", "XXX_NoUnkeyedLiteral", "created_at", "updated_at"}




// prepare for outher
type BaseModel struct {
	*gorm.DB
	ctx *context.Context
}

// SetCtx set context
func (obj *BaseModel) SetCtx(c *context.Context) {
	obj.ctx = c
}

// GetDB get gorm.DB info
func (obj *BaseModel) GetDB() *gorm.DB {
	return obj.DB
}

// UpdateDB update gorm.DB info
func (obj *BaseModel) UpdateDB(db *gorm.DB) {
	obj.DB = db
}


