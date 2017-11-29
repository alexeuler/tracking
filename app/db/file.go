package db

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/up-finder/silk.web/app"
	"github.com/up-finder/silk.web/app/setup"
	"io"
	"os"
	"path"
	"time"
)

// Singleton file database for the package
var File *FileDB = NewFileDB(app.Env)

// Initialization of the singleton file database
func NewFileDB(env *setup.Env) *FileDB {
	dir := path.Join(env.Root, env.DB.File.Path)
	exportDir := path.Join(env.Root, env.DB.File.ExportPath)
	capacity := env.DB.File.NumberOfFiles
	retries := env.DB.File.NumberOfRetries
	retryDelay := env.DB.File.RetryDelay
	return NewCustomFileDB(dir, exportDir, capacity, retries, retryDelay)
}

// Interface that could be saved to file db
type Serializer interface {
	Serialize() string
}

//   File database class
//   The files are stored in folder, specified by the path field.
//   The db writes concurrently to many files (number of files are specified by the capacity field)
//   If write to file fails, the db tries to reopen the file (number of retries is specified
//   by the retries field)
//   The delay between retries is specified by the retryDelay field (measured in milleseconds)
type FileDB struct {
	path       string
	exportPath string
	capacity   int
	retries    int
	retryDelay int
	files      chan *file
}

// ---------- Public FileDB methods

// File db constructor
func NewCustomFileDB(path string, exportPath string, capacity int, retries int, retryDelay int) (res *FileDB) {
	res = &FileDB{path: path, exportPath: exportPath, capacity: capacity, retries: retries, retryDelay: retryDelay}
	res.files = make(chan *file, capacity)
	res.fillChannel()
	return res
}

// Path to db files
func (f *FileDB) Path() string {
	return f.path
}

// Number of files in db for concurrent writes
func (f *FileDB) Capacity() int {
	return f.capacity
}

// Number of retries if write fails
func (f *FileDB) Retries() int {
	return f.retries
}

// The delay between retries (milliseconds)
func (f *FileDB) RetryDelay() int {
	return f.retryDelay
}

// Async save entity implementing serialize method to file db
func (f *FileDB) Save(s Serializer) {
	go f.save(s)
}

// Exports all files to one in the exportPath folder, returns the path of the exported file
func (f *FileDB) Export() string {
	// emptying the channel, blocking any writes
	files := make([]*file, 0, f.capacity)
	for i := 0; i < f.capacity; i++ {
		fl := <-f.files
		fl.stream.Close()
		files = append(files, fl)
	}

	// refilling the channel with new files after export is done
	defer func() {
		f.fillChannel()
	}()

	// merging all files into one in export path
	to := f.filenameForExport()
	for _, fromFile := range files {
		from := fromFile.filename
		err := appendToFile(to, from)
		if err != nil {
			log.Errorf("FileDB: Unable to append from file %s to %s: %v", from, to, err)
			return "" //stop export on fail
		}
		os.Remove(from) // remove file if append is successful
	}
	return to
}

// ---------- Private FileDB methods

// Generates filename with timestamp in export dir
func (f *FileDB) filenameForExport() string {
	t := time.Now()
	return path.Join(f.exportPath, t.Format("2006-01-02_15-04-05.log"))
}

// Filename for the i-th file in channel
func (f *FileDB) filename(i int) string {
	name := fmt.Sprintf("%d.log", i)
	return path.Join(f.path, name)
}

// Appends the content of fromPath file to toPath file
func appendToFile(toPath, fromPath string) (err error) {
	const CHUNK_SIZE = 10 * 1024 * 1024

	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(toPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	buf := make([]byte, CHUNK_SIZE)
	for {
		n, err := from.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := to.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

// Fills the FileDB channel with open files
func (f *FileDB) fillChannel() {
	for i := 0; i < cap(f.files); i++ {
		d, err := newFile(f.filename(i), f.retries, f.retryDelay)
		if err != nil {
			log.Fatalf("FileDB: Unable to init channel: %v", err)
		}
		f.files <- d
	}
}

// Sync save
func (f *FileDB) save(s Serializer) (err error) {
	data := s.Serialize()
	bytes := []byte(data)
	bytes = append(bytes, 10) // adding \n at the end of line
	file := <-f.files
	err = file.write(bytes, f.retries)
	f.files <- file
	if err != nil {
		log.Errorf("File DB: Failed to write data %s to file %s. Dropping data.", data, file.filename)
	}
	return err
}

// Structure for open files in FileDB, don't mess it with os.File

type file struct {
	stream     io.WriteCloser //stream for writing data
	opener     openable       //type for opening files, used as an interface to stub in tests
	filename   string         //path to file
	retries    int            //number of retries if write fails
	retryDelay int            //delay between retries in milliseconds
}

// Interface for opening a file
type openable interface {
	open(name string) (io.WriteCloser, error)
}

// Constructor for the new file
func newFile(filename string, retries int, retryDelay int) (f *file, err error) {
	f = &file{filename: filename, retries: retries, retryDelay: retryDelay, opener: &fileOpener{}}
	err = f.open()
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Method for opening a file, using the opener field
func (f *file) open() (err error) {
	f.stream, err = f.opener.open(f.filename)
	return err
}

// Close current file
func (f *file) close() (err error) {
	return f.stream.Close()
}

// Write string to file with retries in case of fail
func (f *file) write(data []byte, retries int) (err error) {
	_, err = f.stream.Write(data)
	if err != nil {
		log.Errorf("File DB: Error writing data %s to file %s: %v", data, f.filename, err)
		for retries > 0 {
			time.Sleep(time.Millisecond * time.Duration(f.retryDelay))
			retries--
			log.Errorf("File DB: Retry %d: Trying to reopen the file %s", retries, f.filename)
			err = f.close()
			if err != nil {
				continue
			}
			err = f.open()
			if err != nil {
				continue
			}
			return f.write(data, retries)
		}
	}
	return err
}

// Default implementation for openable interface
type fileOpener struct{}

//Opens or appends a file
func (f *fileOpener) open(name string) (io.WriteCloser, error) {
	stream, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("Could not open file %s to append or write with permissions 666: %s", name, err)
	}
	return stream, nil
}
