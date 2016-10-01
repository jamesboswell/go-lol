package loldoc

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kdy1997/go-lol/go-lol-generator/htmlutil"
	"github.com/kdy1997/go-lol/go-lol-generator/patcher"
	"github.com/luci/luci-go/common/logging"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const docURL = `https://developer.riotgames.com/api/methods`

func NewGoQueryDoc() (*goquery.Document, error) {
	return goquery.NewDocument(docURL)
}

func Parse(c context.Context, d *goquery.Document) (*Doc, error) {
	doc := &Doc{}

	s := d.ChildrenFiltered("html").
		ChildrenFiltered("body").
		ChildrenFiltered("#wrap").
		ChildrenFiltered(".body.container")

	removeIfUseless(0, s)

	for _, s := range htmlutil.WrapAll(s.Find(".row .span12 #api_detail")) {
		s.Ensure("div#api_detail")
		ul := s.Children().First().Ensure("ul#resources")

		for _, li := range ul.Children() {
			res, err := parseResource(c, li)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse resource %q \n", res.ID)
			} else if res.ID == "" { // skip
				continue
			}

			doc.Resources = append(doc.Resources, res)
		}
	}

	return doc, nil
}

func parseResource(c context.Context, s htmlutil.Sel) (Resource, error) {
	s.Ensure("li.resource")
	s.ChildrenFiltered(".heading").MustBeSingle().
		ChildrenFiltered("ul.options").MustBeSingle().Remove()
	res := Resource{
		Definitions: make(map[string]Schema),
	}

	heading := s.ChildrenFiltered(".heading").MustBeSingle()
	{
		h2 := heading.Children().First().Ensure("h2")
		a := h2.Children().First().Ensure("a")
		idVer := a.Children().First().Ensure("span").Text()
		res.ID, res.Version = parseResourceIDVersion(idVer)

		// [BR, EUNE, .... , TR]
		regionText := a.Children().Last().Ensure("span").Text()
		res.Regions = parseRegions(regionText)
	}

	// this client doesn't support tournament api.
	if res.ID == "tournament-provider" {
		logging.Warningf(c, "skipping resource %q", res.ID)
		return Resource{}, nil
	}
	c = withResID(c, res.ID)

	for _, s := range s.ChildrenFiltered("ul.endpoints").Children() {
		ops, err := parseEndpoint(c, &res, s)
		if err != nil {
			return res, errors.Wrap(err, "parse endpoint")
		}
		res.Operations = append(res.Operations, ops...)
	}

	logging.Infof(c, "parsed resource %q", res.ID)

	return res, nil
}

func parseEndpoint(c context.Context, res *Resource, s htmlutil.Sel) ([]*Operation, error) {
	s.Ensure("li.endpoint")
	var ops []*Operation

	for _, s := range s.ChildrenFiltered("ul.operations").Children() {
		op, err := parseOperation(c, res, s)
		if err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}

	return ops, nil
}

// ".heading > .http_method": HTTP method to use
// ".heading > .path": HTTP request path
// ".heading > .options": description
func parseOperation(c context.Context, res *Resource, s htmlutil.Sel) (*Operation, error) {
	s.Ensure("li.operation")

	op := Operation{}
	op.res = res

	{ // parse: .heading
		heading := s.ChildrenFiltered("div.heading").MustBeSingle()

		op.RequestPath = heading.
			ChildrenFiltered(".path").MustBeSingle().
			Children().MustBeSingle().EatText()
		op.Description = heading.
			ChildrenFiltered("ul.options").MustBeSingle().
			Children().MustBeSingle().Ensure("li").
			Children().MustBeSingle().Ensure("a").
			EatText()

		heading.Remove() // remove: .heading
	}
	{ // handle special path parameters like 'region' and 'platformId'
		for _, keyword := range []string{"region", "platformId"} {
			if strings.Contains(op.RequestPath, "{"+keyword+"}") {
				typ, _ := patcher.Type(res.ID, "Region")
				op.PathParams = append(op.PathParams, Parameter{
					Name: keyword,
					Type: typ,
				})
			}
		}
	}
	patch, err := patcher.ForOperation(res.ID, op.RequestPath)
	if err != nil {
		return nil, err
	}
	op.MethodName = patch.Name
	op.OverridedMapKey = patch.MapKey

	switch {
	case s.HasClass("get"):
		op.HTTPMethod = "GET"
	case s.HasClass("post"):
		op.HTTPMethod = "POST"
	case s.HasClass("put"):
		op.HTTPMethod = "PUT"
	default:
		return nil, s.WithDump(errors.New("unknown operation method"))
	}

	for _, block := range s.Find(".content .api_block") {
		if err := parseAPIBlock(c, res, &op, block); err != nil {
			return nil, errors.Wrap(err, "parse api block")
		}
	}

	return &op, nil
}

