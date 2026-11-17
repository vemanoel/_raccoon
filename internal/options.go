package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
)

var cli struct {
	Zip     bool      `short:"z"`
	Dedup   bool      `short:"d"`
	Output  string    `short:"o" required:"" env:"RACCOON_OUTPUT"`
	User    string    `short:"u" required:""`
	After   time.Time `short:"a" format:"15:04:05"`
	Before  time.Time `short:"b" format:"15:04:05"`
    Subs    []string  `short:"s" sep:","`
	Titles  []string  `short:"t" sep:"none"`
    Flairs  []string  `short:"f" name:"flair" sep:"none"`
	Ignore  []string  `short:"i" sep:"," enum:"gif,image,video,gallery"`
	Arctic  bool      `default:"true" negatable:""`
}

func main() {
	ctx := kong.Parse(&cli)

	fmt.Printf("User: %s\n", cli.User)
	fmt.Printf("Formats selecionados: %v\n", cli.Formats)

	_ = ctx
}
