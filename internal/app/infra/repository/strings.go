package repository

func contains(s string) string {
	return "%" + s + "%"
}

func startWith(s string) string {
	return s + "%"
}

func endWith(s string) string {
	return "%" + s
}
