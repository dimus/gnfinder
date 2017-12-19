package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/util"
)

func main() {
	dict := dict.LoadDictionary()
	bayes := flag.Bool("bayes", false, "Use Bayes-based name-finding")
	flag.Parse()
	var data []byte
	var err error
	var output []byte
	var opts []util.Opt
	if *bayes {
		opts = append(opts, util.WithBayes(true))
	}
	switch flag.NArg() {
	case 0:
		data, err = ioutil.ReadAll(os.Stdin)
		util.Check(err)
	case 1:
		data, err = ioutil.ReadFile(flag.Arg(0))
		util.Check(err)
	default:
		fmt.Printf("input must be from stdin or file\n")
		os.Exit(1)
	}
	output = gnfinder.FindNamesJSON(data, &dict, opts...)
	fmt.Println(string(output))
}
