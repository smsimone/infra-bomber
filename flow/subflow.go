package flow

import (
	"context"
	"log"

	"it.toduba/bomber/utils"
)

type SubFlow struct {
	BaseBlock   `yaml:"-"`
	Environment *map[string]string `yaml:"environment"`
	Flow        string             `yaml:"flow"`
}

func (s *SubFlow) Exec(ctx context.Context) (*map[string]interface{}, error) {
	f, err := ParseFromYaml(s.Flow)
	if err != nil {
		log.Printf("Failed to parse sub flow: %v", err.Error())
		return nil, err
	}

	ctxVal := utils.GetContextValues(ctx)

	cleanedEnv := make(map[string]string)
	for k, v := range *s.Environment {
		cleanedEnv[k] = ReplacePlaceholders(ctxVal, v)
	}

	if err := f.Execute(&cleanedEnv); err != nil {
		log.Printf("Failed to execute subflow: %v", err.Error())
		return nil, err
	}

	return nil, nil
}
