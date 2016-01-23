package main

type RTMStartError struct {
	message string
}

func (e *RTMStartError) Error() string {
	return e.message
}
