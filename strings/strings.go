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
func BuildSnakeName(s string) string {
	s2 := strings.ToLower(s)
	s3 := ""
	for i := range s {
		if s2[i] != s[i] {
			s3 += "_" + string(s2[i])
		} else {
			s3 += string(s2[i])
		}
	}
	if string(s3[0]) == "_" {
		return s3[1:]
	}
	return s3
}
func UnBuildSnakeName(s string) string {
	s2 := strings.ToUpper(s)
	s1 := string(s[0])
	for i := 1; i < len(s); i++ {
		if string(s[i-1]) == "_" {
			s1 = s1[:len(s1)-1]
			s1 += string(s2[i])
		} else {
			s1 += string(s[i])
		}
	}
	return s1
}