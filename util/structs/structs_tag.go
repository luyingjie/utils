package structs

import (
	"reflect"

	"github.com/gqcn/structs"
)

func TagFields(pointer interface{}, priority []string, recursive bool) []*Field {
	return doTagFields(pointer, priority, recursive, map[string]struct{}{})
}

func doTagFields(pointer interface{}, priority []string, recursive bool, tagMap map[string]struct{}) []*Field {
	var fields []*structs.Field
	if v, ok := pointer.(reflect.Value); ok {
		fields = structs.Fields(v.Interface())
	} else {
		var (
			rv   = reflect.ValueOf(pointer)
			kind = rv.Kind()
		)
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}

		if kind == reflect.Ptr && (!rv.IsValid() || rv.IsNil()) {
			fields = structs.Fields(reflect.New(rv.Type().Elem()).Elem().Interface())
		} else {
			fields = structs.Fields(pointer)
		}
	}
	var (
		tag  = ""
		name = ""
	)
	tagFields := make([]*Field, 0)
	for _, field := range fields {
		name = field.Name()
		// Only retrieve exported attributes.
		if name[0] < byte('A') || name[0] > byte('Z') {
			continue
		}
		tag = ""
		for _, p := range priority {
			tag = field.Tag(p)
			if tag != "" {
				break
			}
		}
		if tag != "" {
			// Filter repeated tag.
			if _, ok := tagMap[tag]; ok {
				continue
			}
			tagFields = append(tagFields, &Field{
				Field: field,
				Tag:   tag,
			})
		}
		if recursive {
			var (
				rv   = reflect.ValueOf(field.Value())
				kind = rv.Kind()
			)
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			if kind == reflect.Struct {
				tagFields = append(tagFields, doTagFields(rv, priority, recursive, tagMap)...)
			}
		}
	}
	return tagFields
}

func TagMapName(pointer interface{}, priority []string, recursive bool) map[string]string {
	fields := TagFields(pointer, priority, recursive)
	tagMap := make(map[string]string, len(fields))
	for _, v := range fields {
		tagMap[v.Tag] = v.Name()
	}
	return tagMap
}

func TagMapField(pointer interface{}, priority []string, recursive bool) map[string]*Field {
	fields := TagFields(pointer, priority, recursive)
	tagMap := make(map[string]*Field, len(fields))
	for _, v := range fields {
		tagMap[v.Tag] = v
	}
	return tagMap
}
