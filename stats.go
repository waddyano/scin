//  Copyright (c) 2020 The Bluge Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"github.com/blugelabs/bluge/index"

	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats [path]",
	Short: "lists the contents of the bluge index",
	Long:  `The list command will list the contents of the Bluge index.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return fmt.Errorf("must specify path to index")
		}

		config := index.DefaultConfig(args[0])
		w, _ := index.OpenWriter(config)
		r, _ := w.Reader()
		fields, _ := r.Fields()
		for _, field := range fields {
			fmt.Printf("Field: %s\n", field)
			s, _ := r.CollectionStats(field)
			fmt.Printf("Doc count %d %d\n", s.TotalDocumentCount(), s.DocumentCount())
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(statsCmd)
}
