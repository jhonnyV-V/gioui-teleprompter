package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
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
	var focusBarSize unit.Dp = 50
	var textWitdh unit.Dp = 550
	var fontSize unit.Sp = 35
	var autoScroll bool = false
	var autoScrollSpeed unit.Dp = 1
	theme := material.NewTheme()
	var list layout.List

	var ops op.Ops

	for {
		wEvent := window.Event()
		switch evenType := wEvent.(type) {
		case app.FrameEvent:

			context := app.NewContext(&ops, evenType)
			contextEvent, _ := context.Source.Event(key.Filter{
				Optional: key.ModShift,
			})
			var stepSize unit.Dp = 1

			switch contextEventType := contextEvent.(type) {

			case key.Event:
				if contextEventType.State != key.Press {
					break
				}

				if contextEventType.Modifiers == key.ModShift {
					stepSize += 5
				}

				if contextEventType.Name == "+" {
					fontSize += unit.Sp(stepSize)
					focusBarSize += unit.Dp(unit.Sp(stepSize)) * 1.5
				}

				if contextEventType.Name == "-" {
					fontSize -= unit.Sp(stepSize)
					focusBarSize -= unit.Dp(unit.Sp(stepSize)) * 1.5
				}

				if contextEventType.Name == "Space" {
					autoScroll = !autoScroll
					if autoScrollSpeed == 0 {
						autoScroll = true
						autoScrollSpeed++
					}
				}

				if contextEventType.Name == "K" {
					scrollY = scrollY - (stepSize * 4)
					if scrollY < 0 {
						scrollY = 0
					}
				}

				if contextEventType.Name == "J" {
					scrollY = scrollY + (stepSize * 4)
				}

				if contextEventType.Name == "F" {
					autoScroll = true
					autoScrollSpeed++
				}

				if contextEventType.Name == "S" {
					if autoScrollSpeed > 0 {
						autoScrollSpeed--
					} else {
						autoScroll = false
					}
				}

				if contextEventType.Name == "W" {
					textWitdh = textWitdh + (stepSize * 10)
				}

				if contextEventType.Name == "N" {
					textWitdh = textWitdh - (stepSize * 10)
				}

				if contextEventType.Name == "U" {
					focusBarY = focusBarY - stepSize
				}

				if contextEventType.Name == "D" {
					focusBarY = focusBarY - stepSize
				}

			}

			//TODO: check why i get nil events
			contextEvent, _ = context.Source.Event(pointer.Filter{})
			switch contextEventType := contextEvent.(type) {
			case pointer.Event:
				fmt.Printf("pointer Event is %s\n", contextEventType.Kind.String())
				if contextEventType.Kind != pointer.Scroll {
					break
				}

				if contextEventType.Modifiers == key.ModShift {
					stepSize = 3
				}

				scrollY = scrollY + (unit.Dp(contextEventType.Scroll.Y) * stepSize)
				if scrollY < 0 {
					scrollY = 0
				}

			default:
				fmt.Printf("pointer Event is %#v\n", contextEventType)
			}

			paint.Fill(&ops, color.NRGBA{R: 0xff, G: 0xfe, B: 0xe0, A: 0xff})

			focusBar := clip.Rect{
				Min: image.Pt(0, int(focusBarY)),
				Max: image.Pt(context.Constraints.Max.X, int(focusBarY)+int(focusBarSize)),
			}.Push(&ops)
			paint.ColorOp{Color: color.NRGBA{R: 0xff, A: 0x66}}.Add(&ops)
			paint.PaintOp{}.Add(&ops)
			focusBar.Pop()

			if autoScroll {
				scrollY = scrollY + autoScrollSpeed
				inv := op.InvalidateCmd{At: context.Now.Add(time.Second * 2 / 100)}
				context.Execute(inv)
			}

			list = layout.List{
				Axis: layout.Vertical,
				Position: layout.Position{
					Offset: int(scrollY),
				},
			}

			marginWidth := (unit.Dp(context.Constraints.Max.X) - textWitdh) / 2
			margins := layout.Inset{
				Left:   marginWidth,
				Right:  marginWidth,
				Top:    0,
				Bottom: 0,
			}

			margins.Layout(context, func(context layout.Context) layout.Dimensions {
				return list.Layout(context, len(paragrahpList),
					func(context layout.Context, index int) layout.Dimensions {
						paragraph := material.Label(theme, fontSize, paragrahpList[index])
						paragraph.Alignment = text.Middle
						return paragraph.Layout(context)
					})
			})
			evenType.Frame(context.Ops)

		case app.DestroyEvent:
			return evenType.Err
		}
	}
}
