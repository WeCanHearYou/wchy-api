package web_test

import (
	"net/url"
	"testing"

	"github.com/getfider/fider/app/models/entity"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/web"
)

func TestReactRenderer_FileNotFound(t *testing.T) {
	RegisterT(t)

	r, err := web.NewReactRenderer("unknown.js")
	Expect(err).IsNotNil()
	Expect(r).IsNil()
}

func TestReactRenderer_EmptyFile(t *testing.T) {
	RegisterT(t)

	r, err := web.NewReactRenderer("/app/pkg/web/testdata/empty.js")
	Expect(err).IsNil()
	Expect(r).IsNotNil()

	u, _ := url.Parse("https://github.com")
	html, err := r.Render(u, web.Map{})
	Expect(html).Equals("")
	Expect(err).IsNil()
}

func TestReactRenderer_RenderEmptyHomeHTML(t *testing.T) {
	RegisterT(t)

	r, err := web.NewReactRenderer("ssr.js")
	Expect(err).IsNil()

	u, _ := url.Parse("https://demo.test.fider.io")
	html, err := r.Render(u, web.Map{
		"tenant": &entity.Tenant{
			Locale: "en",
		},
		"settings": web.Map{
			"locale": "en",
		},
		"props": web.Map{
			"posts":          make([]web.Map, 0),
			"tags":           make([]web.Map, 0),
			"countPerStatus": web.Map{},
		},
	})
	Expect(html).ContainsSubstring(`<div class="c-dev-banner">DEV</div>`)
	Expect(html).ContainsSubstring(`<input type="text" class="c-input" id="input-title"`)
	Expect(html).ContainsSubstring(`What can we do better? This is the place for you to vote, discuss and share ideas.`)
	Expect(html).ContainsSubstring(`No posts have been created yet.`)
	Expect(html).ContainsSubstring(`Powered by Fider`)
	Expect(err).IsNil()
}

func TestReactRenderer_RenderEmptyHomeHTML_Portuguese(t *testing.T) {
	RegisterT(t)

	r, err := web.NewReactRenderer("ssr.js")
	Expect(err).IsNil()

	u, _ := url.Parse("https://demo.test.fider.io")
	html, err := r.Render(u, web.Map{
		"tenant": &entity.Tenant{
			Locale: "pt-BR",
		},
		"settings": web.Map{
			"locale": "pt-BR",
		},
		"props": web.Map{
			"posts":          make([]web.Map, 0),
			"tags":           make([]web.Map, 0),
			"countPerStatus": web.Map{},
		},
	})
	Expect(html).ContainsSubstring(`<div class="c-dev-banner">DEV</div>`)
	Expect(html).ContainsSubstring(`<input type="text" class="c-input" id="input-title"`)
	Expect(html).ContainsSubstring(`O que podemos fazer melhor? Este é o lugar para você votar, discutir e compartilhar ideias.`)
	Expect(html).ContainsSubstring(`Nenhuma postagem foi criada ainda.`)
	Expect(html).ContainsSubstring(`Powered by Fider`)
	Expect(err).IsNil()
}
