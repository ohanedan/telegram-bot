package telegram_bot

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/fatih/structs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func runScript(message *tgbotapi.Message, query string) (bool, error) {

	src := fmt.Sprintf("result := %v", query)

	script := tengo.NewScript([]byte(src))

	tengoMap, err := mapToTengoMap(structs.Map(message))
	if err != nil {
		return false, err
	}

	err = script.Add("message", tengoMap)
	if err != nil {
		return false, err
	}
	err = script.Add("query", query)
	if err != nil {
		return false, err
	}

	compiled, err := script.Run()
	if err != nil {
		return false, err
	}

	return compiled.Get("result").Bool(), nil
}

func mapToTengoMap(input map[string]interface{}) (*tengo.Map, error) {

	resultMap := make(map[string]tengo.Object)
	for key, value := range input {

		if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr &&
			reflect.ValueOf(value).IsNil()) {
			continue
		}

		switch val := value.(type) {
		case map[string]interface{}:
			res, err := mapToTengoMap(val)
			if err != nil {
				return nil, err
			}
			resultMap[key] = res
		case string:
			resultMap[key] = &tengo.String{Value: val}
		case int64:
			resultMap[key] = &tengo.Int{Value: val}
		case int:
			intVal := int64(val)
			resultMap[key] = &tengo.Int{Value: intVal}
		case time.Time:
			resultMap[key] = &tengo.Time{Value: val}
		case bool:
			if val {
				resultMap[key] = tengo.TrueValue
			} else {
				resultMap[key] = tengo.FalseValue
			}
		default:
			return nil, errors.New(fmt.Sprintf("unknown type %T", value))
		}
	}

	result := &tengo.Map{
		Value: resultMap,
	}
	return result, nil
}
