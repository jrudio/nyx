package nyx

import (
	"bytes"
	"golang.org/x/net/html"
)

type Lineup struct {
	Artists []string `json:"artists"`
	Size    int      `json:"size"`
}

// ParseLineup takes in an io.Reader and outputs the Lineup or error
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
func ParseLineup(reader []byte) (Lineup, error) {
	var lineup Lineup

	buff := bytes.NewBuffer(reader)

	doc, err := html.Parse(buff)

	if err != nil {
		return lineup, err
	}

	var f func(*html.Node)

	f = func(n *html.Node) {
		// if n.Type == html.ElementNode && n.Data == "a" {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "favorite-item" {

					// if (n.FirstChild.Type == html.ElementNode && n.FirstChild.NextSibling.Type == html.ElementNode ) {
						// imgNode := n.FirstChild.NextSibling
						// divFavItem := imgNode.NextSibling.NextSibling
						// h4Headline := divFavItem.FirstChild.NextSibling
						// log.Println(h4Headline.FirstChild.Data)
						// if (n.FirstChild.Type == html.ElementNode ) {
							// log.Println(getArtist(n.FirstChild.NextSibling))
							// }
						artist := getFromArtistParentNode(n)

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
func getFromArtistParentNode(n *html.Node) string {
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