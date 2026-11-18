package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"bufio"
	"strconv"
	"strings"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"
	_ "modernc.org/sqlite"
)

type cookiesReader interface {
	read() ([]*http.Cookie, error)
}


type browserReader struct {
	profile string
}

type netscapeReader struct {
	path string
}

func (self *netscapeReader) read([]*http.Cookie, error) {
	data, err := os.ReadFile(self.path)
	if err != nil {
		return nil, err
	}

	var result []*http.Cookie

	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line, "\t")
		if len(fields) != 7 {
			continue
		}

		cookie := &http.Cookie{
			Domain:   fields[0],
			Path:     fields[2],
			Secure:   fields[3] == "TRUE",
			Name:     fields[5],
			Value:    fields[6],
			HttpOnly: false,
		}

		expires, err := strconv.ParseInt(fields[4], 10, 64)
		if err != nil {
			return nil, err
		}

		if expires > 0 {
			cookie.Expires = time.Unix(expires, 0)
		}

		result = append(result, cookie)
	}

	return result, scanner.Err()
}

func (self *browserReader) copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	err = os.MkdirAll(filepath.Dir(dst))
	if err != nil {
		return nil, err
	}

	out, err := os.OpenFile(
	    dst,
	    os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
	    0600,
	)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return nil, err
	}

	return out.Close()
}

func (self *browserReader) cloneDB() (string, error) {
	src := filepath.Join(self.profile, "cookies.sqlite")

	dst, err := os.MkdirTemp("/dev/shm", "*")
	if err != nil {
		return nil, err
	}

	dst = filepath.Join(dst, filepath.Base(src))

	err := self.copyFile(src, dst)
	if err != nil {
		return nil, err
	}

	err = self.copyFile(src+"-wal", dst+"-wal")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	err = self.copyFile(src+"-shm", dst+"-shm")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return  nil, err
	}

	return dst, nil
}

func (self *browserReader) openClonedDB() (*sql.DB, func(), error) {
	path, err := self.cloneDB()
	if err != nil {
		return nil, nil, err
	}

	clone := sql.Open(
		"sqlite",
		fmt.Sprintf("file:%s?mode=ro&immutable=1", path),
	)

	cleanup := func() {
		clone.Close()
		os.RemoveAll(filepath.Dir(path))
	}

	return clone, cleanup, nil
}

func (self *browserReader) read() ([]*http.Cookie, error) {
	db, cleanup, err := self.openClonedDB()
	if err != nil {
		return nil, err
	}
	defer cleanup()

	rows, err := db.Query(`
		SELECT
			name,
			value,
			host,
			path,
			isSecure,
			expiry
		FROM
			moz_cookies
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*http.Cookie

	for rows.Next() {
		var (
			c http.Cookie

			expiry int64
		)

		err = rows.Scan(
			&c.Name,
			&c.Value,
			&c.Domain,
			&c.Path,
			&c.Secure,
			&expiry,
		)
		if err != nil {
			return nil, err
		}

		c.Expires = time.Unix(expiry, 0)

		result = append(result, &c)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ReadCookies(s string) ([]*http.Cookie, error) {
	info, err := os.Stat(s)
	if err != nil {
		return nil, err
	}

	var reader cookiesReader

	if info.IsDir() {
		reader = browserReader{
			profile: s,
		}
	} else {
		reader = netscapeReader{
			path: s,
		}
	} 

	return reader.read()
}
