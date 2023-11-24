package context_utilities

import (
	"context"

	"it.toduba/bomber/enums"
)

type ContextValue struct {
	StepName   string
	OutputName *string
	BaseUrl    string
	Variables  *map[string]interface{}
}

func GetContextValues(ctx context.Context) ContextValue {
	return ctx.Value(enums.Values).(ContextValue)
}
