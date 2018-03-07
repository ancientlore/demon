package main

// multiError is a list of errors.
type multiError []error

// Error gets the text of a multiError.
func (m multiError) Error() string {
	var e string
	for i, s := range m {
		if i > 0 {
			e += "\n"
		}
		e += s.Error()
	}
	return e
}
