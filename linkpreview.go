package linkpreview

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type LinkPreview struct {
	Title       string 
	Description string 
	Image       string 
}

func GetLinkPreview(ctx context.Context, previewUrl string) (*LinkPreview, error) {
	urlPreview, err := url.Parse(previewUrl)
	if err != nil {
		return nil, err
	}

	if urlPreview.Scheme != "http" && urlPreview.Scheme != "https" {
		return nil, &ErrNonHTTPScheme{url: previewUrl}
	}

	// init Client
	client, err := http.NewRequestWithContext(ctx, "GET", previewUrl, nil)
	if err != nil {
		return nil, err
	}

	// TODO: Add support for custom clients
	r, err := http.DefaultClient.Do(client)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	// Check if the response code is 200
	// TODO: Need support for redirects
	if r.StatusCode != 200 {
		return nil, &ErrNonSupportedResponse{respCode: r.StatusCode}
	}

	c := r.Header.Get("Content-Type")

	// parse Mime type
	c, _, _ = strings.Cut(c, ";")
	contentType := strings.TrimSpace(strings.ToLower(c))

	if contentType != "text/html" {
		return nil, &ErrNonSupportedContentType{contentType: contentType}
	}

	gq, err := goquery.NewDocumentFromReader(r.Body)

	if err != nil {
		return nil, err
	}

	title := getTitle(gq)
	description := getDescription(gq)
	image := getImage(gq)	

	if image != "" && !(strings.HasPrefix("http://", image) || strings.HasPrefix("https://", image)) {
		parsedURL, _ := url.Parse(previewUrl)

		joinedURL := url.URL{
			Scheme: parsedURL.Scheme,
			Host:   parsedURL.Host,
			Path:   image,
		}

		image = joinedURL.String()
	}

	return &LinkPreview{
		Title:      title,
		Description: description,
		Image:       image,
	}, nil

}

func getTagsContents(s *goquery.Document, tags []string) string {
	for _, t := range tags {
		e := s.Find(t)
		if e.Size() != 0 {
			if c, b := e.Attr("content"); b {
				return c
			}
		}
	}

	return ""
}

func getTitle(s *goquery.Document) string {
	titleQuery := []string{
		"meta[name='title']",
		"meta[property='og:title']",
		"meta[property='twitter:title']",
	}

	r := getTagsContents(s, titleQuery)
	
	if r == "" {
		return s.Find("title").Text()
	}

	return r
}

func getDescription(s *goquery.Document) string {
	descriptionQuery := []string{
		"meta[name='description']",
		"meta[property='og:description']",
		"meta[property='twitter:description']",
	}

	return getTagsContents(s, descriptionQuery)
}

func getImage(s *goquery.Document) string {
	imageQuery := []string{
		"meta[property='og:image']",
		"link[rel='image_src']",
		"meta[property='twitter:image']",
	}
	i := getTagsContents(s, imageQuery)

	if i == "" {
		i = s.Find("img").First().AttrOr("src", "")
	}

	return i
}