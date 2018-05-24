package youi

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseYUMLSimple(t *testing.T) {
	const simple = `
<Page xmlns="https://yuml.ovo.ovh/schema/components/1.0">
	<Canvas X="10" Y="10" Width="100" Height="100">
		<Image src="out.png" />
	</Canvas>
</Page>
`
	out, err := parseYUML(strings.NewReader(simple))
	if err != nil {
		t.Error(err)
		return
	}

	const expected = `https://yuml.ovo.ovh/schema/components/1.0:Page (xmlns=https://yuml.ovo.ovh/schema/components/1.0)
  ˪ https://yuml.ovo.ovh/schema/components/1.0:Canvas (X=10, Y=10, Width=100, Height=100)
    ˪ https://yuml.ovo.ovh/schema/components/1.0:Image (src=out.png)`

	// Check generated tree
	if out.String() != expected {
		t.Error("generated tree doesn't match expected result")
		fmt.Printf("Expected:\n%s\n\nGot:\n%s\n\n", expected, out)
		return
	}
}
