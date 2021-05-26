package model

type Container interface {
	Add(key string, value Value)
	Get(key string) (value Value, ok bool)
	Remove()
	Len() int
}