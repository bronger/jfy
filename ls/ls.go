package ls

func Handle(stdout, stderr []byte, args ...string) (any, any, error) {
	return map[string]string{"Hallö": "Süpi"}, nil, nil
}
