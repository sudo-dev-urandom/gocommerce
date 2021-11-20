package core

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

type (
	Model struct {
		ID int `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key,column:id"`
	}

	DBFunc func(tx *gorm.DB) error

	PagedFindResult struct {
		TotalData   int         `json:"total_data"`
		Rows        int         `json:"rows"`
		CurrentPage int         `json:"current_page"`
		LastPage    int         `json:"last_page"`
		From        int         `json:"from"`
		To          int         `json:"to"`
		Data        interface{} `json:"data"`
	}
	CompareFilter struct {
		Value1 interface{} `json:"value1"`
		Value2 interface{} `json:"value2"`
	}
)

type tagOptions string

// tag if filter condtion with comma conditionA, ConditionB
func parseTag(tag string) tagOptions {
	return tagOptions(tag)
}

func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

func WithinTransaction(fn DBFunc) (err error) {
	tx := App.DB.Begin()
	defer tx.Commit()
	err = fn(tx)

	return err
}

// Create
func Create(i interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		if !App.DB.NewRecord(i) {
			return fmt.Errorf("cannot create row. not a new record")
		}
		if err = tx.Create(i).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

// Update
func Save(i interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		if App.DB.NewRecord(i) {
			return fmt.Errorf("cannot save row. it is a new record")
		}
		if err = tx.Save(i).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

// Delete
func Delete(i interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		if err = tx.Delete(i).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

// FindbyID finds row by id.
func FindbyID(i interface{}, id int) (err error) {
	return WithinTransaction(func(tx *gorm.DB) error {
		if err = tx.Last(i, id).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

func PagedFindFilter(i interface{}, page int, rows int, order []string, sort []string, filter interface{}, field []string, allfieldcondition ...string) (result PagedFindResult, err error) {

	if page <= 0 {
		page = 1
	}

	query := App.DB

	query = ConditionQuery(query, filter)
	query = OrderSortQuery(query, order, sort)

	temp := query
	var totalRows int

	temp.Find(i).Count(&totalRows)

	if len(field) > 0 {
		query = query.Select(field)
	}

	var (
		offset   int
		lastPage int
	)

	if rows > 0 {
		offset = (page * rows) - rows
		lastPage = int(math.Ceil(float64(totalRows) / float64(rows)))

		query = query.Limit(rows).Offset(offset)
	}

	query.Find(i)

	result = PagedFindResult{
		TotalData:   totalRows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        i,
	}

	return result, err
}

func ConditionQuery(query *gorm.DB, filter interface{}) *gorm.DB {
	refFilter := reflect.ValueOf(filter).Elem()
	refType := refFilter.Type()
	for x := 0; x < refFilter.NumField(); x++ {
		field := refFilter.Field(x)
		// check if empty
		if !reflect.DeepEqual(field.Interface(), reflect.Zero(reflect.TypeOf(field.Interface())).Interface()) {
			con := strings.Split(refType.Field(x).Tag.Get("condition"), ",")
			tags := parseTag(refType.Field(x).Tag.Get("condition"))
			switch con[0] {
			default:
				format := fmt.Sprintf("%s IN (?)", refType.Field(x).Tag.Get("json"))
				if tags.Contains("optional") {
					query = query.Or(format, field.Interface())
				} else {
					query = query.Where(format, field.Interface())
				}
			case "LIKE":
				format := fmt.Sprintf("LOWER(%s) %s ?", refType.Field(x).Tag.Get("json"), con[0])
				field := "%" + strings.ToLower(field.Interface().(string)) + "%"
				if tags.Contains("optional") {
					query = query.Or(format, field)
				} else {
					query = query.Where(format, field)
				}
			case "BETWEEN":
				if values, ok := field.Interface().(CompareFilter); ok && values.Value1 != "" {
					format := fmt.Sprintf("%s %s ? %s ?", refType.Field(x).Tag.Get("json"), con[0], "AND")
					if tags.Contains("optional") {
						query = query.Or(format, values.Value1, values.Value2)
					} else {
						query = query.Where(format, values.Value1, values.Value2)
					}
				}
			case "OR":
				var e []string
				for _, v := range field.Interface().([]string) {
					e = append(e, refType.Field(x).Tag.Get("json")+" = '"+v+"'")
				}
				if tags.Contains("optional") {
					query = query.Or(strings.Join(e, " OR "))
				} else {
					query = query.Where(strings.Join(e, " OR "))
				}

			case "NOT IN":
				format := fmt.Sprintf("%s NOT IN (%s)", refType.Field(x).Tag.Get("json"), field.Interface())
				query = query.Where(format)
			}
		}
	}

	return query
}

func OrderSortQuery(query *gorm.DB, order []string, sort []string) *gorm.DB {
	for k, v := range order {
		q := v
		if len(sort) > k {
			value := sort[k]
			if strings.ToUpper(value) == "ASC" || strings.ToUpper(value) == "DESC" {
				q = v + " " + strings.ToUpper(value)
			}
		}
		query = query.Order(q)
	}

	return query
}
