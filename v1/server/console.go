package server

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	Config "github.com/vrianta/agai/v1/config"
	Log "github.com/vrianta/agai/v1/log"
)

func (s *instance) ServeConsole() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) // Listen for Ctrl+C

	go func() {
		<-quit // Wait for interrupt signal
		fmt.Println("\nShutting down server...")
		s.stopServer()
		os.Exit(0) // Exit program gracefully
	}()

	printHelp()
	handleOutput() // Start prompt handler for consistent `: ` prompt

	// Command input loop
	scanner := bufio.NewScanner(os.Stdin)
	for ; scanner.Scan(); fmt.Print(": ") {
		input := scanner.Text()

		switch input {
		case "stop":
			s.stopServer()
		case "start":
			s.startServer()
		case "exit":
			fmt.Println("Exiting...")
			s.stopServer()
			os.Exit(0)
		case "r": // restart
			fmt.Println("Restarting...")
			s.stopServer()
			s.startServer()
		case "restart":
			fmt.Println("Restarting...")
			s.stopServer()
			s.startServer()
		default:
			fmt.Println("Please Print -h to get available commands")
			continue
		}
	}
}

func (s *instance) startServer() {
	if s.state {
		fmt.Println("Server is in Starting State so skipping the process...")
		return
	}
	// Define the server configuration
	s.server = &http.Server{
		Addr: Config.GetWebConfig().Host + ":" + Config.GetWebConfig().Port, // Host and port
	}

	Log.WriteLogf("%s", "Server Starting at "+Config.GetWebConfig().Host+":"+Config.GetWebConfig().Port)

	s.server.ListenAndServe()
	s.state = true
}

func (s *instance) stopServer() {
	if !s.state {
		fmt.Println("Server is Already Stopped | can't stop it...")
		return
	}
	// Create a timeout context (5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := s.server.Shutdown(ctx); err != nil {
		fmt.Println("Shutdown Error:", err)
	} else {
		fmt.Println("Server shutdown successfully")
	}

	s.state = false
	// s.Sessions = map[string]Session{}
}

// Function to ensure `:` appears after every print
func handleOutput() {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	fmt.Fprint(oldStdout, ": ")
	scanner := bufio.NewScanner(r)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Fprintln(oldStdout, text) // Original content
			fmt.Fprint(oldStdout, ": ")   // Always show `:` prompt after print
		}
	}()
}

// Help function to display available commands
func printHelp() {
	fmt.Println("\nCommand you can Use ----------------")
	fmt.Println(" start    - Start the server")
	fmt.Println(" stop     - Stop the server")
	fmt.Println(" restart  - Restart the server")
	fmt.Println(" r        - Shortcut for restart")
	fmt.Println(" exit     - Stop the server and exit the program")
	fmt.Println(" -h       - Display available commands")
	fmt.Println("------------------------------------")
}
