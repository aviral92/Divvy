package main;

import (
    "log"

	context "golang.org/x/net/context"
    "github.com/Divvy/src/pb"
)

type CommonFileListRPCResponse struct {
    fileList    *pb.FileList
    err         error
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
                              Key: searchQuery})
            searchResponse <- CommonFileListRPCResponse{fileList: fileList,
                                                        err: err}

        } (peer.Client)
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

        } (peer.Client)
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
