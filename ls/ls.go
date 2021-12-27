package ls

import "github.com/bronger/jfy/lib"

func Handle(settings lib.SettingsType, stdout, stderr []byte, args ...string) (any, any, error) {
	return map[string]string{"Hallö": "Süpi"}, nil, nil
}
