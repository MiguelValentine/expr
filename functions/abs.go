package functions

import (
	"github.com/yjhatfdu/expr/types"
)

func absInt(data1, out []int64) {
	for i := range out {
		if data1[i] < 0 {
			out[i] = -data1[i]
		} else {
			out[i] = data1[i]
		}
	}
}

func absFloat(data1, out []float64) {
	for i := range out {
		if data1[i] < 0 {
			out[i] = -data1[i]
		} else {
			out[i] = data1[i]
		}
	}
}

func init() {
	addFunc, _ := NewFunction("abs")
	addFunc.Overload([]types.BaseType{types.Int}, types.Int, func(vectors []types.INullableVector) (types.INullableVector, error) {
		in := vectors[0].(*types.NullableInt)
		output := types.NullableInt{}
		output.Init(len(in.Values))
		output.FilterArr = in.FilterArr
		output.IsNullArr = in.IsNullArr
		absInt(in.Values, output.Values)
		return &output, nil
	})
	addFunc.Overload([]types.BaseType{types.Float}, types.Float, func(vectors []types.INullableVector) (types.INullableVector, error) {
		in := vectors[0].(*types.NullableFloat)
		output := types.NullableFloat{}
		output.Init(len(in.Values))
		output.FilterArr = in.FilterArr
		output.IsNullArr = in.IsNullArr
		absFloat(in.Values, output.Values)
		return &output, nil
	})
}
