package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func DisplayPeers() {
	if len(Node.netMgr.peers) == 0 {
		fmt.Println("No peers found")
		return
	}
	fmt.Println("List of known peers")
	for _, p := range Node.netMgr.peers {
		fmt.Println(p)
	}
}

func DisplayMyFiles() {
	if len(Node.fileMgr.SharedFiles) == 0 {
		fmt.Println("No files shared")
	}

	fmt.Println("Shared files")
	for _, file := range Node.fileMgr.SharedFiles {
		fmt.Println(file.FileName)
	}
}

func DisplayPeerFiles() {
	fileList, err := PeersGetSharedFiles()
	if err != nil {
		fmt.Println("Error getting file list from peers")
	}
	for _, file := range fileList.Files {
		fmt.Sprintf("%v (%v)", file.Name, file.Hash)
	}
}

func ExecuteCommand(cmdStr string) {
	cmdStr = strings.TrimRight(cmdStr, "\n")
	commands := strings.Split(cmdStr, " ")
	if len(commands) <= 0 {
		return
	}
	switch commands[0] {
	case "": // To consume empty return key
	case "peers":
		DisplayPeers()
	case "files":
		DisplayPeerFiles()
	case "myfiles":
		DisplayMyFiles()
	default:
		log.Printf("[CLI] Command not recognized")
	}
}

func cli() {
	// This is an infinite loop that will keep reading commands
	log.Printf("[CLI] Listening for user input")
	cmdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		cmdStr, _ := cmdReader.ReadString('\n')
		ExecuteCommand(cmdStr)
	}
}
