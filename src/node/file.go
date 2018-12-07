package main

import (
	"os"
	//"time"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type File struct {
	FileName	string
	Path		string
	Size		int64
	IsDir		bool
	Hash		string
}

type FileManager struct {
	directoryPath     string
	files             []*File
}

func NewFileManager(directoryPath string) *FileManager {
	fileMgr := &FileManager{}
	fileMgr.directoryPath = directoryPath

	fileMgr.files = []*File{}
	files, err := ioutil.ReadDir(fileMgr.directoryPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		log.Println(f.Name())

		pathToFile := directoryPath + f.Name()
		fi, err := os.Stat(pathToFile)

		if err != nil {
			log.Fatal(err)
		}
		//fileHash := computeHash(pathToFile)

		f := &File{	FileName :	f.Name(),
				Path	 :	fileMgr.directoryPath,
				IsDir	 :	false,
				Size	 :	fi.Size(),
				//Hash	 :	fileHash
				}

		fileMgr.files = append(fileMgr.files, f)
	}

	return fileMgr
}

//check if file exists
func (file *File) checkIfFileExists(name string) bool {
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

func (file *File) setHash() {
	file.Hash = file.computeHash(file.Path + file.FileName)
}

func (file *File) getHash(filePath string) string{
	return file.Hash
}

func (file *File) computeHash(filePath string) string {
	input := strings.NewReader(filePath)

	hash := sha256.New()
	if _, err := io.Copy(hash, input); err != nil {
		log.Fatal(err)
	}
	sum := hash.Sum(nil)

	return hex.EncodeToString(sum)

}
