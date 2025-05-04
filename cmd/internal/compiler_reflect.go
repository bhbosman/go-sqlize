package internal

func (compiler *Compiler) addReflectFunctions() {
	compiler.GlobalTypes[ValueKey{"reflect", "Type"}] = compiler.registerReflectType()
}
