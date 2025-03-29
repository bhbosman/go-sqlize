package internal

//func (self *Compiler) addSomeType() func(state []ICompileState, types []reflect.Type) reflect.Type {
//	return func(state []ICompileState, types []reflect.Type) reflect.Type {
//		sf := []reflect.StructField{
//			{"Data", "", types[0], reflect.StructTag(""), 0, nil, false},
//			{"Assigned", "", reflect.TypeFor[bool](), reflect.StructTag(""), 0, nil, false},
//		}
//		rt := reflect.StructOf(sf)
//		self.GlobalMethodHandlers[GlobalMethodHandlerKey{rt, "Set"}] =
//			func(types []reflect.Type, rt reflect.Type) reflect.Value {
//				rtIn := []reflect.Type{rt, types[0]}
//				rtOut := []reflect.Type{rt}
//				rtFunc := reflect.FuncOf(rtIn, rtOut, false)
//				return reflect.MakeFunc(
//					rtFunc,
//					func(args []reflect.Value) (results []reflect.Value) {
//						resultRt := self.GlobalTypes[someTypeValueKey](state, types)
//						resultRv := reflect.New(resultRt).Elem()
//						resultRv.Field(0).Set(args[1])
//						resultRv.Field(1).Set(reflect.ValueOf(true))
//						return []reflect.Value{resultRv}
//					},
//				)
//			}(types, rt)
//		self.GlobalMethodHandlers[GlobalMethodHandlerKey{rt, "ToNone"}] =
//			func(types []reflect.Type, rt reflect.Type) reflect.Value {
//				rtIn := []reflect.Type{rt}
//				rtOut := []reflect.Type{rt}
//				rtFunc := reflect.FuncOf(rtIn, rtOut, false)
//				return reflect.MakeFunc(
//					rtFunc,
//					func(args []reflect.Value) (results []reflect.Value) {
//						// Note: state may not be correct here
//						resultRt := self.GlobalTypes[someTypeValueKey](state, types)
//						resultRv := reflect.New(resultRt).Elem()
//						resultRv.Field(0).Set(reflect.Zero(types[0]))
//						resultRv.Field(1).Set(reflect.ValueOf(false))
//						return []reflect.Value{resultRv}
//					},
//				)
//			}(types, rt)
//		self.GlobalMethodHandlers[GlobalMethodHandlerKey{rt, "Get"}] =
//			func(types []reflect.Type, rt reflect.Type) reflect.Value {
//				rtIn := []reflect.Type{rt}
//				rtOut := []reflect.Type{types[0]}
//				rtFunc := reflect.FuncOf(rtIn, rtOut, false)
//				return reflect.MakeFunc(
//					rtFunc,
//					func(args []reflect.Value) (results []reflect.Value) {
//						rv := args[0]
//						if rv.Field(1).Equal(reflect.ValueOf(true)) {
//							return []reflect.Value{rv.Field(0)}
//						}
//						return []reflect.Value{reflect.Zero(types[0])}
//					},
//				)
//			}(types, rt)
//		self.GlobalMethodHandlers[GlobalMethodHandlerKey{rt, "IsAssigned"}] =
//			func(types []reflect.Type, rt reflect.Type) reflect.Value {
//				rtIn := []reflect.Type{rt}
//				rtOut := []reflect.Type{reflect.TypeFor[bool]()}
//				rtFunc := reflect.FuncOf(rtIn, rtOut, false)
//				return reflect.MakeFunc(
//					rtFunc,
//					func(args []reflect.Value) (results []reflect.Value) {
//						return []reflect.Value{args[0].Field(1)}
//					},
//				)
//			}(types, rt)
//		return rt
//	}
//}
