package main

import (
	"log"
	"os/exec"
	"strconv"
)

func HandleLowerBandwidth() {
	inputFilePath := "input.mp4"
	outputFilePath := "output.mp4"
	width := 853
	height := 480

	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-vf", "scale=trunc("+strconv.Itoa(width)+"/2)*2:trunc("+strconv.Itoa(height)+"/2)*2", "-c:v", "libx264", "-preset", "slow", "-crf", "23", "-c:a", "copy", outputFilePath)

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error resizing video: %v", err)
	}
}

func SliceVideoFrames() {
	inputFilePath := "video/lengthy.mov"
	outputDirectory := "output/"
	frameDuration := "300" // 5 minutes in seconds

	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-c:v", "copy", "-f", "segment", "-segment_time", frameDuration, "-reset_timestamps", "1", "-map", "0", outputDirectory+"output_%03d.mp4")

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error slicing video: %v", err)
	}
}
