package filestorage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"sscm/pkg/certs"
	"sync"
)

type FS struct {
	dataFile
	lock     sync.Mutex
	filename string
}

type dataFile struct {
	MaxId  int // keep the maximum number of ID, for autoincrement
	Data   []certs.Data
	Maps   map[int]int   // maps of id to index of FS.data
	Childs map[int][]int // maps parent id to its childs' id
}

func New(filename string) (*FS, error) {
	fs := &FS{filename: filename}
	b, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		fs.dataFile.Maps = make(map[int]int)
		fs.dataFile.Childs = make(map[int][]int)
		return fs, nil
	}
	if err != nil {
		return nil, err
	}
	err = gob.NewDecoder(bytes.NewReader(b)).Decode(&fs.dataFile)
	return fs, err
}

func (fs *FS) Store(cert *certs.Data) (err error) {
	fs.lock.Lock()

	fs.MaxId++
	cert.ID = fs.MaxId
	fs.Data = append(fs.Data, *cert)
	fs.Maps[cert.ID] = len(fs.Data) - 1
	if cert.ParentID > 0 {
		fs.Childs[cert.ParentID] = append(fs.Childs[cert.ParentID], cert.ID)
	}

	fs.lock.Unlock()

	// Sync to file
	f, err := os.Create(fs.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	err = gob.NewEncoder(f).Encode(fs.dataFile)
	return err
}

func (fs *FS) Load(id int) (data certs.Data, err error) {
	i, ok := fs.Maps[id]
	if !ok {
		err = errors.New("data not found")
		return
	}
	if i >= len(fs.Data) {
		err = errors.New("error index out of scope")
		return
	}

	data = fs.Data[i]
	return
}

func (fs *FS) Search(query string, parentId int) ([]certs.Data, error) {
	res := make([]certs.Data, 0)
	for _, v := range fs.Data {
		if (parentId >= 0) && v.ParentID != parentId {
			continue
		}
		// TODO: filter query string
		res = append(res, v)
	}

	return res, nil
}
