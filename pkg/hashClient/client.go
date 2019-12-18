package hashClient

type HashClient interface {
	Set(key, field, value string)
	Get(key, field string) string
}
