package remauth

type RemoteAuth interface {
	Valid(token string) bool
}
