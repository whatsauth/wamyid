package pomokit

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PomodoroReport struct {
    ID                  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
    Name                string             `bson:"name" json:"name"`
    PhoneNumber         string             `bson:"phonenumber,omitempty" json:"phonenumber,omitempty"`
    Cycle               int                `bson:"cycle" json:"cycle"`
    Hostname            string             `bson:"hostname" json:"hostname"`
    IP                  string             `bson:"ip" json:"ip"`
    Screenshots         int                `bson:"screenshots" json:"screenshots"`
    Pekerjaan           string             `bson:"pekerjaan" json:"pekerjaan"`
    Token               string             `bson:"token" json:"token"`
    URLPekerjaan        string             `bson:"urlpekerjaan" json:"urlpekerjaan"`
    WaGroupID           string             `bson:"wagroupid" json:"wagroupid"`
    GTmetrixURLTarget   string             `bson:"gtmetrix_url_target" json:"gtmetrix_url_target"`
    GTmetrixGrade       string             `bson:"gtmetrix_grade" json:"gtmetrix_grade"`
    GTmetrixPerformance string             `bson:"gtmetrix_performance" json:"gtmetrix_performance"`
    GTmetrixStructure   string             `bson:"gtmetrix_structure" json:"gtmetrix_structure"`
    LCP                 string             `bson:"lcp" json:"lcp"`
    TBT                 string             `bson:"tbt" json:"tbt"`
    CLS                 string             `bson:"cls" json:"cls"`
    CreatedAt           time.Time          `bson:"createdAt" json:"createdAt"`
}

type Config struct {
    PublicKeyPomokit    string             `json:"publickeypomokit,omitempty" bson:"publickeypomokit,omitempty"`
    DomyikadoAllUserURL string             `json:"domyikadoalluserurl,omitempty" bson:"domyikadoalluserurl,omitempty"`
}