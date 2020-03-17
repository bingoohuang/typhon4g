// Code generated by "enumer -type=ConfFmt -json"; DO NOT EDIT.

//
package base

import (
	"encoding/json"
	"fmt"
)

const _ConfFmtName = "PropertiesFmtTxtFmt"

var _ConfFmtIndex = [...]uint8{0, 13, 19}

func (i ConfFmt) String() string {
	if i < 0 || i >= ConfFmt(len(_ConfFmtIndex)-1) {
		return fmt.Sprintf("ConfFmt(%d)", i)
	}
	return _ConfFmtName[_ConfFmtIndex[i]:_ConfFmtIndex[i+1]]
}

var _ConfFmtValues = []ConfFmt{0, 1}

var _ConfFmtNameToValueMap = map[string]ConfFmt{
	_ConfFmtName[0:13]:  0,
	_ConfFmtName[13:19]: 1,
}

// ConfFmtString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ConfFmtString(s string) (ConfFmt, error) {
	if val, ok := _ConfFmtNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to ConfFmt values", s)
}

// ConfFmtValues returns all values of the enum
func ConfFmtValues() []ConfFmt {
	return _ConfFmtValues
}

// IsAConfFmt returns "true" if the value is listed in the enum definition. "false" otherwise
func (i ConfFmt) IsAConfFmt() bool {
	for _, v := range _ConfFmtValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for ConfFmt
func (i ConfFmt) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for ConfFmt
func (i *ConfFmt) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ConfFmt should be a string, got %s", data)
	}

	var err error
	*i, err = ConfFmtString(s)
	return err
}