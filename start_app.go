package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/vrianta/agai/v1/config"
	"github.com/vrianta/agai/v1/log"
)

var clientChan chan string = make(chan string)
var mu sync.Mutex
var clients = make(map[chan string]bool)
var run_app *exec.Cmd
var migrate_models *exec.Cmd
var migrate_components *exec.Cmd
var browser_proc *exec.Cmd
var output_app_name string
var wait_for_hot_reloader sync.WaitGroup

// this falag is needed because after we restart the application on componenent file change it again pdate the components so we have to detect
// if we are not detecting the component after a application start or an application restart
var FLAG_restarted_application_after_component_change bool = false

func init() {
	if runtime.GOOS == "windows" {
		output_app_name = "app.exe"
	} else {
		output_app_name = "app"
	}

	run_app = new_app_cmd()
	migrate_models = new_migrate_model_cmd()
	migrate_components = new_migrate_component_cmd()
	// config.init()
}

func start_app() {
	if !f.start_app {
		log.Debug("start_app flag is false, exiting.")
		return
	}

	log.Info("Starting server process...")
	if err := run_app.Start(); err != nil {
		panic("Failed to Start the Server - " + err.Error())
	}

	go start_hot_reload()
	go WatchFolders(1 * time.Second)

	url := "http://" + func() string {
		if config.GetHost() == "" {
			return "localhost"
		}
		return config.GetHost()
	}() + ":" + config.GetPort()

	log.Info("Opening browser at: %s", url)
	fmt.Println("Browser URL:", url)
	if err := openBrowser(url); err != nil {
		log.Error("Failed to open browser: %s", err.Error())
	}

	defer func() {
		log.Info("Killing server process...")

		run_app.Process.Kill()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Info("Running. Press Ctrl+C to exit.")

	<-sigChan

	broadcast("close") // clossing all the tabs which are running the app
	wait_for_hot_reloader.Wait()

	fmt.Println("Interrupt received. Exiting...")
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		// The "" is a placeholder for the window title
		cmd = "cmd"
		args = []string{"/c", "start", "", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}

	browser_proc = exec.Command(cmd, args...)
	return browser_proc.Start()
}

func start_hot_reload() {
	log.Info("Starting hot reload server on port 8888")
	http.HandleFunc("/hot-reload", sseHandler)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Error("Hot reload server failed: %s", err.Error())
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("SSE client connected")

	// SSE headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Each client gets its own channel
	client := make(chan string)

	mu.Lock()
	clients[client] = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(clients, client)
		mu.Unlock()
		close(client)
		log.Debug("SSE client disconnected")
	}()

	// Listen and send events
	for {
		select {
		case msg, ok := <-client:
			if !ok {
				return
			}
			if msg == "close" {
				fmt.Println("clossing tabs")
				wait_for_hot_reloader.Done()
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return // Client closed connection
		}
	}
}

func broadcast(msg string) {
	mu.Lock()
	defer mu.Unlock()

	for ch := range clients {
		select {
		case ch <- msg:
			if msg == "close" {
				wait_for_hot_reloader.Add(1)
			}

		default:
			// Drop slow or dead clients
			close(ch)
			delete(clients, ch)
		}
	}
}

func WatchFolders(interval time.Duration) {
	log.Info("Starting file watcher...")

	modTimestamps := make(map[string]time.Time)
	componentTimestamps := make(map[string]time.Time)
	generalTimestamps := make(map[string]time.Time)

	for {
		moduleChanged := checkDir("models", modTimestamps)
		componentChanged := checkDir("components", componentTimestamps)
		generalChanged := checkGeneral(".", generalTimestamps)

		if moduleChanged {
			log.Info("Model file changed")
			onModuleChange()
			time.Sleep(interval)
			broadcast("reload")
		}
		if componentChanged {
			log.Info("Component files changed")
			onComponentChange()
			time.Sleep(interval)
			// log.Warn("We Deteceted a Component Change Please -restart the application so | curretnly hot reload is not supported for components")
			broadcast("reload")
		}
		if generalChanged {
			log.Info("General files changed")
			onGeneralChange()
			time.Sleep(interval)
			broadcast("reload")
		}

	}
}

func checkDir(root string, lastMod map[string]time.Time) bool {
	changed := false

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		modTime := info.ModTime()

		// If we've never seen this file before, add it to the map but don't trigger change
		if _, ok := lastMod[path]; !ok {
			lastMod[path] = modTime
			return nil
		}

		if modTime.After(lastMod[path]) {
			log.Debug("Change detected in %s", path)
			lastMod[path] = modTime
			changed = true
		}
		return nil
	})

	if err != nil {
		log.Error("Error watching %s: %s", root, err.Error())
	}

	return changed
}

