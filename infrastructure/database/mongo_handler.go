package database

import (
	"context"
	mongo "gopkg.in/mgo.v2"
)

//mongoHandler armazena a estrutura para MongoDB
type mongoHandler struct {
	database *mongo.Database
	session  *mongo.Session
}

//NewMongoHandler constrói um novo handler de banco para MongoDB
func NewMongoHandler(c *config) (*mongoHandler, error) {
	session, err := mongo.DialWithTimeout(c.host, c.connTimeout)
	if err != nil {
		return &mongoHandler{}, err
	}

	handler := new(mongoHandler)
	handler.session = session
	handler.database = handler.session.DB(c.database)

	return handler, nil
}

//Store realiza uma inserção no banco de dados
func (mgo mongoHandler) Store(ctx context.Context, collection string, data interface{}) error {
	session := mgo.session.Clone()
	defer session.Close()

	return mgo.database.C(collection).With(session).Insert(data)
}

//Update realiza uma atualização no banco de dados
func (mgo mongoHandler) Update(ctx context.Context, collection string, query interface{}, update interface{}) error {
	session := mgo.session.Clone()
	defer session.Close()

	return mgo.database.C(collection).With(session).Update(query, update)
}

//FindAll realiza uma busca por todos os registros no banco de dados
func (mgo mongoHandler) FindAll(ctx context.Context, collection string, query interface{}, result interface{}) error {
	session := mgo.session.Clone()
	defer session.Close()

	return mgo.database.C(collection).With(session).Find(query).All(result)
}

//FindOne realiza a busca de um item específico no banco de dados
func (mgo mongoHandler) FindOne(
	ctx context.Context,
	collection string,
	query interface{},
	selector interface{},
	result interface{},
) error {
	session := mgo.session.Clone()
	defer session.Close()

	return mgo.database.C(collection).With(session).Find(query).Select(selector).One(result)
}
