package main

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/kchugalinskiy/pulseaudio/simple"
)

func main() {
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Не удалось создать окно:", err)
	}

	exitCh := make(chan interface{})
	defer close(exitCh)

	win.SetTitle("Простой пример")
	win.Connect("destroy", func() {
		gtk.MainQuit()
		close(exitCh)
	})

	f := simple.SampleSpec{
		Format:   simple.SampleFormatS16LE,
		Rate:     44100,
		Channels: 1,
	}
	play, err := simple.New("", "echo example play", simple.StreamDirectionPlayback, "", "Music", &f, nil, nil)
	if err != nil {
		fmt.Printf("creating sample play: %v\n", err)
		return
	}
	defer play.Close()

	rec, err := simple.New("", "echo example rec", simple.StreamDirectionRecord, "", "Music", &f, nil, nil)
	if err != nil {
		fmt.Printf("creating sample rec: %v\n", err)
		return
	}
	defer rec.Close()

	go func() {
		for {
			select {
			case <-exitCh:
				return
			default:
				b, err := rec.Read16(1024)
				if err != nil {
					fmt.Printf("reading audio: %v\n", err)
					continue
				}

				if err = play.Write16(b); err != nil {
					fmt.Printf("writing audio: %v\n", err)
					continue
				}
			}
		}
	}()

	gtk.Main()
}
