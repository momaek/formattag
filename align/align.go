package align

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
)

// Manager align manager
type Manager struct {
	node   ast.Node
	fset   *token.FileSet
	groups []group
}

// group every struct node
type group struct {
	maxTagNum int
	lines     []*line
}

type line struct {
	field     *ast.Field
	tags      []string
	lens      []int
	spaceLens []int
	result    string
}

var mgr *Manager

// Init init align tag manager
func Init(content any) {
	mgr = &Manager{
		fset: token.NewFileSet(),
	}
	switch t := content.(type) {
	case string:
		node, err := parser.ParseFile(mgr.fset, t, nil, parser.ParseComments)
		if err != nil {
			log.Fatal("parse file failed ", err)
			return
		}
		mgr.node = node

	case io.Reader, []byte:
		node, err := parser.ParseFile(mgr.fset, "", t, parser.ParseComments)
		if err != nil {
			log.Fatal("parse file failed ", err)
			return
		}
		mgr.node = node
	default:
		log.Fatal("Unsupported content type")
	}
}

// Do do align
func Do() ([]byte, error) {
	return mgr.Align()
}

// Align align job
func (m *Manager) Align() ([]byte, error) {
	ast.Inspect(m.node, m.processStructTags)
	m.calcTagPosition()

	for _, grp := range m.groups {
		for _, line := range grp.lines {
			if len(line.result) == 0 {
				line.result = fmt.Sprintf("`%s`", strings.Join(line.tags, " "))
			}

			// write back to struct tag field
			line.field.Tag.Value = line.result
		}
	}

	var buf bytes.Buffer
	err := format.Node(&buf, m.fset, m.node)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *Manager) calcTagPosition() {
	for _, grp := range m.groups {
		if len(grp.lines) <= 1 {
			continue
		}

		for i := 0; i < grp.maxTagNum; i++ {
			max := getMaxTagLength(grp.lines, i)
			updateResult(grp.lines, max, i)
		}

		for _, line := range grp.lines {
			line.result = "`" + line.result + "`"
		}
	}
}

func getMaxTagLength(lines []*line, idx int) int {
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

func (m *Manager) processStructTags(node ast.Node) bool {
	// if not a struct node return
	st, ok := node.(*ast.StructType)
	if !ok {
		return true
	}
	// if struct do not have elements return
	if len(st.Fields.List) == 0 {
		return true
	}

	lastLineNum := m.fset.Position(st.Fields.List[0].Pos()).Line
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
				lastLineNum = m.fset.Position(st.Fields.List[idx+1].Pos()).Line
			}

			m.groups = append(m.groups, grp)
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
			lens = append(lens, length(t.String()))
			ln.tags = append(ln.tags, t.String())
		}

		ln.lens = lens

		lineNum := m.fset.Position(field.Pos()).Line
		if lineNum-lastLineNum >= 2 {
			lastLineNum = lineNum
			m.groups = append(m.groups, grp)
			grp = group{
				maxTagNum: tags.Len(),
			}
		}

		lastLineNum = lineNum

		grp.lines = append(grp.lines, ln)
	}

	if len(grp.lines) > 0 {
		m.groups = append(m.groups, grp)
	}

	return true
}

func length(s string) int {
	return len([]rune(s))
}
