package contract

import (
	"gorm.io/gorm"
)

const ORMKey = "web:orm"

// ORM ORMService input parameter
type ORM interface {
	GetDB(opt ...DBOption) (*gorm.DB, error)
}

// DBOption initialize option
type DBOption func(orm ORM) error
