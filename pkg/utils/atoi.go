package utils

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Return string to int and ignore error
func Atoi(value string) int {
	num, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return num
}

func StrToUint64(str string) (uint64, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	return uint64(i), err
}

func StrToUint64NoError(str string) uint64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return uint64(i)
}

func Uint64ToStr(id uint64) string {
	return strconv.FormatUint(id, 10)
}

func Uint64ToStrPtr(id uint64) *string {
	strPtr := strconv.FormatUint(id, 10)
	return &strPtr
}

func ConvertToStringIfNeeded(input interface{}) (interface{}, error) {
    jsonData, err := json.Marshal(input)
    if err != nil {
        return nil, err
    }

    var result interface{}
    err = json.Unmarshal(jsonData, &result)
    if err != nil {
        return nil, err
    }
    return process(result), nil
}

// Recursively process all types, including map, slice, and other basic types
func process(value interface{}) interface{} {
    switch v := value.(type) {
    case map[string]interface{}:
        // Process each element in the map
        for key, val := range v {
            // Check if the key contains "id" (case insensitive)
            if strings.Contains(strings.ToLower(key), "id") {
                // Convert to string based on the value type
                switch val := val.(type) {
                case float64:
                    // Check if the number is negative or has decimal places
                    if val < 0 || val != float64(uint64(val)) {
                        // Keep original value if it's not a valid unsigned integer
                        v[key] = val
                    } else {
                        v[key] = strconv.FormatUint(uint64(val), 10)
                    }
                case int64:
                    if val >= 0 {
                        v[key] = strconv.FormatInt(val, 10)
                    }
                case int:
                    if val >= 0 {
                        v[key] = strconv.Itoa(val)
                    }
                case string:
                    // Already a string, keep as is
                    v[key] = val
                default:
                    // Process other types normally
                    v[key] = process(val)
                }
            } else {
                // Process non-ID fields normally
                v[key] = process(val)
            }
        }
        return v
    case []interface{}:
        // Process each element in the slice
        for i, val := range v {
            v[i] = process(val)
        }
        return v
    default:
        // For other types, return the value as is
        return value
    }
}
