package nyx

import (
	"bytes"
	// "fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Artist struct {
	Img     string    `json:"img"`
	Name    string    `json:"name"`
	Stage   string    `json:"stage"`
	SetTime time.Time `json:"setTime"`
}

type Lineup struct {
	Artists []Artist `json:"artists"`
	Size    int      `json:"size"`
}

const (
	defaultYear = 2024
	// defaultTime =
)

// Get()
//
// # Get returns the artists listed at the given url
//
// Parameters:
// - url (string): URL to the favorited artists from the Insomniac app. Example: https://insom.app/<unique-id>
//
// Returns:
// - Lineup: Returns the artists listed at the Insomniac URL
// - error: Is nil if request and parsing is successful
func Get(url string) (Lineup, error) {
	var lineup Lineup

	resp, err := http.Get(url)

	if err != nil {
		return lineup, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return lineup, err
	}

	return ParseLineup(body)
}

// ParseLineup takes in raw html bytes and outputs the Lineup or errors out
//
// example node:
// <div class="favorite-item">
//
//	<img class="favorite-item__img" src="./lineup_files/07b05349-da5b-11ed-b991-0ee6b8365494.jpg" alt="Da Tweekaz">
//	<div class="favorite-item__info">
//	    <h4 class="headline">Da Tweekaz</h4>
//	    <p></p>
//	    <!-- <p></p> -->
//	</div>
//	<!-- <div class="favorite-item__action">
//	    <img src="/wp-content/assets/splashpages/app-sharing/lineup/img/heart.svg" alt="heart">
//	</div> -->
//
// </div>
func ParseLineup(rawHTML []byte) (Lineup, error) {
	var lineup Lineup
	var setTime time.Time

	buff := bytes.NewBuffer(rawHTML)

	doc, err := html.Parse(buff)

	if err != nil {
		return lineup, err
	}

	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "favorite-panel" {
					setTime = getTimePerforming(n)
				}

				if a.Key == "class" && a.Val == "favorite-item" {

					artist := Artist{
						Img:     getArtistImg(n),
						Name:    getArtist(n),
						Stage:   getStage(n),
						SetTime: setTime,
					}

					lineup.Artists = append(lineup.Artists, artist)
					lineup.Size += 1

					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return lineup, nil
}

// return empty string if not found
func getArtist(n *html.Node) string {
	imgNode := n.FirstChild.NextSibling

	if imgNode == nil {
		return ""
	}

	divFavItem := imgNode.NextSibling.NextSibling

	if divFavItem == nil {
		return ""
	}

	h4Headline := divFavItem.FirstChild.NextSibling

	if h4Headline == nil {
		return ""
	}

	artist := h4Headline.FirstChild

	if artist == nil {
		return ""
	}

	return artist.Data
}

func getArtistImg(n *html.Node) string {
	imgSrc := ""

	if n == nil || n.FirstChild == nil {
		return imgSrc
	}

	if n.FirstChild.NextSibling == nil {
		return imgSrc
	}

	imgNode := n.FirstChild.NextSibling

	for _, attr := range imgNode.Attr {
		if attr.Key != "src" {
			continue
		}

		imgSrc = attr.Val

		break
	}

	return imgSrc
}

func getTimePerforming(n *html.Node) time.Time {
	// verify we are at the right starting point
	hasAttr := false

	for _, attr := range n.Attr {
		if attr.Key == "class" && attr.Val == "favorite-panel" {
			hasAttr = true

			break
		}
	}

	if !hasAttr {
		return time.Time{}
	}

	// verify we are on the <p> tag
	dayNode := n.FirstChild.NextSibling

	if dayNode == nil {
		return time.Time{}
	}

	if dayNode.Data != "p" {
		return time.Time{}
	}

	for _, attr := range dayNode.Attr {
		// find "favorite-date" and convert "<DAY OF WEEK>, <Month> <DAY>"" to time.Time

		if attr.Key == "class" && attr.Val == "favorite-date" {
			dayText := dayNode.FirstChild.Data

			// fmt.Println(dayText)

			date, err := parseDate(dayText, defaultYear)

			if err != nil {
				return time.Time{}
			}

			return date
		}
	}

	// fmt.Println("didn't find day")

	return time.Time{}
}

func parseDate(dateStr string, year int) (time.Time, error) {
	if year == 0 {
		year = defaultYear
	}

	dateSplit := strings.Split(dateStr, " ")

	month := getMonth(dateSplit[1])
	day, err := strconv.ParseInt(dateSplit[2], 10, 32)

	if err != nil {
		return time.Time{}, err
	}

	loc := time.FixedZone("UTC-8", -8*60*60)

	return time.Date(year, month, int(day), 16, 0, 0, 0, loc), nil
}

func getMonth(monthStr string) time.Month {
	switch monthStr {
	case "JANUARY":
		return time.January
	case "FEBRUARY":
		return time.February
	case "MARCH":
		return time.March
	case "APRIL":
		return time.April
	case "MAY":
		return time.May
	case "JUNE":
		return time.June
	case "JULY":
		return time.July
	case "AUGUST":
		return time.August
	case "SEPTEMBER":
		return time.September
	case "OCTOBER":
		return time.October
	case "NOVEMBER":
		return time.November
	case "DECEMBER":
		return time.December
	default:
		return time.Month(0)
	}
}

func getStage(n *html.Node) string {
	// confirm we are at div.favorite-item
	hasAttr := false

	for _, attr := range n.Attr {
		if attr.Key == "class" && attr.Val == "favorite-item" {
			hasAttr = true

			break
		}
	}

	if !hasAttr {
		return ""
	}

	artistNode := n.FirstChild.NextSibling

	if artistNode == nil {
		return ""
	}

	infoNode := artistNode.NextSibling

	if infoNode == nil {
		return ""
	}

	infoNode = infoNode.NextSibling

	if infoNode == nil {
		return ""
	}

	hasAttr = false

	for _, attr := range infoNode.Attr {
		if attr.Key == "class" && attr.Val == "favorite-item__info" {
			hasAttr = true

			break
		}
	}

	if !hasAttr {
		return ""
	}

	artistNameNode := infoNode.FirstChild

	if artistNameNode == nil {
		return ""
	}

	artistNameNode = artistNameNode.NextSibling

	if artistNameNode == nil {
		return ""
	}

	stageNode := artistNameNode.NextSibling

	if stageNode == nil {
		return ""
	}

	stageNode = stageNode.NextSibling

	if stageNode == nil {
		return ""
	}

	if stageNode.DataAtom != atom.P {
		return ""
	}

	stageNode = stageNode.FirstChild

	if stageNode == nil {
		return ""
	}

	return stageNode.Data
}
