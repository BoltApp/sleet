package response

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const xmlHeader string = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"

func TestNewGatewayResponse(t *testing.T) {
	response := NewGatewayResponse()
	fmt.Println(response)
	assert.Equal(t, "", response.Get("invalid_key"), "Expect empty for invalid_key")
	assert.Equal(t, 2, response.GetIntOrDefault("invalid_key", 2), "Expect 2 value ")
	assert.Equal(t, -1, response.GetInt("invalid_key"), "Expect -1 value ")
	assert.Equal(t, -1, response.GetResponseCode(), "Expect -1 value ")
}

func TestParseResponse(t *testing.T) {
	testInput := "<gatewayResponse><par1>a</par1><par2>b</par2><par3></par3></gatewayResponse>"
	response := NewGatewayResponse()
	response.SetFromXML(testInput)
	assert.Equal(t, "a", response.Get("par1"), "Expect a for par1")
	assert.Equal(t, "b", response.Get("par2"), "Expect b for par2")
	assert.Equal(t, "", response.Get("par3"), "Expect par3 empty")
	assert.Equal(t, "", response.Get("par4"), "Expect par4 empty")
}

func TestParseInvalidResponseRoot(t *testing.T) {
	testInput := "<root><par1>a</par1></root>"
	response := NewGatewayResponse()
	response.SetFromXML(testInput)
	fmt.Println(response)
	assert.Equal(t, "4", response.Get(RESPONSE_CODE), "Expect responseCode: 4")
	assert.Equal(t, "400", response.Get(REASON_CODE), "Expect reasonCode: 400")
	assert.Equal(t, "invalid xml", response.Get(EXCEPTION), "Expect exception: invalid xml")
}

func TestParseInvalidResponseParElement(t *testing.T) {
	testInput := "<gatewayResponse><par1>a</par2></gatewayResponse>"
	response := NewGatewayResponse()
	response.SetFromXML(testInput)
	fmt.Println(response)
	assert.Equal(t, "4", response.Get(RESPONSE_CODE), "Expect responseCode: 4")
	assert.Equal(t, "400", response.Get(REASON_CODE), "Expect reasonCode: 400")
	assert.True(t, strings.HasPrefix(response.Get(EXCEPTION), "invalid xml"), "Expect invalid xml prefix")
}

func TestParseXMLHeader(t *testing.T) {
	testInput := xmlHeader + "\n<gatewayResponse><par1>x</par1></gatewayResponse>"
	response := NewGatewayResponse()
	response.SetFromXML(testInput)
	fmt.Println(response)
	assert.Equal(t, "x", response.Get("par1"), "Expect x for par1")
}

func TestParseXMLWithSpaces(t *testing.T) {
	testInput := xmlHeader + " \n <gatewayResponse>\n\n<par1> s     \t</par1>\n  <par2> 5\t \n </par2>  </gatewayResponse>    \n\n"
	response := NewGatewayResponse()
	response.SetFromXML(testInput)
	fmt.Println(response)
	assert.Equal(t, "s", response.Get("par1"), "Expect s for par1")
	assert.Equal(t, 5, response.GetInt("par2"), "Expect 5 for par2")
}
