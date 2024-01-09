package err

import "fmt"

type ConfigParseError struct {
	missingFields []string
}

func ConfigParseErrorInit(input string) *ConfigParseError {
	return &ConfigParseError{
		missingFields: []string{input},
	}
}

func (c *ConfigParseError) AppendOrCreate(input string) *ConfigParseError {
	if c == nil {
		c = ConfigParseErrorInit(input)
	} else {
		c.AddMissingField(input)
	}
	return c
}

func (c *ConfigParseError) AddMissingField(name string) {
	c.missingFields = append(c.missingFields, name)
}

func (c *ConfigParseError) Error() string {
	res := "Missing/empty mandatory fields: \n"
	for _, field := range c.missingFields {
		res += fmt.Sprintf("-> %s\n", field)
	}
	return res
}
