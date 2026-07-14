package communication

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// This test pins the on-the-wire shape of the Manager<->Executor message set to
// the current WireProtocolVersion. It computes a structural fingerprint over the
// message envelope, the MessageType enum values, and every payload struct
// (recursing into the referenced pkg/common types, since those cross the wire
// too), then compares it against a golden fingerprint recorded for the current
// WireProtocolVersion.
//
// Any change to a wire struct's fields, JSON tags, field types, field order, or
// to the MessageType enum values changes the fingerprint and fails this test.
// That is intentional: an incompatible wire change MUST be accompanied by a bump
// of WireProtocolVersion (and a new golden entry below).
//
// If this test fails:
//   - If you did NOT mean to change the wire format, revert your struct change.
//   - If the change is intentional and breaks compatibility, bump
//     WireProtocolVersion in message.go and add a new entry to
//     goldenWireFingerprints with the fingerprint printed by this test.
//   - If the change is intentional and backward compatible (e.g. an added
//     `omitempty` field), you must still decide whether to bump; if you keep the
//     version, update the golden entry for the current version to the new value.
var goldenWireFingerprints = map[uint32]string{
	1: "2834d26567348bec72b0e08d9b44a0b3a58989f2179ad240ec70bf4b26fe2488",
}

// wireMessageTypes is the ordered list of MessageType enum values that travel on
// the wire. Their integer values are part of the wire contract: reordering or
// inserting a value shifts the enum and breaks compatibility. Add new message
// types here so the fingerprint covers them.
var wireMessageTypes = []struct {
	name string
	val  MessageType
}{
	{"ProcessRequestMessage", ProcessRequestMessage},
	{"ProcessResponseMessage", ProcessResponseMessage},
	{"DeployAppRequestMessage", DeployAppRequestMessage},
	{"DeployAppResponseMessage", DeployAppResponseMessage},
	{"GetKeysetRecoveryRequestMessage", GetKeysetRecoveryRequestMessage},
	{"GetKeysetRecoveryResponseMessage", GetKeysetRecoveryResponseMessage},
	{"SetKeysetRecoveryRequestMessage", SetKeysetRecoveryRequestMessage},
	{"SetKeysetRecoveryResponseMessage", SetKeysetRecoveryResponseMessage},
	{"KeysetRecoveryResultMessage", KeysetRecoveryResultMessage},
	{"AdminCommandRequestMessage", AdminCommandRequestMessage},
	{"AdminCommandResponseMessage", AdminCommandResponseMessage},
	{"ErrorMessage", ErrorMessage},
}

// wirePayloadTypes are the message envelope plus every payload struct exchanged
// between Manager and Executor. Keep this list complete: a payload struct that is
// not listed here is not covered by the fingerprint.
func wirePayloadTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(Message{}),
		reflect.TypeOf(ProcessRequestData{}),
		reflect.TypeOf(ProcessResponseData{}),
		reflect.TypeOf(DeployAppRequestData{}),
		reflect.TypeOf(DeployAppResponseData{}),
		reflect.TypeOf(ErrorData{}),
		reflect.TypeOf(GetKeysetRecoveryRequestData{}),
		reflect.TypeOf(GetKeysetRecoveryResponseData{}),
		reflect.TypeOf(SetKeysetRecoveryRequestData{}),
		reflect.TypeOf(SetKeysetRecoveryResponseData{}),
		reflect.TypeOf(KeysetRecoveryResultData{}),
		reflect.TypeOf(AdminCommandRequestData{}),
		reflect.TypeOf(AdminCommandResponseData{}),
	}
}

func TestWireFingerprintPinnedToProtocolVersion(t *testing.T) {
	golden, ok := goldenWireFingerprints[WireProtocolVersion]
	require.Truef(t, ok,
		"no golden wire fingerprint recorded for WireProtocolVersion=%d; add one to goldenWireFingerprints",
		WireProtocolVersion)

	got := computeWireFingerprint()

	require.Equalf(t, golden, got,
		"\n\nThe Manager<->Executor wire format changed but WireProtocolVersion (=%d) was not bumped.\n"+
			"If this change is intentional and breaks compatibility, bump WireProtocolVersion in\n"+
			"message.go and add this fingerprint to goldenWireFingerprints:\n\n\t%d: %q,\n\n"+
			"If it was accidental, revert the struct/enum change.\n",
		WireProtocolVersion, WireProtocolVersion, got)
}

// computeWireFingerprint builds a canonical, deterministic description of the
// wire contract and returns its SHA-256 hex digest.
func computeWireFingerprint() string {
	var sb strings.Builder

	sb.WriteString("messageTypes:\n")
	for _, mt := range wireMessageTypes {
		sb.WriteString(fmt.Sprintf("  %s=%d\n", mt.name, int(mt.val)))
	}

	sb.WriteString("payloads:\n")
	for _, tp := range wirePayloadTypes() {
		stack := map[reflect.Type]bool{}
		var fb strings.Builder
		describeType(tp, stack, &fb)
		sb.WriteString("  " + fb.String() + "\n")
	}

	sum := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(sum[:])
}

