package main

import (
	b64 "encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

const (
	currentTrack    = "name of current track as string"
	currentArtist   = "artist of current track as string"
	currentAlbum    = "album of current track as string"
	playerState     = "player state as string"
	currentPosition = "player position as real"
	trackDuration   = "duration of current track as integer"
	trackImage      = "artwork url of current track"
	//--
	fastSpeed = "5s"
	slowSpeed = "15s"
	//--
	showTime      = true
	showArtist    = false
	showTitleMenu = false
	fontSize      = "12"
	truncateTo    = "36"
	pauseImg      = "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAAlwSFlzAAAN1wAADdcBQiibeAAAActpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IlhNUCBDb3JlIDUuNC4wIj4KICAgPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4KICAgICAgPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIKICAgICAgICAgICAgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIgogICAgICAgICAgICB4bWxuczp0aWZmPSJodHRwOi8vbnMuYWRvYmUuY29tL3RpZmYvMS4wLyI+CiAgICAgICAgIDx4bXA6Q3JlYXRvclRvb2w+d3d3Lmlua3NjYXBlLm9yZzwveG1wOkNyZWF0b3JUb29sPgogICAgICAgICA8dGlmZjpPcmllbnRhdGlvbj4xPC90aWZmOk9yaWVudGF0aW9uPgogICAgICA8L3JkZjpEZXNjcmlwdGlvbj4KICAgPC9yZGY6UkRGPgo8L3g6eG1wbWV0YT4KGMtVWAAAAUVJREFUOBGlk0tKQ1EQRBOjEifRsSgBXUEka/Mz12Wp6FjIAqKDGHABjhT8nXO99fIIz5EFlVvdr7vuN/3eCgPkZw2HjFO4X+MXxgf4VuN2bUmZENvwHM7h9xrNncEtKNLTiD2S9zCNX2hXJNXJ36F3oWhMNgn8YNE7tEltY5rN+c38LbSnwQUqzY4Wf9ScBuqYxsTtFOzw+wRt9JAcr+AxfK1UX8J2zSPx0GWcwCMo+r9Db8GoqbMK9XNRqxpNJxv8HNQPLjUGXqMw7spZKw41+Bc0WFYHtXsUeTDGXblMvFTMoHsUKR6j3eOgUm1OpMZDtLfAK/FDrsjD8+rMSbW5ds0pcQNv4wamIMUxcDSXCa7Rrq4gwucZkzTatG5m86h0dpho5nb++jO57ExYxtyxhiacTfgOJnBsABbQA8vtePjlLfwAb99671Xj0CcAAAAASUVORK5CYII="
	playImg       = "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAAlwSFlzAAAN1wAADdcBQiibeAAAActpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IlhNUCBDb3JlIDUuNC4wIj4KICAgPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4KICAgICAgPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIKICAgICAgICAgICAgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIgogICAgICAgICAgICB4bWxuczp0aWZmPSJodHRwOi8vbnMuYWRvYmUuY29tL3RpZmYvMS4wLyI+CiAgICAgICAgIDx4bXA6Q3JlYXRvclRvb2w+d3d3Lmlua3NjYXBlLm9yZzwveG1wOkNyZWF0b3JUb29sPgogICAgICAgICA8dGlmZjpPcmllbnRhdGlvbj4xPC90aWZmOk9yaWVudGF0aW9uPgogICAgICA8L3JkZjpEZXNjcmlwdGlvbj4KICAgPC9yZGY6UkRGPgo8L3g6eG1wbWV0YT4KGMtVWAAAAU5JREFUOBGNkk0uB0EQR+fvKySCFYkQZ5DYYCPu4Ro4idsgJKwkDsDGRuxF4pv3Wv/GZDISlbzpquqqmuquHjW/Mo76Uc1p1g1YrvY96xU8V7sbW1w6lCk4hBv46qHvACZBSU6rLOA8hyR+otuRxOd6BvOgtEUmMNwwwDaT9Fp1i73DCxhzCua0YttuJMDEp+rT/9bRE+NxiszwvQUDTXR9hF04qnaK2JmdaHsnXnazDTokmx5jERQLXUI/RntzjM8KKJ5zVLQffbbqJ6xbsAcP4MX5I2XVAkOSQt09/zjkL9X77XmEpZrtES6gH1OOYMxfl7jD3r8u0SKOxIoZkdNwEvnr0Bj32W/FR+FlpUgekonq/Yd0jK99hVF8niliIZNMTrF0Y/IcKMltFR0ex0eShKz6bDtJZe2ORYd/U3xh67CmgdzBNTgdxfHbYfMNlDyRQrUR6jUAAAAASUVORK5CYII="
	//--
	debug = false
)

var execDir string

func main() {
	execDir = getExecDir()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "slow":
			exec.Command("mv", execDir+"/spotify."+fastSpeed+".go", execDir+"/spotify."+slowSpeed+".go").Run()
			tell("Bitbar", "quit")
			exec.Command("open", "-a", "BitBar").Run()
		case "fast":
			exec.Command("mv", execDir+"/spotify."+slowSpeed+".go", execDir+"/spotify."+fastSpeed+".go").Run()
			tell("Bitbar", "quit")
			exec.Command("open", "-a", "BitBar").Run()
		case "open":
			tell("Spotify", "activate")
		}
	} else {
		track := tell("Spotify", currentTrack)
		artist := tell("Spotify", currentArtist)
		album := tell("Spotify", currentAlbum)
		state := tell("Spotify", playerState)
		imageURL := tell("Spotify", trackImage)
		pos, _ := strconv.ParseFloat(tell("Spotify", currentPosition), 32)

		if state == "playing" {
			if showTime {
				fmt.Print("(" + parseSeconds(int(pos)) + ") ")
			}
			fmt.Print(track)
			if showArtist {
				fmt.Print(" - " + artist)
			}
			fmt.Print(" | size=" + fontSize)
			fmt.Print(" length=" + truncateTo)
			fmt.Print(" templateImage=" + playImg)
			toggleFastMode(true)
		} else {
			fmt.Print("| templateImage=" + pauseImg)
			toggleFastMode(false)
		}
		fmt.Println("\n---")
		if showTitleMenu {
			fmt.Print(track + " - ")
		}
		fmt.Print(artist)
		fmt.Println("|bash='" + os.Args[0] + "' param1=open terminal=false")
		fmt.Print(album)
		fmt.Println("|bash='" + os.Args[0] + "' param1=open terminal=false")

		fmt.Println("---")
		fmt.Println("| image=" + getTempImage(imageURL) + " bash='" + os.Args[0] + "' param1=open terminal=false")

		if debug == true {
			fmt.Println("execDir:", execDir)
			fmt.Println("speed:", getUpdSpeed())
		}
	}
}

