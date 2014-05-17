package request

import (
	"github.com/franela/goblin"
	"testing"
)

func Test(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Get", func() {
		g.It("Errors when connection refused", func() {
			_, err := Get("http://localhost:8000", map[string]string{
				"foo": "bar",
			})
			g.Assert(err.Error()).Equal("Get http://localhost:8000: dial tcp 127.0.0.1:8000: connection refused")
		})
	})
}
