package main

import (
	"errors"
	"io"
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

func SearchHandler(query *pb.SearchQuery) (*pb.FileList, error) {
	fileList := &pb.FileList{}
	fileList.NodeID = Node.netMgr.ID.String()

	var file *File

	if query.IsHash {
		file = Node.fileMgr.searchFileByHash(query.Key)
	} else {
		file = Node.fileMgr.searchFileByName(query.Key)
	}

    if file != nil {
        fileList.Files = append(fileList.Files, &pb.File{
            Name: file.FileName,
            Hash: file.Hash})
    }
	return fileList, nil
}

func GetSharedFilesHandler() (*pb.FileList, error) {
	fileList := Node.fileMgr.GetSharedFilesList()
	return &fileList, nil
}

// TODO: Function too complex. Try to break into smaller functions
func DownloadFileRequestHandler(request *pb.DownloadRequest, responseChan chan DownloadFileResponse) {
	/*
	 * 1. Check if the file exists
	 * 2. Start sending the file to the client from the specified offset
	 */
	var (
		err     error
		success *pb.Success
		fileBuf = make([]byte, Node.config.ChunkSizeInt)
	)

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

	peer, err := GetPeerFromID(request.NodeID)
	if err != nil {
		log.Printf("[Core] %v", err)
		return
	}

	// Create a stream
	stream, err := peer.Client.ReceiveFile(context.Background())
	if err != nil {
		log.Printf("[Core] Unable to open stream %v", err)
		return
	}

	defer stream.CloseSend()

	// Start sending the file
	for {
		lenRead, err := file.Read(fileBuf)
		if err != nil {
			goto FINISH
		}
		err = stream.Send(&pb.FileChunk{Hash: requestedFile.Hash,
			Content: fileBuf[:lenRead],
			Offset:  0})
		if err != nil {
			goto FINISH
		}
	}

FINISH:
	if err != nil {
		if err == io.EOF {
			status, err := stream.CloseAndRecv()
			if err != nil {
				log.Printf("[Core] %v", err)
			}
			log.Printf("[Core] Transfer Status: %v", status)
			return
		}
		log.Printf("[Core] Unable to send file %v", err)
	}
}

/*
*  CLI Handlers
 */

func PeersSearchFile(searchQuery string, isHash bool) (pb.FileList, error) {
	// Send a search RPC to all peers and wait for their responses
	searchResponse := make(chan CommonFileListRPCResponse)
	remainingResponses := len(Node.netMgr.peers)
	var peerFiles pb.FileList

	for _, peer := range Node.netMgr.peers {
		go func(client pb.DivvyClient) {
			fileList, err := client.Search(context.Background(),
				&pb.SearchQuery{
					IsHash: isHash,
					Key:    searchQuery})
			searchResponse <- CommonFileListRPCResponse{fileList: fileList,
				err: err}

		}(peer.Client)
	}

	// Collecting responses
	for {
		if remainingResponses <= 0 {
			break
		}

		resp := <-searchResponse
		if resp.err != nil {
			log.Printf("[Core] Error in response %v", resp.err)
		}
		if resp.fileList.Files != nil {
			peerFiles.Files = append(peerFiles.Files, resp.fileList.Files...)
		}

		/*
		 *  This could be a BUG. Not sure what will happen when when a grpc fails
		 */
		remainingResponses--
	}

	return peerFiles, nil
}

func PeersGetSharedFiles() (pb.FileList, error) {
	fileListResponse := make(chan CommonFileListRPCResponse)
	remainingResponses := len(Node.netMgr.peers)
	var peerFiles pb.FileList

	for _, peer := range Node.netMgr.peers {
		go func(client pb.DivvyClient) {
			if client == nil {
				log.Printf("[Core] Client is nil")
			}
			fileList, err := client.GetSharedFiles(context.Background(),
				&pb.Empty{})
			if err != nil {
				log.Printf("[Core] Unable to get files from a peer %v", err)
			}
			fileListResponse <- CommonFileListRPCResponse{fileList: fileList,
				err: err}

		}(peer.Client)
	}

	// Collecting responses
	for {
		/*
		 *  This could be a BUG. Not sure what will happen when when a grpc fails
		 */
		if remainingResponses <= 0 {
			break
		}

		resp := <-fileListResponse
		if resp.err != nil {
			log.Printf("[Core] Error in response %v", resp.err)
		}
		if resp.fileList.Files != nil {
			peerFiles.Files = append(peerFiles.Files, resp.fileList.Files...)
		}

		remainingResponses--
	}

	return peerFiles, nil
}
