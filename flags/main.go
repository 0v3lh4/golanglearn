package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	flagsNameFile = "flags.txt"
	targetDir     = "FlagsImage/"
	originURL     = "https://www.cia.gov/library/publications/the-world-factbook/graphics/flags/large/"
	maxDownload   = 249
)

func main() {
	defer timeTrack(time.Now(), "SaveImages")
	createtargetDirIfNotExists()
	saveFiles(getFlagNames())
}

func saveFiles(flagNames []string) {
	name := make(chan string)
	defer close(name)

	var wg sync.WaitGroup
	wg.Add(maxDownload)

	for idx := 0; idx < maxDownload; idx++ {
		go saveImageToURLForTargetDirectory(flagNames[idx], name, &wg)
	}

	go func() {
		for flagSaved := range name {
			fmt.Printf("%s salvo! \n", flagSaved)
		}
	}()

	wg.Wait()
}

func saveImageToURLForTargetDirectory(flagName string, name chan string, wg *sync.WaitGroup) {
	url := originURL + flagName
	resp, err := http.Get(url)

	defer resp.Body.Close()

	checkError(err)

	flagImg, _ := os.Create(targetDir + flagName)
	defer flagImg.Close()

	_, errCopy := io.Copy(flagImg, resp.Body)

	checkError(errCopy)

	name <- flagName

	defer wg.Done()
}

func getFlagNames() []string {
	flagsFile, err := os.Open(flagsNameFile)
	defer flagsFile.Close()
	checkError(err)

	scannerFlags := bufio.NewScanner(flagsFile)
	flagNames := []string{}

	for scannerFlags.Scan() {
		flagNames = append(flagNames, scannerFlags.Text())
	}

	return flagNames
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func createtargetDirIfNotExists() {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		os.Mkdir(targetDir, os.ModeDir)
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s tempo transcorrido %s", name, elapsed)
}
