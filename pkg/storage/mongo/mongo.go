package mongo

import (
	"context"
	"log"
	"skillfactory/GoNews/pkg/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	db *mongo.Database
}

func New(connstr string) (*Storage, error) {
	cl, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connstr))
	if err != nil {
		return nil, err
	}

	db := cl.Database("posts")

	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Posts() ([]storage.Post, error) {
	ps := make([]storage.Post, 0)
	log.Println("mongo.Posts")
	rows, err := s.db.Collection("posts").Find(context.Background(), bson.M{})
	if err != nil {
		return ps, err
	}
	defer rows.Close(context.Background())

	if err = rows.All(context.Background(), &ps); err != nil {
		return ps, err
	}
	// выбираем наменование автора для статьи
	for k, v := range ps {
		row, err := s.db.Collection("authors").Find(context.Background(), bson.M{"id": v.AuthorID})
		if err == nil {
			a := storage.Author{}
			row.Next(context.Background())
			err = row.Decode(&a)
			if err == nil {
				ps[k].AuthorName = a.Name
			}
		}
	}
	return ps, nil
}

func (s *Storage) MaxID() (int, error) {
	o := options.FindOne().SetSort(bson.M{"id": -1})
	row := s.db.Collection("posts").FindOne(context.Background(), bson.M{}, o)
	p := storage.Post{}
	row.Decode(&p)
	return 0, nil
}

func (s *Storage) AddPost(p storage.Post) error {
	maxid, err := s.MaxID()
	if err == nil {
		maxid++
		p.ID = maxid
	}
	b, err := bson.Marshal(p)
	if err != nil {
		return err
	}

	_, err = s.db.Collection("posts").InsertOne(context.Background(), b)
	return err
}

func (s *Storage) UpdatePost(p storage.Post) error {
	_, err := bson.Marshal(p)
	if err != nil {
		return err
	}
	f := bson.M{"id": p.ID}
	_, err = s.db.Collection("posts").UpdateOne(context.Background(), f, bson.M{"$set": p})
	log.Println(err)
	return err
}

func (s *Storage) DeletePost(p storage.Post) error {
	f := bson.M{"id": p.ID}
	_, err := s.db.Collection("posts").DeleteOne(context.Background(), f)
	return err
}
