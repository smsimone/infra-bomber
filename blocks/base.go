package blocks

import (
	"context"
	"fmt"
	"it.toduba/bomber/utils"
	"log"
	"regexp"
	"strings"
)

type BaseBlock interface {
	Exec(ctx context.Context) (*map[string]interface{}, error)
}

// ReplacePlaceholders Rimpiazza tutti i placeholder contenuti in s con i valori salvati nel contesto
func ReplacePlaceholders(ctx utils.ContextValue, s string) string {
	placeholders := getPlaceholders(s)

	tmp := s
	for _, placeholder := range placeholders {
		value := (*ctx.Variables)[placeholder]
		if value == nil {
			log.Fatalf("Missing variable '%v' from flow environment", placeholder)
		} else {
			tmp = strings.ReplaceAll(tmp, fmt.Sprintf("{{%v}}", placeholder), value.(string))
		}
	}

	return tmp
}

func getPlaceholders(s string) []string {
	pFinder := regexp.MustCompile(`{{[a-zA-Z0-9-_]*}}`)
	matches := pFinder.FindAllString(s, -1)
	found := make([]string, 0)
	for _, match := range matches {
		match = strings.ReplaceAll(match, "{{", "")
		match = strings.ReplaceAll(match, "}}", "")
		found = append(found, match)
	}
	return found
}
