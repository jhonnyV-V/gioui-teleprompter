package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
)

var (
	filename      *string
	paragrahpList []string
)

func readText(filename string) []string {
	file, err := os.ReadFile(filename)
	text := []string{}
	if err != nil {
		log.Fatal("Error when reading file:\n  ", err)
	}
	if err == nil {
		// Convert text to a slice of strings.
		text = strings.Split(string(file), "\n")
		// Add extra empty lines a the end. Simple trick to ensure
		// the last line of the speech scrolls out of the screen
		for i := 1; i <= 10; i++ {
			text = append(text, "")
		}
	}

	return text
}

func main() {

	filename = flag.String("file", "speech.txt", "Which .txt file shall I present?")
	flag.Parse()

	paragrahpList = readText(*filename)

	go func() {
		window := new(app.Window)
		window.Option(app.Title("TelePrompter"))
		window.Option(app.Size(unit.Dp(650), unit.Dp(600)))
		window.Option(app.MinSize(unit.Dp(650), unit.Dp(600)))

		if err := loop(window); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)

	}()

	app.Main()
}

func loop(window *app.Window) error {
	var scrollY unit.Dp = 0
	var focusBarY unit.Dp = 170
	var textWitdh unit.Dp = 550
	var fontSize unit.Sp = 35
	var autoScroll bool = false
	var autoScrollSpeed unit.Dp = 1

	var op op.Ops

	for {
		event := window.Event()
		switch evenType := event.(type) {
		case app.FrameEvent:
			context := app.NewContext(&op, evenType)
			context.Dp(unit.Dp(100))

		case app.DestroyEvent:
			return evenType.Err
		}
	}
}
