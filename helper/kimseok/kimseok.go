package kimseok

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCursorFromRegex(db *mongo.Database, patttern string) (cursor *mongo.Cursor, err error) {
	filter := bson.M{"question": primitive.Regex{Pattern: patttern, Options: "i"}}
	cursor, err = db.Collection("qna").Find(context.TODO(), filter)
	return
}

func GetCursorFromString(db *mongo.Database, question string) (cursor *mongo.Cursor, err error) {
	filter := bson.M{"question": question}
	cursor, err = db.Collection("qna").Find(context.TODO(), filter)
	return
}

func GetRandomFromQnASlice(qnas []Datasets) Datasets {
	// Inisialisasi sumber acak dengan waktu saat ini
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// Pilih elemen acak dari slice
	randomIndex := rng.Intn(len(qnas))
	return qnas[randomIndex]
}

func GetQnAfromSliceWithJaro(q string, qnas []Datasets) (dt Datasets) {
	var score float64
	for _, qna := range qnas {
		//fmt.Println(data)
		str2 := qna.Question
		scorex := jaroWinkler(q, str2)
		if score < scorex {
			dt = qna
			score = scorex
		}

	}
	return

}

func GetMessage(msg itmodel.IteungMessage, botname string, db *mongo.Database) string {
	dt, err := QueriesDataRegexpALL(db, msg.Message)
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(dt.Answer)
}

func QueriesDataRegexpALL(db *mongo.Database, queries string) (dest Datasets, err error) {
	//kata akhiran imbuhan
	//queries = SeparateSuffixMu(queries)

	//ubah ke kata dasar
	queries = Stemmer(queries)
	filter := bson.M{"question": queries}
	qnas, err := atdb.GetAllDoc[[]Datasets](db, "qna", filter)
	if err != nil {
		return
	}
	if len(qnas) > 0 {
		dest = GetRandomFromQnASlice(qnas)
		return
	}

	words := strings.Fields(queries) // Tokenization
	//pencarian dengan pengurangan kata dari belakang
	for len(words) > 0 {
		// Join remaining elements back into a string
		remainingMessage := strings.Join(words, " ")
		filter := bson.M{"question": primitive.Regex{Pattern: remainingMessage, Options: "i"}}
		qnas, err = atdb.GetAllDoc[[]Datasets](db, "qna", filter)
		if err != nil {
			return
		} else if len(qnas) > 0 {
			dest = GetQnAfromSliceWithJaro(queries, qnas)
			return
		}
		// Remove the last element
		words = words[:len(words)-1]
	}
	// Reset words untuk pencarian dengan pengurangan kata dari depan
	words = strings.Fields(queries)

	// Pencarian dengan pengurangan kata dari depan
	for len(words) > 0 {
		// Gabungkan elemen yang tersisa menjadi satu string
		remainingMessage := strings.Join(words, " ")
		filter := bson.M{"question": primitive.Regex{Pattern: remainingMessage, Options: "i"}}
		qnas, err = atdb.GetAllDoc[[]Datasets](db, "qna", filter)
		if err != nil {
			return
		} else if len(qnas) > 0 {
			dest = GetQnAfromSliceWithJaro(queries, qnas)
			return
		}
		// Hapus elemen pertama
		words = words[1:]
	}

	return
}
