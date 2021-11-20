package models

import (
	"gocommerce/core"
)

type (
	User struct {
		core.Model
		Name     string `json:"name" gorm:"column:name"`
		Username string `json:"username" gorm:"column:username"`
		Email    string `json:"email" gorm:"column:email"`
		Password string `json:"password" gorm:"column:password"`
		Address  string `json:"address" gorm:"column:address"`
	}
)

func (User) TableName() string {
	return "users"
}

func (p *User) Create() error {
	err := core.Create(&p)
	return err
}

func (p *User) Save() error {
	err := core.Save(&p)
	return err
}

func (p *User) Delete() error {
	err := core.Delete(&p)
	return err
}

func (p *User) FindbyID(id int) error {
	err := core.FindbyID(&p, id)
	return err
}
func (b *User) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result core.PagedFindResult, err error) {
	Question := []User{}
	orders := []string{orderby}
	sorts := []string{sort}
	result, err = core.PagedFindFilter(&Question, page, rows, orders, sorts, filter, []string{})

	return result, err
}
