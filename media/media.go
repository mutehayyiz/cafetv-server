package media

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func fatalf(format string, err error) {
	if err != nil {
		logrus.Fatalf(format, err)
	}
}

type Type string
const(
	video Type= "video"
	photo Type= "photo"
	music Type= "music"
)

type Category string
const(
	howToDo Type= "howtodo"
	add Type= "add"
	didYouKnow Type= "didyouknow"
)

type Media struct {
	ID        	primitive.ObjectID     	`bson:"_id,omitempty" json:"id"`
	Name 		string					`bson:"name,omitempty" json:"name"`
	Description string 					`bson:"description,omitempty" json:"description"`
	MediaType	Type					`bson:"mediaType,omitempty" json:"mediaType"`
	Category 	Category				`bson:"category,omitempty" json:"category"`
	Online 		bool 					`bson:"online,omitempty" json:"online"`
	Hash 		string					`bson:"hash,omitempty" json:"hash"`
}