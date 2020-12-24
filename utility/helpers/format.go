package helpers

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatFloatNumber(value interface{}) string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%.2f", value)
}
