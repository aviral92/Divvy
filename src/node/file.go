package main

import (
	"os"
	//"time"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type File struct {
	//Delete  bool        `json:"de,omitempty"`
	Path string `json:"p,omitempty"`
	Size int64  `json:"s,omitempty"`
	//Mode    os.FileMode `json:"m,omitempty"`
	//ModTime *time.Time  `json:"t,omitempty"`
	IsDir   bool   `json:"d,omitempty"`
	Hash    uint64 `json:"h,omitempty",hash:"ignore"`
	Content []byte `json:"c,omitempty"`
}

type FileManager struct {
	directoryPath     string
	files             []*File
	availableToOthers bool
	//grpcServer        *grpc.Server
}

func NewFileManager(pathToDirectory string) *FileManager {
	fileMgr := &FileManager{}
	fileMgr.directoryPath = pathToDirectory

	fileMgr.files = []*File{}
	f := &File{Path: "/home/vagrant/go/src/github.com/Divvy/README.md"}

	fileMgr.files = append(fileMgr.files, f)

	return fileMgr
}

//check if file exists
func (file *File) exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (fileMgr *FileManager) displayDirectory() {

	files, err := ioutil.ReadDir(fileMgr.directoryPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		log.Println(f.Name())
	}
}

/*func (file *File) GetHashes(names []string) []string {
hashes := make([]string, len(names))
for i, name := range names {
hashes[i] = GetHash(name)
}
return hashes
}*/

func (file *File) GetHash(filePath string) string {
	input := strings.NewReader(filePath)

	hash := sha256.New()
	if _, err := io.Copy(hash, input); err != nil {
		log.Fatal(err)
	}
	sum := hash.Sum(nil)

	return base64.URLEncoding.EncodeToString(sum)

}
