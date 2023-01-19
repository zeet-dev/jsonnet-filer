package cabinet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zeet-dev/jsonnet-filer/pkg/api/v1alpha1"
)

type MockOker struct {
	ok bool
}

func (m MockOker) Ok() bool {
	return m.ok
}

func TestFind(t *testing.T) {
	type testCase[T Oker] struct {
		name        string
		input       any
		wantResults []T
	}
	tests := []testCase[MockOker]{
		{
			name: "given empty input, expect empty results",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults := Find[MockOker](tt.input)
			assert.Equal(t, tt.wantResults, gotResults)
		})
	}
}

func TestIs_string(t *testing.T) {
	type testCase[T any] struct {
		name    string
		input   any
		want    T
		wantErr bool
	}

	tests := []testCase[string]{
		{
			name:  "given string, expect string and no error",
			input: "foo",
			want:  "foo",
		},
		{
			name:    "given int, expect empty string and an error",
			input:   123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Is[string](tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIs_File(t *testing.T) {
	type testCase[T any] struct {
		name    string
		input   any
		want    T
		wantErr bool
	}

	tests := []testCase[v1alpha1.File]{
		{
			name: "given map that looks like a File, expect v1alpha1.File and no error",
			input: map[string]any{
				"apiVersion": v1alpha1.ApiVersion,
				"kind":       v1alpha1.Kind,
				"metadata": map[string]any{
					"name": "foo",
				},
				"content":          "this is a test",
				"encodingStrategy": "yaml",
			},
			want: v1alpha1.NewFile("foo", "this is a test"),
		},
		{
			name: "given map that looks like a File but has extra fields, expect an error",
			input: map[string]any{
				"apiVersion": v1alpha1.ApiVersion,
				"kind":       v1alpha1.Kind,
				"metadata": map[string]any{
					"name": "foo",
				},
				"content":              "this is a test",
				"encodingStrategy":     "yaml",
				"extraShouldNotBeHere": "foo",
			},
			wantErr: true,
		},
		{
			name: "given map that looks like a File, but with missing values, expect File and no error",
			input: map[string]any{
				"apiVersion": v1alpha1.ApiVersion,
				"kind":       v1alpha1.Kind,
				"metadata": map[string]any{
					"name": "foo",
				},
			},
			want: func() v1alpha1.File {
				f := v1alpha1.NewFile("foo", nil)
				f.EncodingStrategy = ""

				return f
			}(),
		},
		{
			name: "given map that is totally different, an error",
			input: map[string]any{
				"foo": "bar",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Is[v1alpha1.File](tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
