package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gavsidhu/notion-hooks/internal/logging"
)

type NotionClient struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

func NewNotionClient(token string) *NotionClient {
	return &NotionClient{
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		baseURL: "https://api.notion.com/v1/",
		token:   token,
	}
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (c *NotionClient) SetToken(token string) {
	c.token = token
}

func (c *NotionClient) GetDatabase(ctx context.Context, databaseID string) (Database, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%sdatabases/%s", c.baseURL, databaseID), nil)
	if err != nil {
		return Database{}, err
	}

	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("error with request", err)
		return Database{}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading body", err)
		return Database{}, err
	}

	var database Database
	err = json.Unmarshal(body, &database)
	if err != nil {
		fmt.Println("error unmarshalling", err)
		return Database{}, err
	}

	return database, nil
}

func (c *NotionClient) GetPage(ctx context.Context, pageID string) (Page, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%spages/%s", c.baseURL, pageID), nil)
	if err != nil {
		return Page{}, err
	}

	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("error with request", err)
		return Page{}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading body", err)
		return Page{}, err
	}

	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		fmt.Println("error unmarshalling", err)
		return Page{}, err
	}

	return page, nil
}

func (c *NotionClient) GetAllDatabasePages(ctx context.Context, databaseID string) (*DatabaseQueryResponse, error) {
	var allPagesResults []Page
	hasMore := true
	nextCursor := ""

	// Rate limit to 3 requests per second. Notion API has a rate limit of 3 requests per second.
	rateLimiter := time.NewTicker(334 * time.Millisecond)
	defer rateLimiter.Stop()

	for hasMore {
		<-rateLimiter.C // Wait for the next tick before making the request

		var requestBody io.Reader
		var req *http.Request
		var err error

		if nextCursor != "" {
			jsonBody := []byte(fmt.Sprintf(`{"start_cursor": "%s"}`, nextCursor))
			requestBody = bytes.NewBuffer(jsonBody)
			req, err = http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%sdatabases/%s/query", c.baseURL, databaseID), requestBody)
		} else {
			req, err = http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%sdatabases/%s/query", c.baseURL, databaseID), nil)
		}
		if err != nil {
			logging.Logger.Error(fmt.Sprintf("Error creating request: %s", err))
			return nil, err
		}

		req.Header.Set("Notion-Version", "2022-06-28")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
		req.Header.Set("Content-Type", "application/json")

		res, err := c.httpClient.Do(req)
		if err != nil {
			logging.Logger.Error(fmt.Sprintf("Error with request: %s", err))
			return nil, err
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("error reading body", err)
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			logging.Logger.Error(fmt.Sprintf("Received non-OK HTTP status code: %d", res.StatusCode))
			return nil, fmt.Errorf("received non-OK HTTP status code: %d", res.StatusCode)
		}

		var pages DatabaseQueryResponse
		err = json.Unmarshal(body, &pages)
		if err != nil {
			logging.Logger.Error(fmt.Sprintf("Error unmarshalling response: %s", err))
			return nil, err
		}

		allPagesResults = append(allPagesResults, pages.Results...)
		hasMore = pages.HasMore
		nextCursor = pages.NextCursor
	}

	finalResponse := DatabaseQueryResponse{
		Object:  "list",
		Results: allPagesResults,
		HasMore: false,
	}

	return &finalResponse, nil
}

func (c *NotionClient) GetAllDatabasePageIDs(ctx context.Context, pages *DatabaseQueryResponse) ([]string, error) {

	var pageIDs []string

	for _, page := range pages.Results {
		pageIDs = append(pageIDs, page.ID)
	}

	return pageIDs, nil
}
