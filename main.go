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

func filter(ss []string, test func(string) bool) []string {
	res := make([]string, 0)
	for _, s := range ss {
		if test(s) {
			res = append(res, s)
		}
	}
	return res
}

func apply(ss []string, f func(string) string) []string {
	res := make([]string, 0)
	for _, s := range ss {
		res = append(res, f(s))
	}
	return res
}

func parseLink(link string) (string, string) {
	title, url, _ := strings.Cut(link, "https://")
	title = strings.TrimSpace(title)
	url = "https://" + url
	return title, url
}

func padStart(s string, length int, value string) string {
	if len(s) < length {
		return strings.Repeat(value, length-len(s)) + s
	}
	return s
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func toString(num int) string {
	return strconv.Itoa(num)
}

func doer(commands <-chan *exec.Cmd, results chan<- error) {
	for job := range commands {
		err := job.Run()
		if err != nil {
			results <- err
		} else {
			results <- nil
		}
	}
}

func main() {
	index := flag.Int("index", 1, "starting index, or if negative disables indexing")
	linksFile := flag.String("file", filepath.Join(".", "links"), "file with links")
	// errorsFile := flag.String("errors", filepath.Join(".", "failed"), "file with errors, if any occur")
	outputDir := flag.String("dir", filepath.Dir(""), "directory to store results in")
	batchSize := flag.Int("batch", 6, "number of parralel downloads")
	// verbose := flag.Bool("verbose", false, "get extra information while downloading")
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
	total := len(links)
	succeeded := 0
	commands := make(chan *exec.Cmd, total)
	results := make(chan error, total)
	for i := 0; i < *batchSize; i += 1 {
		go doer(commands, results)
	}
	padLength := max(2, len(toString(len(links))))
	for i := range links {
		title, url := parseLink(links[i])
		if *index > 0 {
			title = padStart(toString(*index), padLength, "0") + " - " + title + ".mp3"
			title = filepath.Join(*outputDir, title)
			*index += 1
		}
		commands <- exec.Command("yt-dlp", "-f", "bestaudio/best", "--extract-audio",
			"--audio-quality", "0", "--audio-format", "mp3", "-o", title, url)
	}
	close(commands)
	for range links {
		err := <-results
		if err == nil {
			succeeded += 1
		}
	}
	fmt.Printf("Done. Downloaded %d/%d.\n", succeeded, total)
}
