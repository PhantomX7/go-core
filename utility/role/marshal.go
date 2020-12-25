package role

import (
	"reflect"

	"github.com/liip/sheriff"

	"github.com/PhantomX7/go-core/utility/errors"
)

func GetDataJSONByRole(data interface{}, groups ...string) (interface{}, error) {
	kind := reflect.ValueOf(data).Kind()
	if data == nil ||
		((kind == reflect.Ptr) &&
			reflect.ValueOf(data).IsNil()) {
		return nil, errors.ErrUnprocessableEntity
	}

	o := &sheriff.Options{
		Groups: groups,
	}

	filtered, err := sheriff.Marshal(o, data)
	if err != nil {
		return nil, err
	}

	return filtered, nil
}
