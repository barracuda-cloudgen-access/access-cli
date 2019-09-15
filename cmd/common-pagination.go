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
package cmd

import (
	"fmt"
	"math"

	"github.com/spf13/cobra"
)

type pageable interface {
	SetPerPage(perPage *int64)
	SetPage(page *int64)
}

func initPaginationFlags(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations["pagination_flags_init"] = "yes"
	cmd.Flags().Int64("range-start", 1, "start of the range of items to return")
	cmd.Flags().Int64("range-end", 0, "end of the range of items to return (0 to return all items past range-start)")
}

func preRunFlagCheckPagination(cmd *cobra.Command, args []string) error {
	rangeStart, err := cmd.Flags().GetInt64("range-start")
	if err != nil {
		return err
	}
	if rangeStart < 1 {
		return fmt.Errorf("invalid range start %d", rangeStart)
	}

	rangeEnd, err := cmd.Flags().GetInt64("range-end")
	if err != nil {
		return err
	}
	if rangeEnd != 0 && rangeEnd <= rangeStart {
		return fmt.Errorf("invalid range end %d", rangeEnd)
	}

	return nil
}

// forAllPages is a pagination helper
// all int64 usage is because go-swagger really likes int64
// function `do` must return the total number of items
func forAllPages(cmd *cobra.Command, params pageable, do func() (int64, error)) (sliceStart, sliceEnd int64, err error) {
	if _, ok := cmd.Annotations["pagination_flags_init"]; !ok {
		panic("forAllPages called for command where pagination flags were not initialized. This is a bug!")
	}

	rangeStart, err := cmd.Flags().GetInt64("range-start")
	if err != nil {
		return 0, 0, err
	}
	rangeStart-- // user-facing values are 1-based
	rangeEnd, err := cmd.Flags().GetInt64("range-end")
	if err != nil {
		return 0, 0, err
	}
	if rangeEnd == 0 {
		rangeEnd = math.MaxInt64
	} else {
		rangeEnd-- // user-facing values are 1-based
	}

	perPage := int64(50)

	total := int64(math.MaxInt64)
	curPage := rangeStart / perPage
	sliceStart = rangeStart - curPage*perPage
	sliceEnd = rangeEnd - curPage*perPage
	lastPage := rangeEnd/perPage + perPage
	for ; curPage < lastPage && perPage*curPage < total; curPage++ {
		p := curPage + 1
		params.SetPage(&p)
		params.SetPerPage(&perPage)
		total, err = do()
		if err != nil {
			return 0, 0, err
		}
	}
	return sliceStart, sliceEnd, nil
}
