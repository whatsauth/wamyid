package strava

import (
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func StravaActivityUpdateIfEmptyDataHandler(Pesan itmodel.IteungMessage, db *mongo.Database) string {
	reply := "Informasi Stava kamu hari ini: "

	if Pesan.Phone_number != "6282268895372" {
		if isMaintenance {
			reply += "\n\nMaaf kak, sistem sedang maintenance. Coba lagi nanti ya."
			return reply
		}
	}

	// cek apakah akun strava sudah terdaftar di database
	Idata, err := atdb.GetOneDoc[StravaIdentity](db, "strava_identity", bson.M{"phone_number": Pesan.Phone_number})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "Kak, akun Strava kamu belum terdaftar. Silakan daftar dulu!"
		}
		return "\n\nError fetching data dari MongoDB: " + err.Error()
	}

	col := "strava_activity"
	// cek apakah akun strava sudah terdaftar di database
	data, err := atdb.GetOneDoc[StravaActivity](db, col, bson.M{"picture": Idata.Picture})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "Kak, Strava Activity kamu tidak ditemukan di database. Silakan cek kembali."
		}
		return "\n\nError fetching data dari MongoDB: " + err.Error()
	}

	c := colly.NewCollector(
		colly.AllowedDomains(domWeb),
	)

	// ambil link strava activity dari pesan
	rawUrl := extractStravaLink(Pesan.Message)
	if rawUrl == "" {
		return reply + "\n\nMaaf, pesan yang kamu kirim tidak mengandung link Strava. Silakan kirim link aktivitas Strava untuk mendapatkan informasinya."
	}

	stravaActivity := StravaActivity{}
	stravaActivity.ActivityId = data.ActivityId

	c.OnHTML("main", func(e *colly.HTMLElement) {
		e.ForEach("img", func(_ int, imgEl *colly.HTMLElement) {
			imgTitle := imgEl.Attr("title")
			if imgTitle == stravaActivity.Name {
				stravaActivity.Picture = imgEl.Attr("src")
			}
		})
	})

	c.OnHTML("div", func(e *colly.HTMLElement) {
		e.ForEach("span", func(_ int, el *colly.HTMLElement) {
			label := strings.ToLower(strings.TrimSpace(el.Text))
			value := strings.TrimSpace(el.DOM.Next().Text()) // Ambil elemen di sebelahnya

			switch label {
			case "distance":
				stravaActivity.Distance = value
			case "time":
				stravaActivity.MovingTime = value
			case "elevation":
				stravaActivity.Elevation = value
			}
		})
	})

	found := false
	c.OnHTML("div.MapAndElevationChart_mapContainer__VIs6u", func(e *colly.HTMLElement) {
		found = true
	})

	c.OnScraped(func(r *colly.Response) {
		if data.ActivityId == "" {
			reply += "\n\n Strava Activity kak " + Pesan.Alias_name + " tidak di temukan."
			return
		}

		if stravaActivity.Distance == "" || stravaActivity.MovingTime == "" {
			reply += "\n\nMaaf kak, kami tidak dapat mengambil data aktivitas kamu. Coba hubungi admin ya."
			return
		}

		if stravaActivity.Picture == "" {
			reply += "\n\nMaaf kak, sistem tidak dapat mengambil foto profil Strava kamu. Pastikan akun Strava kamu dibuat public(everyone). doc: https://www.do.my.id/mentalhealt-strava"
			return
		}

		if Idata.Picture != stravaActivity.Picture {
			reply += "\n\nAda yang salah nih dengan akun strava kamu, coba lakukan update dengan perintah dibawah yaaa"
			reply += "\n\n *strava update in*"
			reply += "\n\nAtau mungkin link yang kamu share bukan punya kamu 😏"
			return
		}

		if data.ActivityId != stravaActivity.ActivityId {
			reply += "\n\nStrava Activity ID kak " + Pesan.Alias_name + " tidak sama."
			return
		}

		if data.Distance != "" && data.MovingTime != "" {
			reply += "\n\nData Strava kak " + Pesan.Alias_name + " sudah up to date."
			return
		}

		// cek apakah ada map atau tidak di halaman strava
		if !found {
			reply += "\n\nJangan Curang donggg! Silahkan share record aktivitas yang benar dari Strava ya kak, bukan dibikin manual kaya gitu"
			reply += "\nYang semangat dong... yang semangat dong..."
			return
		}

		if stravaActivity.TypeSport == "Ride" {
			reply += "\n\nMaaf kak, sistem hanya dapat mengambil data aktivitas jalan dan lari. Silakan share link aktivitas jalan dan lari Strava kamu."
			return
		}

		if data.Distance == "" && data.MovingTime == "" {
			if data.Status == "Invalid" {
				distanceFloat := parseDistance(stravaActivity.Distance)
				if distanceFloat < 3 {
					stravaActivity.UpdatedAt = time.Now()
					stravaActivity.Status = "Invalid"

					updateData := bson.M{
						"distance":    stravaActivity.Distance,
						"moving_time": stravaActivity.MovingTime,
						"elevation":   stravaActivity.Elevation,
						"updated_at":  stravaActivity.UpdatedAt,
						"status":      stravaActivity.Status,
					}

					// update data ke database jika ada perubahan
					_, err := atdb.UpdateDoc(db, col, bson.M{"activity_id": stravaActivity.ActivityId}, bson.M{"$set": updateData})
					if err != nil {
						reply += "\n\nError updating data to MongoDB: " + err.Error()
						return
					}

					reply += "\n\nWahhh, kamu malas sekali ya, jangan malas lari terus dong kak! 😏"
					reply += "\nSatu hari minimal 3 km, masa kamu cuma " + stravaActivity.Distance + " aja"
					reply += "\n\nJangan lupa jaga kesehatan dan tetap semangat!! 💪🏻💪🏻💪🏻"
					return
				} else {
					// simpan data ke database jika data belum ada
					stravaActivity.CreatedAt = time.Now()
					stravaActivity.Status = "Valid"

					updateData := bson.M{
						"distance":    stravaActivity.Distance,
						"moving_time": stravaActivity.MovingTime,
						"elevation":   stravaActivity.Elevation,
						"updated_at":  stravaActivity.UpdatedAt,
						"status":      stravaActivity.Status,
					}

					// update data ke database jika ada perubahan
					_, err := atdb.UpdateDoc(db, col, bson.M{"activity_id": stravaActivity.ActivityId}, bson.M{"$set": updateData})
					if err != nil {
						reply += "\n\nError updating data to MongoDB: " + err.Error()
						return

					} else {
						reply += "\n\nHaiiiii kak, " + "*" + Pesan.Alias_name + "*" + "! Berikut Progres Aktivitas kamu hari ini yaaa yang di update!! 😀"
						reply += "\n\n- Name: " + stravaActivity.Name
						reply += "\n- Title: " + stravaActivity.Title
						reply += "\n- Date Time: " + stravaActivity.DateTime
						reply += "\n- Type Sport: " + stravaActivity.TypeSport
						reply += "\n- Distance: " + stravaActivity.Distance
						reply += "\n- Moving Time: " + stravaActivity.MovingTime
						reply += "\n- Elevation: " + stravaActivity.Elevation
						reply += "\n\nSemangat terus, jangan lupa jaga kesehatan dan tetap semangat!! 💪🏻💪🏻💪🏻"
					}
				}

			} else {
				reply += "\n\nMaaf kak, Tidak bisa mengambil data aktivitas kamu."
				return
			}
		}
	})

	rawUrl = data.LinkActivity

	err = c.Visit(rawUrl)
	if err != nil {
		return "Link Profile Strava yang anda kirimkan tidak valid. Silakan kirim ulang dengan link yang valid.(3)"
	}

	return reply
}
