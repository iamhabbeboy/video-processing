package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

func main() {
	// HandleLowerBandwidth()
	e := echo.New()

	// t := &Template{
	// 	templates: template.Must(template.ParseGlob("public/views/*.html")),
	// }
	// e.Static("/", "public/assets")
	e.File("/", "media/output")
	// e.Renderer = t
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"name": "Dolly!",
		})
	})

	e.GET("/video-chunk", func(c echo.Context) error {
		f := "media/output/output_000.mp4"
		file, err := os.Open(f)
		if err != nil {
			return err
		}
		defer file.Close()

		// Get the file's information
		fi, err := file.Stat()
		if err != nil {
			return err
		}
		s := fi.Size()
		fx := strconv.FormatInt(s, 10)
		// Set the response headers for streaming
		c.Response().Header().Set("Content-Type", "video/mp4")
		c.Response().Header().Set("Content-Length", fx)
		c.Response().Header().Set("Accept-Ranges", "bytes")

		// Read and stream the video in chunks
		http.ServeContent(c.Response().Writer, c.Request(), "video.mp4", fi.ModTime(), file)

		return nil
	})

	e.GET("/segment-list", func(c echo.Context) error {
		f := "media/output/video.m3u8"
		absolutePath := filepath.Join(".", f)
		// application/vnd.apple.mpegurl
		c.Response().Header().Set("Accept-Ranges", "bytes")
		c.Response().Header().Set("Content-Type", "application/vnd.apple.mpegurlapplication/vnd.apple.mpegurl")
		return c.String(http.StatusOK, absolutePath)
	})
	// ffmpeg -i input.mp4 -c:v copy -c:a copy -hls_time 10 -hls_list_size 0 output.m3u8
	e.GET("/stream/:name", func(c echo.Context) error {
		segmentName := c.Param("name")
		segmentPath := filepath.Join("media/output/", segmentName)

		// Open the video segment file
		file, err := os.Open(segmentPath)
		if err != nil {
			log.Println("Error opening segment file:", err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		defer file.Close()
		fi, _ := file.Stat()
		// Set the response headers for the video segment
		c.Response().Header().Set("Content-Type", "video/MP2T")
		c.Response().Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
		c.Response().WriteHeader(http.StatusOK)

		// Stream the video segment to the response
		_, err = io.Copy(c.Response().Writer, file)
		if err != nil {
			log.Println("Error streaming video segment:", err)
		}
		return c.File(segmentPath)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
