// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Rsspoll polls RSS feeds for updates.
//
// Usage:
//
//	rsspoll [-d days] [file]
//
// Rsspoll reads RSS feeds from a file and prints entries within the
// last d days (default 1).
//
// The "-d" flag specifies the number of days to recall.
//
// The file containing RSS feeds of interest should either be passed as an
// argument or exist at $XDG_CONFIG_HOME/rsspoll/config.txt. Each feed URL
// should be on a separate line within the file.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mmcdole/gofeed"
)

var flagDays = flag.Int("d", 1, "number of days to recall")

func usage() {
	fmt.Fprintf(os.Stderr, "usage: rsspoll [-d days] [file]\n")
	os.Exit(2)
}

func main() {
	log.SetPrefix("rsspoll: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	var name string
	switch flag.NArg() {
	case 0:
		dir, err := os.UserConfigDir()
		if err != nil {
			log.Fatal(err)
		}
		name = filepath.Join(dir, "rsspoll", "config.txt")
	case 1:
		name = flag.Arg(0)
	default:
		usage()
	}
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := process(f, *flagDays); err != nil {
		log.Fatal(err)
	}
}

func process(r io.Reader, d int) error {
	fp := gofeed.NewParser()
	s := bufio.NewScanner(r)
	daysAgo := time.Now().AddDate(0, 0, -d)
	for s.Scan() {
		if err := poll(fp, s.Text(), daysAgo); err != nil {
			return err
		}
	}
	return s.Err()
}

func poll(fp *gofeed.Parser, url string, t time.Time) error {
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}
	for _, it := range feed.Items {
		if it.PublishedParsed != nil && (*it.PublishedParsed).After(t) {
			fmt.Printf("%s: %s\n", it.Title, it.Link)
		}
	}
	return nil
}