func parseAPIBlock(c context.Context, res *Resource, op *Operation, s htmlutil.Sel) (err error) {
	s.Ensure(".api_block")

	if len(s.Children()) == 0 { // status api has an empty .api_block
		return nil
	}

	blockType := s.Children().First().Ensure("h4").EatText() // remove: h4

	switch blockType {
	case "Response Classes":
		for _, respBody := range s.ChildrenFiltered(".block.response_body").Reverse() {
			cnt := len(respBody.Children())
			switch cnt {
			case 1: // <b>Return Value:</b> $class
				respBody.Children().First().Ensure("b").EatExact("Return Value:")

				retValCls := respBody.Text()
				respBody.ReplaceWithHtml(commentf("return: %s", retValCls))
				retValType, err := patcher.Type(res.ID, retValCls)
				if err != nil {
					return errors.Wrapf(err, "failed to parse return value\n")
				}
				op.OrigReturnType = retValType
				logging.Debugf(c, "return value: %q", retValCls)
				continue
			case 3:
				cls, err := parseResponseClass(c, respBody)
				if err != nil {
					return err
				}
				cls.res = res
				res.Definitions[cls.OrigName] = cls

			default:
				return errors.Errorf("unexpected child count %d from response classes\n", cnt)
			}
		}

		return nil

	case "Response Errors":
		respErrors, err := parseResponseErrors(c, s)
		if err != nil {
			return errors.Wrap(err, "failed to parse response errors\n")
		}

		op.ResponseErrors = respErrors
		s.ReplaceWithHtml(commentf("<table> Response Errors </table>"))
		return nil

	case "Query Parameters":
		params, err := parseQueryParams(c, s)
		if err != nil {
			return errors.Wrap(err, "failed to parse query paramters\n")
		}
		op.QueryParams = params
		return nil

	case "Path Parameters":
		params, err := parsePathParams(c, s)
		if err != nil {
			return errors.Wrap(err, "failed to parse path parameters\n")
		}
		op.PathParams = append(op.PathParams, params...)
		return nil

	case "Select Region to Execute Against": // ignore
		return nil

	case "Implementation Notes":
		// has single <p> node as a child
		op.ImplNotes = s.
			Children().MustBeSingle().Ensure("p").Text()
		return nil

	case "Rate Limit Notes":
		{
			// <p> <span> $text </span> </p>
			p := s.Children().MustBeSingle().Ensure("p")
			op.RateLimitNotes = p.Children().MustBeSingle().Ensure("span").Text()
			return nil
		}

	default:
		return s.WithDump(errors.Errorf("unknown api block %q", blockType))
	}
}

func parseResponseErrors(c context.Context, s htmlutil.Sel) ([]*ResponseError, error) {
	s.Ensure(".api_block")

	datas, err := consumeSimpleTable(s.Children().First())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse response error table\n")
	}

	var errs []*ResponseError
	for _, data := range datas {
		code, err := strconv.Atoi(data["HTTP Status Code"])
		if err != nil {
			return nil, err
		}
		reason := data["Reason"]

		errs = append(errs, NewResponseError(code, reason))
	}

	return errs, nil
}

