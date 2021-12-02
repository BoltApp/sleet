package request

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const xmlHeader string = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
const gatewayRequestElem string = document_base // "<gatewayRequest>"

func TestNewGatewayRequest(t *testing.T) {
	request := NewGatewayRequest()

	request2 := GatewayRequest{}
	request2.ToXMLString()

	xml := request.ToXMLString()
	assert.True(t, strings.HasPrefix(xml, xmlHeader), "Expect xml header")
	assert.True(t, strings.HasPrefix(xml, xmlHeader+"\n<"+gatewayRequestElem+">"), "Expect xml header and start root element")
	assert.True(t, strings.HasSuffix(xml, "</"+gatewayRequestElem+">"), "Expect end root element")
	assert.True(t, strings.Contains(xml, "<version>"+version_no+"</version>"))
}

func TestGatewayRequestSet(t *testing.T) {
	request := NewGatewayRequest()
	request.Set("a", "1")
	request.SetInt("b", 2)
	xml := request.ToXMLString()
	assert.True(t, strings.Contains(xml, "<a>1</a>"), "Expect a=1 value")
	assert.True(t, strings.Contains(xml, "<b>2</b>"), "Expect b=2 value")
}

func TestGatewayRequestGet(t *testing.T) {
	request := NewGatewayRequest()
	request.Set("a", "1")
	assert.Equal(t, "1", request.Get("a"), "Expect 1 for key a")
	request.SetInt("b", 2)
	assert.Equal(t, "2", request.Get("b"), "Expect 2 for key b")
	// Get value for invalid key
	assert.Equal(t, "", request.Get("invalid_key"), "Expect empty for invalid_key")
	// Test GetInt values
	assert.Equal(t, -1, request.GetInt("invalid_key"), "Expect -1")
	assert.Equal(t, 0, request.GetIntOrDefault("invalid_key", 0), "Expect 0")
	assert.Equal(t, 3, request.GetIntOrDefault("invalid_key", 3), "Expect 3")
	// Valid int
	request.Set("nr", "1")
	assert.Equal(t, 1, request.GetInt("nr"), "Expect 1")
	request.Set("nr", " 1")
	assert.Equal(t, 1, request.GetInt("nr"), "Expect 1")
	request.Set("nr", "1 ")
	assert.Equal(t, 1, request.GetInt("nr"), "Expect 1")
	// Invalid int
	request.Set("nr", "x")
	assert.Equal(t, -1, request.GetInt("nr"), "Expect -1")
	assert.Equal(t, 5, request.GetIntOrDefault("nr", 5), "Expect 5")
}

func TestGatewayRequestXmlEscape(t *testing.T) {
	request := NewGatewayRequest()
	request.Set("a", "<")
	request.Set("b", ">")
	request.Set("c", "'")
	request.Set("d", "\"")
	request.Set("e", "&")
	xml := request.ToXMLString()
	assert.True(t, strings.Contains(xml, "<a>&lt;</a>"), "Expect a=&lt; value")
	assert.True(t, strings.Contains(xml, "<b>&gt;</b>"), "Expect b=&gt; value")
	assert.True(t, strings.Contains(xml, "<c>&#39;</c>"), "Expect c=&#39; value")
	assert.True(t, strings.Contains(xml, "<d>&#34;</d>"), "Expect d=&#34; value")
	assert.True(t, strings.Contains(xml, "<e>&amp;</e>"), "Expect e=&amp; value")
}
