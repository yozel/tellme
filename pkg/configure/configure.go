package configure

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type info struct {
	v            reflect.Value
	name         string
	env          string
	defaultValue string
	usage        string
	required     bool
}

func Parse(c any) error {
	typeOfC := reflect.TypeOf(c)
	if typeOfC.Kind() != reflect.Ptr {
		return fmt.Errorf("c must be a pointer")
	}

	typeOfC = typeOfC.Elem()

	results := []info{}

	for i := 0; i < typeOfC.NumField(); i++ {
		field := typeOfC.Field(i)
		tagValue, ok := field.Tag.Lookup("configure")
		if !ok {
			continue
		}

		info := info{}
		info.v = reflect.ValueOf(c).Elem().Field(i)
		for _, item := range strings.Split(tagValue, ";") {
			parts := strings.SplitN(item, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid tag %s", tagValue)
			}
			key, value := parts[0], parts[1]
			switch key {
			case "name":
				info.name = value
			case "env":
				info.env = value
			case "default":
				info.defaultValue = value
			case "usage":
				info.usage = value
			case "required":
				info.required = value == "true"
			default:
				return fmt.Errorf("invalid tag %s", tagValue)
			}
		}
		results = append(results, info)
	}

	for _, info := range results {
		if info.env != "" {
			if val, ok := os.LookupEnv(info.env); ok {
				info.defaultValue = val
			}
		}
		switch info.v.Kind() {
		case reflect.String:
			flag.StringVar(info.v.Addr().Interface().(*string), info.name, info.defaultValue, info.usage)
		case reflect.Int:
			if info.defaultValue == "" {
				info.defaultValue = "0"
			}
			def, err := strconv.Atoi(info.defaultValue)
			if err != nil {
				return err
			}
			flag.IntVar(info.v.Addr().Interface().(*int), info.name, def, info.usage)
		case reflect.Int64:
			if info.defaultValue == "" {
				info.defaultValue = "0"
			}
			def, err := strconv.ParseInt(info.defaultValue, 10, 64)
			if err != nil {
				return err
			}
			flag.Int64Var(info.v.Addr().Interface().(*int64), info.name, def, info.usage)
		case reflect.Bool:
			if info.defaultValue == "" {
				info.defaultValue = "false"
			}
			def, err := strconv.ParseBool(info.defaultValue)
			if err != nil {
				return err
			}
			flag.BoolVar(info.v.Addr().Interface().(*bool), info.name, def, info.usage)
		case reflect.Float64:
			if info.defaultValue == "" {
				info.defaultValue = "0"
			}
			def, err := strconv.ParseFloat(info.defaultValue, 64)
			if err != nil {
				return err
			}
			flag.Float64Var(info.v.Addr().Interface().(*float64), info.name, def, info.usage)
		default:
			return fmt.Errorf("unsupported type %s", info.v.Kind())
		}
	}
	flag.Parse()

	for _, info := range results {
		if info.required && info.v.IsZero() {
			return fmt.Errorf("required flag %s is not set", info.name)
		}
	}
	return nil
}
