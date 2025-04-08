package types

import "regexp"

type Rule struct {
	ID        string
	Severity  string
	FileTypes []string
	Pattern   *regexp.Regexp
	Title     string
	Advice    string
}
