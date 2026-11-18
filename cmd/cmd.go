package cmd

import (
	"fmt"
	"time"
	"os"

	"github.com/alecthomas/kong"
)

var cmd struct {
	Config	kong.ConfigFlag `env:"CONFIG_FILE"`
	Output  string    `default:"." env:"OUTPUT"`
	User    string    `required:""`
	After   time.Time `format:"15:04:05"`
	Before  time.Time `format:"15:04:05"`
    Subs    []string  `sep:","`
	Titles  []string  `sep:"none"`
    Flairs  []string  `name:"flair" sep:"none"`
	Only    []string  `sep:"," enum:"gif,image,video,gallery"`
	Profile string	  `env:"BROWSER"`
}

func init() {
  ctx := kong.Parse(
  	&cmd,
  	kong.Name("raccoon"),
  	kong.UsageOnError(),
  	kong.Configuration(kong.JSON),
  	kong.DefaultEnvVars("RACCOON"),
  )

  ctx.FatalIfErrorf(ctx.Run())
}

func main() {
	//
}
