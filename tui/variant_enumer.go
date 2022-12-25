// Code generated by "enumer -type=variant -trimprefix=variant -transform=kebab-case"; DO NOT EDIT.

package tui

import (
	"fmt"
)

const _variantName = "StreamDownload"

var _variantIndex = [...]uint8{0, 6, 14}

func (i variant) String() string {
	i -= 1
	if i < 0 || i >= variant(len(_variantIndex)-1) {
		return fmt.Sprintf("variant(%d)", i+1)
	}
	return _variantName[_variantIndex[i]:_variantIndex[i+1]]
}

var _variantValues = []variant{1, 2}

var _variantNameToValueMap = map[string]variant{
	_variantName[0:6]:  1,
	_variantName[6:14]: 2,
}

// variantString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func variantString(s string) (variant, error) {
	if val, ok := _variantNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to variant values", s)
}

// variantValues returns all values of the enum
func variantValues() []variant {
	return _variantValues
}

// IsAvariant returns "true" if the value is listed in the enum definition. "false" otherwise
func (i variant) IsAvariant() bool {
	for _, v := range _variantValues {
		if i == v {
			return true
		}
	}
	return false
}