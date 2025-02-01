package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func Load(s reflect.Value) error {
	for i := 0; i < s.NumField(); i++ {
		field := s.Type().Field(i)
		name := field.Tag.Get("name")

		if name == "" {
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			if err := Load(s.Field(i)); err != nil {
				return fmt.Errorf("load config struct %s: %v", name, err)
			}
			continue
		}

		value := os.Getenv(name)
		if value == "" {
			return fmt.Errorf("the config %s was not found", name)
		}

		switch s.Field(i).Kind() {
		case reflect.String:
			s.Field(i).SetString(value)
		case reflect.Int:
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("failed to convert %s to int: %v", name, err)
			}
			s.Field(i).SetInt(int64(intValue))
		default:
			return fmt.Errorf("unsupported field type: %s", s.Field(i).Kind())
		}
	}
	return nil
}
