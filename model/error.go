package model

type ErrNotFound struct {
	N   int
	Err error
}

func (e ErrNotFound) Error() string {
	return e.Err.Error()
}
