package request

import (
	"github.com/franela/goblin"
	"testing"
)

func Test(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Get", func() {
		g.It("Errors when connection refused", func() {
			_, err := Get("http://localhost", map[string]string{
				"foo": "bar",
			})
			g.Assert(err.Error()).Equal("Get http://localhost: dial tcp 127.0.0.1:80: connection refused")
		})
	})
}
