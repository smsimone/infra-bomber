package flow

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"it.toduba/bomber/utils"
)

type ScriptBlock struct {
	BaseBlock `yaml:"-"`
	Env       *map[string]string `yaml:"env"`
	Command   string             `yaml:"command"`
	Args      []string           `yaml:"args"`
}

func (s *ScriptBlock) Exec(ctx context.Context) (*map[string]interface{}, error) {
	ctxVal := utils.GetContextValues(ctx)

	stepName := ctxVal.StepName
	outputName := ctxVal.OutputName

	command := exec.Command(s.Command)
	for _, arg := range s.Args {
		command.Args = append(command.Args, ReplacePlaceholders(ctxVal, arg))
	}

	if s.Env != nil {
		for k, v := range *s.Env {
			command.Env = append(command.Env, fmt.Sprintf("%v=%v", k, ReplacePlaceholders(ctxVal, v)))
		}
	}
	out, err := command.Output()
	if err != nil {
		log.Fatalf("[%v] Failed to execute command: %v -> %v", stepName, err.Error(), command.Err)
		return nil, err
	}

	if outputName == nil {
		log.Printf("[%v] Executed script", stepName)
		return nil, nil
	}

	otp := strings.TrimSpace(string(out))
	log.Printf("[%v] Got '%v' output", stepName, *outputName)

	return &map[string]interface{}{*outputName: otp}, nil
}
