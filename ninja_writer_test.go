// Copyright 2014 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blueprint

import (
	"bytes"
	"strings"
	"testing"
)

func ck(err error) {
	if err != nil {
		panic(err)
	}
}

var ninjaWriterTestCases = []struct {
	input  func(w *ninjaWriter)
	output string
}{
	{
		input: func(w *ninjaWriter) {
			ck(w.Comment("foo"))
		},
		output: "# foo\n",
	},
	{
		input: func(w *ninjaWriter) {
			ck(w.Pool("foo"))
		},
		output: "pool foo\n",
	},
	{
		input: func(w *ninjaWriter) {
			ck(w.Rule("foo"))
		},
		output: "rule foo\n",
	},
	{
		input: func(w *ninjaWriter) {
			ck(w.Build("foo comment", "foo", testNinjaStrings("o3", "o4"),
				testNinjaStrings("io3", "io4"), testNinjaStrings("e3", "e4"),
				testNinjaStrings("i3", "i4"), testNinjaStrings("oo3", "oo4"),
				testNinjaStrings("v3", "v4"), []string{"o1", "o2"},
				[]string{"io1", "io2"}, []string{"e1", "e2"},
				[]string{"i1", "i2"}, []string{"oo1", "oo2"},
				[]string{"v1", "v2"}, nil))
		},
		output: "# foo comment\nbuild o1 o2 o3 o4 | io1 io2 io3 io4: foo e1 e2 e3 e4 | i1 i2 i3 i4 || oo1 oo2 $\n        oo3 oo4 |@ v1 v2 v3 v4\n",
	},
	{
		input: func(w *ninjaWriter) {
			ck(w.Build("foo comment", "foo",
				testNinjaStrings(strings.Repeat("o", lineWidth)),
				nil,
				testNinjaStrings(strings.Repeat("i", lineWidth)),
				nil, nil, nil, nil, nil, nil, nil, nil, nil, nil))
		},
		output: "# foo comment\nbuild $\n        " + strings.Repeat("o", lineWidth) + ": $\n        foo $\n        " + strings.Repeat("i", lineWidth) + "\n",
	},
	{
		input: func(w *ninjaWriter) {
			ck(w.Default(nil, testNinjaStrings("foo"), []string{"bar"}))
		},
		output: "default bar foo\n",
	},
	{
		input: func(w *ninjaWriter) {
			ck(w.Assign("foo", "bar"))
		},
		output: "foo = bar\n",
	},
	{
		input: func(w *ninjaWriter) {
			ck(w.ScopedAssign("foo", "bar"))
		},
		output: "    foo = bar\n",
	},
	{
		input: func(w *ninjaWriter) {
			ck(w.Subninja("build.ninja"))
		},
		output: "subninja build.ninja\n",
	},
	{
		input: func(w *ninjaWriter) {
			ck(w.BlankLine())
		},
		output: "\n",
	},
	{
		input: func(w *ninjaWriter) {
			ck(w.Pool("p"))
			ck(w.ScopedAssign("depth", "3"))
			ck(w.BlankLine())
			ck(w.Comment("here comes a rule"))
			ck(w.Rule("r"))
			ck(w.ScopedAssign("command", "echo out: $out in: $in _arg: $_arg"))
			ck(w.ScopedAssign("pool", "p"))
			ck(w.BlankLine())
			ck(w.Build("r comment", "r", testNinjaStrings("foo.o"),
				nil, testNinjaStrings("foo.in"), nil, nil, nil, nil,
				nil, nil, nil, nil, nil, nil))
			ck(w.ScopedAssign("_arg", "arg value"))
		},
		output: `pool p
    depth = 3

# here comes a rule
rule r
    command = echo out: $out in: $in _arg: $_arg
    pool = p

# r comment
build foo.o: r foo.in
    _arg = arg value
`,
	},
}

func TestNinjaWriter(t *testing.T) {
	for i, testCase := range ninjaWriterTestCases {
		buf := bytes.NewBuffer(nil)
		w := newNinjaWriter(buf)
		testCase.input(w)
		if buf.String() != testCase.output {
			t.Errorf("incorrect output for test case %d", i)
			t.Errorf("  expected: %q", testCase.output)
			t.Errorf("       got: %q", buf.String())
		}
	}
}

func testNinjaStrings(s ...string) []*ninjaString {
	ret, _ := parseNinjaStrings(nil, s)
	return ret
}
