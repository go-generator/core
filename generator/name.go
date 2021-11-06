package generator

import "strings"

func BuildNames(name string, options ...func(string) string) map[string]string {
	var toPlural func(string) string
	if len(options) > 0 {
		toPlural = options[0]
	}
	n := make(map[string]string)
	var raw string
	if !strings.Contains(name, "_") {
		raw = BuildSnakeName(name)
	} else {
		raw = strings.ToLower(name)
		name = UnBuildSnakeName(raw)
	}
	n = map[string]string{
		"raw":      raw,
		"name":     strings.ToLower(string(name[0])) + name[1:],
		"Name":     strings.ToUpper(string(name[0])) + name[1:],
		"NAME":     strings.ToUpper(name),
		"constant": strings.ToUpper(raw),
		"lower":    strings.ToLower(name),
	}
	if toPlural == nil {
		return n
	}
	raws := toPlural(raw)
	names := UnBuildSnakeName(raws)
	n["raws"] = raws
	n["names"] = strings.ToLower(string(names[0])) + names[1:]
	n["Names"] = strings.ToUpper(string(names[0])) + names[1:]
	n["NAMES"] = strings.ToUpper(names)
	n["constants"] = strings.ToUpper(raws)
	n["lowers"] = strings.ToLower(names)
	return n
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
func InitEnv(env map[string]string, projectName string) map[string]string {
	init, ok := env["init"]
	if ok {
		if init == "true" {
			outMap := make(map[string]string)
			for k, v := range env {
				tmp := buildEnvNames(k, v)
				for k1, v1 := range tmp {
					outMap[k1] = v1
				}
			}
			tmp := buildProjectName(projectName)
			for k1, v1 := range tmp {
				outMap[k1] = v1
			}
			return outMap
		}
	}
	env["project"] = projectName
	return env
}
func buildProjectName(name string) map[string]string {
	var raw string
	if !strings.Contains(name, "_") {
		raw = BuildSnakeName(name)
	} else {
		raw = strings.ToLower(name)
		name = UnBuildSnakeName(raw)
	}
	return map[string]string{
		"project_raw":      raw,
		"project":          strings.ToLower(string(name[0])) + name[1:],
		"Project":          strings.ToUpper(string(name[0])) + name[1:],
		"project_lower":    strings.ToLower(name),
		"project_name":     BuildSnakeName(name),
		"project_constant": strings.ToUpper(raw),
	}
}
func buildEnvNames(name, v string) map[string]string {
	names := map[string]string{
		name:            v,
		name + "_name":  strings.ToLower(string(v[0])) + v[1:],
		name + "_Name":  strings.ToUpper(string(v[0])) + v[1:],
		name + "_NAME":  strings.ToUpper(v),
		name + "_lower": strings.ToLower(v),
	}
	return names
}
