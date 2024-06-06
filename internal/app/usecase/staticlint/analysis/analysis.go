// Package analysis works with standard analysers from go/analysis
package analysis

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"

	osexit "github.com/avGenie/url-shortener/internal/app/usecase/staticlint/analysis/os_exit_error"
)

// GetAnalyzers Returns all standard analyzers
//
// appends      check for missing values after append
// asmdecl      report mismatches between assembly files and Go declarations
// assign       check for useless assignments
// atomic       check for common mistakes using the sync/atomic package
// atomicalign  check for non-64-bits-aligned arguments to sync/atomic functions
// bools        check for common mistakes involving boolean operators
// buildssa     build SSA-form IR for later passes
// buildtag     check //go:build and // +build directives
// cgocall      detect some violations of the cgo pointer passing rules
// composites   check for unkeyed composite literals
// copylocks    check for locks erroneously passed by value
// ctrlflow     build a control-flow graph
// deepequalerrors check for calls of reflect.DeepEqual on error values
// defers       report common mistakes in defer statements
// directive    check Go toolchain directives such as //go:debug
// errorsas     report passing non-pointer or non-error values to errors.As
// fieldalignment find structs that would use less memory if their fields were sorted
// findcall     find calls to a particular function
// framepointer report assembly that clobbers the frame pointer before saving it
// httpmux      report using Go 1.22 enhanced ServeMux patterns in older Go versions
// httpresponse check for mistakes using HTTP responses
// ifaceassert  detect impossible interface-to-interface type assertions
// inspect      optimize AST traversal for later passes
// loopclosure  check references to loop variables from within nested functions
// lostcancel   check cancel func returned by context.WithCancel is called
// nilfunc      check for useless comparisons between functions and nil
// nilness      check for redundant or impossible nil comparisons
// os_exit_error check for os.Exit calls
// pkgfact      gather name/value pairs from constant declarations
// printf       check consistency of Printf format strings and arguments
// reflectvaluecompare check for comparing reflect.Value values with == or reflect.DeepEqual
// shadow       check for possible unintended shadowing of variables
// shift        check for shifts that equal or exceed the width of the integer
// sigchanyzer  check for unbuffered channel of os.Signal
// slog         check for invalid structured logging calls
// sortslice    check the argument type of sort.Slice
// stdmethods   check signature of methods of well-known interfaces
// stringintconv check for string(int) conversions
// structtag    check that struct field tags conform to reflect.StructTag.Get
// testinggoroutine report calls to (*testing.T).Fatal from goroutines started by a test
// tests        check for common mistaken usages of tests and examples
// timeformat   check for calls of (time.Time).Format or time.Parse with 2006-02-01
// unmarshal    report passing non-pointer or non-interface values to unmarshal
// unreachable  check for unreachable code
// unsafeptr    check for invalid conversions of uintptr to unsafe.Pointer
// unusedresult check for unused results of calls to some functions
// unusedwrite  checks for unused writes
// usesgenerics detect whether a package uses generics features
//
// List of analyzers: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes
func GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		appends.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpmux.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		slog.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
		osexit.OSExitAnalyzer,
	}
}
