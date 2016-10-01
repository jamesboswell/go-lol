package loldoc

import (
	"fmt"
	"go/types"

	"github.com/kdy1997/go-lol/go-lol-generator/patcher"
)

type Doc struct {
	Resources []Resource
}

type Resource struct {
	ID          string
	Version     string
	Regions     []string
	Definitions map[string]Schema
	Operations  []*Operation
}

type child struct {
	res *Resource
}

type Operation struct {
	child
	PathParams  Parameters
	QueryParams Parameters
	Description string

	OrigReturnType  types.Type
	MethodName      string          // from override
	OverridedMapKey types.BasicKind // from override

	HTTPMethod  string
	RequestPath string

	ResponseErrors []*ResponseError
	ImplNotes      string // text from "Implementation Notes" block
	RateLimitNotes string // text from "Rate Limit Notes"
}

type Schema struct {
	*types.Named
	child
	Description string
	Fields      []Field
	OrigName    string
	StructName  string // from override
}

type Field struct {
	origName    string // for json
	Type        types.Type
	goName      string
	Description string
}

type Parameters []Parameter

type Parameter struct {
	Name        string
	Description string
	Required    bool
	Type        types.Type
}

func (ps Parameters) Has(name string) bool {
	for _, p := range ps {
		if p.Name == name {
			return true
		}
	}
	return false
}

type ResponseError struct {
	Code   int
	Reason string
}

func NewField(origName string, typ types.Type, description string) Field {
	return Field{
		origName:    origName,
		goName:      patcher.FieldName(origName),
		Type:        typ,
		Description: description,
	}
}

func (f Field) OrigName() string { return f.origName }
func (f Field) GoName() string   { return f.goName }

func NewResponseError(code int, reason string) *ResponseError {
	return &ResponseError{
		Code:   code,
		Reason: reason,
	}
}

// APIBase returns empty string if it's not a special operation.
func (res Resource) APIBase() string {
	switch res.ID {
	case "lol-static-data":
		return "https://global.api.pvp.net"
	case "lol-status":
		return "https://status.leagueoflegends.com"
	default:
		return ""
	}
}

// APIBase returns empty string if it's not a special operation.
func (op Operation) APIBase() string { return op.res.APIBase() }

func (res Resource) NeedAPIKey() bool {
	switch res.ID {
	case "lol-status":
		return false
	default:
		return true
	}
}

func (op Operation) NeedAPIKey() bool { return op.res.NeedAPIKey() }

func (op Operation) SupportedRegions() []string {
	return op.res.Regions
}

// IsRegional returns true if this operation has a path parameter named 'region' or 'platformId'.
func (op *Operation) IsRegional() bool {
	if op.res.APIBase() == "" { // normal api
		return true
	}

	for _, p := range op.PathParams { // typically api to global region
		if IsRegion(p) {
			return true
		}
	}
	return false
}

func (c child) ResID() string { return c.res.ID }

func (res Resource) GoString() string {
	return fmt.Sprintf("%q", res.ID)
}

func (f Field) GoString() string {
	return fmt.Sprintf("%q %s (%q)", f.GoName(), f.Type, f.OrigName())
}
