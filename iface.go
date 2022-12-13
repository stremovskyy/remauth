package remauth

type RemoteAuth interface {
	Valid(token string) bool
	Check(token string, result func(response interface{}, err error))
}
