package server

import (
	schema "CloudPricingAPI/schema"
	"CloudPricingAPI/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GraphQLQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// Setting up filters for the resource.
func queryFilters() (schema.ProductFilter, schema.PriceFilter) {
	ProductFilter := schema.ProductFilter{
		VendorName:    "azure",
		Region:        "westus",
		Service:       "Virtual Machines",
		ProductFamily: "Compute",
		AttributeFilters: []schema.AttributeFilter{
			{Key: "skuName", Value: "M128m"},
			{Key: "armSkuName", Value: "Standard_M128m"},
		},
	}
	PriceFilter := schema.PriceFilter{
		PurchaseOption: "Consumption",
	}
	return ProductFilter, PriceFilter
}

// Building query based on the product and price filters for the resource, so that we can call GraphQL API.
func buildQuery() GraphQLQuery {

	product, price := queryFilters()
	v := map[string]interface{}{}
	v["productFilter"] = product
	v["priceFilter"] = price

	query := `
		query($productFilter: ProductFilter!, $priceFilter: PriceFilter) {
			products(filter: $productFilter) {
			    attributes { key, value }
				prices(filter: $priceFilter) {
					priceHash
					USD
					purchaseOption
					startUsageAmount
					unit
				}
			}
		}
	`

	return GraphQLQuery{query, v}
}

// Adding header to the request
func AddHeaders(req *http.Request) {

	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Api-Key", utils.ApiKey)
}

func TestHandler(c *gin.Context) {

	query := buildQuery()
	reqBody, err := json.Marshal(query)
	if err != nil {
		fmt.Println("Error while generating request body")
		return
	}

	request, err := http.NewRequest("POST", utils.EndPoint+utils.Path, bytes.NewBuffer(reqBody))

	if err != nil {
		fmt.Println("Error while generating request")
		return
	}

	AddHeaders(request)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Error while sending API request")
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Invalid response from the API")
		return
	}
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(string(respBody)), &jsonMap)

	c.JSON(http.StatusOK, jsonMap)
}
