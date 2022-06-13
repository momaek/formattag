package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
)

type config struct {
	fset   *token.FileSet
	file   string
	groups []group
}

var Version string

func main() {
	err := do()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func do() (err error) {
	c, err := parseConfig(os.Args[1:])
	if err != nil {
		return
	}

	node, err := c.parse()
	if err != nil {
		return
	}

	c.format(node)
	err = c.write(node)
	return
}

func (c *config) write(node ast.Node) error {
	for _, grp := range c.groups {
		for _, line := range grp.lines {
			if len(line.result) == 0 {
				line.result = fmt.Sprintf("`%s`", strings.Join(line.tags, " "))
			}

			line.field.Tag.Value = line.result
		}
	}

	var buf bytes.Buffer
	err := format.Node(&buf, c.fset, node)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.file, buf.Bytes(), 0)
	return err
}

func parseFlag(args []string) string {
	if len(args) == 0 {
		printHelp()
	}

	for idx, cmd := range args {
		if cmd == "version" {
			printVersion()
		}

		if cmd == "-file" {
			if idx != 0 || len(args) != 2 {
				printHelp()
			}
			return args[1]
		}
	}
	return ""
}

func printVersion() {
	fmt.Printf("formattag verion %s %s\n", Version, "darwin/amd64")
	os.Exit(0)
}

func printHelp() {
	fmt.Printf(`formattag is a tool for formatting Go source code.

Usage:

	formattag <command> [flags]

The commands are:
	version     print formattag version

The flags are:
	-file filename to be format
`)
	os.Exit(0)
}

func parseConfig(args []string) (*config, error) {
	c := &config{
		file:   parseFlag(args),
		groups: []group{},
	}

	return c, nil
}

func (c *config) parse() (ast.Node, error) {
	c.fset = token.NewFileSet()
	var content interface{}
	return parser.ParseFile(c.fset, c.file, content, parser.ParseComments)
}

func (c *config) format(node ast.Node) (ast.Node, error) {

	ast.Inspect(node, c.rewrite)
	c.process()
	return nil, nil
}

type line struct {
	field     *ast.Field
	tags      []string
	lens      []int
	spaceLens []int
	result    string
}

var runID = 0

func (c *config) rewrite(node ast.Node) bool {
	st, ok := node.(*ast.StructType)
	if !ok {
		return true
	}

	if len(st.Fields.List) == 0 {
		return true
	}

	c.preProcessStruct(st)

	return true
}

type group struct {
	maxTagNum int
	lines     []*line
}

func (c *config) preProcessStruct(st *ast.StructType, inline ...bool) {
	lastLineNum := c.fset.Position(st.Fields.List[0].Pos()).Line
	grp := group{}
	l := len(st.Fields.List)
	for idx, field := range st.Fields.List {
		if field.Tag == nil {
			continue
		}

		tag, err := strconv.Unquote(field.Tag.Value)
		if err != nil {
			continue
		}

		tag = strings.TrimLeft(tag, " ")
		tag = strings.TrimRight(tag, " ")

		tags, err := structtag.Parse(tag)
		if err != nil {
			continue
		}

		if _, ok := field.Type.(*ast.StructType); ok {
			if idx+1 < l {
				lastLineNum = c.fset.Position(st.Fields.List[idx+1].Pos()).Line
			}

			c.groups = append(c.groups, grp)
			grp = group{}
			continue
		}

		if grp.maxTagNum < tags.Len() {
			grp.maxTagNum = tags.Len()
		}

		ln := &line{
			field: field,
		}

		lens := make([]int, 0, tags.Len())
		for _, key := range tags.Keys() {
			t, _ := tags.Get(key)
			lens = append(lens, len(t.String()))
			ln.tags = append(ln.tags, t.String())
		}

		ln.lens = lens

		lineNum := c.fset.Position(field.Pos()).Line
		if lineNum-lastLineNum >= 2 {
			lastLineNum = lineNum
			c.groups = append(c.groups, grp)
			grp = group{
				maxTagNum: tags.Len(),
			}
		}

		lastLineNum = lineNum

		grp.lines = append(grp.lines, ln)
	}

	if len(grp.lines) > 0 {
		c.groups = append(c.groups, grp)
	}
}

func (c *config) process() {
	for _, grp := range c.groups {
		if len(grp.lines) <= 1 {
			continue
		}

		for i := 0; i < grp.maxTagNum; i++ {
			max := process(grp.lines, i)
			updateResult(grp.lines, max, i)
		}

		for _, line := range grp.lines {
			line.result = "`" + line.result + "`"
		}
	}
}

func process(lines []*line, idx int) int {
	max := 0
	for _, line := range lines {
		if len(line.lens) > idx {
			if line.lens[idx] > max {
				max = line.lens[idx]
			}
		}
	}

	return max
}

func updateResult(lines []*line, max, idx int) {
	for _, line := range lines {
		if len(line.tags) > idx {
			if l := len(line.lens); l > idx && idx < l-1 {
				line.result += line.tags[idx] + strings.Repeat(" ", max-line.lens[idx]+1)
			} else {
				line.result += line.tags[idx]
			}
		}
	}
}
