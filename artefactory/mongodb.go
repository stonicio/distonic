package artefactory

import "gopkg.in/mgo.v2"

type MongoDBArtefactory struct {
	session *mgo.Session
}

func NewMongoDBArtefactory(url string) (*MongoDBArtefactory, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &MongoDBArtefactory{session: session}, nil
}
