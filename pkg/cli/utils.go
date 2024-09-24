package cli

import (
	"fmt"
	"log"
	"reflect"
)

func FlagMustBeMandatory(flagName string, flagValue string) {
	if flagValue == "" {
		log.Fatalf("Flag \"%s\" is mandatory", flagName)
	}
}

func StringToInt(str string) (int, error) {
	var i int
	_, err := fmt.Sscanf(str, "%d", &i)
	return i, err
}

func DisplayPrettyStruct(data any) {
	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Slice {
		log.Fatalf("Data is not a slice")
	}

	for i := 0; i < val.Len(); i++ {
		fmt.Printf("%d. %+v\n", i+1, val.Index(i))
	}
}
