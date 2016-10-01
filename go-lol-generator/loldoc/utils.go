package loldoc

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kdy1997/go-lol/go-lol-generator/htmlutil"
	"github.com/luci/luci-go/common/logging"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type resIDCtxKeyType struct{}

var resIDCtxKey resIDCtxKeyType

func withResID(c context.Context, resID string) context.Context {
	c = logging.SetField(c, "resource", resID)
	return context.WithValue(c, resIDCtxKey, resID)
}

// resID returns resource id configured by WithResID()
func resID(c context.Context) string {
	val := c.Value(resIDCtxKey)
	if val == nil {
		panic(`loldoc: context must contain resource id`)
	}

	return val.(string)
}

func consumeSelect(s htmlutil.Sel) (vals map[string]string, err error) {
	s.Ensure("select")

	vals = make(map[string]string)

	for _, os := range s.Children() {
		os.Ensure("option")

		text := os.Text()
		val, ok := os.Attr("value")
		if !ok {
			return nil, s.WithDump(errors.Errorf("consumeSelect: no key for value %q", text))
		}

		vals[val] = text
	}
	s.ReplaceWithHtml(commentf("select: %v", vals))
	return
}

func consumeSimpleTable(table htmlutil.Sel) ([]map[string]string, error) {
	table.Ensure("table")

	thead := table.ChildrenFiltered("thead").MustBeSingle()
	tbody := table.ChildrenFiltered("tbody").MustBeSingle()

	columns := consumeTableRow(thead.Children().First(), true)
	rows := consumeTableRows(tbody)

	var data []map[string]string
	for _, row := range rows {
		if len(row) != len(columns) {
			return nil, table.WithDump(errors.Errorf(`row.len(%d) != columns.len(%d)`, len(row), len(columns)))
		}
		rowData := make(map[string]string)
		for i := range columns {
			rowData[columns[i]] = row[i]
		}
		data = append(data, rowData)
	}
	return data, nil
}

func consumeTableRows(tbody htmlutil.Sel) [][]string {
	tbody.Ensure("tbody")
	var rows [][]string

	for _, tr := range tbody.Children() {
		row := consumeTableRow(tr, false)
		rows = append(rows, row)
	}

	return rows
}

func consumeTableRow(tr htmlutil.Sel, isHead bool) []string {
	tr.Ensure("tr")
	var vals []string

	var tagName string
	if isHead {
		tagName = "th"
	} else {
		tagName = "td"
	}

	for _, cs := range tr.Children() {
		cs.Ensure(tagName)
		vals = append(vals, cs.Text())
	}

	tr.ReplaceWithHtml(commentf("Row: %v", vals))

	return vals
}

func removeIfUseless(_ int, s *goquery.Selection) {
	if isUseless(s) {

		s.Remove()
		return
	}
	s.RemoveAttr("onclick").RemoveAttr("style")
	s.Children().Each(removeIfUseless)
	// check twice.
	if isUseless(s) {
		s.Remove()
		return
	}
}

func isUseless(s *goquery.Selection) bool {
	if s.Is("head") || s.Is("script") || s.Is("style") || s.Is("title") {
		return true
	}

	if rel, _ := s.Attr("rel"); s.Is("link") && (rel == "stylesheet" || rel == "shortcut icon") {
		return true
	}

	if s.Is("div#footer") ||
		s.Is("div.navbar") ||
		s.Is("div.header.container.ezreal") ||
		s.Is("div#inputs-link") ||
		s.Is(".sandbox_header") {
		return true
	}

	if s.Is("div.push") && len(s.Children().Nodes) == 0 {
		return true
	}

	return false
}

//
func isEmpty(s *goquery.Selection) bool {
	if len(s.Nodes) == 0 {
		return true
	}
	if len(s.Children().Nodes) != 0 {
		return false
	}
	if s.Text() != "" {
		return false
	}

	if a, _ := s.Attr("class"); a != "" {
		return false
	}
	if a, _ := s.Attr("href"); a != "" {
		return false
	}
	if a, _ := s.Attr("id"); a != "" {
		return false
	}
	return true
}

func commentf(format string, args ...interface{}) string {
	return "<!-- " + fmt.Sprintf(format, args...) + " -->"
}

func parseResourceIDVersion(s string) (id, ver string) {
	if len(s) < 2 {
		return
	}
	lastIdx := strings.LastIndexByte(s, '-')
	if lastIdx == -1 {
		return s, ""
	}

	id = s[0:lastIdx]
	ver = s[lastIdx+1:] // exclude '-' from version
	return
}

// [BR, EUNE, EUW, JP, KR, LAN, LAS, NA, OCE, RU, TR]
// => []string{"BR", "EUNE", "EUW", "JP", "KR", "LAN", "LAS", "NA", "OCE", "RU", "TR"}
func parseRegions(src string) []string {
	src = strings.Replace(src, " ", "", -1)
	if len(src) <= 2 {
		return nil
	}
	if src[0] != '[' || src[len(src)-1] != ']' {
		return nil
	}

	src = src[1 : len(src)-1]

	return strings.Split(src, ",")
}

func IsRegion(p Parameter) bool {
	switch p.Name {
	case "region", "platformId":
		return true
	}
	return false
}
