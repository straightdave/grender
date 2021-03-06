package grender

import (
	"embed"
	"strings"
	"testing"
)

//go:embed templates/*
var testFS embed.FS

func TestCreateGrender(t *testing.T) {
	r := New(OptionMissingKeyZero(true))
	if r.missingKeyZero != true {
		t.Logf(">> missingkey should be true")
		t.Fail()
	}
}

func TestLoadFromFS(t *testing.T) {
	r := New(
		OptionTemplateDir("templates"),
		OptionLayoutDir("templates/layouts"),
	)

	if err := r.LoadFromFS(testFS); err != nil {
		t.Logf("failed to load templates from filesystem: %v", err)
		t.FailNow()
	}

	if r.get("hehe", false) == nil {
		t.Logf("hehe.tmpl not loaded")
		t.Fail()
	}

	res, err := r.Render("haha", "hehe", nil)
	if err != nil {
		t.Logf("failed to render: %v", err)
		t.FailNow()
	}

	// test fixtures may contain blank line or spaces due to IDE configuration.
	if strings.TrimSpace(res.String()) != "haha hehe" {
		t.Logf("> actual: %v", res)
		t.Fail()
	}
}

func TestNotLoadByExt(t *testing.T) {
	r := New(
		OptionLayoutDir("templates/layouts"),
		OptionTemplateDir("templates"),
	)

	if err := r.LoadFromFS(testFS); err != nil {
		t.Logf("failed to load templates from filesystem: %v", err)
		t.FailNow()
	}

	if r.get("not_a_tmpl", false) != nil {
		t.Logf("not_a_tmpl should not be loaded")
		t.Fail()
	}
}

func TestLoadByExt(t *testing.T) {
	r := New(
		OptionLayoutDir("templates/layouts"),
		OptionTemplateDir("templates"),
		OptionTemplateExt([]string{"no"}),
	)

	if err := r.LoadFromFS(testFS); err != nil {
		t.Logf("failed to load templates from filesystem: %v", err)
		t.FailNow()
	}

	if r.get("not_a_tmpl", false) == nil {
		t.Logf("not_a_tmpl should be loaded")
		t.Fail()
	}
}

func TestAddTemplates(t *testing.T) {
	r := New()
	if err := r.Add("P1", `page`); err != nil {
		t.Logf(">> failed to add template P1: %v", err)
		t.FailNow()
	}

	if err := r.Add("P1", `page`); err == nil {
		t.Logf(">> cannot add new templates with duplicate names")
		t.Fail()
	}

	if err := r.AddLayout("L1", `layout1`); err != nil {
		t.Logf(">> failed to add template L1: %v", err)
		t.FailNow()
	}

	if err := r.AddLayout("L1", `layout1`); err == nil {
		t.Logf(">> cannot add new templates with duplicate names")
		t.Fail()
	}
}

func TestInvalidTemplate(t *testing.T) {
	r := New()
	if err := r.Add("P1", `{{ .Name`); err == nil {
		t.Logf(">> not triggering errors when input is invalid")
		t.Fail()
	}

	if err := r.AddLayout("L1", `{{ .Name`); err == nil {
		t.Logf(">> not triggering errors when input is invalid")
		t.Fail()
	}
}

func TestNoTemplates(t *testing.T) {
	r := New()
	if _, err := r.Render("", "", nil); err == nil {
		t.Logf(">> no page name, should fail")
		t.Fail()
	}

	if _, err := r.Render("", "P1", nil); err == nil {
		t.Logf(">> no such page, should fail")
		t.Fail()
	}

	if _, err := r.Render("L1", "", nil); err == nil {
		t.Logf(">> no such layout, should fail")
		t.Fail()
	}
}

func TestNoShare(t *testing.T) {
	r := New()
	r.Add("P1", `{{ share "S1" }}`)

	// panic in custom functions would be transformed into errors
	// by golang templating engine

	_, err := r.Render("", "P1", nil)
	if err != nil {
		if !strings.Contains(err.Error(), "no shared template S1") {
			t.Logf(">> should have 'no shared template S1' error")
			t.Fail()
		}
	}

	r.AddLayout("L1", `{{ share "S1" }}: {{ yield }}`)
	r.Add("P2", `hello`)
	_, err = r.Render("L1", "P2", nil)
	if err != nil {
		if !strings.Contains(err.Error(), "no shared template S1") {
			t.Logf(">> should have 'no shared template S1' error")
			t.Fail()
		}
	}
}

func TestCurrent(t *testing.T) {
	r := New()
	r.AddLayout("L1", `This is layout {{ current }}: {{ yield }}`)
	r.Add("P1", `This is page {{ current }}`)

	out, err := r.Render("L1", "P1", nil)
	if err != nil {
		t.Fatalf(">> failed to render: %v", err)
	}

	if out.String() != "This is layout L1: This is page P1" {
		t.Logf(">> actual: %s", out)
		t.Fail()
	}
}

func TestNoYield(t *testing.T) {
	r := New()
	r.AddLayout("L1", `No yield`)
	r.Add("P1", `hello`)

	out, err := r.Render("L1", "P1", nil)
	if err != nil {
		t.Fatalf(">> failed to render: %v", err)
	}

	if out.String() != "No yield" {
		t.Logf(">> actual: %s", out)
		t.Fail()
	}
}

func TestRenderWithData(t *testing.T) {
	r := New(OptionMissingKeyZero(true))

	r.AddLayout("L1", `NoValue: {{.Subtitle}} Page: {{ yield }}`)
	r.Add("P1", `DownCase: {{.name}} UpperCase: {{.Age}} Object: {{.obj.p1}} ObjectNotExists: {{.obj.notExists}}`)

	out, err := r.Render("L1", "P1", map[string]interface{}{
		"Name": "dave",
		"name": "mike",
		"Age":  18,
		"obj": map[string]int{
			"p1": 1,
		},
	})

	if err != nil {
		t.Fatalf(">> failed to render: %v", err)
	}

	if out.String() != "NoValue: <no value> Page: DownCase: mike UpperCase: 18 Object: 1 ObjectNotExists: 0" {
		t.Logf(">> actual: %s", out)
		t.Fail()
	}
}

func TestShared(t *testing.T) {
	r := New()

	r.AddLayout("L1", `Layout {{ share "S1" }} => {{ yield }}`)
	r.Add("S1", `Shared Content {{.name}}`)
	r.Add("P1", `{{ share "S1" }} => {{ .name }}`)

	out, err := r.Render("L1", "P1", map[string]interface{}{
		"name": "dave",
	})

	if err != nil {
		t.Fatalf(">>failed to render: %v", err)
	}

	if out.String() != "Layout Shared Content dave => Shared Content dave => dave" {
		t.Logf(">> actual: %s", out)
		t.Fail()
	}
}
