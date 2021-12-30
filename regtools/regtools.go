package regtools

import "regexp"

func GetParams(reg *regexp.Regexp, text string) (paramsMap map[string]string) {
	match := reg.FindStringSubmatch(text)

	paramsMap = make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}
