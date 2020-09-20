package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/mutehayyiz/cafetv-server/config"
	"github.com/mutehayyiz/cafetv-server/media"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type Handler interface {
	AddIfNotExists(m * media.Media) error

	GetByID(id string, m * media.Media) error
	GetByCategory(category string, medias *[]*media.Media) error
	GetAll(medias *[]*media.Media) error

	Update(id string, l *media.Media) error

	DeleteByID(id string) error
	DeleteAll() error
	DropDatabase() error
}

var MediaHandler Handler

func Connect() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	MongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Global.MongoURL))
	fatalf("Problem while connecting to Mongo: %s", err)
	MediaHandler = mediaMongoHandler{MongoClient.Database(config.Global.DBName).Collection("media")}
}

func fatalf(format string, err error) {
	if err != nil {
		logrus.Fatalf(format, err)
	}
}

type mediaMongoHandler struct {
	col *mongo.Collection
}


func (m mediaMongoHandler) AddIfNotExists(l *media.Media)  error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	filter := bson.M{"hash": l.Hash }
	res := m.col.FindOne(ctx, filter)
	err := res.Err()
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
	} else {
		var existingMedia media.Media
		_ = res.Decode(&existingMedia)

		return  errors.New(fmt.Sprintf("there is already such media with ID: %s", existingMedia.ID.Hex()))
	}


	update := bson.M{"$set": l}
	_, err = m.col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return errors.New(fmt.Sprintf("error while inserting license: %s", err))
	}

	return nil
}

func (m mediaMongoHandler) GetByID(id string, l *media.Media) error {
	licenseID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(fmt.Sprintf("ID format error: %s", err))
	}

	filter := bson.M{"_id": licenseID}
	res := m.col.FindOne(context.Background(), filter)
	err = res.Err()
	if err != nil {
		return err
	}
	_ = res.Decode(l)

	return nil
}

func (m mediaMongoHandler) GetByCategory(category string, medias *[]*media.Media) error {
	cur, err := m.col.Find(context.Background(), bson.M{"category": media.Category(category)})
	if err != nil {
		return err
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {

		var l media.Media
		err := cur.Decode(&l)
		if err != nil {
			return err
		}

		*medias = append(*medias, &l)

	}

	return cur.Err()
}

func (m mediaMongoHandler) GetAll(medias *[]*media.Media) error {
	cur, err := m.col.Find(context.Background(), bson.D{})
	if err != nil {
		return err
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {

		var l media.Media
		err := cur.Decode(&l)
		if err != nil {
			return err
		}

		*medias = append(*medias, &l)

	}

	return cur.Err()
}

func (m mediaMongoHandler) Update(id string, l *media.Media) error {
	licenseID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(fmt.Sprintf("ID format error: %s", err))
	}

	filter := bson.M{"_id": licenseID}
	res := m.col.FindOne(context.Background(), filter)
	err = res.Err()
	if err != nil {
		return err
	}
	_ = res.Decode(l)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	update := bson.M{"online": l.Online}
	_, err = m.col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return errors.New(fmt.Sprintf("error while inserting license: %s", err))
	}

	return nil
}


func (m mediaMongoHandler) DeleteByID(id string) error {
	licenseID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(fmt.Sprintf("ID format error: %s", err))
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.M{"_id": licenseID}
	res, err := m.col.DeleteOne(ctx, filter)
	if res == nil {
		return errors.New("there is no license to delete")

	}
	if res.DeletedCount == 0 {
		return errors.New(fmt.Sprintf("there is no license to delete"))
	}

	if err != nil {
		return errors.New("license cannot be deleted")
	}

	logrus.Info("License successfully deleted")

	return nil
}

func (m mediaMongoHandler) DeleteAll() error {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.M{}
	res, err := m.col.DeleteMany(ctx, filter)

	if res == nil {
		return errors.New("there is no license to delete")
	}

	if res.DeletedCount == 0 {
		return errors.New(fmt.Sprintf("there is no license to delete"))
	}
	if err != nil {
		return errors.New("media cannot be deleted")
	}

	logrus.Info("Medias successfully deleted")

	return nil
}

func (m mediaMongoHandler) DropDatabase() error {
	return m.col.Database().Drop(context.Background())
}
