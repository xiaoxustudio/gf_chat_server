package verifiy

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/gogf/gf/v2/frame/g"
)

var FieldsValidateRules = map[string]map[string]interface{}{
	"username": g.Map{
		"limit": [2]int{2, 10},
		"regex": "^[a-zA-Z0-9]{2,10}",
	},
	"password": g.Map{
		"limit": [2]int{2, 15},
		"regex": "^[a-zA-Z0-9]{2,15}",
	},
	"repassword": g.Map{
		"limit": [2]int{2, 15},
		"regex": "^[a-zA-Z0-9]{2,15}",
	},
	"email": g.Map{
		"limit": [2]int{4, 20},
		"regex": `^[a-z0-9]+@[a-z0-9]+\.[a-z]{2,4}$`,
	},
}

func Exec(_r map[string]interface{}, exclude []string) (bool, error) {
	if _r == nil {
		return false, errors.New("no map data provided")
	}

	requiredFields := map[string]bool{
		"username": false,
		"password": false,
		"email":    false,
	}
	if reflect.TypeOf(exclude).Kind() == reflect.Slice {
		for _, val := range exclude {
			delete(requiredFields, val)
		}
	}

	for key, value := range _r {
		str, ok := value.(string)
		if !ok {
			return false, fmt.Errorf("value for key '%s' is not a string", key)
		}

		if _, exists := requiredFields[key]; exists {
			requiredFields[key] = true

			if err := validateField(key, str); err != nil {
				return false, err
			}
			if key == "repassword" && _r["password"] != value {
				return false, fmt.Errorf("%s is not equal for password", key)
			}
		}
	}

	for field, checked := range requiredFields {
		if !checked {
			return false, fmt.Errorf("required field '%s' not found or is empty", field)
		}
	}

	return true, nil
}

func validateField(field, value string) error {
	if len(value) == 0 {
		return fmt.Errorf("%s is empty", field)
	}
	fieldData := FieldsValidateRules[field]
	limit := fieldData["limit"].([2]int)
	if len(value) < limit[0] || len(value) > limit[1] {
		return fmt.Errorf("%s length is not in range (%d-%d)", field, limit[0], limit[1])
	}
	re := regexp.MustCompile(fieldData["regex"].(string))
	if !re.MatchString(value) {
		return fmt.Errorf("%s cannot access the rules", field)
	}
	return nil
}
