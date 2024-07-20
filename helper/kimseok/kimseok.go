package kimseok

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func QueriesDataRegexpALL(db *mongo.Database, ctx context.Context, queries string) (dest Datasets, score float64, err error) {
	//kata akhiran imbuhan
	queries = SeparateSuffixMu(queries)
	var cursor *mongo.Cursor
	//ubah ke kata dasar
	queries = Stemmer(queries)
	splits := strings.Split(queries, " ")
	if len(splits) >= 5 {
		queries = splits[len(splits)-3] + " " + splits[len(splits)-2] + " " + splits[len(splits)-1]
		filter := bson.M{"question": primitive.Regex{Pattern: queries, Options: "i"}}
		cursor, err = db.Collection("datasets").Find(ctx, filter)

		if err != nil && err != mongo.ErrNoDocuments {
			queries = splits[len(splits)-2] + " " + splits[len(splits)-1]
			filter = bson.M{"question": primitive.Regex{Pattern: queries, Options: "i"}}
			cursor, err = db.Collection("datasets").Find(ctx, filter)
			if err != nil && err != mongo.ErrNoDocuments {
				queries = splits[len(splits)-1]
				filter = bson.M{"question": primitive.Regex{Pattern: queries, Options: "i"}}
				cursor, err = db.Collection("datasets").Find(ctx, filter)
				if err != nil && err != mongo.ErrNoDocuments {
					filter = bson.M{"question": primitive.Regex{Pattern: queries, Options: "i"}}
					cursor, err = db.Collection("datasets").Find(ctx, filter)
					if err != nil && err != mongo.ErrNoDocuments {
						return dest, score, err
					}
				}
			}
		}
	} else if len(splits) == 1 {
		queries = splits[0]
		filter := bson.M{"question": primitive.Regex{Pattern: queries, Options: "i"}}
		cursor, err = db.Collection("datasets").Find(ctx, filter)
	} else if len(splits) <= 4 {
		queries = splits[len(splits)-2] + " " + splits[len(splits)-1]
		filter := bson.M{"question": primitive.Regex{Pattern: queries, Options: "i"}}
		cursor, err = db.Collection("datasets").Find(ctx, filter)

		if err != nil && err != mongo.ErrNoDocuments {
			queries = splits[len(splits)-1]
			filter = bson.M{"question": primitive.Regex{Pattern: queries, Options: "i"}}
			cursor, err = db.Collection("datasets").Find(ctx, filter)
			if err != nil && err != mongo.ErrNoDocuments {
				return dest, score, err
			}
		}
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var data Datasets
		err := cursor.Decode(&data)
		if err != nil {
			return data, score, err
		}
		//fmt.Println(data)
		str2 := data.Question
		scorex := jaroWinkler(queries, str2)
		if score < scorex {
			dest = data
			score = scorex
		}
	}
	return dest, score, err
}
