// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"text/template"
)

var Impls = []Impl{
	{"MPMCc_go", Blocking | Nonblocking},
	{"MPMCq_go", Blocking | Nonblocking},
	{"MPMCqp_go", Blocking | Nonblocking},

	// {"SPSCr_mc", Broken | Blocking | Batched},
	// {"SPSCrs_mc", Broken | Blocking | Batched},
	// {"MPSCr_mc", Broken | Blocking | Batched},
	// {"MPSCrs_mc", Blocking | Batched},

	{"SPSCns_dv", Blocking | Nonblocking | Unbounded},
	{"MPSCns_dv", Blocking | Nonblocking | Unbounded},
	// {"MPSCnsi_dv", Blocking | Nonblocking | Unbounded},

	{"MPMCqs_dv", Blocking | Nonblocking},
	{"MPMCqsp_dv", Blocking | Nonblocking},
	{"SPMCqs_dv", Blocking | Nonblocking},
	{"SPMCqsp_dv", Blocking | Nonblocking},
	{"SPSCqs_dv", Blocking | Nonblocking},
	{"SPSCqsp_dv", Blocking | Nonblocking},
}

type Flag int

const (
	Unbounded = Flag(1 << iota)
	Batched
	Blocking
	Nonblocking
)

type Impl struct {
	Name  string
	Flags Flag
}

func (impl *Impl) Faces() []string {
	faces := []string{}
	suffix := impl.Name[:4]

	if impl.Blocking() {
		faces = append(faces, suffix)
	}
	if impl.Nonblocking() {
		faces = append(faces, "Nonblocking"+suffix)
	}

	return faces
}

func (impl *Impl) New() string {
	batched, bounded := impl.Batched(), impl.Bounded()
	switch {
	case !batched && bounded:
		return "New" + impl.Name + "(size)"
	case batched && bounded:
		return "New" + impl.Name + "(batchSize, size)"
	case !batched && !bounded:
		return "New" + impl.Name + "()"
	case batched && !bounded:
		return "New" + impl.Name + "(batchSize)"
	}
	return "New" + impl.Name + "()"
}

func (impl *Impl) Bounded() bool     { return !impl.Unbounded() }
func (impl *Impl) Unbounded() bool   { return impl.Flags&Unbounded == Unbounded }
func (impl *Impl) Batched() bool     { return impl.Flags&Batched == Batched }
func (impl *Impl) Blocking() bool    { return impl.Flags&Blocking == Blocking }
func (impl *Impl) Nonblocking() bool { return impl.Flags&Nonblocking == Nonblocking }

func main() {
	outname := flag.String("out", "", "")
	flag.Parse()

	var b bytes.Buffer
	t := template.Must(template.New("").Parse(T))
	err := t.Execute(&b, Impls)
	check(err)

	dst, err := format.Source(b.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed format: %v\n", err)
		dst = b.Bytes()
	}

	var out io.Writer = os.Stdout
	if *outname != "" {
		f, err := os.Create(*outname)
		check(err)
		defer f.Close()
		out = f
	}

	out.Write(dst)
}

func check(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

const T = `package queue

// Code generated by all_gen.go; DO NOT EDIT.

//go:generate go run all_gen.go -out all_test.go

import (
	"testing"
	"strconv"

	"github.com/loov/queue/testsuite"
)

{{ range . }}
{{ $impl := . }}

{{ range .Faces -}}
var _ {{ . }} = (*{{$impl.Name}})(nil)
{{ end }}

func Test{{.Name}}(t *testing.T) {
	{{- if .Batched -}}
	for _, batchSize := range testsuite.BatchSizes {
	{{- else -}}
	batchSize := 0;
	{{- end -}}
		{{- if .Bounded -}}
		for _, size := range testsuite.TestSizes {
		{{- else -}}
		size := 0;
		{{- end -}}
			name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
			t.Run(name, func(t *testing.T){ testsuite.Tests(t, func() testsuite.Queue { return {{.New}} }) })
		{{- if .Bounded -}}
		}
		{{- end -}}
	{{- if .Batched -}}
	}
	{{- end -}}
}

func Benchmark{{.Name}}(b *testing.B) {
	{{- if .Batched -}}
	for _, batchSize := range testsuite.BatchSizes {
	{{- else -}}
	batchSize := 0;
	{{- end -}}
		{{- if .Bounded -}}
		for _, size := range testsuite.BenchSizes {
		{{- else -}}
		size := 0;
		{{- end -}}
			name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
			b.Run(name, func(b *testing.B){ testsuite.Benchmarks(b, func() testsuite.Queue { return {{.New}} }) })
		{{- if .Bounded -}}
		}
		{{- end -}}
	{{- if .Batched -}}
	}
	{{- end -}}
}
{{ end }}
`
