package server

import (
	"github.com/gin-gonic/gin"
	"net/url"
	"strconv"
)

type LinksList struct {
	Next *string `json:"next,omitempty"`
	Prev *string `json:"prev,omitempty"`
	Self string  `json:"self"`
}

type Query struct {
	Limit      int
	Offset     *int
	Q          *string
	Sort       *[]string
	Parameters map[string]string
	Host       string
}

type Paginate struct {
	Links  LinksList `json:"_links"`
	Count  int       `json:"count"`
	Limit  int       `json:"limit"`
	Offset *int      `json:"offset,omitempty"`
	Size   int       `json:"size"`
}

func GenerateUrlLink(u *url.URL, params Query) string {
	q := u.Query()
	q.Set("limit", strconv.Itoa(params.Limit))
	if params.Offset != nil {
		q.Set("offset", strconv.Itoa(*params.Offset))
	}
	if params.Q != nil {
		q.Set("q", *params.Q)
	}
	if params.Sort != nil {
		for _, pp := range *params.Sort {
			q.Add("sort", pp)
		}
	}
	for parameter, value := range params.Parameters {
		q.Add(parameter, value)
	}

	u.RawQuery = q.Encode()
	return params.Host + u.RequestURI()
}

func GetLinksList(ctx *gin.Context, count int, params Query) LinksList {
	offset := 0

	if params.Offset != nil {
		offset = *params.Offset
	}
	prevOffset := offset - params.Limit

	if prevOffset < 0 {
		prevOffset = 0
	}

	var prevUrl string
	var nextUrl string
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	host := scheme + "://" + ctx.Request.Host
	if offset != 0 {
		prev, _ := url.Parse(ctx.Request.URL.Path)
		prevUrl = GenerateUrlLink(prev, Query{
			Limit:      params.Limit,
			Offset:     &prevOffset,
			Q:          params.Q,
			Sort:       params.Sort,
			Parameters: params.Parameters,
			Host:       host,
		})
	}

	if count > params.Limit+offset {
		prev, _ := url.Parse(ctx.Request.URL.Path)
		offset = offset + params.Limit
		nextUrl = GenerateUrlLink(prev, Query{
			Limit:      params.Limit,
			Offset:     &offset,
			Q:          params.Q,
			Sort:       params.Sort,
			Parameters: params.Parameters,
			Host:       host,
		})
	}

	var next *string = &nextUrl

	if nextUrl == "" {
		next = nil
	}

	var pre *string = &prevUrl

	if prevUrl == "" {
		pre = nil
	}

	return LinksList{
		Next: next,
		Prev: pre,
		Self: host + ctx.Request.URL.RequestURI(),
	}
}

func GetPaginate[T any](ctx *gin.Context, query Query, result []T, count int) Paginate {
	links := GetLinksList(ctx, count, query)
	return Paginate{
		Links: LinksList{
			Next: links.Next,
			Prev: links.Prev,
			Self: links.Self,
		},
		Count:  count,
		Limit:  query.Limit,
		Offset: query.Offset,
		Size:   len(result),
	}
}
