package functions

import (
	"fmt"
	"github.com/yjhatfdu/expr/types"
)

func init() {
	// multiIf(if,then,else)
	// multiIf(case1,then1,case2,then2....)
	// multiIf(case1,then1,case2,then2....,default)
	multiIf, _ := NewFunction("multiIf")
	multiIf.Generic(func(inputTypes []types.BaseType) (baseType types.BaseType, e error) {
		if len(inputTypes) < 2 {
			return 0, fmt.Errorf("multiIf should have at least 2 arguments got (%d) argument", len(inputTypes))
		}
		outType := inputTypes[1]
		for i := 0; i < len(inputTypes)/2; i++ {
			t := inputTypes[2*i+1]
			if t != outType && t+types.ScalaOffset != outType && t-types.ScalaOffset != outType {
				return 0, fmt.Errorf("argument #%d should be %s, got %s", 2*i+1, types.GetTypeName(outType), types.GetTypeName(t))
			}
		}

		//todo 使用函数进行类型比较
		lastType := inputTypes[len(inputTypes)-1]
		if lastType != outType && lastType+types.ScalaOffset != outType && lastType-types.ScalaOffset != outType {
			return 0, fmt.Errorf("argument #%d should be %s, got %s", len(inputTypes)-1, types.GetTypeName(outType), types.GetTypeName(lastType))
		}
		return outType, nil
	}, func(vectors []types.INullableVector) (vector types.INullableVector, e error) {
		var hasDefault bool
		if len(vectors)%2 == 1 {
			hasDefault = true
		}
		l := len(vectors)
		filterMasks := make([][]bool, l)
		for i := range vectors {
			filterMasks[i] = vectors[i].GetFilterArr()
		}
		out, err := BroadCastMultiGeneric(vectors, vectors[1].Type(), func(values []interface{}, index int) (i interface{}, e error) {
			for i := 0; i < l/2; i++ {
				v := values[2*i]
				var truthy bool
				if v == nil {
					continue
				}
				switch vi := v.(type) {
				case int64:
					truthy = vi != 0
				case float64:
					truthy = vi != 0
				case string:
					truthy = vi != ""
				case bool:
					truthy = vi
				}
				if truthy {
					return values[2*i+1], nil
				}
			}
			if hasDefault {
				return values[len(values)-1], nil
			} else {
				return nil, nil
			}
		})
		if err != nil {
			return nil, err
		}
		out.SetFilterArr(CalFilterMask(filterMasks))
		return out, nil
	})
}
