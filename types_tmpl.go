// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var typesTmpl = `

	{{define "SimpleType"}}
		{{$type := replaceReservedWords .Name | makePublic}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
		{{if ne .List.ItemType ""}}
			type {{$type}} []{{toGoType .List.ItemType false }}
		{{else if ne .Union.MemberTypes ""}}
			type {{$type}} string
		{{else if .Union.SimpleType}}
			type {{$type}} string
		{{else}}
			type {{$type}} {{toGoTypeNoPointer .Restriction.Base}}
		{{end}}
        {{if .Restriction.SimpleType}}
        {{template "SimpleType" .Restriction.SimpleType}}
        {{end}}
		{{if .Restriction.Enumeration}}
			const (
					{{with .Restriction}}
						{{range .Enumeration}}
						{{if .Doc}} {{.Doc | comment}} {{end}}
						{{$type}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$type}} = "{{goString .Value}}" {{end}}
					{{end}}
					)
		{{end}}
	{{end}}
	{{define "ComplexContent"}}
		{{$baseType := toGoType .Extension.Base false }}
		{{ if $baseType }}
			{{$baseType}}
		{{end}}
		{{template "Elements" .Extension.SequenceChoice}}
        {{template "Elements" .Extension.SequenceChoiceSequence}}
		{{template "Elements" .Extension.Sequence}}
		{{template "Elements" .Extension.Choice}}
		{{template "Elements" .Extension.SequenceChoice}}
		{{template "Attributes" .Extension.Attributes}}
		{{template "Elements" .Extension.Choice}}
		{{template "AttributeGroups" .Extension.AttributeGroup}}
	{{end}}
	{{define "Attributes"}}
		{{ $targetNamespace := getNS }}
		{{range .}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
			{{ if ne .Type "" }}
				{{ normalize .Name | makeFieldPublic}} {{toGoType .Type false}} ` + "`" + `xml:"{{with $targetNamespace}}{{.}} {{end}}{{.Name}},attr,omitempty" json:"{{.Name}},omitempty"` + "`" + `
			{{ else }}
				{{ normalize .Name | makeFieldPublic}} string ` + "`" + `xml:"{{with $targetNamespace}}{{.}} {{end}}{{.Name}},attr,omitempty" json:"{{.Name}},omitempty"` + "`" + `
			{{ end }}
		{{end}}
	{{end}}
	{{define "AttributeGroups"}}
		{{range .}}
			{{if ne .Ref ""}}
				{{template "Attributes" getAttributesFromGroup .Ref}}
			{{ end }}
		{{end}}
	{{end}}
	{{define "SimpleContent"}}
		Value {{toGoTypeNoPointer .Extension.Base }} ` + "`xml:\",chardata\" json:\"-,\"`" + `
		{{template "Attributes" .Extension.Attributes}}
		{{template "AttributeGroups" .Extension.AttributeGroup}}
	{{end}}
	{{define "ComplexTypeInline"}}
		{{replaceReservedWords .Key | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}struct {
		{{with .ComplexType}}
			{{if ne .ComplexContent.Extension.Base ""}}
				{{template "ComplexContent" .ComplexContent}}
			{{else if ne .SimpleContent.Extension.Base ""}}
				{{template "SimpleContent" .SimpleContent}}
			{{else}}
				{{template "Elements" .Sequence}}
				{{template "Elements" .SequenceSequence}}
				{{template "Elements" .SequenceSequenceSequence}}
				{{template "Elements" .Choice}}
				{{template "Elements" .SequenceChoice}}
                {{template "Elements" .SequenceChoiceSequence}}
				{{template "Elements" .All}}
				{{template "Attributes" .Attributes}}
				{{template "Elements" .ChoiceSequence}}
				{{template "AttributeGroups" .AttributeGroup}}
			{{end}}
		{{end}}
		} ` + "`" + `xml:"{{.Key}},omitempty" json:"{{.Key}},omitempty"` + "`" + `
	{{end}}
{{define "Elements"}}
	{{range .}}
		{{if ne .Ref ""}}
			{{removeNS .Ref | replaceReservedWords  | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{toGoType .Ref .Nillable }} ` + "`" + `xml:"{{.Ref }},omitempty" json:"{{.Ref }},omitempty"` + "`" + `
		{{else}}
		{{if not .Type}}
			{{if .SimpleType}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .SimpleType.List.ItemType ""}}
					{{ normalize .Name | makeFieldPublic}} []{{toGoType .SimpleType.List.ItemType false}} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
				{{else}}
					{{ normalize .Name | makeFieldPublic}} {{toGoType .SimpleType.Restriction.Base false}} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
				{{end}}
			{{else}}
                {{ $complexInlineName := setElementInComplexInlineCache .}}
				{{replaceReservedWords .Name | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}*  {{$complexInlineName}} ` + "`" + `xml:"{{.Name | removeNS}},omitempty json:"{{.Name | removeNS }},omitempty""` + "`" + `
			{{end}}
		{{else}}
			{{if .Doc}}{{.Doc | comment}} {{end}}
			{{replaceAttrReservedWords .Name | makeFieldPublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{toGoType .Type .Nillable }} ` + "`" + `xml:"{{.Name | removeNS }},omitempty" json:"{{.Name | removeNS }},omitempty"` + "`" + ` {{end}}
		{{end}}
	{{end}}
{{end}}
{{range .Schemas}}
	{{ $targetNamespace := .TargetNamespace }}
	{{range .SimpleType}}
		{{template "SimpleType" .}}
	{{end}}
	{{range .Elements}}
  		{{$name := .Name}}
		{{if not .Type}}
			{{with .SimpleType}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .List.ItemType ""}}
					type {{$name}} []{{toGoType .List.ItemType false }}
				{{else if ne .Union.MemberTypes ""}}
					type {{$name}} string
				{{else if .Union.SimpleType}}
					type {{$name}} string
				{{else}}
					type {{$name}} {{toGoTypeNoPointer .Restriction.Base}}
				{{end}}
				{{if .Restriction.SimpleType}}
				{{template "SimpleType" .Restriction.SimpleType}}
				{{end}}
				{{if .Restriction.Enumeration}}
					const (
							{{with .Restriction}}
								{{range .Enumeration}}
								{{if .Doc}} {{.Doc | comment}} {{end}}
								{{$name}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$name}} = "{{goString .Value}}" {{end}}
							{{end}}
							)
				{{end}}
          	{{end}}
		    {{/* ComplexTypeLocal */}}
            {{with .ComplexType}}
                type {{$name | replaceReservedWords | makePublic}} struct {
					XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$name}}\"`" + `
					{{if ne .ComplexContent.Extension.Base ""}}
						{{template "ComplexContent" .ComplexContent}}
					{{else if ne .SimpleContent.Extension.Base ""}}
						{{template "SimpleContent" .SimpleContent}}
					{{else}}
						{{template "Elements" .Sequence}}
                        {{template "Elements" .SequenceSequence}}
						{{template "Elements" .SequenceSequenceSequence}}
						{{template "Elements" .Choice}}
						{{template "Elements" .SequenceChoice}}
                        {{template "Elements" .SequenceChoiceSequence}}
						{{template "Elements" .All}}
						{{template "Attributes" .Attributes}}
                        {{template "Elements" .ChoiceSequence}}
                        {{template "AttributeGroups" .AttributeGroup}}
					{{end}}
				}
            {{end}}
		{{else}}
			type {{$name | replaceReservedWords | makePublic}} {{toGoType .Type false | removePointerFromType}}
		{{end}}
    {{end}}
	{{range .ComplexTypes}}
		{{/* ComplexTypeGlobal */}}
		{{$name := replaceReservedWords .Name | makePublic}}
		type {{$name}} struct {
			{{$typ := findNameByType .Name}}
			{{if ne .Name $typ}}
					XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$typ}}\"`" + `
			{{end}}
			{{if ne .ComplexContent.Extension.Base ""}}
				{{template "ComplexContent" .ComplexContent}}
			{{else if ne .SimpleContent.Extension.Base ""}}
				{{template "SimpleContent" .SimpleContent}}
			{{else}}
				{{template "Elements" .Sequence}}
                {{template "Elements" .SequenceSequence}}
				{{template "Elements" .SequenceSequenceSequence}}
                {{template "Elements" .Choice}}
				{{template "Elements" .SequenceChoice}}
                {{template "Elements" .SequenceChoiceSequence}}
				{{template "Elements" .All}}
				{{template "Attributes" .Attributes}}
                {{template "Elements" .ChoiceSequence}}
                {{template "AttributeGroups" .AttributeGroup}}
			{{end}}
		}
	{{end}}
{{end}}
`