func parsePathParams(c context.Context, s htmlutil.Sel) (Parameters, error) {
	s.Ensure(".api_block")
	table := s.Children().First().Ensure("table.table")

	var params Parameters

	tbody := table.ChildrenFiltered("tbody").MustBeSingle().Ensure(".operation-params")

	for _, tr := range tbody.Children() {
		tr.Ensure("tr")

		param := Parameter{}
		{
			td := tr.Children().First().Ensure("td.code")
			cs := td.ChildrenFiltered("div.required").MustBeSingle()
			if cs.EatText() == "required" { // remove: div.required
				param.Required = true
			}
			param.Name = td.Text() // remove: td.code
		}

		param.Description = tr.Children().Last().Ensure("td").EatText()

		{
			td := tr.Children().Last().Ensure("td")

			typeStr := td.Children().First().Ensure("span.model-signature").Text()
			typeStr = patcher.PathParamType(param.Name, typeStr)
			typ, err := patcher.Type(resID(c), typeStr)
			if err != nil {
				return nil, err
			}
			param.Type = typ
			td.Remove()
		}
		params = append(params, param)
		logging.Debugf(c, "path param %q", param.Name)
	}
	return params, nil
}

// <table> query params </table>
// <select> REGION </select>
func parseQueryParams(c context.Context, s htmlutil.Sel) (Parameters, error) {
	s.Ensure(".api_block")
	var params Parameters

	table := s.Children().First().Ensure("table")
	thead := table.Children().First().Ensure("thead")
	tbody := table.Children()[1].Ensure("tbody.operation-params")

	columns := consumeTableRow(thead.Children().MustBeSingle().Ensure("tr"), true)

	for _, tr := range tbody.Children() {
		tr.Ensure("tr")

		param := Parameter{}
		{
			param.Required = true

			codeBlock := tr.Children().First().Ensure("td.code")
			optionalBlock := codeBlock.Children().First().Ensure("div.required")
			optional := optionalBlock.EatText() // remove: div.required
			if optional == "optional" {
				param.Required = false
			}
			param.Name = codeBlock.EatText() // remove: td.code
		}

		// description
		param.Description = tr.Children().Last().Ensure("td").EatText() // remove: td

		// get parameter type
		{
			td := tr.Children().Last().Ensure("td")
			typeStr := td.Children().First().Ensure("span.model-signature").Text()
			typ, err := patcher.Type(resID(c), typeStr)
			if err != nil {
				return nil, err
			}
			param.Type = typ
			td.Remove() // remove: td
		}

		params = append(params, param)
	}

	// ignore: <h4> $REGION </h4>

	_ = columns
	return params, nil
}

// <b>$class</b> - $description
// <br>
// <table> $fields... </table>
func parseResponseClass(c context.Context, s htmlutil.Sel) (cls Schema, err error) {
	s.Ensure(".response_body")
	cls.Fields = make([]Field, 0)
	defer func() {
		if err != nil {
			err = errors.Wrapf(err, "failed to parse class %q \n", cls.OrigName)
		}
	}()

	cls.OrigName = s.Children().First().Ensure("b").EatText() // remove: b
	cp, err := patcher.ForClass(resID(c), cls.OrigName)
	if err != nil {
		return cls, err
	}
	cls.StructName = cp.Name
	c = logging.SetField(c, "class", cls.StructName)

	s.Children().First().Ensure("br").Remove() // remove: <br>

	{
		rows, err := consumeSimpleTable(s.Children().First().Ensure("table"))
		if err != nil {
			return cls, err
		}

		for _, row := range rows {
			rawName := row["Name"]
			typeStr := row["Data Type"]
			description := row["Description"]

			typeStr = patcher.FieldTypeString(resID(c), cls.OrigName, rawName, typeStr)

			typ, err := patcher.Type(resID(c), typeStr)
			if err != nil {
				err = errors.Wrapf(err, "failed to parse type of field %q\n", rawName)
				return cls, err
			}

			cls.Fields = append(cls.Fields, NewField(rawName, typ, description))
		}
	}

	s.Children().First().Remove() // remove: table
	cls.Description = s.Text()

	return
}
