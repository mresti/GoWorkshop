package config

import (
	pb "github.com/wizelineacademy/GoWorkshop/proto/list"
	"gopkg.in/mgo.v2"
)

// Struct used for maintaining HTTP Request Context
type Context struct {
	MongoSession *mgo.Session
	ListService  pb.ListClient
}

// Close mgo.Session
func (c *Context) Close() {
	c.MongoSession.Close()
}

// Returns mgo.collection for the given name
func (c *Context) DbCollection(name string) *mgo.Collection {
	return c.MongoSession.DB(AppConfig.Database).C(name)
}

// Create a new Context object for each HTTP request
func NewContext() *Context {
	session := GetSession().Copy()
	context := &Context{
		MongoSession: session,
		ListService:  GetListService(),
	}
	return context
}
