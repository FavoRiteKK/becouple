package main

func ConcateErrorWith(errs []error, delim string) string {
	var s string
	for i, e := range errs {
		s += e.Error()
		if i < len(errs)-1 {
			s += "\n"
		}
	}

	return s
}
