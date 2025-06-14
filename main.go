// Copyright © 2023- Luka Ivanovic
// This code is licensed under the terms of the MIT licence (see LICENCE for details).

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func filter[T any](ss []T, test func(T) bool) []T {
	res := make([]T, 0)
	for _, s := range ss {
		if test(s) {
			res = append(res, s)
		}
	}
	return res
}

func apply[T, U any](ss []T, f func(T) U) []U {
	res := make([]U, 0)
	for _, s := range ss {
		res = append(res, f(s))
	}
	return res
}

func parseLink(link string) (string, string, error) {
	title, url, found := strings.Cut(link, "https://")
	if !found {
		return "", "", fmt.Errorf("no link found in line: %s", link)
	}
	title = strings.TrimSpace(title)
	url = "https://" + url // return protocol part (removed by strings.Cut)
	return title, url, nil
}

func padLeft(s string, length int, value string) string {
	if len(s) < length {
		return strings.Repeat(value, length-len(s)) + s
	}
	return s
}

func toString(num int) string {
	return strconv.Itoa(num)
}

func doer(verbose bool, commands <-chan *exec.Cmd, results chan<- error) {
	for job := range commands {
		if verbose {
			fmt.Println(job)
		}
		output, err := job.Output()
		// err := job.Run()
		if err == nil {
			if verbose {
				p := strings.Split(job.String(), " -o ")
				fmt.Println(p[1], "downloaded successfully")
			}
			results <- nil
		} else {
			// fmt.Println("bad", strings.Split(job.String(), " -o "))
			if verbose {
				fmt.Println(string(output))
			}
			results <- err
		}
	}
}

func main() {
	index := flag.Int("index", 1, "starting index, or if negative disables indexing")
	linksFile := flag.String("file", filepath.Join(".", "links"), "file with links")
	// errorsFile := flag.String("errors", filepath.Join(".", "failed"), "file with errors, if any occur")
	outputDir := flag.String("dir", filepath.Dir(""), "directory to store results in")
	batchSize := flag.Int("batch", 6, "number of parralel downloads")
	verbose := flag.Bool("verbose", false, "get extra information while downloading")
	additional_command := flag.String("command", "", "additional command to pass to yt-dlp")
	flag.Parse()
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		log.Fatalln("yt-dlp not found")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		log.Fatalln("ffpeg not found")
	}
	rawLinks, err := os.ReadFile(*linksFile)
	if err != nil {
		panic(err)
	}
	links := strings.Split(strings.TrimSpace(string(rawLinks)), "\n")
	links = apply(links, func(s string) string {
		return strings.TrimSpace(s)
	})
	links = filter(links, func(s string) bool {
		return !strings.HasPrefix(s, "--")
	})
	links = filter(links, func(s string) bool {
		return len(s) != 0
	})
	total := len(links)
	if total == 0 {
		fmt.Fprintln(os.Stderr, "There is no songs to download :(")
		os.Exit(1)
	}
	succeeded := 0
	commands := make(chan *exec.Cmd, total)
	results := make(chan error, total)
	for i := 0; i < *batchSize; i += 1 {
		go doer(*verbose, commands, results)
	}
	padLength := max(2, len(toString(len(links))))
	go func() {
		for i := range links {
			title, url, err := parseLink(links[i])
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
			if *index > 0 {
				title = padLeft(toString(*index), padLength, "0") + " - " + title + ".opus"
				title = filepath.Join(*outputDir, title)
				*index += 1
			}
			arguments := []string{"-f", "bestaudio/best", "--extract-audio",
				"--audio-quality", "0", "--audio-format", "opus"}
			if (*additional_command) != "" {
				arguments = append(arguments, strings.Fields(*additional_command)...)
			}
			arguments = append(arguments, "-o", title, url)
			command := exec.Command("yt-dlp", arguments...)
			commands <- command
		}
	}()
	// close(commands)
	for range links {
		err := <-results
		if err == nil {
			succeeded += 1
		}
	}
	fmt.Printf("Done. Downloaded %d/%d.\n", succeeded, total)
}
