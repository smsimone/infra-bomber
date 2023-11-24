package flow

import (
	"context"
	"log"
	"os"

	"it.toduba/bomber/blocks"
	"it.toduba/bomber/enums"
	"it.toduba/bomber/utils"

	"gopkg.in/yaml.v2"
)

type Flow struct {
	Name        string             `yaml:"name"`
	BaseUrl     string             `yaml:"baseUrl"`
	Environment *map[string]string `yaml:"environment"`
	Steps       []Item             `yaml:"steps"`
}

func ParseFromYaml(yamlPath string) (*Flow, error) {
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Printf("Failed to read yaml file: %v", err.Error())
		return nil, err
	}

	flow := new(Flow)
	if err := yaml.Unmarshal(data, flow); err != nil {
		log.Printf("Failed to unmarshal yaml file: %v", err.Error())
		return nil, err
	}

	return flow, nil
}

func (f *Flow) Execute(envVars *map[string]string) {
	ctxVal := utils.ContextValue{
		BaseUrl:   f.BaseUrl,
		Variables: &map[string]interface{}{},
	}

	if envVars != nil {
		for k, v := range *envVars {
			(*ctxVal.Variables)[k] = v
		}
	}

	if f.Environment != nil {
		for k, v := range *f.Environment {
			(*ctxVal.Variables)[k] = blocks.ReplacePlaceholders(ctxVal, v)
		}
	}

	ctx := context.WithValue(context.Background(), enums.Values, ctxVal)

	for _, item := range f.Steps {
		ctxVal := utils.GetContextValues(ctx)
		ctxVal.StepName = item.Name
		ctxVal.OutputName = item.Output
		ctx = context.WithValue(context.Background(), enums.Values, ctxVal)

		out, err := item.Request.Exec(ctx)
		if err != nil {
			log.Fatalf("Failed to run step '%v': %v", item.Name, err.Error())
		}
		if out != nil {
			ctxVal := utils.GetContextValues(ctx)
			for k, v := range *out {
				(*ctxVal.Variables)[k] = v
			}
			ctx = context.WithValue(context.Background(), enums.Values, ctxVal)
		}
	}
}
