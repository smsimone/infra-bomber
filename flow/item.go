package flow

import (
	"encoding/json"
	"fmt"
	"log"
)

type Item struct {
	Request BaseBlock `yaml:"request"`
	Output  *string   `yaml:"output"`
	Name    string    `yaml:"name"`
	CanFail bool      `yaml:"can_fail"`
}

func (i *Item) UnmarshalYAML(unmarshal func(interface{}) error) error {
	generic := new(map[string]interface{})
	if err := unmarshal(generic); err != nil {
		log.Printf("Failed to unmarshal item to map interface: %v", err.Error())
		return err
	}

	i.Name = (*generic)["name"].(string)
	if val, ok := (*generic)["output"]; ok {
		tmp := val.(string)
		i.Output = &tmp
	}

	if val, ok := (*generic)["can_fail"]; ok {
		tmp := val.(bool)
		i.CanFail = tmp
	} else {
		i.CanFail = false
	}

	request := (*generic)["request"].(map[interface{}]interface{})

	var block BaseBlock
	if _, ok := request["command"]; ok {
		tmp := new(ScriptBlock)
		if err := convertToObj(request, tmp); err != nil {
			log.Printf("Failed to convert to script block: %v", err.Error())
			return err
		}
		block = tmp
	} else if _, ok := request["method"]; ok {
		tmp := new(HttpBlock)
		if err := convertToObj(request, tmp); err != nil {
			log.Printf("Failed to convert to http block: %v", err.Error())
			return err
		}
		block = tmp
	} else if _, ok := request["flow"]; ok {
		tmp := new(SubFlow)
		if err := convertToObj(request, tmp); err != nil {
			log.Printf("Failed to convert to http block: %v", err.Error())
			return err
		}
		block = tmp
	}

	if block == nil {
		panic("Failed to parse block")
	}

	i.Request = block

	return nil
}

func convertToObj(data map[interface{}]interface{}, out interface{}) error {
	convertible := convert(data)

	marshaled, err := json.Marshal(convertible)
	if err != nil {
		log.Printf("Failed to marshal item: %v", err.Error())
		return err
	}

	err = json.Unmarshal(marshaled, out)
	if err != nil {
		log.Printf("Failed to unmarshal item: %v", err.Error())
	}

	return err
}

func convert(m map[interface{}]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			res[fmt.Sprint(k)] = convert(v2)
		case []interface{}:
			var data []interface{}
			for _, item := range v2 {
				switch v3 := item.(type) {
				case map[interface{}]interface{}:
					data = append(data, convert(v3))
				case string:
				case int:
					data = append(data, v3)
				default:
					log.Fatalf("Invalid type: %+v", v3)
				}
			}
			res[fmt.Sprint(k)] = data
		default:
			res[fmt.Sprint(k)] = v
		}
	}
	return res
}
