package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// HandleLowerBandwidth()
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
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

	e.GET("/segment", func(c echo.Context) error {
		f := "media/output/video.m3u8"
		// Read the M3U8 playlist file
		content, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		c.Response().Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		return c.String(http.StatusOK, string(content))
	})
	// ffmpeg -i input.mp4 -c:v copy -c:a copy -hls_time 10 -hls_list_size 0 output.m3u8
	e.GET("/:name", func(c echo.Context) error {
		segmentName := c.Param("name")
		// segmentPath := filepath.Join("media/output/", segmentName)
		absolutePath, err := filepath.Abs(fmt.Sprintf("media/output/%s", segmentName))
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		}
		return c.File(absolutePath)
	})

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
