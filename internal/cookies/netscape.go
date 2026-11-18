package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"bufio"
	"strconv"
	"bytes"
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
	content, err := os.ReadFile(self.path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))

	var result []*http.Cookie

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || line[0] == '#' {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 7 {
			continue
		}

		cookie := &http.Cookie{
			Domain:   fields[0],
			HttpOnly: strings.HasPrefix(
				fields[0],
				"#HttpOnly_",
			),
			Path:     fields[2],
			Secure:   fields[3] == "TRUE",
			Name:     fields[5],
			Value:    fields[6],
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

