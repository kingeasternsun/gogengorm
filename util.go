package main

import (
	"strings"

	"github.com/fatih/camelcase"
)

func convertFieldName(transform, fieldName string) (name string, unknown bool) {
	name = ""
	unknown = false

	splitted := camelcase.Split(fieldName)

	switch transform {
	case "snakecase":
		var lowerSplitted []string
		for _, s := range splitted {
			lowerSplitted = append(lowerSplitted, strings.ToLower(s))
		}

		name = strings.Join(lowerSplitted, "_")
	case "lispcase":
		var lowerSplitted []string
		for _, s := range splitted {
			lowerSplitted = append(lowerSplitted, strings.ToLower(s))
		}

		name = strings.Join(lowerSplitted, "-")
	case "camelcase":
		var titled []string
		for _, s := range splitted {
			titled = append(titled, strings.Title(s))
		}

		titled[0] = strings.ToLower(titled[0])

		name = strings.Join(titled, "")
	case "pascalcase":
		var titled []string
		for _, s := range splitted {
			titled = append(titled, strings.Title(s))
		}

		name = strings.Join(titled, "")
	case "keep":
		name = fieldName
	default:
		unknown = true
	}

	return
}
