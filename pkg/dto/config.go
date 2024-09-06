package dto

type GlobalConf struct {
	Verbose          bool
	SuppressWarnings bool
	ErrorOnWarning   bool
	Files            []string
	Command          string
}
