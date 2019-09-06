package schema

import (
	"reflect"
	"strings"
)

type UnsupportedKindError struct {
	kind reflect.Kind
}

func (e *UnsupportedKindError) Error() string {
	return "schema: schema cannot contain kind " + e.kind.String()
}

func getJSONName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	name := strings.Split(tag, ",")[0]
	if name == "" {
		return field.Name
	} else if name == "-" {
		return ""
	}
	return name
}

func isOmitEmpty(field reflect.StructField) bool {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return false
	}
	tags := strings.Split(jsonTag, ",")[1:]
	for _, tag := range tags {
		if tag == "omitempty" {
			return true
		}
	}
	return false
}

type schemaTags map[string]string

func newSchemaTags(tag reflect.StructTag) schemaTags {
	schemaTags := schemaTags{}
	tags := strings.Split(tag.Get("schema"), ";")
	for _, t := range tags {
		if t == "" {
			continue
		}
		if strings.Contains(t, ":") {
			parts := strings.SplitN(t, ":", 2)
			schemaTags[parts[0]] = parts[1]
		} else {
			schemaTags[t] = "true"
		}
	}
	return schemaTags
}

func (t schemaTags) itemsTags() schemaTags {
	itemsTags := schemaTags{}
	for k, v := range t {
		if strings.HasPrefix(k, "items_") {
			itemsTags[strings.TrimPrefix(k, "items_")] = v
		}
	}
	return itemsTags
}
