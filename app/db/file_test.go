package db

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

var testDir = os.TempDir()

// Save string "1" to file
func TestSyncWrite(t *testing.T) {
	fileStub := newFileDBStub(1, 0)
	fileStub.save(unit("1"))
	file, _ := os.Open(path.Join(testDir, "0.log"))
	bytes, _ := ioutil.ReadAll(file)
	content := string(bytes)
	if content != "1\n" {
		t.Fatalf("File DB: TestSyncWrite: expected content of file to be %s, but got %s", "1", content)
	}
}

// Save string "1" 100 times and check that it saves to all 3 files in FileDB with capacity 3
func TestAsyncWrite(t *testing.T) {
	fileStub := newFileDBStub(3, 0)
	data := unit("1")
	for i := 0; i < 100; i++ {
		fileStub.Save(data)
	}
	time.Sleep(time.Millisecond * 10)
	for i := 0; i < 3; i++ {
		file, _ := os.Open(path.Join(testDir, fmt.Sprintf("%d.log", i)))
		bytes, _ := ioutil.ReadAll(file)
		if len(bytes) == 0 {
			t.Fatalf("File DB: TestAsyncWrite: Zero bytes in file %d", i)
		}
	}
}

// Save "1" to three files, export them and check that exported version has 10 lines of "1"
func TestExport(t *testing.T) {
	stub := newFileDBStub(3, 0)
	for i := 0; i < 10; i++ {
		stub.save(unit("1"))
	}
	exportPath := stub.Export()
	data, _ := ioutil.ReadFile(exportPath)
	str := string(data)
	if str != "1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n" {
		t.Errorf("File DB: TestAppendFile: expected 1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n, got %s", str)
	}
}

func TestAppendFile(t *testing.T) {
	f1Path := path.Join(testDir, "test1")
	f2Path := path.Join(testDir, "test2")
	exportPath := path.Join(testDir, "result")
	os.Remove(f1Path)
	os.Remove(f2Path)
	os.Remove(exportPath)
	ioutil.WriteFile(f1Path, []byte("123"), 0666)
	ioutil.WriteFile(f2Path, []byte("asd"), 0666)
	appendToFile(exportPath, f1Path)
	appendToFile(exportPath, f2Path)
	data, _ := ioutil.ReadFile(exportPath)
	str := string(data)
	if str != "123asd" {
		t.Errorf("File DB: TestAppendFile: expected 123asd, got %s", str)
	}
}

// Setup the stream that fails x times and FileDB that retries y times.
// Check that write fails if x > y and succeds otherwise
func TestFailedStream(t *testing.T) {
	// ignore error logs
	level := log.GetLevel()
	log.SetLevel(log.PanicLevel)
	defer func() { log.SetLevel(level) }()

	cases := []struct {
		Failures int
		Retries  int
	}{
		{2, 3},
		{3, 3},
		{4, 0},
		{0, 4},
	}

	for _, testCase := range cases {
		stub := newUnreliableFileDBStub(3, testCase.Retries, testCase.Failures)
		err := stub.save(unit("1"))
		if testCase.Retries-testCase.Failures >= 0 {
			if err != nil {
				t.Fatalf("File DB: FailedStream: received error with %d retries and %d failures",
					testCase.Retries, testCase.Failures)
			}
		} else {
			if err == nil {
				t.Fatalf("File DB: FailedStream: didn't receive error with %d retries and %d failures",
					testCase.Retries, testCase.Failures)
			}
		}
	}
}

// -------- Private functions

func newFileDBStub(capacity int, retries int) *FileDB {
	for i := 0; i < capacity; i++ {
		os.Remove(path.Join(testDir, fmt.Sprintf("%d.log", i)))
	}
	return NewCustomFileDB(testDir, testDir, capacity, retries, 0)
}

func newUnreliableFileDBStub(capacity, numberOfRetries, numberOfFails int) *FileDB {
	fileDB := newFileDBStub(capacity, numberOfRetries)
	for i := 0; i < capacity; i++ {
		file := <-fileDB.files
		unreliable := &unreliableWriteCloser{failCount: numberOfFails}
		file.stream = unreliable
		file.opener = unreliable
		fileDB.files <- file
	}
	return fileDB
}

// String implementing Serializer interface
type unit string

func (u unit) Serialize() string {
	return string(u)
}

// io.WriteCloser that fails predefined number of times.
// Implements open interface to stub open methods
type unreliableWriteCloser struct {
	failCount int
}

func (s *unreliableWriteCloser) Write(data []byte) (int, error) {
	if s.failCount > 0 {
		s.failCount--
		return 0, fmt.Errorf("Sample error")
	}
	return 1, nil
}

func (s *unreliableWriteCloser) Close() error {
	return nil
}

func (s *unreliableWriteCloser) open(name string) (io.WriteCloser, error) {
	return s, nil
}
