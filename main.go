package main

import (
	"net/http"
	"os"
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

	e.Logger.Fatal(e.Start(":1323"))
}
