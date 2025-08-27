package model

type CodeParams struct {
	Language     string   `json:"language" jsonschema:"language to run"`
	Code         string   `json:"code" jsonschema:"source code to execute"`
	Stdin        string   `json:"stdin,omitempty" jsonschema:"optional standard input"`
	Dependencies []string `json:"dependencies,omitempty" jsonschema:"list of packages to install"`
}
