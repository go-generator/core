package strings

import "strings"

var ends = []string{"ies", "ees", "ses", "xes", "zes", "shes", "ches", "aes", "oes", "ues"}
var replaces = []string{"y", "ee", "s", "x", "z", "sh", "ch", "ae", "oe", "ue"}
var plural = []string{"people", "women", "men", "fungus", "feet", "teeth"}
var singular = []string{"person", "woman", "man", "fungi", "foot", "tooth"}

func ToSingular(s string) string {
	if len(s) <= 1 {
		return s
	}
	for i, si := range plural {
		if strings.HasSuffix(s, si) {
			return s[0:len(s) - len (si)] + singular[i]
		}
	}
	for i, si := range ends {
		if strings.HasSuffix(s, si) {
			return s[0:len(s) - len (si)] + replaces[i]
		}
	}
	if strings.HasSuffix(s, "s") {
		return s[0: len(s) - 1]
	}
	return s
}
func ToPlural(s string) string {
	if len(s) <= 1 {
		return s
	}
	for i, si := range singular {
		if strings.HasSuffix(s, si) {
			return s[0:len(s) - len (si)] + plural[i]
		}
	}
	x := s[len(s) - 1:]
	if x == "y" {
		return s[0: len(s) - 1] + "ies"
	}
	if x == "s" || x == "x" || strings.HasSuffix(s, "sh") || strings.HasSuffix(s, "ch") {
		return s[0:] + "es"
	}
	return s[0:] + "s"
}
