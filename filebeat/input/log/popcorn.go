package log

import (
"bytes"
"io"
"os"
"sync"
)

// define the file structure
type Popcorn struct {
sync.Mutex
path string
file *os.File
}

func (f *Popcorn) Initial(path string) {
f.path = path
}

func (f *Popcorn) close() error {
return f.file.Close()
}

func (f *Popcorn) open() error {
file, e := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE, 0666)
if e != nil {
return e
} else {
f.file = file
}
return nil
}

func (f *Popcorn) PopLine(size int64) ([]byte, error) {
// lock the critical section
f.Lock()
// unlock after checking out the data
defer f.Unlock()

// read the file
err := f.open()
if err != nil {
	return nil, err
}

buf := bytes.NewBuffer(make([]byte, 0, size))

_, err = f.file.Seek(0, os.SEEK_SET)
if err != nil {
	return nil, err
}
_, err = io.Copy(buf, f.file)
if err != nil {
	return nil, err
}
line, err := buf.ReadString('\n')
if err != nil && err != io.EOF {
	return nil, err
}

_, err = f.file.Seek(0, os.SEEK_SET)
if err != nil {
	return nil, err
}
nw, err := io.Copy(f.file, buf)
if err != nil {
	return nil, err
}
err = f.file.Truncate(nw)
if err != nil {
	return nil, err
}
err = f.file.Sync()
if err != nil {
	return nil, err
}

_, err = f.file.Seek(0, os.SEEK_SET)
if err != nil {
	return nil, err
}
err = f.close()
if err != nil {
	return nil, err
}
return []byte(line), nil
}
