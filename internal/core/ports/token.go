package ports

type TokenGenerator interface {
	Generate(claims map[string]interface{}) (string, error)
	Validate(token string) (map[string]interface{}, error)
}
