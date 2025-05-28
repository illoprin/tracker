package genreType

import (
	"log/slog"
	"reflect"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

type GenreModel map[string]string

var AllowedGenres = []string{
	"classical",
	"baroque",
	"romantic",
	"opera",
	"ballet",
	"jazz",
	"blues",
	"swing",
	"bebop",
	"cool jazz",
	"free jazz",
	"soul",
	"funk",
	"rhythm and blues",
	"rock",
	"hard rock",
	"psychedelic rock",
	"progressive rock",
	"punk rock",
	"post-punk",
	"garage rock",
	"grunge",
	"alternative rock",
	"indie rock",
	"shoegaze",
	"dream pop",
	"noise rock",
	"gothic rock",
	"industrial rock",
	"metal",
	"heavy metal",
	"thrash metal",
	"black metal",
	"death metal",
	"doom metal",
	"progressive metal",
	"nu metal",
	"metalcore",
	"folk",
	"country",
	"bluegrass",
	"americana",
	"pop",
	"synthpop",
	"electropop",
	"indie pop",
	"dream pop",
	"dance pop",
	"teen pop",
	"eurodance",
	"house",
	"deep house",
	"progressive house",
	"tech house",
	"electro house",
	"acid house",
	"techno",
	"minimal techno",
	"detroit techno",
	"trance",
	"progressive trance",
	"uplifting trance",
	"psytrance",
	"hard trance",
	"drum and bass",
	"jungle",
	"dubstep",
	"brostep",
	"trap",
	"electro",
	"breakbeat",
	"big beat",
	"idm",
	"glitch",
	"ambient",
	"dark ambient",
	"new age",
	"chillout",
	"lo-fi",
	"trip hop",
	"hip hop",
	"rap",
	"boom bap",
	"trap rap",
	"drill",
	"gangsta rap",
	"conscious hip hop",
	"r&b",
	"neo soul",
	"k-pop",
	"j-pop",
	"c-pop",
	"reggae",
	"ska",
	"dub",
	"afrobeat",
	"latin",
	"bossa nova",
	"salsa",
	"tango",
	"flamenco",
	"world music",
}

func ValidateGenres(fl validator.FieldLevel) bool {
	// configure logger
	logger := slog.With(slog.String("function", "genreType.ValidateGenres"))

	// check field type
	field := fl.Field()
	if field.Kind() != reflect.Slice {
		return false
	}

	// check min length
	if field.Len() < 1 {
		return false
	}

	// check each element
	for i := 0; i < field.Len(); i++ {
		genre := field.Index(i).String()

		// to lower case, trim
		genre = strings.ToLower(strings.TrimSpace(genre))

		logger.Info("validate genre", slog.String("genre", genre))

		// check genre in allowed genres slice
		if !slices.Contains(AllowedGenres, genre) {
			return false
		}
	}

	return true
}
