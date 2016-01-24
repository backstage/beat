package schemas

import (
	"fmt"
	"net/url"
)

type Link struct {
	Rel          string                 `json:"rel" bson:"rel"`
	Href         string                 `json:"href" bson:"href"`
	Title        string                 `json:"title,omitempty" bson:"title,omitempty"`
	TargetSchema map[string]interface{} `json:"targetSchema,omitempty" bson:"targetSchema,omitempty"`
	MediaType    string                 `json:"mediaType,omitempty" bson:"mediaType,omitempty"`
	Method       string                 `json:"method,omitempty" bson:"method,omitempty"`
	EncType      string                 `json:"encType,omitempty" bson:"encType,omitempty"`
	Schema       map[string]interface{} `json:"schema,omitempty" bson:"schema,omitempty"`
}

type Links []*Link

func (l Links) ApplyBaseUrl(baseUrl string) {
	for _, link := range l {
		if isRelativeLink(link.Href) {
			link.Href = fmt.Sprintf("%s%s", baseUrl, link.Href)
		}
	}
}

// ConcatenateLinks generate new links with merge with tailLinks
func (l Links) ConcatenateLinks(tailLinks *Links) *Links {
	currentSize := len(l)
	expandSize := len(*tailLinks)

	newLinks := make(Links, currentSize+expandSize)
	copy(newLinks, l)

	for i, link := range *tailLinks {
		newLinks[currentSize+i] = link
	}

	return &newLinks
}

func BuildDefaultLinks(collectionName string) Links {
	collectionUrl := fmt.Sprintf("/%s", collectionName)
	itemUrl := fmt.Sprintf("/%s/{id}", collectionName)

	return Links{
		&Link{Rel: "self", Href: itemUrl},
		&Link{Rel: "item", Href: itemUrl},
		&Link{Rel: "create", Method: "POST", Href: collectionUrl},
		&Link{Rel: "update", Method: "PUT", Href: itemUrl},
		&Link{Rel: "delete", Method: "DELETE", Href: itemUrl},
		&Link{Rel: "parent", Href: collectionUrl},
	}
}

func isRelativeLink(link string) bool {
	url, err := url.Parse(link)

	if err != nil {
		return false
	}

	return url.Host == "" && url.Scheme == "" && !isUriTemplate(link)
}

func isUriTemplate(link string) bool {
	return len(link) > 0 && link[0] == '{'
}