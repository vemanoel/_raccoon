package internal

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
)

var cli struct {
	Browser struct {
		Name string `env:"NAME"`
		Prof string `env:"PROFILE"`
	} `embed:"" envprefix:"BROWSER_" prefix:"browser-"`
	Output  string    `default:"." env:"OUTPUT"`
	User    string    `required:""`
	After   time.Time `format:"15:04:05"`
	Before  time.Time `format:"15:04:05"`
    Subs    []string  `sep:","`
	Titles  []string  `sep:"none"`
    Flairs  []string  `name:"flair" sep:"none"`
	Only    []string  `sep:"," enum:"gif,image,video,gallery"`
} `embed:"" envprefix:"RACCOON_"`
