package main

type SlacksocError struct {
	message string
}

func (e *SlacksocError) Error() string {
	return e.message
}
