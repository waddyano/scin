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
	"lexers"
	"log"
	"os"
	"path/filepath"

	"github.com/blugelabs/bluge"

	"github.com/spf13/cobra"
	"github.com/zeebo/xxh3"
)

func buildLineMap(input []byte) []uint32 {
	m := make([]uint32, 0, len(input)/1024)
	line := 1
	last_off := 0
	for off, by := range input {
		if by == '\n' {
			if off > last_off+1024 {
				m = append(m, uint32(off), uint32(line))
				last_off = off
			}
			line += 1
		}
	}

	return m
}

const MAX_BATCH = 20

type batchHandler struct {
	writer     *bluge.Writer
	filenames  []string
	addedCount int
}

func (bh *batchHandler) addFiles() error {
	if len(bh.filenames) == 0 {
		return nil
	}
	batch := bluge.NewBatch()
	count := 0
	for _, filename := range bh.filenames {
		if !lexers.CanLex(filename) {
			continue
		}
		f, err := os.Open(filename)
		if err != nil {
			continue
		}
		count++
		data, _ := io.ReadAll(f)
		f.Close()

		h128 := xxh3.Hash128(data)
		id := fmt.Sprintf("%8x%8x", h128.Hi, h128.Lo)

		idq := bluge.NewTermQuery(id).SetField("_id")
		sr := bluge.NewAllMatches(idq).IncludeLocations()
		r, _ := bh.writer.Reader()
		res, _ := r.Search(context.Background(), sr)
		r.Close()
		match, err := res.Next()
		if match != nil {
			//existing := ""
			//match.VisitStoredFields(func(field string, value []byte) bool {
			//	if field == "filename" {
			//		existing = string(value)
			//	}
			//	return true
			//})
			//fmt.Printf("existing match %s %s %v\n", filename, existing, err)
			new_id := fmt.Sprintf("%s%8x", id, xxh3.Hash([]byte(filename)))
			doc := bluge.NewDocument(new_id)
			doc.AddField(bluge.NewKeywordField("filepath", filename).StoreValue())
			doc.AddField(bluge.NewKeywordField("filename", filepath.Base(filename)).StoreValue())
			doc.AddField(bluge.NewKeywordField("duplicate", id).StoreValue())
			batch.Update(doc.ID(), doc)
			continue
		}

		//fmt.Printf("%s %d bytes %d lines\n", filename, len(data), len(buildLineMap(data))/2)
		lineMap := buildLineMap(data)
		buf := make([]byte, len(lineMap)*4)
		for i, d := range lineMap {
			binary.BigEndian.PutUint32(buf[i*4:], d)
		}

		//fmt.Printf("%s\n", h128)
		doc := bluge.NewDocument(id)
		fieldObj := bluge.NewTextFieldBytes("content", data).WithAnalyzer(SourceAnalyzer{filename: filename, comments: false}).SearchTermPositions()
		doc.AddField(fieldObj)
		fieldObj = bluge.NewTextFieldBytes("comments", data).WithAnalyzer(SourceAnalyzer{filename: filename, comments: true}).SearchTermPositions()
		doc.AddField(fieldObj)
		doc.AddField(bluge.NewKeywordField("filepath", filename).StoreValue())
		doc.AddField(bluge.NewKeywordField("filename", filepath.Base(filename)).StoreValue())
		doc.AddField(bluge.NewStoredOnlyField("linemap", buf))
		batch.Update(doc.ID(), doc)
	}
	bh.addedCount += count
	err := bh.writer.Batch(batch)
	bh.filenames = bh.filenames[:0]
	return err
}

func (bh *batchHandler) addDirectory(dirname string) error {
	//fmt.Printf("process dir %s\n", dirname)
	entries, _ := os.ReadDir(dirname)
	for _, de := range entries {
		ste, _ := de.Info()
		file := filepath.Join(dirname, de.Name())
		if ste.IsDir() {
			bh.addDirectory(file)
		} else {
			bh.filenames = append(bh.filenames, file)
			if len(bh.filenames) == MAX_BATCH {
				err := bh.addFiles()
				if err != nil {
					log.Printf("error updating document: %v", err)
				}
			}
		}
	}
	return nil
}

var createCmd = &cobra.Command{
	Use:   "index [path] [file]",
	Short: "adds file to the scin index",
	Long:  `The index command will scan directories and add their data into the scin index.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) < 2 {
			return fmt.Errorf("must specify db and path to index")
		}

		config := bluge.DefaultConfig(args[0])
		fmt.Printf("%v\n", config)
		writer, err := bluge.OpenWriter(config)

		defer func() {
			err = writer.Close()
			if err != nil {
				log.Fatalf("error closing writer: %v", err)
			}
		}()

		name := args[1]
		st, err := os.Stat(name)
		bh := batchHandler{filenames: make([]string, 0, MAX_BATCH), writer: writer}
		if err == nil && st.IsDir() {
			err = bh.addDirectory(name)
			bh.addFiles()
		} else {
			bh.filenames = append(bh.filenames, name)
			bh.addFiles()
		}
		fmt.Printf("added %d\n", bh.addedCount)
		if err != nil {
			log.Printf("error updating document: %v", err)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}
