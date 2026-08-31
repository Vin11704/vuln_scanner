package main

type Finding struct{
	RuleID string
	Filename string
	Severity string
	Line int // which line of the file
	Message string // explanation of the finding
}

type Rule interface{
	ID() string
	AppliesTo(path string) bool
	Check(path string, content []byte) []Finding
}

