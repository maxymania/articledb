package articledb

import (
	"sync"
	"fmt"
	"encoding/json"
	"path/filepath"
	"github.com/maxymania/articledb/godbm"
)

const PostingPermitted = "y"
const PostingNotPermitted = "n"
const PostingModerated = "m"

type Group struct{
	Name string
	Description string
	Begin int64
	End int64
	Offset int64
	Count int64
	Posting bool
}

type GroupDB struct{
	glm sync.Mutex
	GroupList *godbm.HashDB
	GroupAssoc *godbm.HashDB
}

func NewGroupDB (path string) (*GroupDB,error){
	l := filepath.Join(path,"group.list")
	a := filepath.Join(path,"group.assoc")
	dbl,e := godbm.Open(l)
	if e!=nil{
		dbl,e = godbm.Create(l,10)
		if e!=nil { return nil,e }
	}
	dba,e := godbm.Open(a)
	if e!=nil{
		dba,e = godbm.Create(a,10)
		if e!=nil { dbl.Close(); return nil,e }
	}
	return &GroupDB{GroupList:dbl,GroupAssoc:dba},nil
}

// Adds a Group inside a GroupDB. It returns false,nil if the group was already existing,
// it returns true,err where err!=nil if creating the group did not work and true,nil otherwise.
func(gdb *GroupDB) CreateGroup(grp *Group) (bool,error){
	gdb.glm.Lock(); defer gdb.glm.Unlock()
	grp.Begin=0
	grp.End=0
	grp.Offset=1000
	grp.Count=0
	t,e := gdb.GroupList.Get([]byte(grp.Name))
	// fmt.Println(t)
	if e==nil && len(t)>0{ return false,nil }
	j,e := json.Marshal(grp)
	// fmt.Println("debug:",string(j))
	if e!=nil{ return true,e }
	gdb.GroupList.Set([]byte(grp.Name),j)
	return true,nil
}

// this is for Test-Purposes, use CreateGroup instead
func(gdb *GroupDB) CreateGroupStr(grp string) (bool,error){
	return gdb.CreateGroup(&Group{Name:grp})
}

func(gdb *GroupDB) getNumberInGroup(group string) (int64,error){
	var grp Group
	gdb.glm.Lock(); defer gdb.glm.Unlock()
	j,e := gdb.GroupList.Get([]byte(group))
	if e!=nil{ return 0,e }
	e = json.Unmarshal(j,&grp)
	if e!=nil{ return 0,e }
	grp.End++
	if(grp.End>=grp.Offset){ grp.Offset+=1000 }
	j,e = json.Marshal(&grp)
	if e!=nil{ return 0,e }
	gdb.GroupList.Set([]byte(group),j)
	return grp.End,nil
}

// Iterate over all Groups
// 
// isok is a pointer to a boolean, with wich the consumer can signalize
// to not continue iterating. If isok is nil it is treated as always true.
func(gdb *GroupDB) IterateGroups(isok *bool, grpc chan <- *Group) {
	defer close(grpc)
	kvc := make(chan godbm.KeyValuePair)
	go gdb.GroupList.Iterate(isok,kvc)
	for kvp := range kvc {
		grp := new(Group)
		if json.Unmarshal(kvp.Value,grp)!=nil { continue }
		grpc <- grp
	}
}

// Adds an article to the Group
func(gdb *GroupDB) AddArticle(id ,group string) error {
	n,e := gdb.getNumberInGroup(group)
	if e!=nil { return e }
	e = gdb.GroupAssoc.Set([]byte(fmt.Sprintf("%s %d",group,n)),[]byte(id))
	if e!=nil { return e }
	return nil
}

// Gets a Internal Number of a Number from a group
func(gdb *GroupDB) GetNumber(group string,rn int64) (int64,error){
	var grp Group
	j,e := gdb.GroupList.Get([]byte(group))
	if e!=nil{ return 0,e }
	e = json.Unmarshal(j,&grp)
	if e!=nil{ return 0,e }
	return grp.Offset-rn,nil
}

// Gets a Article-ID Number of a Number from a group
func(gdb *GroupDB) GetID(group string,rn int64) (string,error){
	i,e := gdb.GetNumber(group,rn)
	if e!=nil { return "",e }
	id,e := gdb.GroupAssoc.Get([]byte(fmt.Sprintf("%s %d",group,i)))
	if e!=nil { return "",e }
	return string(id),nil
}

// Gets an Group object. This is for GetIDFast
func(gdb *GroupDB) GetGroupObj(group string) (*Group,error){
	grp := new(Group)
	j,e := gdb.GroupList.Get([]byte(group))
	if e!=nil{ return nil,e }
	e = json.Unmarshal(j,grp)
	if e!=nil{ return nil,e }
	return grp,nil
}

// this is like GetID, but it takes an Group object instead of a group name, so it is more
// efficient if many GetID/GetIDFast calls are done.
func(gdb *GroupDB) GetIDFast(group *Group,rn int64) (string,error){
	i := group.Offset-rn
	id,e := gdb.GroupAssoc.Get([]byte(fmt.Sprintf("%s %d",group,i)))
	if e!=nil { return "",e }
	return string(id),nil
}

// Gets the low and high value for a group (Defined in NNTP)
// Note that this is the external representation
func(gdb *GroupDB) GetRange(group string) (low int64,high int64,err error){
	var grp Group
	j,e := gdb.GroupList.Get([]byte(group))
	if e!=nil{ return 0,0,e }
	e = json.Unmarshal(j,&grp)
	if e!=nil{ return 0,0,e }
	low = grp.Offset-grp.End
	high = grp.Offset-grp.Begin
	err = nil
	return
}

