package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"github.com/liviosoares/go-watson-sdk/watson"
	"github.com/liviosoares/go-watson-sdk/watson/speech_to_text"
	//"github.com/liviosoares/go-watson-sdk/watson/tone_analyzer" bugs with this, using node to rune tone analyzer
	"io"
	"log"
	"sync"
	"strings"
	"os/exec"
	"github.com/gorilla/websocket"
)

var tpl *template.Template

type Client struct {
	conn *websocket.Conn
}

var client Client

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

type Obj struct {
	Analysis string `json:"analysis"`
	Result []speech_to_text.Result `json:"results,omitempty"`
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {

	var s string = ""
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if req.Method == http.MethodPost {

		f, h, err := req.FormFile("q")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		bs, err := ioutil.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		dst, err := os.Create(filepath.Join("./audio/", h.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = dst.Write(bs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}


		/**
		* Watson Client API
		 */
		c, err := speech_to_text.NewClient(watson.Config{})
		if err != nil {
			log.Printf("client failed %#v\n", err)
		}
		output, stream, err := c.NewStream("", "audio/wav", map[string]interface{}{"continuous": true, "interim_results": false, "timestamps": true})
		if err != nil {
			log.Printf("stream failed %#v %s\n", err, err.Error())
		}

		fn := mp3ToWav(h.Filename)

		g, err := os.Open("audio/" + fn)
		if err != nil {
			log.Printf("stream failed to open audio file %s %s\n", "test_data/speech.wav", err)
		}

		go func() {
			_, err = io.Copy(stream, g)
			if err != nil {
				log.Printf("io failed to copy audio file to API %s\n", err.Error())
			}
		}()

		done := false
		for (!done){
			select {
			case event, ok := <-output:
				if !ok || len(event.Error) > 0 {
					log.Printf("failed to transcribe %#v %s\n", ok, event.Error)
					done = true;
					return
				}
				if len(event.Results) == 0 {
					log.Printf("failed to transcribe, empty results %#v\n", event)
					done = true;
					return
				}

				command := "node test.js '" + event.Results[0].Alternatives[0].Transcript + "'"
				output := execNode(command)

				client.conn.WriteJSON(Obj{Analysis:output, Result:event.Results})
			}
		}
	}

	tpl.ExecuteTemplate(w, "index.hbs", s)
}

func exe_cmd(cmd string, wg *sync.WaitGroup) {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:]
	_, err := exec.Command(head,parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	wg.Done()
}

func execNode(cmd string) string {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:]
	log.Println(head, parts)
	err := os.Chdir("/Users/adamwolf/Desktop")
	if err != nil {
		log.Println(err.Error())
	}

	res, err := exec.Command(head, parts[0], parts[1]).CombinedOutput()
	if err != nil {
		fmt.Printf("%s", err)
	}
	return string(res);
}

func wsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Origin") != "http://"+req.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}

	conn, err := websocket.Upgrade(w, req, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	client.conn = conn
}

func mp3ToWav(mp3 string) (fn string) {
	fn = mp3
	var ext = filepath.Ext(mp3)
	var name = fn[0:len(fn) - len(ext)]
	if ext == ".mp3" {
		newFilename := name + ".wav"
		wg := new(sync.WaitGroup)
		commands := []string{"mpg123 -w ./audio/" + newFilename + " ./audio/" + mp3}
		for _, str := range commands {
			wg.Add(1)
			go exe_cmd(str, wg)
		}
		wg.Wait()
		fn = newFilename
	}
	return
}
