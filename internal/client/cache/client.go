package cache

type Client interface {
	Set(key string, value interface{}) bool
	Get(key string) interface{}
	Remove(key string) bool
}
