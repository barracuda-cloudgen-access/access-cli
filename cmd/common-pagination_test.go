// Package cmd implements fyde-cli commands
package cmd

/*
Copyright © 2019 Fyde, Inc. <hello@fyde.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/nbio/st"
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
		{
			rangeStart: 1, rangeEnd: 10, pageSize: 5,
			data:            []int{},
			expectedFetches: 1,
			expected:        []int{},
		},
		{
			rangeStart: 1, rangeEnd: 10, pageSize: 5,
			data:            genData[0:20],
			expectedFetches: 2,
			expected:        genData[0:9],
		},
		{
			rangeStart: 1, rangeEnd: 10, pageSize: 50,
			data:            genData[0:20],
			expectedFetches: 1,
			expected:        genData[0:9],
		},
		{
			rangeStart: 1, rangeEnd: 10, pageSize: 10,
			data:            genData[0:20],
			expectedFetches: 1,
			expected:        genData[0:9],
		},
		{
			rangeStart: 1, rangeEnd: 0, pageSize: 7,
			data:            genData[0:20],
			expectedFetches: 3,
			expected:        genData[0:20],
		},
		{
			rangeStart: 99, rangeEnd: 104, pageSize: 10,
			data:            genData,
			expectedFetches: 2,
			expected:        genData[98:103],
		},
		{
			rangeStart: 10, rangeEnd: 0, pageSize: 7,
			data:            genData[0:20],
			expectedFetches: 2,
			expected:        genData[9:20],
		},
		{
			rangeStart: 1, rangeEnd: 0, pageSize: 34,
			data:            genData,
			expectedFetches: 7,
			expected:        genData,
		},
		{
			rangeStart: 1, rangeEnd: 11, pageSize: 10,
			data:            genData,
			expectedFetches: 1,
			expected:        genData[0:10],
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
		st.Expect(t, err, nil)
		result = result[cutStart:cutEnd]

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Error("Result should be", testCase.expected, "but was", result)
		}
		if fPageable.fetchCount != testCase.expectedFetches {
			t.Error("Expected result in", testCase.expectedFetches, "fetches but", fPageable.fetchCount, "were made")
		}
	}
}