func toggleFastMode(active bool) {
	if active && getUpdSpeed() != "fast" {
		exec.Command(execDir+"/spotify."+slowSpeed+".go", "fast").Run()
	} else if !active && getUpdSpeed() != "slow" {
		exec.Command(execDir+"/spotify."+fastSpeed+".go", "slow").Run()
	}
}

func getUpdSpeed() string {
	speed := strings.Split(os.Args[0], ".")[1]

	if speed == fastSpeed {
		return "fast"
	}

	return "slow"
}

func getExecDir() string {
	if execDir == "" {
		execDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	return execDir
}

func parseSeconds(seconds int) string {
	return fmt.Sprintf("%d:%02d", seconds/60, seconds%60)
}

func tell(to, cmd string) string {
	out, err := exec.Command("osascript", "-e", "tell application \""+to+"\" to "+cmd).Output()
	if err != nil {
		fmt.Println("(E)\n---\nError: Spotify not running or no song selected")
		toggleFastMode(false)
		os.Exit(0)
	}
	return strings.TrimSpace(string(out))
}

func getTempImage(url string) string {
	paths := strings.Split(url, "/")
	fileName := "/tmp/" + paths[len(paths)-1] + ".png"

	base64Image := getBase64FromImageString(fileName)

	if base64Image != "" {
		return base64Image
	}

	imgBytes, format := getImageFromURL(url)

	saveImage(fileName, imgBytes)
	resizeImage(fileName, format)
	return saveBase64StringFromImage(fileName)
}

func getBase64FromImageString(fileName string) string {
	bytesRead, err := ioutil.ReadFile(strings.Replace(fileName, ".png", ".txt", 1))

	if err != nil || len(bytesRead) == 0 {
		return ""
	}

	return string(bytesRead)
}

func saveBase64StringFromImage(fileName string) string {
	bytesRead, err := ioutil.ReadFile(fileName)

	if err == nil && len(bytesRead) != 0 {
		file, err := os.Create(strings.Replace(fileName, ".png", ".txt", 1))

		if err != nil {
			fmt.Println("Error", err)
			return ""
		}
		defer file.Close()

		encoded := b64.StdEncoding.EncodeToString(bytesRead)
		file.Write([]byte(encoded))

		return encoded
	}

	return ""
}

func saveImage(fileName string, imageBytes []byte) {
	file, err := os.Create(fileName)

	if err != nil {
		fmt.Println("Error", err)
		return
	}
	defer file.Close()

	file.Write(imageBytes)
}

func resizeImage(fileName, format string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var img image.Image

	if format == "image/jpeg" {
		img, err = jpeg.Decode(file)
	} else if format == "image/png" {
		img, err = png.Decode(file)
	}
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	m := resize.Thumbnail(200, 200, img, resize.Lanczos3)

	out, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if format == "image/jpeg" {
		jpeg.Encode(out, m, nil)
	} else if format == "image/png" {
		png.Encode(out, m)
	}
}

func getImageFromURL(url string) ([]byte, string) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("No image!", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, http.DetectContentType(body)
}
