/**
ExMachina driver
----------------
Connects to a running Kubernetes master and observe Pod events. Upon reception, it sends
musical bits to a SonicPi instance
*/
package main

import (
	"bytes"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/template"
	"time"

	"./moles"
	"./sonic"
	"k8s.io/kubernetes/pkg/api"
)

const (
	timeout   = 2 * time.Second
	masterURL = "http://127.0.0.1:8001"
	sonicURL  = "127.0.0.1:4557"
)

var (
	evtChan      chan int
	upTpl, dwTpl *template.Template
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool)

	setupSigs(sigs, done)

	loadTemplates()

	sonic := setupSonic()
	defer func() {
		sonic.Conn.Close()
		close(evtChan)
	}()

	evtChan = make(chan int)
	setupTimer(sonic, evtChan)

	moles.Watch(masterURL, up, down)

	<-done
}

func setupSonic() *sonic.Conn {
	sonic, err := sonic.Dial(sonicURL)
	if err != nil {
		log.Fatalf("Unable to open Sonic UDP connection")
	}
	return sonic
}

func setupTimer(sonic *sonic.Conn, evtChan <-chan int) {
	timer := time.NewTimer(timeout)
	go func(timer *time.Timer, eChan <-chan int) {
		var (
			tunes            []string
			upCount, dwCount int
		)

		for {
			select {
			case v := <-evtChan:
				switch v {
				case -1:
					dwCount++
				case 1:
					upCount++
				}
				if upCount > 0 {
					tunes = append(tunes, upNote(upCount))
				}
				if dwCount > 0 {
					tunes = append(tunes, dwNote(dwCount))
				}

			case <-timer.C:
				if len(tunes) > 0 {
					playList := strings.Join(tunes, "\n")
					sonic.Blast(playList)
					upCount, dwCount, tunes = 0, 0, []string{}
				}
				timer.Reset(timeout)
			}
		}
	}(timer, evtChan)
}

func upNote(inc int) string {
	return hydrate(upTpl, inc)
}
func dwNote(inc int) string {
	return hydrate(dwTpl, inc)
}
func hydrate(tpl *template.Template, inc int) string {
	buff := bytes.NewBufferString("")
	err := tpl.Execute(buff, struct{ Count int }{Count: inc})
	if err != nil {
		log.Fatalf("Unable to hydrate template %s", tpl.Name())
	}
	return buff.String()
}

func up(pod *api.Pod) {
	evtChan <- 1
}
func down(pod *api.Pod) {
	evtChan <- -1
}

func setupSigs(sigs chan os.Signal, done chan<- bool) {
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()
}

func loadTemplates() {
	var err error

	upTpl, err = template.ParseFiles("templates/up.rb")
	if err != nil {
		log.Fatalf("Unable to load up template")
	}

	dwTpl, err = template.ParseFiles("templates/down.rb")
	if err != nil {
		log.Fatalf("Unable to load down template")
	}
}
