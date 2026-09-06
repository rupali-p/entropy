package main

import (
	"fmt"
	"os"
	"sort"
)

func main() {
	videoDir := getMostRecentVideo()
	fmt.Println("Most recent video directory:", videoDir)
}

func getMostRecentVideo() string {

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return ""
	}
	fmt.Println("Current working directory:", cwd)

	fmt.Println("\nTraversing to the Screencaptures directory...")
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return ""
	}
	screencapturesDir := userHomeDir + "/screencaptures"

	// fmt.Println("Screencaptures directory:", screencapturesDir)

	// this is retrieving files based on alphabetical order
	allfiles, err := os.ReadDir(screencapturesDir)
	if err != nil {
		fmt.Println("Error reading screencaptures directory:", err)
		return ""
	}
	// fmt.Println("Files in Screencaptures directory:", allfiles)

	filesSortedByModificationTime := sortFilesByModificationTime(allfiles)
	mostRecentFile, err := filesSortedByModificationTime[len(filesSortedByModificationTime)-1].Info()
	if err != nil {
		fmt.Println("Error getting most recent file info:", err)
		return ""
	}
	
	recentVideo := struct {
		title string
		size  int64
	}{
		title: mostRecentFile.Name(),
		size: mostRecentFile.Size(),
	}


	fmt.Printf("\n\nMost recent video file: %+v\n\n", recentVideo)

	return screencapturesDir + "/" + recentVideo.title

}

func sortFilesByModificationTime(files []os.DirEntry) []os.DirEntry {

	sort.Slice(files, func(i, j int) bool {
		infoI, _ := files[i].Info()
		infoJ, _ := files[j].Info()
		return infoI.ModTime().Before(infoJ.ModTime())
	})
	fmt.Println("\nSorted files by modification time:", files)

	return files
}