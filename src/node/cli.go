package main

import (
		//"fmt"
		"net"
		"log"
		"github.com/urfave/cli"
		"os"
       )

func Run() {
	app := cli.NewApp()
	app.Name = "divvy"
	app.Usage = "fight the loneliness!"
	app.Commands = []cli.Command{
	{
	Name:    "setup",
	Aliases: []string{"s"},
	Usage:   "Initial setup",
	Action: func(c *cli.Context) error {
		 //fmt.Println("added task: ", c.Args().First())
		 // Initialize node
		Node = newNodeT()

		initNode(&Node)

		// Discovery listener. Do this before sending the discovery messages
		go Node.netMgr.ListenForDiscoveryMessages()

		Node.netMgr.DiscoverPeers()
		//Node.fileMgr.displayDirectory()
		// Once everything is setup start listening. This call is blocking
		// Do not put any logic after gRPC serve

		// gRPC server
		conn, err := net.Listen("tcp", controlPort)
		if err != nil {
			log.Fatalf("[Node] Failed to open port %v because %v", controlPort, err)
		}
		log.Printf("[Node] Listening on port %v", controlPort)

		if Node.netMgr.address == nil {
			log.Printf("[Node] Network manager has no address")
			goto EXIT
		}

		err = Node.netMgr.grpcServer.Serve(conn)
		if err != nil {
			 log.Fatalf("[Node] Failed to serve %v", err)
		}

		EXIT:
			log.Printf("[Node] Bye from Divvy!")
			return nil
		 },
	   },
	   {
		Name:    "show",
		Aliases: []string{"ls"},
		Usage:   "show shared files in directory",
		Action: func(c *cli.Context) error {
			log.Println("njnjnj")
			Node.fileMgr.displayDirectory()
			return nil
		},
	   },
	}

	err := app.Run(os.Args)
	if err != nil {
	     log.Println(err)
	}
}
