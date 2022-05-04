package anaml

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
var identifierPattern = regexp.MustCompile(`^[0-9]+$`)

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandIdentifierList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vv, _ := strconv.Atoi(v.(string))
			vs = append(vs, vv)
		}
	}
	return vs
}

func expandSingleMap(value interface{}) (map[string]interface{}, error) {
	if value == nil {
		return nil, errors.New("Value is null")
	}

	array, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Value is not an array. Value: %v", value)
	}

	if len(array) == 0 {
		return nil, errors.New("Array is empty")
	}

	single, ok := array[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Value at index 0 of array is not a map. Value: %v", array[0])
	}

	return single, nil
}

func getNullableInt(d *schema.ResourceData, key string) *int {
	value, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	intValue, ok := value.(int)
	if !ok {
		return nil
	}
	return &intValue
}

func getNullableString(d *schema.ResourceData, key string) *string {
	value, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	stringValue, ok := value.(string)
	if !ok {
		return nil
	}
	return &stringValue
}

func getNullableMapString(d map[string]interface{}, key string) *string {
	rawValue, ok := d[key]
	if ok {
		stringValue := rawValue.(string)
		return &stringValue
	} else {
		return nil
	}
}

func identifierList(ints []int) []string {
	vs := make([]string, 0, len(ints))
	for _, v := range ints {
		vs = append(vs, strconv.Itoa(v))
	}
	return vs
}

func validateAnamlName() schema.SchemaValidateFunc {
	return validation.StringMatch(namePattern, "Names must start with a lowercase a-z and contain only a-z, underscores, and digits.")
}

func validateAnamlIdentifier() schema.SchemaValidateFunc {
	return validation.StringMatch(identifierPattern, "Must be parsable as an integer")
}

func ValidateDuration() schema.SchemaValidateFunc {
	return func(i interface{}, k string) ([]string, []error) {
		_, err := time.ParseDuration(i.(string))
		if err != nil {
			return nil, []error{err}
		}
		return nil, nil
	}
}

func validateMapKeysAnamlIdentifier() schema.SchemaValidateDiagFunc {
	return validation.MapKeyMatch(identifierPattern, "Map keys must be parsable as an integer")
}

type IdAndVersion struct {
	ID      int
	Version string
}

func (n *IdAndVersion) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&n.ID, &n.Version}
	expectedLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), expectedLen; g != e {
		return fmt.Errorf("Wrong number of fields in response tuple: %d != %d", g, e)
	}
	return nil
}

func unmarshalIdAndVersion(input []byte, idAndVersion *IdAndVersion) error {
	if err := json.Unmarshal(input, idAndVersion); err != nil {
		return err
	}
	return nil
}
