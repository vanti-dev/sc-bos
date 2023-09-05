package hpd3

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

// for use as a mock Client for tests
type memoryClient struct {
	points PointData
}

func (c *memoryClient) FetchSensorData(_ context.Context, pointName string, out *PointData) error {
	return copyFieldByJsonName(out, c.points, pointName)
}

// copies a single field from src to dest identified by its JSON tag (or field name if no JSON tag exists)
func copyFieldByJsonName(dest, src any, fieldName string) error {
	destR := reflect.ValueOf(dest).Elem()
	srcR := reflect.ValueOf(src)
	if srcR.Kind() == reflect.Ptr || srcR.Kind() == reflect.Interface {
		srcR = srcR.Elem()
	}

	srcField, ok := findJsonField(srcR.Type(), fieldName)
	if !ok {
		return fmt.Errorf("src has no field with JSON name %q", fieldName)
	}
	dstField, ok := findJsonField(destR.Type(), fieldName)
	if !ok {
		return fmt.Errorf("dest has no field with JSON name %q", fieldName)
	}

	srcFieldValue := srcR.FieldByIndex(srcField.Index)
	dstFieldValue := destR.FieldByIndex(dstField.Index)
	dstFieldValue.Set(srcFieldValue)
	return nil
}

func findJsonField(t reflect.Type, fieldName string) (reflect.StructField, bool) {
	for idx := 0; idx < t.NumField(); idx++ {
		field := t.Field(idx)
		// find the effective JSON name of the field
		jsonName, _, _ := strings.Cut(field.Tag.Get("json"), ",") // remove stuff like ",omitempty"
		if jsonName == "-" {
			// ignored field
			continue
		} else if jsonName == "" {
			// json name defaults to field name
			jsonName = field.Name
		}

		if jsonName == fieldName {
			return field, true
		}
	}
	return reflect.StructField{}, false
}
