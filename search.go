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
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/blugelabs/bluge"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [path]",
	Short: "lists the contents of the bluge index",
	Long:  `The list command will list the contents of the Bluge index.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return fmt.Errorf("must specify path to index")
		}

		config := bluge.DefaultConfig(args[0])
		r, _ := bluge.OpenReader(config)
		term := strings.ToLower(args[1])
		qcontent := bluge.NewTermQuery(term).SetField("content")
		qcomments := bluge.NewTermQuery(term).SetField("comments")
		q := bluge.NewBooleanQuery()
		q.AddShould(qcontent)
		q.AddShould(qcomments)
		sr := bluge.NewAllMatches(q).IncludeLocations()
		res, _ := r.Search(context.Background(), sr)
		for {
			match, err := res.Next()
			if match == nil {
				break
			}

			filename := ""
			var linemap []uint32
			match.VisitStoredFields(func(field string, value []byte) bool {
				if field == "linemap" {
					linemap = make([]uint32, len(value)/4)
					for i := 0; i < len(value)/4; i++ {
						linemap[i] = binary.BigEndian.Uint32(value[i*4:])
					}
				}
				if field == "filepath" {
					filename = string(value)
				}
				return true
			})

			f, err := os.Open(filename)
			if err != nil {
				fmt.Printf("can not open %s - %s\n", filename, err.Error())
				continue
			}
			data, _ := io.ReadAll(f)
			f.Close()
			last_line := -1
			for _, ftl := range match.FieldTermLocations {
				fmt.Printf("%s: %d\n", ftl.Field, ftl.Location.Start)
				idx := sort.Search(len(linemap)/2, func(i int) bool { return uint32(ftl.Location.Start) >= linemap[2*i] })
				off := 0
				if idx < len(linemap)/2 {
					off = int(linemap[2*idx+1])
				}
				line := 1
				start := 0
				for off < ftl.Location.Start {
					if data[off] == '\n' {
						start = off + 1
						line++
					}
					off++
				}
				end := start
				for end < len(data) && data[end] != '\n' {
					end += 1
				}
				if line != last_line {
					fmt.Printf("%s: %d: %s\n", filename, line, data[start:end])
					last_line = line
				}
			}

			if err != nil {
				break
			}
		}

		//pi, _ := r.PostingsIterator()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)
}
