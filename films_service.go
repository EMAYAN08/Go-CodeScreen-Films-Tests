package films

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

const filmsEndpointUrl string = "https://toolbox.palette-adv.spectrocloud.com:5002/films"

// Your API token. Needed to successfully authenticate when calling the films endpoint.
// Must be included in the Authorization header in the request sent to the films endpoint.
const apiToken string = "8c5996d5-fb89-46c9-8821-7063cfbc18b1"

type Film struct {
	Name         string  `json:"name"`
	Length       int     `json:"length"`
	Rating       float64 `json:"rating"`
	ReleaseDate  string  `json:"releaseDate"`
	DirectorName string  `json:"directorName"`
}

var (
	films     []Film
	filmsOnce sync.Once
)

func fetchFilms() ([]Film, error) {
	req, _ := http.NewRequest("GET", filmsEndpointUrl, nil)

	req.Header.Set("Authorization", "Bearer "+apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var films []Film
	if err := json.NewDecoder(resp.Body).Decode(&films); err != nil {
		return nil, err
	}
	return films, nil
}

// GetFilms retrieves the data for all films by calling the https://toolbox.palette-adv.spectrocloud.com:5002/films endpoint.
func GetFilms() []Film {
	filmsOnce.Do(func() {
		films, _ = fetchFilms()
	})

	return films
}

// BestRatedFilm retrieves the name of the best rated film that was directed by the director with the given name.
// If there are no films directed by the given director, return an empty string.
// Note: there will only be one film with the best rating.
func BestRatedFilm(directorName string) string {
	films := GetFilms()
	var bestFilm Film
	found := false

	for _, film := range films {
		if film.DirectorName == directorName {
			if !found || film.Rating > bestFilm.Rating {
				bestFilm = film
				found = true
			}
		}
	}

	if found {
		return bestFilm.Name
	}
	return ""
}

// DirectorWithMostFilms retrieves the name of the director who has directed the most films
// in the CodeScreen Film service.
func DirectorWithMostFilms() string {

	films := GetFilms()
	count := make(map[string]int)

	for _, film := range films {
		count[film.DirectorName]++
	}
	var maxDir string
	max := 0

	for dir, count := range count {
		if count > max {
			maxDir = dir
			max = count
		}
	}
	//TODO Implement
	return maxDir
}

// AverageRating retrieves the average rating for the films directed by the given director, rounded to 1 decimal place.
// If there are no films directed by the given director, return 0.0.
func AverageRating(directorName string) float64 {
	// TODO: Implement the actual retrieval of films
	films := GetFilms()
	var totalRating float64
	var count int

	for _, film := range films {
		// Check if the film's director matches the specified director
		if film.DirectorName == directorName {
			totalRating += film.Rating
			count++
		}
	}

	// If no films were found for the director, return 0.0
	if count == 0 {
		return 0.0
	}

	// Calculate the average rating
	value := totalRating / float64(count)
	return float64(int(value*10+0.5)) / 10
}

// func roundToOneDecimal(value float64) float64 {
// 	return float64(int(value*10+0.5)) / 10
// }

/*
ShortestFilmReleaseGap retrieves the shortest number of days between any two film releases directed by the given director.
If there are no films directed by the given director, return 0.
If there is only one film directed by the given director, return 0.
Note: no director released more than one film on any given day.

For example, if the service returns the following 3 films:

	{
	    "name": "Batman Begins",
	    "length": 140,
	    "rating": 8.2,
	    "releaseDate": "2006-06-16",
	    "directorName": "Christopher Nolan"
	},

	{
	    "name": "Interstellar",
	    "length": 169,
	    "rating": 8.6,
	    "releaseDate": "2014-11-07",
	    "directorName": "Christopher Nolan"
	},

	{
	    "name": "Prestige",
	    "length": 130,
	    "rating": 8.5,
	    "releaseDate": "2006-11-10",
	    "directorName": "Christopher Nolan"
	}

Then this method should return 147 for Christopher Nolan, as Prestige was released 147 days after Batman Begins.
*/
func ShortestFilmReleaseGap(directorName string) int {
	//TODO Implement
	films := GetFilms()
	var releaseDates []time.Time
	for _, film := range films {
		// Check if the film's director matches the specified director
		if film.DirectorName == directorName {
			// Parse the release date of the film
			date, err := time.Parse("2006-01-02", film.ReleaseDate)
			if err == nil {
				// If the date is parsed successfully, add it to the releaseDates slice
				releaseDates = append(releaseDates, date)
			}
		}
	}
	if len(releaseDates) < 2 {
		return 0
	}
	sort.Slice(releaseDates, func(i, j int) bool {
		return releaseDates[i].Before(releaseDates[j])
	})
	minGap := int(^uint(0) >> 1) // max int
	for i := 1; i < len(releaseDates); i++ {
		// Calculate the gap in days between consecutive release dates
		gap := int(releaseDates[i].Sub(releaseDates[i-1]).Hours() / 24)
		// Update minGap if the current gap is smaller
		if gap < minGap {
			minGap = gap
		}
	}
	return minGap
}
