package cookies

type reader interface {
	read() ([]*http.Cookie, error)
}

func Read(s string) ([]*http.Cookie, error) {
	info, err := os.Stat(s)
	if err != nil {
		return nil, err
	}

	var r reader

	if info.Mode().IsRegular() {
		r = browserReader{
			profile: s,
		}
	} else if info.IsDir() {
		r = netscapeReader{
			path: s,
		}
	} else {
		return nil, fmt.Errorf("")
	}

	return r.read()
}

