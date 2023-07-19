package client // import "github.com/docker/docker/client"

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/errdefs"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestContainerExportError(t *testing.T) {
	client := &Client{
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	}
	_, err := client.ContainerExport(context.Background(), "nothing")
	assert.Check(t, is.ErrorType(err, errdefs.IsSystem))
}

func TestContainerExport(t *testing.T) {
	expectedURL := "/containers/container_id/export"
	client := &Client{
		client: newMockClient(func(r *http.Request) (*http.Response, error) {
			if !strings.HasPrefix(r.URL.Path, expectedURL) {
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, r.URL)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte("response"))),
			}, nil
		}),
	}
	body, err := client.ContainerExport(context.Background(), "container_id")
	if err != nil {
		t.Fatal(err)
	}
	defer body.Close()
	content, err := io.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "response" {
		t.Fatalf("expected response to contain 'response', got %s", string(content))
	}
}
