package midjourneyapi

import (
	"fmt"
	"os"
	"testing"
)

var apiKey = os.Getenv("API_KEY")

func TestClient_Describe(t *testing.T) {
	client := NewClient(apiKey)
	f, err := os.Open("testdata/example.jpg")
	if err != nil {
		t.Fatal(err)
	}

	result, err := client.Describe(f)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(result)
}
