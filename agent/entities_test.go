package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergedImagesKeepsOgENV(t *testing.T) {
	depB := DeploymentBlueprint{
		Template: Template{
			Images: []Image{
				{
					Name: "wp",
					Environment: []Environment{
						{Variable: "FOO", Value: "bar"},
					},
				},
			},
		},
		Override: Template{
			Images: []Image{
				{
					Name: "wp",
				},
			},
		},
	}

	mImgs := depB.MergedImages()

	e := []Environment{{Variable: "FOO", Value: "bar"}}

	assert.Equal(t, e, mImgs[0].Environment)
}

func TestMergedImagesOverridesENV(t *testing.T) {
	depB := DeploymentBlueprint{
		Template: Template{
			Images: []Image{
				{Name: "wp",
					Environment: []Environment{
						{Variable: "FOO", Value: "bar"},
					},
				},
			},
		},
		Override: Template{
			Images: []Image{
				{
					Name: "wp",
					Environment: []Environment{
						{Variable: "FOO", Value: "overridden"},
					},
				},
			},
		},
	}

	mImgs := depB.MergedImages()

	e := []Environment{{Variable: "FOO", Value: "overridden"}}

	assert.Equal(t, e, mImgs[0].Environment)
}

func TestMergedImagesExtraENVs(t *testing.T) {
	depB := DeploymentBlueprint{
		Template: Template{
			Images: []Image{
				{
					Name: "wp",
					Environment: []Environment{
						{Variable: "FOO", Value: "bar"},
					},
				},
			},
		},
		Override: Template{
			Images: []Image{
				{
					Name: "wp",
					Environment: []Environment{
						{Variable: "FOO", Value: "overridden"},
						{Variable: "MORE", Value: "stuff"},
					},
				},
			},
		},
	}

	mImgs := depB.MergedImages()

	e := []Environment{
		{Variable: "FOO", Value: "overridden"},
		{Variable: "MORE", Value: "stuff"},
	}

	assert.Equal(t, e, mImgs[0].Environment)
}
