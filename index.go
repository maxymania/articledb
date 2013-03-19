package articledb

import (
	"sync"
	"encoding/json"
	"github.com/maxymania/articledb/godbm"
	"path/filepath"
	"net/textproto"
	"os"
)

type Entry struct{
	ArticleID string
	Address Position
	Header textproto.MIMEHeader
}

type Index struct{
	jem sync.Mutex
	encf *os.File
	Enc *json.Encoder
	Dbm *godbm.HashDB
}
func NewIndex(path string) (*Index,error) {
	j := filepath.Join(path,"index.json")
	d := filepath.Join(path,"index.dbm")
	db,e := godbm.Open(d)
	if e!=nil {
		db,e = godbm.Create(d,10)
		if e!=nil { return nil,e }
		jr,e := os.Open(j)
		if e==nil{
			ent := new(Entry)
			jd := json.NewDecoder(jr)
			for {
				if jd.Decode(ent)!=nil { break }
				jm,e := json.Marshal(&(ent.Address))
				if e!=nil { continue }
				db.Set([]byte(ent.ArticleID),jm)
			}
			jr.Close()
		}
	}
	jw,e := os.OpenFile(j,os.O_WRONLY|os.O_APPEND|os.O_CREATE,0644)
	if e!=nil {
		db.Close()
		return nil,e
	}
	return &Index{encf:jw,Enc:json.NewEncoder(jw),Dbm:db},nil
}

// sets a Tuple of Message-Id, Position descriptor and MIME-Header in the index.
func (i *Index) Put(id string,p *Position,h textproto.MIMEHeader) error {
	var ent Entry
	ent.ArticleID=id
	ent.Address = *p
	ent.Header = h
	j,e := json.Marshal(&ent)
	if e!=nil { return e }
	i.Dbm.Set([]byte(id),j)
	
	i.jem.Lock()
	i.Enc.Encode(&ent)
	i.jem.Unlock()
	return nil
}

// gets a Tuple of Position descriptor and MIME-Header from a specific Message-Id from the index.
func (i *Index) Get(id string) (*Position,textproto.MIMEHeader,error) {
	var ent Entry
	j,e := i.Dbm.Get([]byte(id))
	if e!=nil { return nil,nil,e }
	e = json.Unmarshal(j,&ent)
	if e!=nil { return nil,nil,e }
	p := new(Position)
	*p = ent.Address
	h := ent.Header
	return p,h,nil
}

