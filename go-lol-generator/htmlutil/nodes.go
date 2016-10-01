package htmlutil

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

type Sel struct {
	*goquery.Selection
}

type Sels []Sel

func Wrap(s *goquery.Selection) Sel {
	if s.Length() != 1 {
		panic(errors.Errorf("htmlutil: invalid selector %s", Dump(s)))
	}
	return Sel{s}
}

func WrapAll(s *goquery.Selection) Sels {
	var ret Sels
	for i := range s.Nodes {
		ret = append(ret, Wrap(s.Eq(i)))
	}
	return ret
}

// NOTE: text replaces 's' with comment
func (sel Sel) Text() string {
	if cnt := len(sel.Selection.Children().Nodes); cnt != 0 {
		err := errors.Errorf("cannot read text: node must not have a child, but has %d", cnt)
		panic(sel.WithDump(err))
	}

	s := sel.Selection.Contents()
	for _, node := range s.Nodes {
		if node.Type == html.CommentNode { // ignore comment node
			continue
		}

		if node.Type == html.TextNode {
			t := node.Data
			return t
		}
	}

	return ""
}

func (sel Sel) EatText() string {
	if cnt := len(sel.Selection.Children().Nodes); cnt != 0 {
		err := errors.Errorf("cannot read text: node must not have a child, but has %d", cnt)
		panic(sel.WithDump(err))
	}

	s := sel.Selection.Contents()
	for _, node := range s.Nodes {
		if node.Type == html.CommentNode { // ignore comment node
			continue
		}

		if node.Type == html.TextNode {
			t := node.Data
			sel.ReplaceWithHtml(fmt.Sprintf("<!-- %q -->", t))
			return t
		}
	}

	return ""
}

func (sel Sel) EatExact(text string) {
	got := sel.Text()
	if got != text {
		panic(errors.Errorf("EatText: want %q, but got %q", text, got))
	}
	sel.Remove()
}

func (sel Sel) Ensure(selector string) Sel {
	if !sel.Selection.Is(selector) {
		err := errors.Errorf("expected selector: %q\n", selector)
		panic(sel.WithDump(err))
	}
	return sel
}

func (sel Sel) Find(selector string) Sels {
	return WrapAll(sel.Selection.Find(selector))
}

func (sel Sel) Children() Sels {
	return WrapAll(sel.Selection.Children())
}

func (sel Sel) ChildrenFiltered(selector string) Sels {
	return WrapAll(sel.Selection.ChildrenFiltered(selector))
}

func (sel Sel) Dump() string {
	return Dump(sel.Selection)
}

// WithDump wraps err with error with dump data.
func (sel Sel) WithDump(err error) error {
	return errors.Wrap(err, sel.Dump())
}

func (ss Sels) MustBeSingle() Sel {
	switch len(ss) {
	case 1:
		return ss[0]
	default:
		panic(errors.Errorf("must have single node, but has %d", len(ss)))
	}
}

func (ss Sels) First() Sel {
	return ss[0]
}

func (ss Sels) Last() Sel {
	return ss[len(ss)-1]
}

func (ss Sels) Children() Sels {
	var ret Sels
	for _, s := range ss {
		ret = append(ret, s.Children()...)
	}
	return ret
}

func (ss Sels) Reverse() Sels {
	var ret Sels
	for i := len(ss) - 1; i >= 0; i-- {
		ret = append(ret, ss[i])
	}
	return ret
}

// Dump creates an error containing dumped html of s.
func Dump(s *goquery.Selection) string {
	const prefix = "\nDUMP HTML: \n"

	if len(s.Nodes) == 0 {
		return prefix + "<!-- EMPTY NODE -->" + "\n"
	}

	str, err := goquery.OuterHtml(s)
	if err == nil {
		str = gohtml.FormatWithLineNo(str)
		return prefix + str + "\n"
	}

	str, _ = s.Html()
	str = gohtml.FormatWithLineNo(str)
	return prefix + "(INNER)\n" + str + "\n"
}
