package component

type (
	component map[string]any

	// map[string]map[string]any -> "[component_key/field_key value] => { "tableheading" : "value" } "
	components map[string]component

	// [table_name](all the components)
	storage map[string]components
)
