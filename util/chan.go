package util


func NewChan()chan interface{}{
	c := make(chan interface{})
	return c
}
