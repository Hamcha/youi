package youi

func init() {
	namespaces = make(map[string]componentList)
	initBuiltinComponents()
}
