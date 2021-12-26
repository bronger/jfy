package ls

func Handle(stdout, stderr []byte, args ...string) any {
	return map[string]string{"Hallö": "Süpi"}
}
