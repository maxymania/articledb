package articledb

import (
	"fmt"
	"io"
	"os"
	"sync"
	"log"
	"compress/flate"
	"path/filepath"
)

const (
    NoCompression      = flate.NoCompression
    BestSpeed          = flate.BestSpeed
    BestCompression    = flate.BestCompression
    DefaultCompression = flate.DefaultCompression
)

// A position Descriptor stores the filename, the 
type Position struct{
	Name string
	Begin int64
	Length int64
}

type fileWriter struct{
	Name string
	File *os.File
	Offset int64
}

func(f *fileWriter) WriteStream(s io.Reader,level int) (*Position,error) {
	begin,e := f.File.Seek(f.Offset,os.SEEK_SET)
	if e!=nil { return nil,e }
	w,e := flate.NewWriter(f.File,level)
	if e!=nil { return nil,e }
	_,e = io.Copy(w,s)
	w.Close()
	if e!=nil {
		f.File.Seek(f.Offset,os.SEEK_SET)
		return nil,e
	}
	end,e := f.File.Seek(0,os.SEEK_CUR)
	if e!=nil { return nil,e }
	f.Offset = end
	return &Position{f.Name,begin,end-begin},nil
}
// func(f *fileWriter) 


// Do not modify any Field at all
type StorageWriter struct{
	cincr sync.Mutex
	Counter uint
	Path string
	Queue chan *fileWriter
}

// creates a new Storage writer for a specific folder
func NewStorageWriter(path string,size int) *StorageWriter{
	return &StorageWriter{
		Counter : 0,
		Path : path,
		Queue : make(chan *fileWriter,size),
	}
}


func(s *StorageWriter) createFileWriter() (*fileWriter,error){
	s.cincr.Lock()
		c := s.Counter
		s.Counter++
	s.cincr.Unlock()
	fn := fmt.Sprintf("data._%d_",c)
	f,e := os.OpenFile(filepath.Join(s.Path,fn),os.O_RDWR|os.O_CREATE,0644)
	if e!=nil { return nil,e }
	l,e := f.Seek(0,os.SEEK_END)
	if e!=nil { return nil,e }
	return &fileWriter{fn,f,l},nil
}
func(s *StorageWriter) getFileWriter() (*fileWriter,error) {
	select {
	case f := <- s.Queue:
		return f,nil
	default:
		return s.createFileWriter()
	}
	panic("unreachable")
}
func(s *StorageWriter) putFileWriter(f *fileWriter) {
	select {
	case s.Queue <- f:return
	default:
		f.File.Close()
		log.Println("queue overflow:",f.Name) // Queue full
	}
}

// writes an article into the Storage and returns the Position descriptor
func(s *StorageWriter) WriteStream(src io.Reader,level int) (*Position,error) {
	f,e := s.getFileWriter()
	if e!=nil { return nil,e }
	defer s.putFileWriter(f)
	return f.WriteStream(src,level)
}

// Do not modify any Field at all
type StorageReader struct{
	mu sync.Mutex
	Path string
	Fobjs map[string]*os.File
}
// creates an reader for a Storage Folder
func NewStorageReader(path string) *StorageReader{
	return &StorageReader{
		Path : path,
		Fobjs : make(map[string]*os.File),
	}
}

// reads a stored thing at the given Position
func (s *StorageReader) Read(pos *Position) (io.ReadCloser,error) {
	s.mu.Lock(); defer s.mu.Unlock()
	fobj,ok := s.Fobjs[pos.Name]
	if !ok {
		f,e := os.Open(filepath.Join(s.Path,pos.Name))
		if e!=nil { return nil,e }
		s.Fobjs[pos.Name]=f
		fobj=f
	}
	r := io.NewSectionReader(fobj,pos.Begin,pos.Length)
	r2 := flate.NewReader(r)
	return r2,nil
}
