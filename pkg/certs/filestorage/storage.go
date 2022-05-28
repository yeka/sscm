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
	maxId  int // keep the maximum number of ID, for autoincrement
	data   []certs.Data
	maps   map[int]int   // maps of id to index of FS.data
	childs map[int][]int // maps parent id to its childs' id
}

func New(filename string) (*FS, error) {
	fs := &FS{filename: filename}
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = gob.NewDecoder(bytes.NewReader(b)).Decode(&fs.dataFile)
	return fs, err
}

func (fs *FS) Store(cert *certs.Data) (err error) {
	fs.lock.Lock()

	fs.maxId++
	cert.ID = fs.maxId
	fs.data = append(fs.data, *cert)
	fs.maps[cert.ID] = len(fs.data) - 1
	if cert.ParentID > 0 {
		fs.childs[cert.ParentID] = append(fs.childs[cert.ParentID], cert.ID)
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

func (fs *FS) Load(id int, isRoot bool) (data certs.Data, err error) {
	if (isRoot && fs.data[fs.maps[id]].ParentID > 0) || (!isRoot && fs.data[fs.maps[id]].ParentID == 0) {
		err = errors.New("data not found")
		return
	}

	data = fs.data[fs.maps[id]]
	return
}

func (fs *FS) Search(query string, mode certs.SearchMode) ([]certs.Data, error) {
	res := make([]certs.Data, 0)
	for _, v := range fs.data {
		if v.ParentID > 0 && mode == certs.RootOnly {
			continue
		}
		if v.ParentID == 0 && mode == certs.NonRootOnly {
			continue
		}
		// TODO: filter query string
		res = append(res, v)
	}

	return res, nil
}
