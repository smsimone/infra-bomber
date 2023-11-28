package blocks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"it.toduba/bomber/enums"

	"it.toduba/bomber/utils"
)

type HttpBlock struct {
	BaseBlock          `yaml:"-"`
	Body               *map[string]interface{} `yaml:"body"`
	Headers            *map[string]string      `yaml:"headers"`
	BodySelector       *string                 `yaml:"bodySelector"`
	Path               string                  `yaml:"path"`
	Method             string                  `yaml:"method"`
	ExpectedStatusCode int                     `yaml:"expectedStatusCode"`
}

func (s *HttpBlock) Exec(ctx context.Context) (*map[string]interface{}, error) {
	ctxVal := utils.GetContextValues(ctx)

	stepName := ctxVal.StepName

	url := fmt.Sprintf("%v%v", ctxVal.BaseUrl, s.Path)

	var content io.Reader
	if s.Body != nil {
		body := prepareRequestBody(ctx, s.Body)
		jsonVal, err := json.Marshal(body)
		if err != nil {
			log.Printf("Failed to sanitize input: %v", err)
			return nil, err
		}
		content = bytes.NewBuffer(jsonVal)
	}

	req, err := http.NewRequest(s.Method, ReplacePlaceholders(ctxVal, url), content)
	if err != nil {
		log.Printf("[%v] Failed to build request: %v", stepName, err.Error())
		return nil, err
	}

	if s.Headers != nil {
		for k, val := range *s.Headers {
			req.Header.Add(k, ReplacePlaceholders(ctxVal, val))
		}
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if err != nil {
		log.Printf("[%v] Failed to send http request: %v", stepName, err.Error())
		return nil, err
	} else if resp.StatusCode != s.ExpectedStatusCode {
		log.Printf("[%v] Received invalid status code: expected %v - got %v", stepName, s.ExpectedStatusCode, resp.StatusCode)
		return nil, fmt.Errorf("received invalid status code")
	}

	var j interface{}

	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		log.Printf("[%v] Failed to convert body to json: %v", stepName, err.Error())
		return nil, err
	} else {
		log.Printf("[%v] Got response", stepName)
		if content == nil || s.BodySelector == nil {
			return nil, nil
		}
		return getOutput(ctx, j, s.BodySelector), nil
	}
}

func prepareRequestBody(ctx context.Context, body *map[string]interface{}) map[string]interface{} {
	if body == nil {
		return nil
	}

	sanitized := make(map[string]interface{})
	for key, val := range *body {
		if parsed, ok := val.(string); ok {
			sanitized[key] = ReplacePlaceholders(ctx.Value(enums.Values).(utils.ContextValue), parsed)
		} else if parsed, ok := val.(map[string]interface{}); ok {
			sanitized[key] = prepareRequestBody(ctx, &parsed)
		} else {
			sanitized[key] = val
		}
	}
	return sanitized
}

func getOutput(ctx context.Context, body interface{}, selector *string) *map[string]interface{} {
	ctxVal := utils.GetContextValues(ctx)

	outputName := ctxVal.OutputName
	stepName := ctxVal.StepName

	var data map[string]interface{}
	if decoded, ok := body.(map[string]interface{}); ok {
		data = decoded
	} else {
		log.Printf("[%v] Failed to unmarshal body", stepName)
		return nil
	}

	if selector != nil && outputName != nil {
		return &map[string]interface{}{
			*outputName: data[*selector],
		}
	}
	return nil
}
