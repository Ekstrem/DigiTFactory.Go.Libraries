// Package lifecycle содержит утилиты жизненного цикла объектов.
package lifecycle

import "reflect"

// GetValueObjects извлекает объект-значения из структуры через рефлексию.
// Возвращает словарь с именами типов в качестве ключей и значениями полей.
// Включаются только поля-указатели и интерфейсы, которые не являются nil.
func GetValueObjects(obj any) map[string]any {
	result := make(map[string]any)
	if obj == nil {
		return result
	}

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return result
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Пропускаем приватные поля
		if !fieldType.IsExported() {
			continue
		}

		// Пропускаем примитивные типы
		kind := field.Kind()
		if kind == reflect.Bool || kind == reflect.Int || kind == reflect.Int8 ||
			kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 ||
			kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 ||
			kind == reflect.Uint32 || kind == reflect.Uint64 || kind == reflect.Float32 ||
			kind == reflect.Float64 || kind == reflect.String {
			continue
		}

		// Пропускаем nil-значения
		if (kind == reflect.Ptr || kind == reflect.Interface || kind == reflect.Map ||
			kind == reflect.Slice) && field.IsNil() {
			continue
		}

		result[fieldType.Type.Name()] = field.Interface()
	}

	return result
}
