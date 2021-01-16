package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var CustomValidator *cValidator

type cValidator struct {
	db *gorm.DB
}

func NewValidator(db *gorm.DB) {
	CustomValidator = &cValidator{
		db: db,
	}
}

// check if value of request is unique in database
// tag format : unique=tablename.columnname
func (cv *cValidator) Unique() validator.Func {
	return func(fl validator.FieldLevel) bool {
		arr := strings.Split(fl.Param(), ".")
		rows, err := cv.db.Table(arr[0]).Select("*").Where("`"+arr[1]+"` = ?", fl.Field().Interface()).Rows()
		if err != nil {
			return false
		}

		// loop through all item and check for deleted
		for rows.Next() {
			// get all columns
			cols, _ := rows.Columns()

			//find deleted at position
			deletedIndex := -1
			for i, v := range cols {
				if v == "deleted_at" {
					deletedIndex = i
					break
				}
			}

			// return invalid if there is no deleted_at field
			if deletedIndex == -1 {
				_ = rows.Close()
				return false
			}

			// Result is your slice string.
			rawResult := make([][]byte, len(cols))
			result := make([]string, len(cols))

			dest := make([]interface{}, len(cols)) // A temporary interface{} slice
			for i := range rawResult {
				dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
			}

			// scan all to the result array
			err = rows.Scan(dest...)
			if err != nil {
				_ = rows.Close()
				return false
			}

			// iterate all raw and convert to string
			for i, raw := range rawResult {
				if raw == nil {
					result[i] = "nil"
				} else {
					result[i] = string(raw)
				}
			}

			// return not valid if deleted value is nil
			if result[deletedIndex] == "nil" {
				_ = rows.Close()
				return false
			}
		}
		_ = rows.Close()
		return true
	}
}

// check if value of request exist in database
// tag format : exist=tablename.columnname
func (cv *cValidator) Exist() validator.Func {
	return func(fl validator.FieldLevel) bool {
		arr := strings.Split(fl.Param(), ".")
		return cv.db.Table(arr[0]).
			Select("*").
			Where("`"+arr[1]+"` = ?", fl.Field().Interface()).
			Take(&struct{}{}).Error != gorm.ErrRecordNotFound
	}
}
