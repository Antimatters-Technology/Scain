package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// CanonicalJSON converts an object to canonical JSON string for deterministic hashing
func CanonicalJSON(obj interface{}) (string, error) {
	if obj == nil {
		return "null", nil
	}

	v := reflect.ValueOf(obj)
	return canonicalJSONValue(v)
}

// canonicalJSONValue recursively processes values to create canonical JSON
func canonicalJSONValue(v reflect.Value) (string, error) {
	// Handle pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return "null", nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return "true", nil
		}
		return "false", nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint()), nil

	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", v.Float()), nil

	case reflect.String:
		bytes, err := json.Marshal(v.String())
		if err != nil {
			return "", fmt.Errorf("failed to marshal string: %w", err)
		}
		return string(bytes), nil

	case reflect.Array, reflect.Slice:
		if v.IsNil() {
			return "null", nil
		}
		
		var elements []string
		for i := 0; i < v.Len(); i++ {
			element, err := canonicalJSONValue(v.Index(i))
			if err != nil {
				return "", fmt.Errorf("failed to serialize array element %d: %w", i, err)
			}
			elements = append(elements, element)
		}
		
		return "[" + strings.Join(elements, ",") + "]", nil

	case reflect.Map:
		if v.IsNil() {
			return "null", nil
		}
		
		// Get all keys and sort them
		keys := v.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			keyI, err := canonicalJSONValue(keys[i])
			if err != nil {
				return false
			}
			keyJ, err := canonicalJSONValue(keys[j])
			if err != nil {
				return false
			}
			return keyI < keyJ
		})
		
		var pairs []string
		for _, key := range keys {
			keyStr, err := canonicalJSONValue(key)
			if err != nil {
				return "", fmt.Errorf("failed to serialize map key: %w", err)
			}
			
			valueStr, err := canonicalJSONValue(v.MapIndex(key))
			if err != nil {
				return "", fmt.Errorf("failed to serialize map value: %w", err)
			}
			
			pairs = append(pairs, keyStr+":"+valueStr)
		}
		
		return "{" + strings.Join(pairs, ",") + "}", nil

	case reflect.Struct:
		return canonicalJSONStruct(v)

	default:
		// For other types, try to marshal to JSON and unmarshal to generic interface{}
		jsonBytes, err := json.Marshal(v.Interface())
		if err != nil {
			return "", fmt.Errorf("failed to marshal value of type %s: %w", v.Type(), err)
		}
		
		var generic interface{}
		if err := json.Unmarshal(jsonBytes, &generic); err != nil {
			return "", fmt.Errorf("failed to unmarshal to generic interface: %w", err)
		}
		
		return canonicalJSONValue(reflect.ValueOf(generic))
	}
}

// canonicalJSONStruct processes struct values for canonical JSON
func canonicalJSONStruct(v reflect.Value) (string, error) {
	t := v.Type()
	
	// Collect field name-value pairs
	type fieldPair struct {
		name  string
		value string
	}
	
	var pairs []fieldPair
	
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		
		// Skip unexported fields
		if !field.IsExported() {
			continue
		}
		
		// Get JSON field name
		fieldName := field.Name
		omitEmpty := false
		
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" && parts[0] != "-" {
				fieldName = parts[0]
			}
			
			// Check for omitempty and skip zero values
			if len(parts) > 1 {
				for _, option := range parts[1:] {
					if option == "omitempty" {
						omitEmpty = true
						break
					}
				}
			}
		}
		
		// Skip zero values if omitempty is set
		if omitEmpty && isZeroValue(fieldValue) {
			continue
		}
		
		// Serialize field value
		valueStr, err := canonicalJSONValue(fieldValue)
		if err != nil {
			return "", fmt.Errorf("failed to serialize field %s: %w", fieldName, err)
		}
		
		// Create JSON key string
		keyBytes, err := json.Marshal(fieldName)
		if err != nil {
			return "", fmt.Errorf("failed to marshal field name %s: %w", fieldName, err)
		}
		
		pairs = append(pairs, fieldPair{
			name:  string(keyBytes),
			value: valueStr,
		})
	}
	
	// Sort pairs by field name
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].name < pairs[j].name
	})
	
	// Build JSON object string
	var pairStrings []string
	for _, pair := range pairs {
		pairStrings = append(pairStrings, pair.name+":"+pair.value)
	}
	
	return "{" + strings.Join(pairStrings, ",") + "}", nil
}

// isZeroValue checks if a reflect.Value represents the zero value of its type
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return v.Len() == 0
	case reflect.Struct:
		return v.IsZero()
	default:
		return false
	}
}

// SortObjectKeys recursively sorts object keys in any data structure
func SortObjectKeys(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}

	v := reflect.ValueOf(obj)
	return sortObjectKeysValue(v).Interface()
}

// sortObjectKeysValue recursively sorts keys in reflect.Value
func sortObjectKeysValue(v reflect.Value) reflect.Value {
	// Handle pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return v
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		if v.IsNil() {
			return v
		}

		// Create new map with same type
		newMap := reflect.MakeMap(v.Type())
		
		// Copy all key-value pairs, recursively sorting values
		for _, key := range v.MapKeys() {
			value := sortObjectKeysValue(v.MapIndex(key))
			newMap.SetMapIndex(key, value)
		}
		
		return newMap

	case reflect.Array, reflect.Slice:
		if v.IsNil() {
			return v
		}

		// Create new slice with same type
		newSlice := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		
		// Copy all elements, recursively sorting
		for i := 0; i < v.Len(); i++ {
			element := sortObjectKeysValue(v.Index(i))
			newSlice.Index(i).Set(element)
		}
		
		return newSlice

	case reflect.Struct:
		// Create new struct
		newStruct := reflect.New(v.Type()).Elem()
		
		// Copy all fields, recursively sorting
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanSet() {
				field := sortObjectKeysValue(v.Field(i))
				newStruct.Field(i).Set(field)
			}
		}
		
		return newStruct

	default:
		// For primitive types, return as-is
		return v
	}
} 