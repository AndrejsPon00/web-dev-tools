package pp

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/AndrejsPon00/web-dev-tools/backend/module"
	"github.com/goccy/go-json"
)

const (
	BASE_PP_URL       = "https://apipub.pp.lv/lv/api_user/v1/search/lots?"
	BASE_IMAGE_URL    = "https://img.pp.lv/"
	DEFAULT_IMAGE_URL = "https://st2.depositphotos.com/38069286/47731/v/450/depositphotos_477315358-stock-illustration-picture-isolated-background-gallery-symbol.jpg"
	BASE_SEARCH_QUERY = "&query="
	BASE_PAGE_QUERY   = "currentPage="
	FILTER_MIN_PRICE  = "minPrice="
	FILTER_MAX_PRICE  = "maxPrice="
	POSTS_IN_ONE_PAGE = 20
)

func ScrapPosts(input string, currentPage uint8, filter *module.Filter, wg *sync.WaitGroup, paginationChan chan *module.Pagination, result chan *module.PreviewPost, errorChan chan error) {
	defer wg.Done()

	url := getFullURL(input, currentPage, filter)
	log.Println(url)
	rawResponse, err := FetchResponse(url)
	if err != nil {
		errorChan <- err
		return
	}

	response, err := DecodeResponse(rawResponse)
	if err != nil {
		errorChan <- err
		return
	}

	SendPaginationPostsToChannel(currentPage, response, paginationChan)
	SendPreviewPostsToChannel(response, result)
}

func FetchResponse(input string) ([]byte, error) {
	resp, err := http.Get(input)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	return body, err
}

func DecodeResponse(response []byte) (*Response, error) {
	if isNil(response) {
		return nil, nil
	}

	res := &Response{}
	err := json.Unmarshal(response, res)
	return res, err
}

func SendPreviewPostsToChannel(response *Response, resultChan chan *module.PreviewPost) {
	if isNil(response) {
		return
	}

	for _, item := range response.Content.Data {
		if containsAds(item) {
			continue
		}

		resultChan <- &module.PreviewPost{
			Title:        item.Title,
			URL:          item.RedirectURL,
			PreviewImage: getPreviewImageURL(item),
			Price:        getPrice(item),
		}
	}
}

func SendPaginationPostsToChannel(currentPage uint8, response *Response, paginationChan chan *module.Pagination) {
	paginationChan <- &module.Pagination{
		Source:  module.SOURCE_PP,
		HasNext: hasNextPage(currentPage, response),
	}
}

func addFilter(url string, f *module.Filter) string {
	if f.PriceMin > 0 {
		url = fmt.Sprintf("%s&%s%d", url, FILTER_MIN_PRICE, f.PriceMin)
	}

	if f.PriceMax > 0 {
		url = fmt.Sprintf("%s&%s%d", url, FILTER_MAX_PRICE, f.PriceMax)
	}

	return url
}

func hasNextPage(currentPage uint8, response *Response) bool {
	if isNil(response) {
		return false
	}
	var totalPages float32 = float32(response.Content.ItemsCount) / POSTS_IN_ONE_PAGE
	return totalPages > float32(currentPage)
}

func getPrice(item *Data) string {
	if len(item.Prices) > 0 {
		return fmt.Sprintf("%s €", item.Prices[0].Value)
	}
	return "0 €"
}

func getPreviewImageURL(item *Data) string {
	if item.FileCount > 0 {
		return fmt.Sprintf("%s%s", BASE_IMAGE_URL, item.Files[0].File.Versions.OriginalFile.Path)
	}
	return DEFAULT_IMAGE_URL
}

func containsAds(item *Data) bool {
	return len(item.AdFilters) == 0
}

func encodeSpacesForURL(query string) string {
	return strings.ReplaceAll(query, " ", "%20")
}

func getFullURL(query string, pageNumber uint8, filter *module.Filter) string {
	if pageNumber == 0 {
		pageNumber = 1
	}

	url := fmt.Sprintf("%s%s%d%s%s", BASE_PP_URL, BASE_PAGE_QUERY, pageNumber, BASE_SEARCH_QUERY, encodeSpacesForURL(query))
	if filter != nil {
		url = addFilter(url, filter)
	}
	return url
}

func isNil(value interface{}) bool {
	return value == nil || reflect.ValueOf(value).IsNil()
}
