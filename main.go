package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	supabase "github.com/nedpals/supabase-go"
)

// write your supabase cofig
const (
	SUPABASE_URL     = ""
	SUPABSE_ANON_KEY = ""
)

func main() {

	islandFile, err := os.Open("./island.csv")
	if err != nil {
		log.Fatalln(err)
	}
	var islands []*Island
	if err := gocsv.UnmarshalFile(islandFile, &islands); err != nil {
		log.Fatalln(err)
	}

	supabase := supabase.CreateClient(SUPABASE_URL, SUPABSE_ANON_KEY)

	var prefecures []*Prefecture
	if err := supabase.DB.From("prefecture").Select("prefecture_id, japanese_name").Execute(&prefecures); err != nil {
		log.Fatalln(err)
	}
	prefectureIDMap := map[string]string{}
	for _, v := range prefecures {
		prefectureIDMap[v.JapaneseName] = v.PrefectureID
	}

	prefixSQL := "INSERT INTO island (japanese_name, japanese_kana_name, english_name, area, town_name, prefecture_id) VALUES"
	format := `('%s', '%s', '%s', %v, '%s', '%s')`

	values := []string{}
	for _, v := range islands {
		values = append(values, fmt.Sprintf(format, v.JapaneseName,
			v.JapaneseKanaName,
			v.EnglishName,
			v.Area,
			v.TownName,
			prefectureIDMap[v.PrefectureID]))
	}
	valueSQL := strings.Join(values, ",")
	fmt.Println(prefixSQL)
	fmt.Println(valueSQL)
}

type Island struct {
	JapaneseName     string  `json:"japanese_name" csv:"japanese_name"`
	JapaneseKanaName string  `json:"japanese_kana_name" csv:"japanese_kana_name"`
	EnglishName      string  `json:"english_name" csv:"english_name"`
	TownName         string  `json:"town_name" csv:"town_name"`
	Area             float64 `json:"area" csv:"area"`
	PrefectureID     string  `json:"prefecture_id" csv:"prefecture_id"`
}

type Prefecture struct {
	PrefectureID string `json:"prefecture_id" csv:"prefecture_id"`
	JapaneseName string `json:"japanese_name" csv:"japanese_name"`
}
