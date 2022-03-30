package manager

import (
	"fmt"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/cloudbeat/resources/fetching"
)

type FetcherCtxCreatorStrategy interface {
	Register(fetching.ContextCreator)
	CreateContext(fetcherName string) (fetching.FetcherCtx, error)
}

type FetcherCtxCreator struct {
	ctxCreators map[string]fetching.ContextCreator
}

func (en *FetcherCtxCreator) Register(enricher fetching.ContextCreator) {
	name := enricher.GetName()
	_, ok := en.ctxCreators[name]
	if ok {
		logp.L().Warnf("fetcher %q factory method overwritten", name)
	}

	en.ctxCreators[name] = enricher
}

func (en *FetcherCtxCreator) CreateContext(fetcherName string) (fetching.FetcherCtx, error) {
	extraElements := make(fetching.FetcherCtx)
	for _, ctxCreator := range en.ctxCreators {
		ctxName := ctxCreator.GetContextName()
		ctx, err := ctxCreator.GetContext(fetcherName)
		if err != nil {
			return nil, fmt.Errorf("could not add context %s to request: %w", ctxName, err)
		}
		extraElements[ctxName] = ctx
	}

	return extraElements, nil
}
