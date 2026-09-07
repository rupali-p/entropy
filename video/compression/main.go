package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

func main() {

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
	}
	fmt.Println("Current working directory:", cwd)

	fmt.Println("\nTraversing to the Screencaptures directory...")
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
	}
	screencapturesDir := userHomeDir + "/screencaptures"

	videoDir, videoTitle := getMostRecentVideo(screencapturesDir)
	fmt.Println("Most recent video directory:", videoDir)
	fmt.Println("Most recent video title:", videoTitle)

	compressedVideo := compressVideo(screencapturesDir, videoTitle, videoDir)
	fmt.Println("Compressed video:", compressedVideo)

}

func compressVideo(screencapturesDir string, title string, path string) os.FileInfo {

	fmt.Println("Compressing video at path:", path)
	
	cmd := exec.Command("ffmpeg", "-i", path, "-c:v", "libx265", "-crf", "28", "-preset", "medium", "-c:a", "aac", "-b:a", "128k", screencapturesDir+"/compressed-"+title)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error compressing video:", err)
		return nil
	}

	compressedVideo, err := os.Stat(screencapturesDir + "/compressed-" + title)
	if err != nil {
		fmt.Println("Error getting compressed video info:", err)
		return nil
	}

	fmt.Printf("\n\nCompressed video file: %+v\n\n", compressedVideo)

	return compressedVideo
}

func getMostRecentVideo(screencapturesDir string) (string, string) {

	// fmt.Println("Screencaptures directory:", screencapturesDir)

	// this is retrieving files based on alphabetical order
	allfiles, err := os.ReadDir(screencapturesDir)
	if err != nil {
		fmt.Println("Error reading screencaptures directory:", err)
		return "", ""
	}
	// fmt.Println("Files in Screencaptures directory:", allfiles)

	filesSortedByModificationTime := sortFilesByModificationTime(allfiles)
	recentVideoFile := getMostRecentVideoFile(filesSortedByModificationTime)
	if recentVideoFile == (os.DirEntry)(nil) {
		fmt.Println("No video files found in the Screencaptures directory.")
		return "", ""
	}
	recentVideoFileInfo, err := recentVideoFile.Info()
	if err != nil {
		fmt.Println("Error getting recent video file info:", err)
		return "", ""
	}

	
	recentVideo := struct {
		title string
		size  string
	}{
		title: recentVideoFileInfo.Name(),
		size: fmt.Sprintf("%dKB", recentVideoFileInfo.Size()/1000),
	}


	fmt.Printf("\n\nMost recent video file: %+v\n\n", recentVideo)

	return screencapturesDir + "/" + recentVideo.title, recentVideo.title
}

func getMostRecentVideoFile(files []os.DirEntry) os.DirEntry {
	for i := 0; i < 5 && i < len(files); i++ {
		info, err := files[i].Info()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}
		if filepath.Ext(info.Name()) == ".mov" {
			return files[i]
		}
	}
	return nil
}

func sortFilesByModificationTime(files []os.DirEntry) []os.DirEntry {

	sort.Slice(files, func(i, j int) bool {
		infoI, _ := files[i].Info()
		infoJ, _ := files[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})
	fmt.Println("\nSorted files by modification time:", files)

	return files
}