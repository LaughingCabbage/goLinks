/*
 *Copyright 2018-2019 Kevin Gentile
 *
 *Licensed under the Apache License, Version 2.0 (the "License");
 *you may not use this file except in compliance with the License.
 *You may obtain a copy of the License at
 *
 *http://www.apache.org/licenses/LICENSE-2.0
 *
 *Unless required by applicable law or agreed to in writing, software
 *distributed under the License is distributed on an "AS IS" BASIS,
 *WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *See the License for the specific language governing permissions and
 *limitations under the License.
 */

package blockmap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

var tmpDir string

func TestMain(m *testing.M) {
	//Generate initial blockmap from the test root
	tmpDir, _ = ioutil.TempDir(".", "test")
	cwd, _ := os.Getwd()
	tmpDir = cwd + string(os.PathSeparator) + tmpDir
	fmt.Println("tmpdir: " + tmpDir)

	for i := 0; i < 2; i++ {
		iStr := strconv.Itoa(i)
		log.Println("Creating Archive " + iStr)
		subTmpdir, err := ioutil.TempDir(tmpDir, "test"+iStr)
		if err != nil {
			panic(err)
		}
		//Write random data to files
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for j := 0; j < 1; j++ {
			jStr := strconv.Itoa(j)
			buff := make([]byte, 100)
			r.Read(buff)
			tmpfile, err := ioutil.TempFile(subTmpdir, "file"+jStr)
			if err != nil {
				panic(err)
			}

			if _, err := tmpfile.Write(buff); err != nil {
				panic(err)
			}
			if err := tmpfile.Close(); err != nil {
				panic(err)
			}
		}
	}

	result := m.Run()
	fmt.Println("Tearing Down")
	os.RemoveAll(tmpDir)
	os.Exit(result)
}

func TestBlockMap_New(t *testing.T) {
	path := tmpDir
	b := New(path)
	if b == nil {
		t.Error(errors.New("blockmap: failed to make new blockmap | " + path))
	}
}

func TestBlockMap_Generate(t *testing.T) {
	t.Log(tmpDir)
	b := New(tmpDir)
	if err := b.Generate(); err != nil {
		t.Error(err, tmpDir)
	}
}

func TestBlockMap_PrintBlockMap(t *testing.T) {
	t.Skip()
	b := New(tmpDir)
	if err := b.Generate(); err != nil {
		b.PrintBlockMap()
		t.Error(err)
	}

}

func TestEqual(t *testing.T) {
	//Initialize A
	a := New(tmpDir)
	fmt.Println(tmpDir)
	if err := a.Generate(); err != nil {
		t.Error(err)
	}
	//Initialize B
	b := New(tmpDir)
	if err := b.Generate(); err != nil {
		t.Error(err)
	}

	if !Equal(a, b) {
		t.Error(errors.New("blockmap: failed to evaluate equal blockmaps"))
	}

	c := New(tmpDir)
	if Equal(a, c) {
		t.Error(errors.New("blockmap: evaluated equality in unequal blockmaps"))
	}
}

func TestBlockMap_IO(t *testing.T) {
	for i := 0; i < 10; i++ {
		b := New(tmpDir)
		if err := b.Generate(); err != nil {
			t.Error(err)
		}

		//Save the blockmap
		if err := b.Save(b.Root); err != nil {
			t.Error(err)
		}

		//Load the blockmap in a new structure
		a := New(tmpDir)
		fmt.Println("loading link file at: " + b.Root)
		if err := a.Load(b.Root); err != nil {
			t.Error(err)
		}

		//Ensure both maps are equal
		if !Equal(b, a) {
			t.Error(errors.New("BlockMapIO failed to reload map"))
		}

		//Re-generate the link and validate it with the current file
		if err := a.Generate(); err != nil {
			t.Error(err)
		}
	}
}

func TestBlockMap_JSON(t *testing.T) {
	b1 := New(tmpDir)
	if err := b1.Generate(); err != nil {
		t.Error(err)
	}

	b1JSON, err := json.Marshal(b1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("b1 JSON")
	fmt.Println(string(b1JSON))

	b2 := New(tmpDir)
	if err := b2.Generate(); err != nil {
		t.Error(err)
	}

	b2JSON, err := json.Marshal(b2)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("b2 JSON")
	fmt.Println(string(b2JSON))

	if !bytes.Equal(b1JSON, b2JSON) {
		t.Error("blockmap: json input and output are not equal")
	}
}
