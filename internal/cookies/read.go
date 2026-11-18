package cookies

type reader interface {
	read(string) ([]*http.Cookie, error)
}

func Read(s string) ([]*http.Cookie, error) {
	info, err := os.Stat(s)
	if err != nil {
		return nil, err
	}

	var r reader

	if info.Mode().IsRegular() {
		r = netscapeReader{}
	} else if info.IsDir() {
		r = firefoxReader{}
	} else {
		return nil, fmt.Errorf("")
	}

	return r.read(s)
}

