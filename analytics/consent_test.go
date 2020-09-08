package analytics

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newConsentTempfile(t *testing.T) string {
	file, err := ioutil.TempFile("", "section-cli-analytics-consent")
	if err != nil {
		t.FailNow()
	}
	return file.Name()
}

func TestConsentDetectsIfConsentNotRecorded(t *testing.T) {
	assert := assert.New(t)

	// Setup
	consentPath = newConsentTempfile(t)

	// Test
	assert.False(IsConsentRecorded())
}

func TestConsentPromptsForConsentIfConsentNotRecorded(t *testing.T) {
	assert := assert.New(t)
	assert.False(IsConsentRecorded())
}

func TestConsentSubmitNoopsIfNoConsent(t *testing.T) {
	assert := assert.New(t)
	var called bool
	ConsentGiven = false
	// Setup
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		body, _ := ioutil.ReadAll(r.Body)
		t.Logf("%s", body)
	}))
	HeapBaseURI = ts.URL

	// Invoke
	e := Event{
		Name: "CLI invoked",
		Properties: map[string]string{
			"Args":       "apps list",
			"Subcommand": "apps list",
			"Errors":     "",
		},
	}
	err := Submit(e)

	// Test
	assert.NoError(err)
	assert.False(called)
}
