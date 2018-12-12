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
	if fileList.Files == nil {
		return
	}

	for _, file := range fileList.Files {
		fmt.Println(file)
	}
}

func DisplaySearchResult(searchParam string, isHash bool) {
	// TODO: Complete this function
	fileList, err := PeersSearchFile(searchParam, isHash)
	if err != nil {
		fmt.Println("Error getting file via search parameter specified from peers")
	}
	if fileList.Files == nil {
		return
	}

	for _, file := range fileList.Files {
		fmt.Println(file)
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
	case "search":
		if len(commands) == 3 {
			if commands[1] == "hash" {
				DisplaySearchResult(commands[2], true)
			} else if commands[1] == "name" {
				DisplaySearchResult(commands[2], false)
			}
		}
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