func checkGeneral(root string, lastMod map[string]time.Time) bool {
	changed := false

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip specific directories
		if info.IsDir() && (info.Name() == "models" || info.Name() == "components" || info.Name() == ".git") {
			return filepath.SkipDir
		}
		if info.Name() == "sessions.data" || info.Name() == output_app_name {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		modTime := info.ModTime()

		// Don't trigger change on first encounter
		if _, ok := lastMod[path]; !ok {
			lastMod[path] = modTime
			return nil
		}

		if modTime.After(lastMod[path]) {
			log.Debug("General change detected in %s", path)
			lastMod[path] = modTime
			changed = true
		}

		return nil
	})

	if err != nil {
		log.Error("Error watching general files: %s", err.Error())
	}

	return changed
}

func app_build() {
	if err := exec.Command("go", "build", "-o", output_app_name, ".").Run(); err != nil {
		log.Error("Failed to Build the app - %s", err.Error())
	}
}

func new_app_cmd() *exec.Cmd {
	app_build()
	r := exec.Command("./"+output_app_name, "-ss")
	r.SysProcAttr = &syscall.SysProcAttr{
		// CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP, // Windows only
	}
	return r
}

func new_migrate_model_cmd() *exec.Cmd {
	return exec.Command("go", "run", ".", "-mm")
}

func new_migrate_component_cmd() *exec.Cmd {
	return exec.Command("go", "run", ".", "-mc")
}

func onModuleChange() {
	log.Info("Restarting server due to module changes...")
	if err := run_app.Process.Kill(); err != nil {
		log.Error("Failed to kill server: %s", err.Error())
	}

	run_app.Wait()
	run_app = new_app_cmd()

	migrate_models = new_migrate_model_cmd()

	if err := migrate_models.Run(); err != nil {
		log.Error("Model migration failed: %s", err.Error())
	}
	if model_migration_output, err := migrate_models.Output(); err == nil {
		fmt.Println("Output for Model Migration : ", model_migration_output)
	} else {
		panic(err.Error())
	}

	if err := run_app.Start(); err != nil {
		panic("Failed to restart server: " + err.Error())
	}
	log.Info("Server restarted after module migration.")
}

func onComponentChange() {
	// now on the first run it will skip the component change
	if !FLAG_restarted_application_after_component_change {
		FLAG_restarted_application_after_component_change = true
		return
	}
	log.Info("Restarting server due to component changes...")
	if err := run_app.Process.Kill(); err != nil {
		log.Error("Failed to kill server: %s", err.Error())
	}

	run_app.Wait()
	run_app = new_app_cmd()

	migrate_components = new_migrate_component_cmd()

	if err := migrate_components.Run(); err != nil {
		log.Error("Component migration failed: %s", err.Error())
	}

	fmt.Println(migrate_components.Output())

	if err := run_app.Start(); err != nil {
		panic("Failed to restart server: " + err.Error())
	}
	log.Info("Server restarted after component migration.")
	// and if we are restarted the application means by default the components will be updated in the disk
	// hence we need to make it false again so that we skip the next change scan
	FLAG_restarted_application_after_component_change = false
}

func onGeneralChange() {
	log.Warn("Restarting server due to general changes...")
	if err := run_app.Process.Kill(); err != nil {
		panic("Failed to kill server: " + err.Error())
	}

	run_app.Wait()          // waiting for the app to close
	run_app = new_app_cmd() // creating new cmd to start the application

	if err := run_app.Start(); err != nil {
		panic("Failed to restart server: " + err.Error())
	}
	// later I should check if I can connect to the started server before proceeding
	log.Info("Server restarted due to general change.")
}

func migrate_model_and_component() {
	if f.migrate_model {
		// start migrating the models

		log.Info("Migrating Models")
		migrate_models := new_migrate_model_cmd()

		if err := migrate_models.Run(); err != nil {
			log.Error("Model migration failed: %s", err.Error())
		}
	}

	if f.migrate_component {

		log.Info("Migrating Components")
		migrate_components := new_migrate_component_cmd()

		if output, err := migrate_components.Output(); err != nil {
			panic("Component migration failed: " + err.Error())
		} else {
			log.Info("%s", string(output))
		}

		migrate_components.Wait()

	}
}
