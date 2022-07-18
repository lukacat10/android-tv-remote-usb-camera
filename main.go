package main

import (
	"android-tv-remote-control/ws"
	"encoding/base64"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
)

func connect_adb(w http.ResponseWriter, req *http.Request) {
	exec.Command(`Z:\Android\Sdk\platform-tools\adb.exe`, "connect", "").Start()
	w.WriteHeader(200)
}

func start_app(w http.ResponseWriter, req *http.Request) {
	exec.Command(`Z:\Android\Sdk\platform-tools\adb.exe`, "shell", "am", "start", "-n", "pl.nextcamera/pl.nextcamera.MainActivity").Start()
	w.WriteHeader(200)
}

func kill_app(w http.ResponseWriter, req *http.Request) {
	exec.Command(`Z:\Android\Sdk\platform-tools\adb.exe`, "shell", "am", "force-stop", "pl.nextcamera").Start()
	w.WriteHeader(200)
}

func enter(w http.ResponseWriter, req *http.Request) {
	exec.Command(`Z:\Android\Sdk\platform-tools\adb.exe`, "shell", "input", "keyevent", "66").Start()
	w.WriteHeader(200)
}

func turn_off_device(w http.ResponseWriter, req *http.Request) {
	exec.Command(`Z:\Android\Sdk\platform-tools\adb.exe`, "shell", "input", "keyevent", "26").Start()
	w.WriteHeader(200)
}

func hello(w http.ResponseWriter, req *http.Request) {
	cmd := exec.Command(`Z:\Android\Sdk\platform-tools\adb.exe`, "exec-out", "screencap", "-p")
	resut, _ := cmd.Output()

	w.Write(resut[16+17:])
}

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	//proxy, err := NewProxy("http://127.0.0.1:3000")
	//if err != nil {
	//	panic(err)
	//}

	hub := ws.NewHub()
	go hub.Run()

	go func() {
		for {
			cmd := exec.Command(`Z:\Android\Sdk\platform-tools\adb.exe`, "exec-out", "screencap", "-p")
			resut, _ := cmd.Output()
			sEnc := base64.StdEncoding.EncodeToString(resut[16+17:])
			hub.Broadcast([]byte(sEnc))
		}
	}()

	// handle all requests to your server using the proxy
	//http.HandleFunc("/", ProxyRequestHandler(proxy))
	fs := http.FileServer(http.Dir("./client/build"))
	http.Handle("/", fs)
	http.HandleFunc("/api/get_image", hello)

	http.HandleFunc("/api/start_app", start_app)
	http.HandleFunc("/api/kill_app", kill_app)
	http.HandleFunc("/api/enter", enter)
	http.HandleFunc("/api/turn_off_device", turn_off_device)
	http.HandleFunc("/api/connect_adb", connect_adb)

	http.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})
	log.Fatal(http.ListenAndServe(":12312", nil))
}
