package f5

type F5Metadata struct {
	Name    string `@( Ident | QF5Name | F5Name | QString ) "{"`
	Value   string `(  "value" @( Ident | QF5Name | F5Name | QString )`
	Persist string ` | "persist" @( "true" | "false" ) )* "}"`
}
