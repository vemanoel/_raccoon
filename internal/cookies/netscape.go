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

