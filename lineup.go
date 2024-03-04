package nyx

import (
	"bytes"
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