package model

type Value interface {
	Len() int	// 返回 value 占用内存的大小
}
