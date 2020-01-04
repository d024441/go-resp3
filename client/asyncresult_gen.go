// Code generated by "converter"; DO NOT EDIT.

package client

// ToBool converts a redis value to a bool.
// In case the conversion is not supported a ConversionError is returned.
func (r *asyncResult) ToBool() (bool, error) {
	if err := r.wait(); err != nil {
		return false, err
	}
	if r.err != nil {
		return false, r.err
	}
	return r.value.ToBool()
}

// ToFloat64 converts a redis value to a float64.
// In case the conversion is not supported a ConversionError is returned.
func (r *asyncResult) ToFloat64() (float64, error) {
	if err := r.wait(); err != nil {
		return 0, err
	}
	if r.err != nil {
		return 0, r.err
	}
	return r.value.ToFloat64()
}

// ToInt64 converts a redis value to an int64.
// In case the conversion is not supported a ConversionError is returned.
func (r *asyncResult) ToInt64() (int64, error) {
	if err := r.wait(); err != nil {
		return 0, err
	}
	if r.err != nil {
		return 0, r.err
	}
	return r.value.ToInt64()
}

// ToInt64Slice returns a slice with values of type int64. In case value conversion to string is not possible
// a ConversitionError is returned.
func (r *asyncResult) ToInt64Slice() ([]int64, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.value.ToInt64Slice()
}

// ToSlice returns a slice with values of type interface{}.
func (r *asyncResult) ToSlice() ([]interface{}, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.value.ToSlice()
}

// ToSlice2 returns a slice with values of type []interface{}. In case value conversion to []interface{} is not possible
// a ConversitionError is returned.
func (r *asyncResult) ToSlice2() ([][]interface{}, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.value.ToSlice2()
}

// ToString converts a redis value to a string.
// In case the conversion is not supported a ConversionError is returned.
func (r *asyncResult) ToString() (string, error) {
	if err := r.wait(); err != nil {
		return "", err
	}
	if r.err != nil {
		return "", r.err
	}
	return r.value.ToString()
}

// ToStringInt64Map returns a map with keys of type string and values of type int64. In case key or value conversion is not possible
// a ConvertionError is returned.
func (r *asyncResult) ToStringInt64Map() (map[string]int64, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.value.ToStringInt64Map()
}

// ToStringMap returns a map with keys of type string. In case key conversion to string is not possible
// a ConvertionError is returned.
func (r *asyncResult) ToStringMap() (map[string]interface{}, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.value.ToStringMap()
}

// ToStringMapSlice returns a slice with values of type map[string]interfcae{}. In case value conversion to map[string]interface{} is not possible
// a ConversitionError is returned.
func (r *asyncResult) ToStringMapSlice() ([]map[string]interface{}, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.value.ToStringMapSlice()
}

// ToStringSet returns a map with keys of type string and boolean true values. In case key conversion to string is not possible
// a ConvertionError is returned.
func (r *asyncResult) ToStringSet() (map[string]bool, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.value.ToStringSet()
}

// ToStringSlice returns a slice with values of type string. In case value conversion to string is not possible
// a ConversitionError is returned.
func (r *asyncResult) ToStringSlice() ([]string, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.value.ToStringSlice()
}

// ToStringStringMap returns a map with keys and values of type string. In case key or value conversion to string is not possible
// a ConvertionError is returned.
func (r *asyncResult) ToStringStringMap() (map[string]string, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.value.ToStringStringMap()
}

// ToTree returns a tree with nodes of type []interface{} and leaves of type interface{}. In case value conversion to []interface{} is not possible
// a ConversitionError is returned.
func (r *asyncResult) ToTree() ([]interface{}, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.value.ToTree()
}

// ToXrange returns a slice with values of type IdMap. In case the conversion is not possible
// a ConversitionError is returned.
func (r *asyncResult) ToXrange() ([]IDMap, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.value.ToXrange()
}

// ToXread returns a map[string] with values of type IdMap. In case the conversion is not possible
// a ConversitionError is returned.
func (r *asyncResult) ToXread() (map[string][]IDMap, error) {
	if err := r.wait(); err != nil {
		return nil, err
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.value.ToXread()
}