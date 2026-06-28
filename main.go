package main

// TODO: add in readme
// To learn about native-messaging protocol (common to browsers like Chrome or Firefox)
// see https://developer.chrome.com/docs/extensions/develop/concepts/native-messaging#native-messaging-host-protocol

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/colangelo/mozeidon-z-messaging/models"
	"github.com/james-barrow/golang-ipc"
	"github.com/rickypc/native-messaging-host"
)

func main() {
	if handled, out := handleFlags(os.Args[1:]); handled {
		fmt.Println(out)
		os.Exit(0)
	}

	if err := webBrowserProxy(); err != nil {
		log.Printf("Error in mozeidon_native_app: %v", err)
	}
}

type IpcIncomingMessage struct {
	Command string `json:"command"        binding:"required"`
	Args    string `json:"args,omitempty"`
}

// isEndOfStream reports whether the browser sent the {"data":"end"} terminator.
// Parses the decoded map instead of byte-comparing marshaled JSON, so key order
// / whitespace can't break streaming.
func isEndOfStream(response *host.H) bool {
	if response == nil {
		return false
	}
	d, ok := (*response)["data"]
	return ok && d == "end"
}

func webBrowserProxy() error {
	browserMessagingClient := (&host.Host{}).Init()

	// Reclaim ipc sockets leaked by previous instances that were hard-killed
	// (SIGKILL/crash) and so never ran their own shutdown cleanup.
	sweepStaleSockets()

	// Step 1. Register this running native-app profile into the ProfileDirectory
	var nativeAppProfile *models.NativeAppProfile
	firstMessage := &host.H{}
	if err := browserMessagingClient.OnMessage(os.Stdin, firstMessage); err != nil {
		return fmt.Errorf("Error receiving message from browser: %w", err)
	}
	response, err := json.Marshal(firstMessage)
	if err != nil {
		return fmt.Errorf("error parsing registration response message: %w", err)
	}
	var registrationData models.RegistrationInfoResponse
	if err := json.Unmarshal(response, &registrationData); err != nil {
		return fmt.Errorf("error parsing registration message: %w", err)
	}
	nativeAppProfile, err = models.GetNativeAppProfile(&registrationData)
	if err != nil {
		return fmt.Errorf("error building native-app profile: %w", err)
	}

	profileDataDir, err := models.GetProfileDirectory()
	if err != nil {
		return fmt.Errorf("Error getting the profile directory: %w", err)
	}

	jsonProfile, err := json.MarshalIndent(nativeAppProfile, "", "  ")
	jsonProfilePath := filepath.Join(profileDataDir, nativeAppProfile.FileName)

	if err := os.WriteFile(jsonProfilePath, jsonProfile, 0644); err != nil {
		return fmt.Errorf("error writing profile file: %w", err)
	}

	// Step 2. Start IPC server
	ipcConfig := &ipc.ServerConfig{

		Encryption:        true,  // allows encryption to be switched off (bool - default is true)
		UnmaskPermissions: false, // single-user: native-app and CLI run as the same user
	}

	ipcServer, err := ipc.StartServer(nativeAppProfile.IpcName, ipcConfig)
	if err != nil {
		return fmt.Errorf("Error starting %s ipc-server: %w", nativeAppProfile.IpcName, err)
	}

	/*
		cleanup removes everything this instance created: the IPC unix socket and the
		profile-registry file, so nothing is left behind. ipcServer.Close() unlinks
		the socket via its listener; the explicit os.Remove is a belt-and-suspenders
		in case a future golang-ipc stops unlinking. Idempotent.
	*/
	cleanup := func() {
		ipcServer.Close()
		os.Remove(socketPath(nativeAppProfile.IpcName))
		os.Remove(jsonProfilePath)
	}

	/*
		Clean up on exits triggered by the browser-extension: handle SIGTERM/SIGINT
		sent from the browser-extension. Won't catch SIGKILL/Windows — those
		leftovers are reclaimed by sweepStaleSockets() on the next startup.
		@see https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Native_messaging#closing_the_native_app
	*/
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		cleanup()
		os.Exit(0)
	}()

	// Also clean up when the proxy returns with an error (e.g. the browser closed
	// stdin, so the next read hits EOF).
	defer cleanup()

	// Listen to client, and handle incoming message
	for {
		message, _ := ipcServer.Read()
		if message.MsgType > 0 {

			// Parse incoming message
			incomingMessage := IpcIncomingMessage{}
			if err := json.Unmarshal(message.Data, &incomingMessage); err != nil {
				log.Printf("skipping malformed ipc message: %v", err)
				continue
			}

			// Send incoming message to browser
			request := &host.H{"payload": incomingMessage}
			if err := browserMessagingClient.PostMessage(os.Stdout, request); err != nil {
				return fmt.Errorf("Error posting message to browser: %w", err)
			}

			for {
				// Wait for browser messages
				// browser may send many messages before the data:end message
				response := &host.H{}
				if err := browserMessagingClient.OnMessage(os.Stdin, response); err != nil {
					return fmt.Errorf("Error receiving message from browser: %w", err)
				}

				// send back browser message to client
				responseMessage, err := json.Marshal(response)
				if err != nil {
					return fmt.Errorf("error marshaling browser response: %w", err)
				}
				if err := ipcServer.Write(1, responseMessage); err != nil {
					return fmt.Errorf("error writing to ipc server: %w", err)
				}
				// end of browser response for the incoming message
				if isEndOfStream(response) {
					break
				}
			}
		}
	}
}
