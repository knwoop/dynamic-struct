package _examples

import (
	"fmt"
	dynamicstruct "github.com/kntaka/dynamic-struct"
	"gopkg.in/yaml.v2"
	"log"
	"reflect"
	"time"
)

func getYamlExample() interface{} {
	integer := 0
	uinteger := uint(0)
	str := ""
	float := 0.0
	boolean := false

	subInstance := dynamicstruct.NewStruct().
		AddField("Integer", integer, "").
		AddField("Text", str, `yaml:"subText"`).
		Build().
		New()

	instance := dynamicstruct.NewStruct().
		AddField("Integer", integer, `yaml:"int" validate:"lt=123"`).
		AddField("Uinteger", uinteger, `validate:"gte=0"`).
		AddField("Text", str, `yaml:"someText"`).
		AddField("Float", float, `yaml:"double"`).
		AddField("Boolean", boolean, "").
		AddField("Time", time.Time{}, "").
		AddField("Slice", []int{}, "").
		AddField("PointerInteger", &integer, "").
		AddField("PointerUinteger", &uinteger, "").
		AddField("PointerText", &str, "").
		AddField("PointerFloat", &float, "").
		AddField("PointerBoolean", &boolean, "").
		AddField("PointerTime", &time.Time{}, "").
		AddField("SubStruct", subInstance, `yaml:"subData"`).
		AddField("Anonymous", "", `yaml:"-" validate:"required"`).
		Build().
		New()

	data := []byte(`
{
	"int": 123,
	"Uinteger": 456,
	"someText": "example",
	"double": 123.45,
	"Boolean": true,
	"Time": "2018-12-27T19:42:31+07:00",
	"Slice": [1, 2, 3],
	"PointerInteger": 345,
	"PointerUinteger": 234,
	"PointerFloat": 567.89,
	"PointerText": "pointer example",
	"PointerBoolean": true,
	"PointerTime": "2018-12-28T01:23:45+07:00",
	"subData": {
		"Integer": 12,
		"subText": "sub example"
	},
	"Anonymous": "avoid to read"
}
`)

	err := yaml.Unmarshal([]byte(data), &instance)
	if err != nil {
		return err
	}

	fmt.Printf("--- type:\n%v\n\n", reflect.TypeOf(instance))
	fmt.Printf("--- instance:\n%v\n\n", instance)

	yml, err := yaml.Marshal(&instance)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return yml
}