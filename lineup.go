package nyx

import (
	"bytes"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

type Artist struct {
	Name string `json:"name"`
	Img string `json:"img"`
}

type Lineup struct {
	Artists []Artist `json:"artists"`
	Size    int      `json:"size"`
}

// Get()
//
// Get returns the artists listed at the given url
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

	defer resp.Body.Close();

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
//     <img class="favorite-item__img" src="./lineup_files/07b05349-da5b-11ed-b991-0ee6b8365494.jpg" alt="Da Tweekaz">
//     <div class="favorite-item__info">
//         <h4 class="headline">Da Tweekaz</h4>
//         <p></p>
//         <!-- <p></p> -->
//     </div>
//     <!-- <div class="favorite-item__action">
//         <img src="/wp-content/assets/splashpages/app-sharing/lineup/img/heart.svg" alt="heart">
//     </div> -->
// </div>
func ParseLineup(rawHTML []byte) (Lineup, error) {
	var lineup Lineup

	buff := bytes.NewBuffer(rawHTML)

	doc, err := html.Parse(buff)

	if err != nil {
		return lineup, err
	}

	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {

				if a.Key == "class" && a.Val == "favorite-item" {

						artist := Artist{
							Name: getArtist(n),
							Img: getArtistImg(n),
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