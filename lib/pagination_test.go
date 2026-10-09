package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaginateRecordLimit(t *testing.T) {
	tests := []struct {
		name          string
		total         int
		limits        []int
		wantPageSizes []int
		wantOffsets   []int
		wantItems     int
	}{
		{
			name:          "uncapped pagination",
			total:         205,
			wantPageSizes: []int{defaultPageSize, defaultPageSize, defaultPageSize},
			wantOffsets:   []int{0, 100, 200},
			wantItems:     205,
		},
		{
			name:          "cap within first page",
			total:         205,
			limits:        []int{25},
			wantPageSizes: []int{25},
			wantOffsets:   []int{0},
			wantItems:     25,
		},
		{
			name:          "cap across pages with a smaller final request",
			total:         205,
			limits:        []int{125},
			wantPageSizes: []int{defaultPageSize, 25},
			wantOffsets:   []int{0, 100},
			wantItems:     125,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pageSizes, offsets []int
			fetch := func(pageLimit, offset int) ([]int, int, error) {
				pageSizes = append(pageSizes, pageLimit)
				offsets = append(offsets, offset)

				count := min(pageLimit, tt.total-offset)
				if count < 0 {
					count = 0
				}
				page := make([]int, count)
				for i := range page {
					page[i] = offset + i
				}
				return page, len(page), nil
			}

			pages, err := paginate(context.Background(), fetch, tt.limits...)
			require.NoError(t, err)
			assert.Equal(t, tt.wantPageSizes, pageSizes)
			assert.Equal(t, tt.wantOffsets, offsets)

			var items []int
			for _, page := range pages {
				items = append(items, page...)
			}
			assert.Len(t, items, tt.wantItems)
		})
	}
}

func TestPaginateRejectsNonPositiveRecordLimit(t *testing.T) {
	for _, limit := range []int{0, -1} {
		_, err := paginate(context.Background(), func(pageLimit, offset int) ([]int, int, error) {
			return nil, 0, nil
		}, limit)
		require.EqualError(t, err, "limit must be a positive integer")
	}
}

func TestGetTopicsRecordLimitCapsReturnedPage(t *testing.T) {
	requestedLimit := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedLimit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
		topics := make([]Topic, 8)
		for i := range topics {
			topics[i].ID = strconv.Itoa(i)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(TopicsResponse{Topics: topics})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	responses, err := client.GetTopics(context.Background(), 3)
	require.NoError(t, err)
	require.Len(t, responses, 1)
	assert.Equal(t, 3, requestedLimit)
	assert.Len(t, responses[0].Topics, 3)
}

func TestGetJobsByDateRecordLimit(t *testing.T) {
	var pageSizes, offsets []int
	baseTime := time.Now().UTC()
	createdAt := baseTime.Format(dateFormat)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageLimit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		pageSizes = append(pageSizes, pageLimit)
		offsets = append(offsets, offset)

		jobs := make([]Job, pageLimit)
		for i := range jobs {
			jobs[i].ID = strconv.Itoa(offset + i)
			jobs[i].CreatedAt = createdAt
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(JobsResponse{Jobs: jobs})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	startDate := baseTime.Add(-time.Hour)
	endDate := baseTime.Add(time.Hour)
	responses, err := client.GetJobsByDate(context.Background(), startDate, endDate, 125)
	require.NoError(t, err)
	assert.Equal(t, []int{defaultPageSize, 25}, pageSizes)
	assert.Equal(t, []int{0, 100}, offsets)

	total := 0
	for _, response := range responses {
		total += len(response.Jobs)
	}
	assert.Equal(t, 125, total)
}

func TestGetJobsRecordLimitPreservesAgeCutoff(t *testing.T) {
	requests := 0
	createdAt := time.Now().UTC().Add(-30 * 24 * time.Hour).Format(dateFormat)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		jobs := make([]Job, defaultPageSize)
		for i := range jobs {
			jobs[i].ID = strconv.Itoa(i)
			jobs[i].CreatedAt = createdAt
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(JobsResponse{Jobs: jobs})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	responses, err := client.GetJobs(context.Background(), 7, 150)
	require.NoError(t, err)
	assert.Equal(t, 1, requests, "the age cutoff should stop pagination after the old page")
	assert.Len(t, responses, 1)
	assert.Len(t, responses[0].Jobs, defaultPageSize)
}
