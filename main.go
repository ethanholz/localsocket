package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"

	"github.com/ethereum/go-ethereum/rpc"
)

type ShellService struct{}

type ExecArgs struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type ExecResult struct {
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}

func (s *ShellService) Exec(ctx context.Context, args ExecArgs) (*ExecResult, error) {
	if args.Command == "" {
		return nil, fmt.Errorf("missing command")
	}
	res := ExecResult{
		ExitCode: 0,
		Stdout:   "",
		Stderr:   "",
	}

	cmd := exec.CommandContext(ctx, args.Command, args.Args...)
	stdout, err := cmd.Output()
	res.Stdout = string(stdout)
	if exitError, ok := err.(*exec.ExitError); ok {
		res.ExitCode = exitError.ExitCode()
		res.Stderr = string(exitError.Stderr)
	}

	return &res, nil
}

func main() {
	rpcServer := rpc.NewServer()

	if err := rpcServer.RegisterName("shell", new(ShellService)); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	// JSON-RPC over WebSocket.
	//
	// Development only: "*" allows all origins.
	// In production, restrict this to your actual origin.
	mux.Handle("/rpc", rpcServer.WebsocketHandler([]string{"*"}))

	// Static files.
	//
	// Files in ./public are served from /.
	// Example: ./public/index.html -> http://localhost:8546/
	// Example: ./public/app.js     -> http://localhost:8546/app.js
	publicDir := "./public"

	if err := os.MkdirAll(publicDir, 0o755); err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.Dir(publicDir)))

	slog.Info("serving static files", "url", "http://localhost:8546/")
	slog.Info("serving JSON-RPC websocket", "url", "ws://localhost:8546/rpc")

	if err := http.ListenAndServe(":8546", mux); err != nil {
		panic(err)
	}
}
