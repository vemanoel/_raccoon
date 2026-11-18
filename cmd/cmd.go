package cmd

import (
	"fmt"
	"time"
	"os"

	"github.com/alecthomas/kong"
)

var cmd struct {
	Output  string    `default:"." env:"OUTPUT"`
	User    string    `required:""`
	After   time.Time `format:"15:04:05"`
	Before  time.Time `format:"15:04:05"`
    Subs    []string  `sep:","`
	Titles  []string  `sep:"none"`
    Flairs  []string  `name:"flair" sep:"none"`
	Only    []string  `sep:"," enum:"gif,image,video,gallery"`
	Profile string	  `env:"BROWSER"`
} `embed:"" envprefix:"RACCOON_"`

func init() {
  err := kong.Parse(&cmd).Run()
  if err != nil {
	  fmt.Println(err)
	  os.Exit(1)
  }
}

func main() {
	//
}
