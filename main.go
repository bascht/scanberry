package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"

	"embed"
	"net/http"

	"github.com/gofiber/template/html"
	"github.com/spf13/viper"

	"git.bascht.space/bascht/scanberry/scan"
)

//go:embed views/*
var viewsfs embed.FS

var documents = make(map[string]*scan.Document)
var basedir string

func main() {
	// Create a new engine by passing the template folder
	// and template extension using <engine>.New(dir, ext string)
	// engine := html.New("./views", ".html")

	// // We also support the http.FileSystem interface
	// // See examples below to load templates from embedded files
	// engine := html.NewFileSystem(http.Dir("./views"), ".html")
	engine := html.NewFileSystem(http.FS(viewsfs), ".html")

	viper.SetConfigName("scanberry")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(filepath.Join(xdg.ConfigHome, "scanberry"))

	viper.SetDefault("listen.host", "0.0.0.0")
	viper.SetDefault("listen.port", "0.0.0.0")
	// viper.SetDefault("filename.prefixes", [""])

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Could not load config")
	}

	basedir = xdg.RuntimeDir + "/himbeerscan"
	os.MkdirAll(filepath.Join(basedir, "downloads"), 0700)
	os.Chdir(basedir)
	fmt.Sprintf("Starting to work in %v", basedir)

	engine.Reload(true)       // Optional. Default: false
	engine.Debug(false)       // Optional. Default: false
	engine.Layout("embed")    // Optional. Default: "embed"
	engine.Delims("{{", "}}") // Optional. Default: engine delimiters

	// engine.AddFunc("greet", func(name string) string {
	// 	return "Hallouhu, " + name + "!"
	// })

	// After you created your engine, you can pass it to Fiber's Views Engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/downloads", basedir+"/downloads")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/scan")
	})

	app.Get("/scan", func(c *fiber.Ctx) error {
		prefixes := viper.GetStringSlice("filename.prefixes")
		timestamp := time.Now().Format("2006-01-02")
		return c.Render("views/scan", fiber.Map{
			"Title":     "Hello, World!",
			"Prefixes":  prefixes,
			"Timestamp": timestamp,
		}, "views/layouts/main")
	})

	app.Post("/scan", func(c *fiber.Ctx) error {

		document := scan.Document{
			Id:     uuid.NewString(),
			Name:   c.FormValue("name"),
			Date:   time.Now(),
			Duplex: c.FormValue("duplex") == "on",
			Events: make(chan scan.Event),
		}

		documents[document.Id] = &document

		go scan.Process(basedir, &document)
		document.Events <- scan.Event{Message: "Geht los", Type: "info"}

		// document.Events <- scan.Event{Message: "Scanning initiated", Type: "info"}
		return c.Render("views/status", fiber.Map{
			"id":                    document.Id,
			"name":                  document.Name,
			"FullName":              document.FullName(),
			"FullNameWithExtension": document.FullNameWithExtension(),
		}, "views/layouts/main")
	})

	app.Get("/show/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if document, ok := documents[id]; ok {
			return c.Render("views/show", fiber.Map{
				"id":                    document.Id,
				"name":                  document.Name,
				"FullName":              document.FullName(),
				"FullNameWithExtension": document.FullNameWithExtension(),
			}, "views/layouts/main")
		} else {
			return c.Status(fiber.StatusNotFound).SendString("Sorry, document no longer exists")
		}
	})

	app.Get("/events/:id", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		id := c.Params("id")
		document := documents[id]

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			fmt.Fprintf(w, "event: status\ndata: started\n\n")
			w.Flush()

			for event := range document.Events {
				msg, err := json.Marshal(event)
				if err != nil {
					log.Fatal("Could not encode event")
				}

				fmt.Fprintf(w, "data: %s \n\n", msg)
				w.Flush()
			}

			fmt.Fprintf(w, "event: status\ndata: stopped\n\n")
			w.Flush()

		}))

		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
