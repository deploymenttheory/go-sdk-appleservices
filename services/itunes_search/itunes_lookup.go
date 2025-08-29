package itunes_search

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type LookupParamsBuilder struct {
	params map[string]string
}

func NewLookupParams() *LookupParamsBuilder {
	return &LookupParamsBuilder{
		params: make(map[string]string),
	}
}

func (lp *LookupParamsBuilder) ID(id int) *LookupParamsBuilder {
	if id > 0 {
		lp.params["id"] = strconv.Itoa(id)
	}
	return lp
}

func (lp *LookupParamsBuilder) IDs(ids []int) *LookupParamsBuilder {
	if len(ids) > 0 {
		idStrings := make([]string, len(ids))
		for i, id := range ids {
			idStrings[i] = strconv.Itoa(id)
		}
		lp.params["id"] = strings.Join(idStrings, ",")
	}
	return lp
}

func (lp *LookupParamsBuilder) UPC(upc string) *LookupParamsBuilder {
	if upc != "" {
		lp.params["upc"] = upc
	}
	return lp
}

func (lp *LookupParamsBuilder) EAN(ean string) *LookupParamsBuilder {
	if ean != "" {
		lp.params["ean"] = ean
	}
	return lp
}

func (lp *LookupParamsBuilder) ISRC(isrc string) *LookupParamsBuilder {
	if isrc != "" {
		lp.params["isrc"] = isrc
	}
	return lp
}

func (lp *LookupParamsBuilder) ISBN(isbn string) *LookupParamsBuilder {
	if isbn != "" {
		lp.params["isbn"] = isbn
	}
	return lp
}

func (lp *LookupParamsBuilder) AMGArtistID(amgArtistID string) *LookupParamsBuilder {
	if amgArtistID != "" {
		lp.params["amgArtistId"] = amgArtistID
	}
	return lp
}

func (lp *LookupParamsBuilder) AMGArtistIDs(amgArtistIDs []string) *LookupParamsBuilder {
	if len(amgArtistIDs) > 0 {
		lp.params["amgArtistId"] = strings.Join(amgArtistIDs, ",")
	}
	return lp
}

func (lp *LookupParamsBuilder) AMGAlbumID(amgAlbumID string) *LookupParamsBuilder {
	if amgAlbumID != "" {
		lp.params["amgAlbumId"] = amgAlbumID
	}
	return lp
}

func (lp *LookupParamsBuilder) AMGAlbumIDs(amgAlbumIDs []string) *LookupParamsBuilder {
	if len(amgAlbumIDs) > 0 {
		lp.params["amgAlbumId"] = strings.Join(amgAlbumIDs, ",")
	}
	return lp
}

func (lp *LookupParamsBuilder) AMGVideoID(amgVideoID string) *LookupParamsBuilder {
	if amgVideoID != "" {
		lp.params["amgVideoId"] = amgVideoID
	}
	return lp
}

func (lp *LookupParamsBuilder) Entity(entity string) *LookupParamsBuilder {
	if entity != "" {
		lp.params["entity"] = entity
	}
	return lp
}

func (lp *LookupParamsBuilder) Limit(limit int) *LookupParamsBuilder {
	if limit > 0 {
		if limit > MaxLimit {
			limit = MaxLimit
		}
		lp.params["limit"] = strconv.Itoa(limit)
	}
	return lp
}

func (lp *LookupParamsBuilder) Sort(sort string) *LookupParamsBuilder {
	if sort != "" {
		lp.params["sort"] = sort
	}
	return lp
}

func (lp *LookupParamsBuilder) Country(country string) *LookupParamsBuilder {
	if country != "" {
		lp.params["country"] = country
	}
	return lp
}

func (lp *LookupParamsBuilder) Build() map[string]string {
	return lp.params
}

func (c *Client) Lookup(params map[string]string) (*SearchResponse, error) {
	if len(params) == 0 {
		c.logger.Error("At least one lookup parameter is required")
		return nil, fmt.Errorf("at least one lookup parameter is required")
	}

	queryString := c.buildQueryString(params)
	c.logger.Debug("Making iTunes lookup request", zap.String("query", queryString))

	resp, err := c.baseClient.HTTP.R().
		SetQueryParams(params).
		Get(BaseLookupURL)

	if err != nil {
		c.logger.Error("Failed to make lookup request", zap.Error(err))
		return nil, fmt.Errorf("failed to make lookup request: %w", err)
	}

	if !resp.IsSuccess() {
		c.logger.Error("Lookup request failed", zap.Int("status", resp.StatusCode()))
		return nil, fmt.Errorf("lookup request failed with status %d", resp.StatusCode())
	}

	var lookupResponse SearchResponse
	err = json.Unmarshal(resp.Body(), &lookupResponse)
	if err != nil {
		c.logger.Error("Failed to unmarshal lookup response", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal lookup response: %w", err)
	}

	c.logger.Debug("Lookup request successful",
		zap.Int("results", lookupResponse.ResultCount),
		zap.Int("status", resp.StatusCode()))

	return &lookupResponse, nil
}
