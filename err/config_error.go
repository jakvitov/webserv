package err

import "fmt"

type ConfigParseError struct {
	missingFields []string
}

func ConfigParseErrorInit(input []string) *ConfigParseError {
	return &ConfigParseError{
		missingFields: input,
	}
}

func (c *ConfigParseError) Error() string {
	res := "Missing/empty mandatory fields: \n"
	for _, field := range c.missingFields {
		res += fmt.Sprintf("-> %s\n", field)
	}
	return res
}
