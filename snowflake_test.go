package snowflake

import "testing"

func BenchmarkGenerate(b *testing.B) {
	node,_ := NewNode(1)
	for i:=0; i<b.N ; i++  {
		node.Generate()
	}
}
