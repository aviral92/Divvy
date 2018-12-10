package main

import (
	"errors"
	"log"
	"os"

	"github.com/Divvy/src/pb"
	context "golang.org/x/net/context"
)

type CommonFileListRPCResponse struct {
	fileList *pb.FileList
	err      error
}

/*
*  Utility functions
 */
func GetPeerFromID(nodeID string) (*PeerT, error) {
	for _, p := range Node.netMgr.peers {
		if nodeID == p.ID.String() {
			return &p, nil
		}
	}

	return nil, errors.New("Peer not found")
}

/*
*  RPC Handlers
 */

// Search is called by the NetworkManager on SEARCH RPC call
func SearchHandler(query *pb.SearchQuery) (*pb.FileList, error) {
	// TODO: Call the File manager to get all the files matching name/hash
	if query.IsHash {
		file := Node.fileMgr.searchFileByHash(query.Key)
		log.Printf("Searched file: %v", file.FileName)
	}

	return &pb.FileList{}, nil
}

func GetSharedFilesHandler() (*pb.FileList, error) {
	// TODO: Call the File manager to get all files
	return &pb.FileList{}, nil
}

func DownloadFileRequestHandler(request *pb.DownloadRequest, responseChan chan DownloadFileResponse) {
	/*
	 * 1. Check if the file exists
	 * 2. Start sending the file to the client from the specified offset
	 */
	var err error
	var success *pb.Success

	requestedFile := Node.fileMgr.searchFileByHash(request.Hash)
	if requestedFile == nil {
		err = errors.New("Request file not found")
		success = nil
	} else {
		success = &pb.Success{}
		err = nil
	}

	responseChan <- DownloadFileResponse{
		success: success,
		err:     err}
	if err != nil {
		// Unsuccessful request
		return
	}

	// Successful request. Send file to the client
	file, fileErr := os.Open(requestedFile.Path)
	if fileErr != nil {
		// TODO: This should be relayed to the peer as well.
		log.Printf("[Core] Unable to open file: %v", err)
		return
	}

	_ = file

	peer, err := GetPeerFromID(request.NodeID)
	if err != nil {
		log.Printf("[Core] %v", err)
	}

	_ = peer

	// Create a stream
}

/*
*  CLI Handlers
 */

func PeersSearchFile(searchQuery string) (*pb.FileList, error) {
	// Send a search RPC to all peers and wait for their responses
	searchResponse := make(chan CommonFileListRPCResponse)
	remainingResponses := len(Node.netMgr.peers)
	var peerFiles *pb.FileList

	for _, peer := range Node.netMgr.peers {
		go func(client pb.DivvyClient) {
			fileList, err := client.Search(context.Background(),
				&pb.SearchQuery{
					IsHash: false,
					Key:    searchQuery})
			searchResponse <- CommonFileListRPCResponse{fileList: fileList,
				err: err}

		}(peer.Client)
	}

	// Collecting responses
	for {
		resp := <-searchResponse
		if resp.err != nil {
			peerFiles.Files = append(peerFiles.Files, resp.fileList.Files...)
		}

		/*
		 *  This could be a BUG. Not sure what will happen when when a grpc fails
		 */
		remainingResponses--
		if remainingResponses <= 0 {
			break
		}
	}

	return peerFiles, nil
}

func PeersGetSharedFiles() (*pb.FileList, error) {
	fileListResponse := make(chan CommonFileListRPCResponse)
	remainingResponses := len(Node.netMgr.peers)
	var peerFiles *pb.FileList

	for _, peer := range Node.netMgr.peers {
		go func(client pb.DivvyClient) {
			fileList, err := client.GetSharedFiles(context.Background(),
				&pb.Empty{})
			fileListResponse <- CommonFileListRPCResponse{fileList: fileList,
				err: err}

		}(peer.Client)
	}

	// Collecting responses
	for {
		resp := <-fileListResponse
		if resp.err != nil {
			peerFiles.Files = append(peerFiles.Files, resp.fileList.Files...)
		}

		/*
		 *  This could be a BUG. Not sure what will happen when when a grpc fails
		 */
		remainingResponses--
		if remainingResponses <= 0 {
			break
		}
	}

	return peerFiles, nil
}
