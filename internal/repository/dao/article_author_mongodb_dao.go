package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/wx-up/go-book/internal/repository/dao/model"
)

const (
	DatabaseName = "go_book"
	// TblNameArticle 制作库表名
	TblNameArticle = "articles"

	// TblNamePublishedArticle 发布库表名
	TblNamePublishedArticle = "published_articles"
)

type MongoDBArticleDAO struct {
	client *mongo.Client
	// node     *snowflake.Node
	idGen IDGen
}

type IDGen func() int64

func NewMongoDBArticleDAO(client *mongo.Client, gen IDGen) *MongoDBArticleDAO {
	return &MongoDBArticleDAO{
		client: client,
		idGen:  gen,
	}
}

// InitArticlesCollection 预先创建好索引
func InitArticlesCollection(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	index := []mongo.IndexModel{
		{
			Keys:    bson.D{{"id", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{"author_id", 1}, {"create_time", 1}},
			Options: options.Index(),
		},
	}
	_, err := db.Collection("articles").Indexes().CreateMany(ctx, index)
	if err != nil {
		return err
	}
	_, err = db.Collection("published_articles").Indexes().CreateMany(ctx, index)
	return err
}

func (m *MongoDBArticleDAO) Insert(ctx context.Context, article model.Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.CreateTime = now
	article.UpdateTime = now
	id := m.idGen()
	article.Id = id
	_, err := m.client.Database(DatabaseName).Collection(TblNameArticle).InsertOne(ctx, article)
	return id, err
}

func (m *MongoDBArticleDAO) UpdateById(ctx context.Context, article model.Article) error {
	filter := bson.M{"id": article.Id, "author_id": article.AuthorId}
	res, err := m.client.Database(DatabaseName).Collection(TblNameArticle).UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		"title":       article.Title,
		"content":     article.Content,
		"update_time": time.Now().UnixMilli(),
	}})
	if err != nil {
		return err
	}
	if res.ModifiedCount <= 0 {
		return ErrArticleNotFound
	}
	return nil
}

func (m *MongoDBArticleDAO) Sync(ctx context.Context, article model.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)
	session, err := m.client.StartSession()
	if err != nil {
		return 0, err
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		if id > 0 {
			err = m.UpdateById(ctx, article)
		} else {
			id, err = m.Insert(ctx, article)
		}
		if err != nil {
			return nil, err
		}

		// 同步到线上库,upsert 语义
		now := time.Now().UnixMilli()
		article.Id = id
		article.CreateTime = now
		article.UpdateTime = now
		filter := bson.M{"id": id}
		return m.client.Database(DatabaseName).
			Collection(TblNamePublishedArticle).
			UpdateOne(ctx, filter, bson.M{
				"$set": bson.M{
					"id":          id,
					"title":       article.Title,
					"content":     article.Content,
					"author_id":   article.AuthorId,
					"update_time": article.UpdateTime,
				},
				// $setOnInsert 表示当插入时需要插入 create_time 字段
				"$setOnInsert": bson.M{
					"create_time": article.CreateTime,
				},
			}, options.Update().SetUpsert(true))
	}, txnOptions)
	return id, err
}

func (m *MongoDBArticleDAO) Transaction(ctx context.Context, f func(dao ArticleDAO, readerDao ReaderArticleDAO) error) error {
	// TODO implement me
	panic("implement me")
}

func (m *MongoDBArticleDAO) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	// TODO implement me
	panic("implement me")
}
