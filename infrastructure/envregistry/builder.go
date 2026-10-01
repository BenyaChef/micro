package envregistry

import envregistrymodel "github.com/BenyaChef/micro/infrastructure/envregistry/model"

type Builder struct {
	source envregistrymodel.Source
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Source(source envregistrymodel.Source) *Builder {
	b.source = source

	return b
}

func (b *Builder) Build() (*Registry, error) {
	source := b.source
	if source == nil {
		source = envregistrymodel.OSSource
	}

	return &Registry{source: source}, nil
}
