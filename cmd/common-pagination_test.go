package cmd

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/spf13/cobra"
)

type fakePageable struct {
	perPage    int64
	page       int64
	data       []int
	fetchCount int
}

// SetPerPage sets how many fake results to return per fakeFetch
func (f *fakePageable) SetPerPage(perPage *int64) {
	f.perPage = *perPage
}

// SetPerPage sets the current page of fake results
func (f *fakePageable) SetPage(page *int64) {
	f.page = *page
}

func (f *fakePageable) fakeFetch() []int {
	f.fetchCount++
	start := (f.page - 1) * f.perPage
	if start >= int64(len(f.data)) {
		return []int{}
	}
	end := start + f.perPage
	if end > int64(len(f.data)) {
		end = int64(len(f.data))
	}
	return f.data[start:end]
}

func TestPagination(t *testing.T) {
	type tcase struct {
		rangeStart, rangeEnd, pageSize int
		data                           []int
		expectedFetches                int
		expected                       []int
	}

	genData := make([]int, 234)
	for i := 0; i < 234; i++ {
		genData[i] = i + 1
	}

	testCases := []tcase{
		tcase{
			rangeStart: 1, rangeEnd: 10, pageSize: 5,
			data:            []int{},
			expectedFetches: 1,
			expected:        []int{},
		},
		tcase{
			rangeStart: 1, rangeEnd: 10, pageSize: 5,
			data:            genData[0:20],
			expectedFetches: 2,
			expected:        genData[0:9],
		},
		tcase{
			rangeStart: 1, rangeEnd: 10, pageSize: 50,
			data:            genData[0:20],
			expectedFetches: 1,
			expected:        genData[0:9],
		},
		tcase{
			rangeStart: 1, rangeEnd: 10, pageSize: 10,
			data:            genData[0:20],
			expectedFetches: 1,
			expected:        genData[0:9],
		},
		tcase{
			rangeStart: 1, rangeEnd: 0, pageSize: 7,
			data:            genData[0:20],
			expectedFetches: 3,
			expected:        genData[0:20],
		},
		tcase{
			rangeStart: 99, rangeEnd: 104, pageSize: 10,
			data:            genData,
			expectedFetches: 2,
			expected:        genData[98:103],
		},
		tcase{
			rangeStart: 10, rangeEnd: 0, pageSize: 7,
			data:            genData[0:20],
			expectedFetches: 2,
			expected:        genData[9:20],
		},
		tcase{
			rangeStart: 1, rangeEnd: 0, pageSize: 34,
			data:            genData,
			expectedFetches: 7,
			expected:        genData,
		},
	}

	for _, testCase := range testCases {
		fakeCmd := &cobra.Command{}
		initPaginationFlags(fakeCmd)
		fakeCmd.Flags().Set("range-start", strconv.Itoa(testCase.rangeStart))
		fakeCmd.Flags().Set("range-end", strconv.Itoa(testCase.rangeEnd))

		fPageable := &fakePageable{
			data: testCase.data,
		}
		global.FetchPerPage = testCase.pageSize

		result := []int{}
		cutStart, cutEnd, err := forAllPages(fakeCmd, fPageable, func() (int, int64, error) {
			r := fPageable.fakeFetch()
			result = append(result, r...)
			return len(r), int64(len(fPageable.data)), nil
		})
		if err != nil {
			t.Error(err)
		}
		result = result[cutStart:cutEnd]

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Error("Result should be", testCase.expected, "but was", result)
		}
		if fPageable.fetchCount != testCase.expectedFetches {
			t.Error("Expected result in", testCase.expectedFetches, "fetches but", fPageable.fetchCount, "were made")
		}
	}
}
