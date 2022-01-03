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

var dumpCmd = &cobra.Command{
	Use:   "dump [path]",
	Short: "lists the contents of the bluge index",
	Long:  `The list command will list the contents of the Bluge index.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return fmt.Errorf("must specify path to index")
		}

		dir := index.NewFileSystemDirectory(args[0])

		snapshotIDs, err := dir.List(index.ItemKindSnapshot)
		if err != nil {
			return fmt.Errorf("error listing snapshots: %v", err)
		}
		for _, snapshotID := range snapshotIDs {
			fmt.Printf("snapshot: %d\n", snapshotID)
		}

		segmentIDs, err := dir.List(index.ItemKindSegment)
		if err != nil {
			return fmt.Errorf("error listing snapshots: %v", err)
		}
		for _, segmentID := range segmentIDs {
			fmt.Printf("segment: %d\n", segmentID)
		}
		config := index.DefaultConfig(args[0])
		fmt.Printf("Config: %v\n", config)
		w, _ := index.OpenWriter(config)
		r, _ := w.Reader()
		fields, _ := r.Fields()
		for _, field := range fields {
			fmt.Printf("Field: %s\n", field)
		}
		field := "_id"
		if len(args) > 1 {
			field = args[1]
		}
		it, _ := r.DictionaryIterator(field, nil, nil, nil)

		defer func() {
			_ = it.Close()
		}()

		//termCount := 0
		curr, err := it.Next()
		for err == nil && curr != nil {
			fmt.Printf("%v\n", it)
			curr, err = it.Next()
		}

		term := "xdata"
		if len(args) > 2 {
			term = args[2]
		}
		itr, _ := r.PostingsIterator([]byte(term), "content", true, true, true)
		fmt.Printf("itr %v\n", itr)
		for {
			x, _ := itr.Next()
			if x == nil {
				break
			}
			fmt.Printf("posting %v num %d\n", x, x.Number())
			filename := ""
			r.VisitStoredFields(x.Number(), func(field string, value []byte) bool {
				//fmt.Printf("fileld %s %s\n", field, string(value))
				if field == "filename" {
					filename = string(value)
				}
				return true
			})
			//fmt.Printf("n locs %d\n", len(x.Locations()))
			for _, loc := range x.Locations() {
				//fmt.Printf("%T %v\n", loc, loc)
				fmt.Printf("%s: start %d end %d\n", filename, loc.Start(), loc.End())
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(dumpCmd)
}
