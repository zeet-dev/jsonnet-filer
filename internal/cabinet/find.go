package cabinet

import (
	"github.com/mitchellh/mapstructure"
)

type Oker interface {
	Ok() bool
}

func Is[T any](input any) (T, error) {
	var output T
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		Result:      &output,
	})

	if err != nil {
		return output, err
	}

	err = decoder.Decode(input)
	return output, err
}

func Find[T Oker](input any) (results []T) {
	switch v := input.(type) {
	case map[string]any:
		result, err := Is[T](input)
		if err == nil && result.Ok() {
			results = append(results, result)
			return
		}

		for _, potential := range v {
			results = append(results, Find[T](potential)...)
		}
	case []any:
		for _, potential := range v {
			results = append(results, Find[T](potential)...)
		}
	default:
		// we don't care about any other values than map types
	}

	return results
}
