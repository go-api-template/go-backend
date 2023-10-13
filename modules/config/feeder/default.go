package feeder

import (
	"errors"
	"fmt"
	"github.com/golobby/cast"
	"reflect"
	"time"
	"unsafe"
)

// Default feeder.
// It feeds the fields using default value if tag exists.
//
// Example:
//
//	type Config struct {
//		Host string `default:"localhost"`
//		Port int `default:"8080"`
//	}
type Default struct{}

// Feed feeds the structure with the default value.
func (f Default) Feed(structure interface{}) error {
	inputType := reflect.TypeOf(structure)
	if inputType != nil {
		if inputType.Kind() == reflect.Ptr {
			if inputType.Elem().Kind() == reflect.Struct {
				return f.fillStruct(reflect.ValueOf(structure).Elem())
			}
		}
	}

	return errors.New("env: invalid structure")
}

// fillStruct sets a reflected struct fields with the default value.
func (f Default) fillStruct(s reflect.Value) error {
	for i := 0; i < s.NumField(); i++ {
		if t, exist := s.Type().Field(i).Tag.Lookup("default"); exist {
			// handle time.Duration
			if s.Type().Field(i).Type.String() == "time.Duration" {
				d, err := time.ParseDuration(t)
				if err != nil {
					return fmt.Errorf("default: cannot set `%v` field; err: %v", s.Type().Field(i).Name, err)
				}
				// #nosec G103
				ptr := reflect.NewAt(s.Field(i).Type(), unsafe.Pointer(s.Field(i).UnsafeAddr())).Elem()
				ptr.Set(reflect.ValueOf(d))
			} else {
				v, err := cast.FromType(t, s.Type().Field(i).Type)
				if err != nil {
					return fmt.Errorf("default: cannot set `%v` field; err: %v", s.Type().Field(i).Name, err)
				}
				// #nosec G103
				ptr := reflect.NewAt(s.Field(i).Type(), unsafe.Pointer(s.Field(i).UnsafeAddr())).Elem()
				ptr.Set(reflect.ValueOf(v))
			}
		} else if s.Type().Field(i).Type.Kind() == reflect.Struct {
			if err := f.fillStruct(s.Field(i)); err != nil {
				return err
			}
		} else if s.Type().Field(i).Type.Kind() == reflect.Ptr {
			if !s.Field(i).IsZero() && s.Field(i).Elem().Type().Kind() == reflect.Struct {
				if err := f.fillStruct(s.Field(i).Elem()); err != nil {
					return err
				}
			}
		}
	}

	/*

		for i := 0; i < s.NumField(); i++ {
			if t, exist := s.Type().Field(i).Tag.Lookup("env"); exist {
				if osv := os.Getenv(t); osv != "" {
					v, err := cast.FromType(osv, s.Type().Field(i).Type)
					if err != nil {
						return fmt.Errorf("env: cannot set `%v` field; err: %v", s.Type().Field(i).Name, err)
					}

					ptr := reflect.NewAt(s.Field(i).Type(), unsafe.Pointer(s.Field(i).UnsafeAddr())).Elem()
					ptr.Set(reflect.ValueOf(v))
				}
			} else if s.Type().Field(i).Type.Kind() == reflect.Struct {
				if err := fillStruct(s.Field(i)); err != nil {
					return err
				}
			} else if s.Type().Field(i).Type.Kind() == reflect.Ptr {
				if s.Field(i).IsZero() == false && s.Field(i).Elem().Type().Kind() == reflect.Struct {
					if err := fillStruct(s.Field(i).Elem()); err != nil {
						return err
					}
				}
			}
		}
	*/

	return nil
}
