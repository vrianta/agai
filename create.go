package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/vrianta/agai/v1/config"
)

func create_controller() {

	for _, controller_name := range f.controller_names_to_create {
		fmt.Printf("🔧 Creating controller: %s\n", controller_name)

		package_name := strings.ToLower(controller_name)
		view_name := ""
		if f.create_view {
			view_name = package_name
			f.view_names_to_create = append(f.view_names_to_create, view_name)
		}

		// Set output location: controller/controller_name/controller_name.controller.go
		controller_output_location := fmt.Sprintf("%s/%s/%s.controller.go", f.controllers_root, controller_name, controller_name)

		if file_info, err := os.Stat(controller_output_location); file_info != nil && err == nil {
			fmt.Printf("⚠️  Skipped: Controller '%s' already exists at %s\n", controller_name, controller_output_location)
			continue
		}

		// Read the template from embed
		controller_template, err := templates.ReadFile("templates/controller.go.template")
		if err != nil {
			fmt.Printf("❌ Error: Failed to read controller template: %v\n", err)
			return
		}

		targetDir := filepath.Join(f.controllers_root, controller_name)

		// Create controller directory
		fmt.Printf("📁 Creating directory: %s\n", targetDir)
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			fmt.Printf("❌ Error: Could not create directory %s: %v\n", targetDir, err)
			return
		}

		// Parse and render the template
		tmpl, err := template.New(controller_name).Parse(string(controller_template))
		if err != nil {
			fmt.Printf("❌ Error: Template parse failed for %s: %v\n", "controller.go.template", err)
			return
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{
			"package_name":    package_name,
			"controller_name": capitalize(controller_name),
			"view_name":       view_name,
		})
		if err != nil {
			fmt.Printf("❌ Error: Template execution failed for %s: %v\n", "controller.go.template", err)
			return
		}

		// Write the final file
		err = os.WriteFile(controller_output_location, buf.Bytes(), 0644)
		if err != nil {
			fmt.Printf("❌ Error: Could not write controller file to %s: %v\n", controller_output_location, err)
			return
		}

		fmt.Printf("✅ Controller '%s' created at %s\n\n", controller_name, controller_output_location)
	}
}

/*
Create View :
it will creat the folder views if that does not exists
then it will create a subfolder with the view name but if the folder exists int will log.Error that the view is already exists and return the function
inside that it will create a file called index.php
*/
func create_view() {
	for _, view_name := range f.view_names_to_create {
		fmt.Printf("🧩 Creating view: %s\n", view_name)

		viewRoot := config.GetWebConfig().ViewFolder
		viewDir := filepath.Join(viewRoot, view_name)
		viewFile := filepath.Join(viewDir, "index.php")

		// Check if view already exists
		if fileInfo, err := os.Stat(viewDir); err == nil && fileInfo.IsDir() {
			fmt.Printf("⚠️  Skipped: View '%s' already exists at %s\n", view_name, viewDir)
			continue
		}

		// Read the view template from embedded FS
		viewTemplate, err := templates.ReadFile("templates/index.php.template")
		if err != nil {
			fmt.Printf("❌ Error: Failed to read view template: %v\n", err)
			return
		}

		// Create the view directory
		fmt.Printf("📁 Creating directory: %s\n", viewDir)
		if err := os.MkdirAll(viewDir, os.ModePerm); err != nil {
			fmt.Printf("❌ Error: Could not create view directory %s: %v\n", viewDir, err)
			return
		}

		// Parse and render the template
		tmpl, err := template.New(view_name).Parse(string(viewTemplate))
		if err != nil {
			fmt.Printf("❌ Error: Template parse failed for %s: %v\n", "index.php.template", err)
			return
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{
			"view_name": capitalize(view_name),
		})
		if err != nil {
			fmt.Printf("❌ Error: Template execution failed for %s: %v\n", "index.php.template", err)
			return
		}

		// Write index.php to view folder
		err = os.WriteFile(viewFile, buf.Bytes(), 0644)
		if err != nil {
			fmt.Printf("❌ Error: Could not write view file to %s: %v\n", viewFile, err)
			return
		}

		fmt.Printf("✅ View '%s' created at %s\n\n", view_name, viewFile)
	}
}