// describeType writes a canonical structural description of t into sb. It
// recurses through pointers, slices, arrays, maps and named structs so that a
// change anywhere in a type that crosses the wire is reflected in the output.
// Unexported fields are skipped (they are never serialized). A per-path stack
// guards against cyclic types.
func describeType(t reflect.Type, stack map[reflect.Type]bool, sb *strings.Builder) {
	for t.Kind() == reflect.Ptr {
		sb.WriteString("*")
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Struct:
		name := t.PkgPath() + "." + t.Name()
		if t.Name() == "" {
			name = "(anon)"
		}
		if stack[t] {
			sb.WriteString("@cycle:" + name)
			return
		}
		stack[t] = true
		sb.WriteString("struct " + name + "{")
		fields := make([]string, 0, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" { // unexported
				continue
			}
			var fsb strings.Builder
			fsb.WriteString(f.Name + " `" + f.Tag.Get("json") + "` ")
			describeType(f.Type, stack, &fsb)
			fields = append(fields, fsb.String())
		}
		// Field declaration order is part of the contract for reflection but not
		// for JSON; keep declaration order to catch reordering explicitly.
		sb.WriteString(strings.Join(fields, "; "))
		sb.WriteString("}")
		delete(stack, t)
	case reflect.Slice:
		sb.WriteString("[]")
		describeType(t.Elem(), stack, sb)
	case reflect.Array:
		sb.WriteString(fmt.Sprintf("[%d]", t.Len()))
		describeType(t.Elem(), stack, sb)
	case reflect.Map:
		sb.WriteString("map[")
		describeType(t.Key(), stack, sb)
		sb.WriteString("]")
		describeType(t.Elem(), stack, sb)
	case reflect.Interface:
		sb.WriteString("interface{" + t.PkgPath() + "." + t.Name() + "}")
	default:
		// Basic/named scalar: record the kind plus the named type (if any) so
		// that e.g. changing uint8 -> uint32 or swapping named types is caught.
		sb.WriteString(t.Kind().String())
		if t.Name() != "" && t.Name() != t.Kind().String() {
			sb.WriteString("(" + t.PkgPath() + "." + t.Name() + ")")
		}
	}
}

// TestWireFingerprintListsAreComplete parses message.go and asserts that every
// exported payload struct (name ending in "Data") and every MessageType constant
// declared there is covered by wirePayloadTypes / wireMessageTypes. Without this,
// a newly added message type appended after ErrorMessage — the most common way
// the protocol evolves — would be silently invisible to the fingerprint until a
// human remembered to add it to the lists, reintroducing exactly the discipline
// failure this guard exists to remove.
func TestWireFingerprintListsAreComplete(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "message.go", nil, 0)
	require.NoError(t, err)

	coveredPayloads := map[string]bool{}
	for _, tp := range wirePayloadTypes() {
		coveredPayloads[tp.Name()] = true
	}
	coveredMsgTypes := map[string]bool{}
	for _, mt := range wireMessageTypes {
		coveredMsgTypes[mt.name] = true
	}

	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		switch gd.Tok {
		case token.TYPE:
			for _, spec := range gd.Specs {
				ts := spec.(*ast.TypeSpec)
				if !ts.Name.IsExported() {
					continue
				}
				if _, isStruct := ts.Type.(*ast.StructType); !isStruct {
					continue
				}
				if !strings.HasSuffix(ts.Name.Name, "Data") {
					continue
				}
				require.Truef(t, coveredPayloads[ts.Name.Name],
					"payload struct %s is declared in message.go but missing from wirePayloadTypes()", ts.Name.Name)
			}
		case token.CONST:
			// A const block declares MessageType values when its first spec is
			// typed `MessageType` (the rest inherit the type via iota).
			if len(gd.Specs) == 0 {
				continue
			}
			first, ok := gd.Specs[0].(*ast.ValueSpec)
			if !ok || first.Type == nil {
				continue
			}
			ident, ok := first.Type.(*ast.Ident)
			if !ok || ident.Name != "MessageType" {
				continue
			}
			for _, spec := range gd.Specs {
				vs := spec.(*ast.ValueSpec)
				for _, name := range vs.Names {
					require.Truef(t, coveredMsgTypes[name.Name],
						"MessageType constant %s is declared in message.go but missing from wireMessageTypes", name.Name)
				}
			}
		}
	}
}

// TestWireFingerprintDetectsChange is a self-check that the fingerprint function
// actually reacts to structural differences, so the guard above cannot silently
// rot into a no-op.
func TestWireFingerprintDetectsChange(t *testing.T) {
	stack := map[reflect.Type]bool{}
	var a strings.Builder
	describeType(reflect.TypeOf(ProcessRequestData{}), stack, &a)

	type altProcessRequestData struct {
		Request          *int   `json:"request"`
		ApplicationState string `json:"applicationState"`
	}
	stack = map[reflect.Type]bool{}
	var b strings.Builder
	describeType(reflect.TypeOf(altProcessRequestData{}), stack, &b)

	require.NotEqual(t, a.String(), b.String(),
		"fingerprint must differ when struct fields differ")
	require.NotEmpty(t, a.String())
	// Guard against accidental non-determinism from map iteration in describeType:
	// two independent computations must agree.
	require.Equal(t, computeWireFingerprint(), computeWireFingerprint())
}
