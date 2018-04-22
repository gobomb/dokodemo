package util

func NewChan(i int) chan interface{} {
	c := make(chan interface{}, i)
	return c
}
