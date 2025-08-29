package itunes_search

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

type SearchParamsBuilder struct {
	params map[string]string
}

func NewSearchParams() *SearchParamsBuilder {
	return &SearchParamsBuilder{
		params: make(map[string]string),
	}
}

func (sp *SearchParamsBuilder) Term(term string) *SearchParamsBuilder {
	if term != "" {
		sp.params["term"] = term
	}
	return sp
}

func (sp *SearchParamsBuilder) Country(country string) *SearchParamsBuilder {
	if country != "" {
		sp.params["country"] = country
	}
	return sp
}

func (sp *SearchParamsBuilder) Media(media string) *SearchParamsBuilder {
	if media != "" {
		sp.params["media"] = media
	}
	return sp
}

func (sp *SearchParamsBuilder) Entity(entity string) *SearchParamsBuilder {
	if entity != "" {
		sp.params["entity"] = entity
	}
	return sp
}

func (sp *SearchParamsBuilder) Attribute(attribute string) *SearchParamsBuilder {
	if attribute != "" {
		sp.params["attribute"] = attribute
	}
	return sp
}

func (sp *SearchParamsBuilder) Callback(callback string) *SearchParamsBuilder {
	if callback != "" {
		sp.params["callback"] = callback
	}
	return sp
}

func (sp *SearchParamsBuilder) Limit(limit int) *SearchParamsBuilder {
	if limit > 0 {
		if limit > MaxLimit {
			limit = MaxLimit
		}
		sp.params["limit"] = strconv.Itoa(limit)
	}
	return sp
}

func (sp *SearchParamsBuilder) Lang(lang string) *SearchParamsBuilder {
	if lang != "" {
		sp.params["lang"] = lang
	}
	return sp
}

func (sp *SearchParamsBuilder) Version(version int) *SearchParamsBuilder {
	if version > 0 {
		sp.params["version"] = strconv.Itoa(version)
	}
	return sp
}

func (sp *SearchParamsBuilder) Explicit(explicit string) *SearchParamsBuilder {
	if explicit != "" {
		sp.params["explicit"] = explicit
	}
	return sp
}

func (sp *SearchParamsBuilder) Build() map[string]string {
	return sp.params
}

func (c *Client) Search(params map[string]string) (*SearchResponse, error) {
	if len(params) == 0 {
		c.logger.Error("At least one search parameter is required")
		return nil, fmt.Errorf("at least one search parameter is required")
	}

	queryString := c.buildQueryString(params)
	c.logger.Debug("Making iTunes search request", zap.String("query", queryString))

	resp, err := c.baseClient.HTTP.R().
		SetQueryParams(params).
		Get(BaseSearchURL)

	if err != nil {
		c.logger.Error("Failed to make search request", zap.Error(err))
		return nil, fmt.Errorf("failed to make search request: %w", err)
	}

	if !resp.IsSuccess() {
		c.logger.Error("Search request failed", zap.Int("status", resp.StatusCode()))
		return nil, fmt.Errorf("search request failed with status %d", resp.StatusCode())
	}

	var searchResponse SearchResponse
	err = json.Unmarshal(resp.Body(), &searchResponse)
	if err != nil {
		c.logger.Error("Failed to unmarshal search response", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	c.logger.Debug("Search request successful",
		zap.Int("results", searchResponse.ResultCount),
		zap.Int("status", resp.StatusCode()))

	return &searchResponse, nil
}
