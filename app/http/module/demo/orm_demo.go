package demo

import (
	"database/sql"
	"time"

	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/gin"
	"github.com/awaketai/goweb/framework/provider/orm"
)

func (api *DemoApi) OrmOperate(c *gin.Context) {
	logger := c.MustMake(contract.LogKey).(contract.Log)
	logger.Info(c, "request start", nil)
	gormService := c.MustMake(contract.ORMKey).(contract.ORM)
	db, err := gormService.GetDB(orm.WithConfigPath("database.default"))
	if err != nil {
		logger.Info(c, "inser err", map[string]any{
			"err": err,
		})
		c.AbortWithError(500, err)
		return
	}
	db.WithContext(c)
	// struct name plural:table name
	err = db.AutoMigrate(&User{})
	if err != nil {
		logger.Info(c, "migrate err", map[string]any{
			"err": err,
		})
		c.AbortWithError(500, err)
		return
	}
	logger.Info(c, "migrate ok", nil)
	// insert
	email := "foo@gmail.com"
	birthday := time.Date(2001, 1, 1, 1, 1, 1, 1, time.Local)
	user := &User{
		Name: "foo",
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
		Age:          uint8(25),
		Birthday:     &birthday,
		MemberNumber: sql.NullString{},
		ActivatedAt:  sql.NullTime{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = db.Create(user).Error
	if err != nil {
		logger.Info(c, "inser err", map[string]any{
			"id":  user.ID,
			"err": err,
		})
		c.AbortWithError(500, err)
		return
	}
	// update
	user.Name = "bar"
	err = db.Save(user).Error
	if err != nil {
		logger.Info(c, "save err", map[string]any{
			"id":  user.ID,
			"err": err,
		})
		c.AbortWithError(500, err)
		return
	}
	// query
	queryUser := &User{ID: user.ID}
	err = db.First(queryUser).Error
	if err != nil {
		logger.Info(c, "query err", map[string]any{
			"id":  user.ID,
			"err": err,
		})
		c.AbortWithError(500, err)
		return
	}
	// del
	err = db.Delete(queryUser).Error
	if err != nil {
		logger.Info(c, "del err", map[string]any{
			"id":  user.ID,
			"err": err,
		})
		c.AbortWithError(500, err)
		return
	}
}
