package gremgo

import (
	"encoding/json"
	"log"
)

// sortResponse ensures that data goes to the function that requested the data
func sortResponse(c *Client, msg []byte) (err error) {
	var data map[string]interface{}
	err = json.Unmarshal(msg, &data) // Unwrap message
	if err != nil {
		log.Fatal(err)
	}
	c.results[data["requestId"].(string)] = data // Prepare data for requester
	c.reschan[data["requestId"].(string)] <- 1   // Notify requester that data is ready
	close(c.reschan[data["requestId"].(string)]) // Close requester channel
	return
}
