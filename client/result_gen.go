// Code generated by "converter"; DO NOT EDIT.

package client

// Result represents a redis command result.
type Result interface {
	// Attr returns the attribute of a Redis value if provided - <nil> otherwise.
	Attr() (Map, error)
	// Err returns any of the following errors:
	// - <nil>: result received successfully.
	// - ErrNotFlushed: result not received yet (pipeline not flushed).
	// - ErrTimeout: timeout reached before result is available.
	// - RedisError: redis error message if a redis command was executed unsuccessfully.
	Err() error
	// IsNull returns <true> if the redis value is null.
	IsNull() (bool, error)
	// Kind returns the type of a Redis value.
	Kind() (RedisKind, error)
	// Map returns a Map if the Redis type is a map.
	Map() (Map, error)
	// Set returns a Set if the Redis type is a set.
	Set() (Set, error)
	// Slice returns a Slice if the Redis type is an array.
	Slice() (Slice, error)
	// ToBool converts a redis value to a bool.
	// In case the conversion is not supported a ConversionError is returned.
	ToBool() (bool, error)
	// ToFloat64 converts a redis value to a float64.
	// In case the conversion is not supported a ConversionError is returned.
	ToFloat64() (float64, error)
	// ToInt64 converts a redis value to an int64.
	// In case the conversion is not supported a ConversionError is returned.
	ToInt64() (int64, error)
	// ToInt64Slice returns a slice with values of type int64. In case value conversion to string is not possible
	// a ConversitionError is returned.
	ToInt64Slice() ([]int64, error)
	// ToSlice returns a slice with values of type interface{}.
	ToSlice() ([]interface{}, error)
	// ToSlice2 returns a slice with values of type []interface{}. In case value conversion to []interface{} is not possible
	// a ConversitionError is returned.
	ToSlice2() ([][]interface{}, error)
	// ToString converts a redis value to a string.
	// In case the conversion is not supported a ConversionError is returned.
	ToString() (string, error)
	// ToStringInt64Map returns a map with keys of type string and values of type int64. In case key or value conversion is not possible
	// a ConvertionError is returned.
	ToStringInt64Map() (map[string]int64, error)
	// ToStringMap returns a map with keys of type string. In case key conversion to string is not possible
	// a ConvertionError is returned.
	ToStringMap() (map[string]interface{}, error)
	// ToStringMapSlice returns a slice with values of type map[string]interfcae{}. In case value conversion to map[string]interface{} is not possible
	// a ConversitionError is returned.
	ToStringMapSlice() ([]map[string]interface{}, error)
	// ToStringSet returns a map with keys of type string and boolean true values. In case key conversion to string is not possible
	// a ConvertionError is returned.
	ToStringSet() (map[string]bool, error)
	// ToStringSlice returns a slice with values of type string. In case value conversion to string is not possible
	// a ConversitionError is returned.
	ToStringSlice() ([]string, error)
	// ToStringStringMap returns a map with keys and values of type string. In case key or value conversion to string is not possible
	// a ConvertionError is returned.
	ToStringStringMap() (map[string]string, error)
	// ToTree returns a tree with nodes of type []interface{} and leaves of type interface{}. In case value conversion to []interface{} is not possible
	// a ConversitionError is returned.
	ToTree() ([]interface{}, error)
	// ToXrange returns a slice with values of type IdMap. In case the conversion is not possible
	// a ConversitionError is returned.
	ToXrange() ([]IDMap, error)
	// ToXread returns a map[string] with values of type IdMap. In case the conversion is not possible
	// a ConversitionError is returned.
	ToXread() (map[string][]IDMap, error)
	// Value returns a Redis value.
	Value() (RedisValue, error)
	// VerbatimString returns a VerbatimString if the Redis type is a verbatim string.
	VerbatimString() (VerbatimString, error)
}