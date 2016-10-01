package loldoc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/kdy1997/go-lol/go-lol-generator/htmlutil"
	"github.com/luci/luci-go/common/logging/memlogger"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
)

func TestParse(t *testing.T) { // dirty, but simple
	c := newTestingContext()

	data, err := ioutil.ReadFile("./methods.html")
	if err != nil {
		panic(err)
	}
	gqDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		panic(errors.Wrap(err, "failed to parse goquery document"))
	}
	defer memlogger.MustDumpStdout(c)

	doc, err := Parse(c, gqDoc)
	if err != nil {
		t.Error(err)
	}

	_ = doc
}

func TestParsingUtils(t *testing.T) {
	Convey("consumeSelect", t, func() {
		s := htmlutil.Wrap(mustParse(`<select class="select any class" id="dnjaf" name="virtual" >
<option value="a" >A</option>
<option value="b">B</option>
<option value="c" >C &amp; C</option>
</select>
`).Find("select"))

		vals, err := consumeSelect(s)
		if err != nil {
			t.Errorf("ParseSelect returns error: %s", err)
		} else {
			_ = vals
		}
		So(err, ShouldBeNil)
		So(vals, ShouldResemble, map[string]string{
			"a": "A",
			"b": "B",
			"c": "C & C",
		})
	})

	Convey("parseRegion", t, func() {
		So(parseRegions("[BR, EUNE, EUW]"), ShouldResemble, []string{"BR", "EUNE", "EUW"})
	})

	Convey("parseResourceIDVersion", t, func() {
		datas := []struct {
			src, id, ver string
		}{
			{
				src: "champion-v1.2",
				id:  "champion", ver: "v1.2",
			},
		}

		for _, data := range datas {
			Convey(fmt.Sprintf("Returns (id=%q,ver=%q) for %s", data.id, data.ver, data.src), func() {
				id, ver := parseResourceIDVersion(data.src)
				So(id, ShouldEqual, data.id)
				So(ver, ShouldEqual, data.ver)
			})
		}
	})

}

func mustParse(src string) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(src))
	if err != nil {
		panic(fmt.Sprintf("failed to parse document: %s", err))
	}
	return doc
}

func newTestingContext() context.Context {
	c := context.Background()
	c = memlogger.Use(c)
	return c
}
