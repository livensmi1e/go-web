package ports

type Validator interface {
	Validate(i interface{}) error
}
