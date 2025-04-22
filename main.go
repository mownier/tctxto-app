package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	browserAppName := "Browser App"
	localIp := ""

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("error getting network interfaces: %v\n", err)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIp = ipnet.IP.String()
				break
			}
		}
	}

	if len(localIp) == 0 {
		log.Fatalln("unable to determine local IP")
	}

	envVars := map[string]string{
		"TCTXTO_SERVER_PORT":      "3232",
		"TCTXTO_ENABLE_RELECTION": "false",
		"TCTXTO_CONSUMERS":        "",

		"TCTXTO_PROXY_BA": "",
		"TCTXTO_PROXY_AO": "",
		"TCTXTO_PROXY_DP": "2121",

		"TCTXTO_PROXY_ORIGIN": "",
		"TCTXTO_SERVER_PK":    "",
		"TCTXTO_CLIENT_PORT":  "2323",
	}

	err = godotenv.Load()
	if err != nil {
		log.Printf("the .env file cannot be loaded: %v\n", err)
		log.Println("will generate .env file")
		vars := []string{}
		for k, v := range envVars {
			vars = append(vars, fmt.Sprintf("%s=%s", k, v))
		}
		s := strings.Join(vars, "\n")
		err = os.WriteFile("./.env", []byte(s), 0644)
		if err != nil {
			log.Printf("unable to generate .env file: %v\n", err)
		}
		err = godotenv.Load()
		if err != nil {
			log.Printf("generated .env file cannot be loaded: %v\n", err)
		}
	}

	var consumersData []byte
	consumersPath := os.Getenv("TCTXTO_CONSUMERS")
	if len(consumersPath) == 0 {
		consumersPath = "./consumers.json"
		if _, err := os.Stat(consumersPath); err == nil {
			consumersData, err = os.ReadFile(consumersPath)
			if err != nil {
				log.Fatalf("error reading existing consumers file at %s: %v\n", consumersPath, err)
			}
		} else if os.IsNotExist(err) {
			c := consumer{
				Name:      browserAppName,
				PublicKey: uuid.New().String(),
			}
			consumersData, err = json.Marshal([]consumer{c})
			if err != nil {
				log.Fatalf("error marshalling generated consumers file at %s: %v\n", consumersPath, err)
			}
			err = os.WriteFile(consumersPath, consumersData, 0644)
			if err != nil {
				log.Fatalf("error writing generated consumers file at %s: %v\n", consumersPath, err)
			}
		} else {
			log.Fatalf("cannot determine existence of consumers file at %s: %v\n", consumersPath, err)
		}
	} else {
		consumersData, err = os.ReadFile(consumersPath)
		if err != nil {
			log.Fatalf("error reading custom consumers file at %s: %v\n", consumersPath, err)
		}
	}

	var consumers []*consumer
	err = json.Unmarshal(consumersData, &consumers)
	if err != nil {
		log.Fatalf("error unmarshalling consumers from %s: %v\n", consumersPath, err)
	}

	if len(consumers) == 0 {
		log.Fatalf("no consumers found in %s\n", consumersPath)
	}

	var browserAppPk = ""
	for _, c := range consumers {
		if c.Name == browserAppName {
			browserAppPk = c.PublicKey
			break
		}
	}
	if len(browserAppPk) == 0 {
		log.Fatalf("public key for consumer with name '%s' is empty in %s\n", browserAppName, consumersPath)
	}

	serverPort := os.Getenv("TCTXTO_SERVER_PORT")
	if len(serverPort) == 0 {
		serverPort = "3232"
	}

	proxyPort := os.Getenv("TCTXTO_PROXY_DP")
	if len(proxyPort) == 0 {
		proxyPort = "2121"
	}

	clientPort := os.Getenv("TCTXTO_CLIENT_PORT")
	if len(clientPort) == 0 {
		clientPort = "2323"
	}

	enableReflection := "false"
	enableReflectionEnvVar := os.Getenv("TCTXTO_ENABLE_RELECTION")
	if len(enableReflectionEnvVar) > 0 &&
		(strings.ToLower(enableReflectionEnvVar) == "false" ||
			strings.ToLower(enableReflectionEnvVar) == "true") {
		enableReflection = enableReflectionEnvVar
	}

	envVars["TCTXTO_SERVER_PORT"] = serverPort
	envVars["TCTXTO_ENABLE_RELECTION"] = enableReflection
	envVars["TCTXTO_CONSUMERS"] = consumersPath

	envVars["TCTXTO_PROXY_BA"] = fmt.Sprintf("%s:%s", localIp, serverPort)
	envVars["TCTXTO_PROXY_AO"] = fmt.Sprintf("http://%s:%s", localIp, clientPort)
	envVars["TCTXTO_PROXY_DP"] = proxyPort

	envVars["TCTXTO_PROXY_ORIGIN"] = fmt.Sprintf("http://%s:%s", localIp, proxyPort)
	envVars["TCTXTO_SERVER_PK"] = browserAppPk
	envVars["TCTXTO_CLIENT_PORT"] = clientPort

	for k, v := range envVars {
		if len(os.Getenv(k)) == 0 {
			if err := os.Setenv(k, v); err != nil {
				log.Printf("problem setting default value %s for %s\n", v, k)
			}
		}
	}

	executables := []executable{
		{
			name:    "server",
			command: tctxtosvCommand,
		},
		{
			name:    "proxy",
			command: tctxtopxCommand,
		},
		{
			name:    "client",
			command: tctxtoclCommand,
		},
	}

	for _, e := range executables {
		fmt.Println("executable")
		fmt.Printf("  name: %s\n", e.name)
		fmt.Printf("  command: %s\n", e.command)
	}

	fmt.Println("environment variables")
	for k, v := range envVars {
		fmt.Printf("  %s=%s\n", k, v)
	}

	var wg sync.WaitGroup
	execProcesses := make(map[string]*exec.Cmd)
	var execProcessesMu sync.Mutex

	ctx, cancel := context.WithCancel(context.Background())

	for _, e := range executables {
		wg.Add(1)
		go func(e executable) {
			defer wg.Done()

			cmd := exec.CommandContext(ctx, e.command)

			err = cmd.Start()
			if err != nil {
				log.Printf("unable to start %s: %v\n", e.name, err)
			}

			execProcessesMu.Lock()
			execProcesses[e.name] = cmd
			execProcessesMu.Unlock()

			err = cmd.Wait()
			if err != nil && ctx.Err() == nil {
				log.Printf("%s exited with error: %v\n", e.name, err)
			} else if err == nil {
				log.Printf("%s exited normally\n", e.name)
			}
		}(e)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Signal received: %v. Shutting down...\n", sig)
		cancel()

		shutdownTimeout := time.Duration(5 * time.Second)
		shutdownTimer := time.NewTimer(shutdownTimeout)

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			log.Println("All executables have terminated")
		case <-shutdownTimer.C:
			log.Println("Timeout reached. Forcefully terminate remaining executable...")
			execProcessesMu.Lock()
			for name, cmd := range execProcesses {
				if cmd.Process != nil {
					log.Printf("Forcefully terminate %s (PID %d)\n", name, cmd.Process.Pid)
					if err := cmd.Process.Kill(); err != nil {
						log.Printf("Error terminating %s: %v\n", name, err)
					}
				}
			}
			execProcessesMu.Unlock()
		}

		os.Exit(0)
	}()

	wg.Wait()
	log.Println("Shutdown")
}

type executable struct {
	name    string
	command string
}

type consumer struct {
	PublicKey string `json:"public_key"`
	Name      string `json:"name"`
}
