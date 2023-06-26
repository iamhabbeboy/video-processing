package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// type Template struct {
// 	templates *template.Template
// }

// func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
// 	return t.templates.ExecuteTemplate(w, name, data)
// }

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

	e.Logger.Fatal(e.Start(":1323"))
}

func parseRangeHeader(rangeHeader string, fileSize int64) (int64, int64, error) {
	// Parse the range header value
	// Example: bytes=500-999
	startEnd := rangeHeader[len("bytes="):]
	startEndSplit := strings.Split(startEnd, "-")
	startStr := startEndSplit[0]
	endStr := startEndSplit[1]

	// Parse the start and end positions
	start, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	end, err := strconv.ParseInt(endStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	// Adjust the end position if it exceeds the file size
	if end >= fileSize {
		end = fileSize - 1
	}

	return start, end, nil
}