var typesTmplComplexInline = `
{{define "SimpleType"}}
	{{$typeName := replaceReservedWords .Name | makePublic}}
	{{if .Doc}} {{.Doc | comment}} {{end}}
	{{if ne .List.ItemType ""}}
		type {{$typeName}} []{{toGoType .List.ItemType false | removePointerFromType}}
	{{else if ne .Union.MemberTypes ""}}
		type {{$typeName}} string
	{{else if .Union.SimpleType}}
		type {{$typeName}} string
	{{else if .Restriction.Base}}
		type {{$typeName}} {{toGoType .Restriction.Base false | removePointerFromType}}
    {{else}}
		type {{$typeName}} interface{}
	{{end}}

	{{if .Restriction.Enumeration}}
	const (
		{{with .Restriction}}
			{{range .Enumeration}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{$typeName}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$typeName}} = "{{goString .Value}}" {{end}}
		{{end}}
	)
	{{end}}
{{end}}

{{define "ComplexContent"}}
	{{$baseType := toGoType .Extension.Base false}}
	{{ if $baseType }}
		{{$baseType}}
	{{end}}

	{{template "Elements" .Extension.SequenceChoice}}
    {{template "Elements" .Extension.SequenceChoiceSequence}}
	{{template "Elements" .Extension.Sequence}}
	{{template "Elements" .Extension.Choice}}
	{{template "Elements" .Extension.SequenceChoice}}
	{{template "Attributes" .Extension.Attributes}}
	{{template "Elements" .Extension.Choice}}
    {{template "AttributeGroups" .Extension.AttributeGroup}}
{{end}}

{{define "Attributes"}}
    {{ $targetNamespace := getNS }}
	{{range .}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
		{{ if ne .Type "" }}
			{{ normalize .Name | makeFieldPublic}} {{toGoType .Type false}} ` + "`" + `xml:"{{with $targetNamespace}}{{.}} {{end}}{{.Name}},attr,omitempty" json:"{{.Name}},omitempty"` + "`" + `
		{{ else }}
			{{ normalize .Name | makeFieldPublic}} string ` + "`" + `xml:"{{with $targetNamespace}}{{.}} {{end}}{{.Name}},attr,omitempty" json:"{{.Name}},omitempty"` + "`" + `
		{{ end }}
	{{end}}
{{end}}

{{define "AttributeGroups"}}
	{{range .}}
		{{if ne .Ref ""}}
          {{template "Attributes" getAttributesFromGroup .Ref}}
         {{ end }}
	{{end}}
{{end}}

{{define "SimpleContent"}}
	Value {{toGoType .Extension.Base false}} ` + "`xml:\",chardata\" json:\"-,\"`" + `
	{{template "Attributes" .Extension.Attributes}}
	{{template "AttributeGroups" .Extension.AttributeGroup}}
{{end}}

{{define "ComplexTypeInline"}}
	{{replaceReservedWords .Key | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}} struct {
	{{with .ComplexType}}
		{{if ne .ComplexContent.Extension.Base ""}}
			{{template "ComplexContent" .ComplexContent}}
		{{else if ne .SimpleContent.Extension.Base ""}}
			{{template "SimpleContent" .SimpleContent}}
		{{else}}
			{{template "Elements" .Sequence}}
			{{template "Elements" .SequenceSequence}}
			{{template "Elements" .SequenceSequenceSequence}}
			{{template "Elements" .Choice}}
			{{template "Elements" .SequenceChoice}}
			{{template "Elements" .SequenceChoiceSequence}}
			{{template "Elements" .All}}
			{{template "Attributes" .Attributes}}
			{{template "Elements" .ChoiceSequence}}
            {{template "AttributeGroups" .AttributeGroup}}
		{{end}}
	{{end}}
	} ` + "`" + `xml:"{{.Key}},omitempty" json:"{{.Key}},omitempty"` + "`" + `
{{end}}

{{define "Elements"}}
	{{range .}}
		{{if ne .Ref ""}}
			{{removeNS .Ref | replaceReservedWords  | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{toGoType .Ref .Nillable }} ` + "`" + `xml:"{{.Ref }},omitempty" json:"{{.Ref }},omitempty"` + "`" + `
		{{else}}
		{{if not .Type}}
			{{if .SimpleType}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .SimpleType.List.ItemType ""}}
					{{ normalize .Name | makeFieldPublic}} []{{toGoType .SimpleType.List.ItemType false}} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
				{{else}}
					{{ normalize .Name | makeFieldPublic}} {{toGoType .SimpleType.Restriction.Base false}} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
				{{end}}
			{{else}}
				{{ $complexInlineName := setElementInComplexInlineCache .}}
				{{replaceReservedWords .Name | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}*  {{$complexInlineName}} ` + "`" + `xml:"{{.Name | removeNS}},omitempty json:"{{.Name | removeNS }},omitempty""` + "`" + `
			{{end}}
		{{else}}
			{{if .Doc}}{{.Doc | comment}} {{end}}
			{{replaceAttrReservedWords .Name | makeFieldPublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{toGoType .Type .Nillable }} ` + "`" + `xml:"{{.Name | removeNS }},omitempty" json:"{{.Name | removeNS }},omitempty"` + "`" + ` {{end}}
		{{end}}
	{{end}}
{{end}}

{{define "Any"}}
	{{range .}}
		Items     []string ` + "`" + `xml:",any" json:"items,omitempty"` + "`" + `
	{{end}}
{{end}}

{{ range $Key, $Value  := getComplexInlineCache }}
 type {{replaceReservedWords $Key | makePublic}} struct {
    {{with  $Value.SimpleType}}
                {{$name := $Key}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .List.ItemType ""}}
					type {{$name}} []{{toGoType .List.ItemType }}
				{{else if ne .Union.MemberTypes ""}}
					type {{$name}} string
				{{else if .Union.SimpleType}}
					type {{$name}} string
				{{else}}
					type {{$name}} {{toGoTypeNoPointer .Restriction.Base}}
				{{end}}
				{{if .Restriction.SimpleType}}
				{{template "SimpleType" .Restriction.SimpleType}}
				{{end}}
				{{if .Restriction.Enumeration}}
					const (
							{{with .Restriction}}
								{{range .Enumeration}}
								{{if .Doc}} {{.Doc | comment}} {{end}}
								{{$name}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$name}} = "{{goString .Value}}" {{end}}
							{{end}}
							)
				{{end}}
    {{end}}
	{{with $Value.ComplexType}}
		{{if ne .ComplexContent.Extension.Base ""}}
			{{template "ComplexContent" .ComplexContent}}
		{{else if ne .SimpleContent.Extension.Base ""}}
			{{template "SimpleContent" .SimpleContent}}
		{{else}}
			{{template "Elements" .Sequence}}
			{{template "Elements" .SequenceSequence}}
			{{template "Elements" .SequenceSequenceSequence}}
			{{template "Elements" .Choice}}
			{{template "Elements" .SequenceChoice}}
			{{template "Elements" .SequenceChoiceSequence}}
			{{template "Elements" .All}}
			{{template "Attributes" .Attributes}}
			{{template "Elements" .ChoiceSequence}}
            {{template "AttributeGroups" .AttributeGroup}}
		{{end}}
	{{end}}
	}
{{end}}
`
