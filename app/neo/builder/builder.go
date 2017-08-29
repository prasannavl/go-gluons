package builder

import (
	"net/http"

	"pvl/apicore/app/neo"

	"pvl/apicore/app/neo/hconv"
	"pvl/apicore/app/neo/mconv"
)

type ChainBuilder struct {
	chain   neo.Chain
	handler neo.Handler
}

func Create(middlewares ...neo.Middleware) ChainBuilder {
	return ChainBuilder{neo.Chain{Middlewares: middlewares}, nil}
}

func (b *ChainBuilder) Add(m ...neo.Middleware) *ChainBuilder {
	b.chain.Middlewares = append(b.chain.Middlewares, m...)
	return b
}

func (b *ChainBuilder) AddSimple(m ...neo.SimpleMiddleware) *ChainBuilder {
	s := make([]neo.Middleware, 0, len(m))
	for _, x := range m {
		s = append(s, mconv.FromSimple(x))
	}
	b.chain.Middlewares = append(b.chain.Middlewares, s...)
	return b
}

func (b *ChainBuilder) Handler(finalHandler neo.Handler) *ChainBuilder {
	b.handler = finalHandler
	return b
}

func (b *ChainBuilder) Build() neo.Handler {
	h := b.handler
	if h == nil {
		h = hconv.FromHttp(http.DefaultServeMux)
	}
	c := b.chain
	mx := c.Middlewares
	mLen := len(mx)
	for i := range mx {
		h = c.Middlewares[mLen-1-i](h)
	}
	return h
}

func (b *ChainBuilder) BuildHttpFactory(errorHandler func(error)) func(neo.Context) http.Handler {
	return hconv.ToHttpFactory(b.Build(), errorHandler)
}