func create_models() {
	for _, model_name := range f.model_names_to_create {
		fmt.Printf("🛠️  Creating model: %s\n", model_name)

		model_output_path := fmt.Sprintf("models/%s.model.go", strings.ToLower(model_name))

		// Skip if model file already exists
		if file_info, err := os.Stat(model_output_path); file_info != nil && err == nil {
			fmt.Printf("⚠️  Skipped: Model '%s' already exists at %s\n", model_name, model_output_path)
			continue
		}

		if f.create_component {
			f.component_names_to_create = append(f.component_names_to_create, model_name)
		}

		// Read model template
		model_template, err := templates.ReadFile("templates/model.go.template")
		if err != nil {
			fmt.Printf("❌ Error: Failed to read model template: %v\n", err)
			return
		}

		// Parse template
		tmpl, err := template.New(model_name).Parse(string(model_template))
		if err != nil {
			fmt.Printf("❌ Error: Template parse failed: %v\n", err)
			return
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{
			"model_name": capitalize(model_name),
		})
		if err != nil {
			fmt.Printf("❌ Error: Template execution failed: %v\n", err)
			return
		}

		// Ensure models directory exists
		if err := os.MkdirAll("models", os.ModePerm); err != nil {
			fmt.Printf("❌ Error: Failed to create models directory: %v\n", err)
			return
		}

		// Write model file
		err = os.WriteFile(model_output_path, buf.Bytes(), 0644)
		if err != nil {
			fmt.Printf("❌ Error: Failed to write model file: %v\n", err)
			return
		}

		fmt.Printf("✅ Model '%s' created at: %s\n\n", model_name, model_output_path)
	}
}

/*
First it will check if the components folder is present or not if not then create it
Then it will check if the component alreadt present or not if present then skip the creation and log error and return
Then check if the same named model is present in the models folder or not
Then eavluate the model in the file and craete component according to that
*/
func create_components() {
	for _, componentName := range f.component_names_to_create {
		fmt.Printf("🧩 Creating component: %s\n", componentName)

		componentFile := filepath.Join("components", fmt.Sprintf("%s.component.json", strings.ToLower(componentName)))
		modelFile := filepath.Join("models", fmt.Sprintf("%s.model.go", strings.ToLower(componentName)))

		// Ensure components/ directory exists
		if err := os.MkdirAll("components", os.ModePerm); err != nil {
			fmt.Printf("❌ Failed to create components directory: %v\n", err)
			continue
		}

		// Check if component already exists
		if _, err := os.Stat(componentFile); err == nil {
			fmt.Printf("⚠️  Skipped: Component already exists at %s\n", componentFile)
			continue
		}

		// Check if model file exists
		modelContent, err := os.ReadFile(modelFile)
		if err != nil {
			fmt.Printf("❌ Error: Model file not found for component '%s' (%s)\n", componentName, modelFile)
			continue
		}

		// Parse model file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, "", modelContent, parser.AllErrors)
		if err != nil {
			fmt.Printf("❌ Error parsing model file for %s: %v\n", componentName, err)
			continue
		}

		// Extract struct field names from the model.New call
		fields := map[string]interface{}{}

		ast.Inspect(node, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) != 2 {
				return true
			}

			ident, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || ident.Sel.Name != "New" {
				return true
			}

			composite, ok := call.Args[1].(*ast.CompositeLit)
			if !ok {
				return true
			}

			for _, elt := range composite.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key := fmt.Sprintf("%s", kv.Key)
				fields[key] = defaultValueForExpr(kv.Key)
			}
			return false
		})

		// Create JSON map
		componentData := map[string]interface{}{
			"0": fields,
		}

		jsonBytes, err := json.MarshalIndent(componentData, "", "  ")
		if err != nil {
			fmt.Printf("❌ Failed to marshal JSON for component '%s': %v\n", componentName, err)
			continue
		}

		if err := os.WriteFile(componentFile, jsonBytes, 0644); err != nil {
			fmt.Printf("❌ Failed to write component file: %v\n", err)
			continue
		}

		fmt.Printf("✅ Component '%s' created at %s\n\n", componentName, componentFile)
	}
}

// defaultValueForExpr returns a Go zero value for a given field expression
func defaultValueForExpr(expr ast.Expr) interface{} {
	switch v := expr.(type) {
	case *ast.CompositeLit:
		return ""
	case *ast.BasicLit:
		if strings.Contains(v.Value, "\"") {
			return ""
		}
		return 0
	case *ast.Ident:
		if v.Name == "true" || v.Name == "false" {
			return false
		}
		return ""
	default:
		return ""
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
