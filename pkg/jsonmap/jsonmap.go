// Copyright © 2022 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsonmap

import (
	"reflect"
	"strings"
)

func StructToJSONMap(i interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	AddJSONFieldsToMap(reflect.ValueOf(i), m)
	return m
}

// AddJSONFieldsToMap is a helper for marshalling struct fields down into a map
func AddJSONFieldsToMap(val reflect.Value, data map[string]interface{}) {
	varType := val.Type()
	if varType.Kind() == reflect.Ptr {
		AddJSONFieldsToMap(val.Elem(), data)
		return
	}
	for i := 0; i < varType.NumField(); i++ {
		f := val.Field(i)
		fType := varType.Field(i)
		if fType.Anonymous {
			AddJSONFieldsToMap(f, data)
			continue
		}
		if !f.CanInterface() {
			continue
		}
		tag, ok := varType.Field(i).Tag.Lookup(`json`)
		var fieldName string
		if ok && len(tag) > 0 {
			tagValues := strings.Split(tag, ",")
			tagName := tagValues[0]
			if tagName == "-" || len(tagValues) > 1 && tagValues[1] == "omitempty" && isEmptyValue(f) {
				continue
			}
			fieldName = tagName
		} else {
			fieldName = fType.Name
		}
		data[fieldName] = f.Interface()
	}
}

// had to copy these rules over from json as not exposed
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
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
	}
	return false
}
