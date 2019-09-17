package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Divvy/src/pb"
	"github.com/radovskyb/watcher"
)

type File struct {
	FileName string
	Path     string
	Size     int64
	IsDir    bool
	Hash     string
}

type FileManager struct {
	DirectoryPath string
	SharedFiles   []*File
}

func NewFileManager(DirectoryPath string) *FileManager {
	fileMgr := &FileManager{}
	fileMgr.DirectoryPath = DirectoryPath

	fileMgr.SharedFiles = []*File{}

	//Read Directory
	files, err := ioutil.ReadDir(fileMgr.DirectoryPath)
	if err != nil {
		log.Fatal(err)
	}

	//Go through all files and add them to struct
	for _, f := range files {
		log.Println(f.Name())

		pathToFile := DirectoryPath + "/" + f.Name()
		fi, err := os.Stat(pathToFile)

		if err != nil {
			log.Fatal(err)
		}

		f := &File{FileName: f.Name(),
			Path:  fileMgr.DirectoryPath,
			IsDir: false,
			Size:  fi.Size(),
		}
		f.setHash()

		fileMgr.SharedFiles = append(fileMgr.SharedFiles, f)
	}

	go fileMgr.createListener(fileMgr.DirectoryPath)
	return fileMgr
}

func (fileMgr *FileManager) createListener(DirectoryPath string) {
	// creates a new file watcher
	w := watcher.New()
	//w.SetMaxEvents(1)

	w.FilterOps(watcher.Remove, watcher.Rename, watcher.Move, watcher.Create, watcher.Write)

	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Println(event.String()) // Print the event's info.
				s := strings.Split(event.String(), " ")
				fileMgr.HandleEvent(s)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()
	if err := w.AddRecursive(DirectoryPath); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func (fileMgr *FileManager) HandleEvent(s []string) {
	if s[0] == "FILE" {
		switch op := s[2]; op {
		case "RENAME":
			for _, f := range fileMgr.SharedFiles {
				if f.FileName == s[1][1:len(s[1])-1] {
					f.FileName = filepath.Base(s[5][:len(s[5])-1])
				}
			}
		case "REMOVE":
			i := 0
			log.Println(s[1][1 : len(s[1])-1])
			if filepath.Ext(s[1][1:len(s[1])-1]) != ".swp" {
				for _, f := range fileMgr.SharedFiles {
					if f.FileName != s[1][1:len(s[1])-1] {
						i++
					}
					break
				}
				log.Println(i)
				i--
				a := fileMgr.SharedFiles
				a = append(a[:i], a[i+1:]...)
				//copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index
				//a[len(a)-1] = nil     // Erase last element (write zero value)
				fileMgr.SharedFiles = a
				//[:len(a)-1]     // Truncate slice
			}
		case "WRITE":
			//i := 0
			log.Println(s[1][1 : len(s[1])-1])
			if filepath.Ext(s[1][1:len(s[1])-1]) != ".swp" {
				for _, f := range fileMgr.SharedFiles {
					if f.FileName == s[1][1:len(s[1])-1] {
						log.Println("recompute hash")
						f.setHash()
						break
					}
				}
			}
		case "CREATE":
			//i := 0
			log.Println(s[1][1 : len(s[1])-1])
			if filepath.Ext(s[1][1:len(s[1])-1]) != ".swp" {
				f := &File{FileName: s[1][1 : len(s[1])-1],
					Path:  fileMgr.DirectoryPath,
					IsDir: false,
				}
				f.setHash()

				fileMgr.SharedFiles = append(fileMgr.SharedFiles, f)
			}
		}
	}
	fileMgr.displayDirectory()

}

//check if file exists
func (fileMgr *FileManager) searchFileByName(name string) *File {
	for _, f := range fileMgr.SharedFiles {
		if name == f.FileName {
			return f
		}
	}
	return nil
}

func (fileMgr *FileManager) searchFileByHash(hash string) *File {
	for _, f := range fileMgr.SharedFiles {
		if hash == f.Hash {
			return f
		}
	}
	return nil
}

func (fileMgr *FileManager) displayDirectory() {

	for _, f := range fileMgr.SharedFiles {
		log.Print(f.FileName)
	}
}

func (fileMgr *FileManager) GetSharedFilesList() pb.FileList {
	fileList := pb.FileList{}
	fileList.NodeID = Node.netMgr.ID.String()
	for _, file := range fileMgr.SharedFiles {
		fileList.Files = append(fileList.Files, &pb.File{
			Name: file.FileName,
			Hash: file.Hash})
	}
	return fileList
}

func (file File) String() string {
	return fmt.Sprintf("%v (%v)", file.FileName, file.Hash)
}

func (file *File) setHash() {
	file.Hash = file.computeHash(file.Path + file.FileName)
}

func (file *File) getHash(filePath string) string {
	if len(file.Hash) == 0 {
		file.Hash = file.computeHash(filePath)
	}

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
