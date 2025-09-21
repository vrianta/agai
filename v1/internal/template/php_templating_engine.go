package template

import (
	"fmt"
	"regexp"
	"strings"
)

// PHPToGoTemplate converts PHP template code to Go html/template syntax
func PHPToGoTemplate(phpTemplate string) string {
	// Step 1: Convert short echo tags <?= ... ?> to Go template syntax
	shortEchoPattern := regexp.MustCompile(`<\?=\s*(.*?)\s*\?>`)
	phpTemplate = shortEchoPattern.ReplaceAllStringFunc(phpTemplate, func(match string) string {
		expr := shortEchoPattern.FindStringSubmatch(match)[1]
		return "{{ " + convertPHPExprToGo(expr) + " }}"
	})

	// Step 2: Convert PHP blocks <?php ... ?>
	phpBlockPattern := regexp.MustCompile(`<\?php\s*(.*?)\s*\?>`)
	phpTemplate = phpBlockPattern.ReplaceAllStringFunc(phpTemplate, func(match string) string {
		code := phpBlockPattern.FindStringSubmatch(match)[1]
		return convertPHPBlockToGo(code)
	})

	// fmt.Println(phpTemplate)
	return phpTemplate
}

func convertPHPExprToGo(expr string) string {
	// Step 1: Replace -> with . for property access
	expr = strings.ReplaceAll(expr, "->", ".")

	// Step 2: Handle function calls first
	functionCallPattern := regexp.MustCompile(`(\w+)\((.*?)\)`)
	expr = functionCallPattern.ReplaceAllStringFunc(expr, func(match string) string {
		parts := functionCallPattern.FindStringSubmatch(match)
		funcName := parts[1]
		args := convertPHPVarsToGo(parts[2])

		switch funcName {
		case "strtoupper":
			return fmt.Sprintf("Upper %s", args)
		case "strtolower":
			return fmt.Sprintf("lower %s", args)
		case "strlen":
			return fmt.Sprintf("(Strlen %s)", args)
		case "count":
			return fmt.Sprintf("(Len %s)", args)
		case "htmlspecialchars":
			return fmt.Sprintf("Html %s", args)
		case "isset":
			return fmt.Sprintf("ne %s nil", args)
		case "empty":
			return fmt.Sprintf("eq %s \"\"", args)
		case "print":
			return fmt.Sprintf("print %s", args)
		default:
			return fmt.Sprintf("call .%s %s", funcName, args)
		}
	})

	// Step 3: Replace PHP-style vars like $x or $$x with .x or .x (same scope)
	expr = convertPHPVarsToGo(expr)

	// Step 4: Optional cleanup (for readability)
	expr = strings.ReplaceAll(expr, "&&", "and")
	expr = strings.ReplaceAll(expr, "||", "or")
	expr = strings.ReplaceAll(expr, "!", "not ")

	return expr
}

// convertPHPVarsToGo converts PHP-style variable expressions into Go template syntax.
//
// It handles the following transformations:
// - Object properties:         $$object->property        -> $object.property
// - Array (map) access:        $array['key']            -> array.key
// - Indexed arrays:            $array[0]                -> index array 0
// - Simple variables:          $var                     -> .var
// - Operators:
//     - ==                    -> eq
//     - !=                    -> ne
//     - <                     -> lt
//     - >                     -> gt
//     - <=                    -> le
//     - >=                    -> ge
//     - &&                    -> and
//     - ||                    -> or
//     - !                     -> not

func convertPHPVarsToGo(expr string) string {
	// Convert object property access: $object->property -> $object.property
	// fmt.Printf("190 - %s\n", expr)
	objPropPattern_doubledollar := regexp.MustCompile(`\$\$(\w+)->(\w+)`)
	if objPropPattern_doubledollar.MatchString(expr) {
		matches := objPropPattern_doubledollar.FindStringSubmatch(expr)
		if len(matches) > 2 {
			objName := matches[1]
			propName := matches[2]
			// fmt.Printf("197 - $%s.%s\n", objName, propName)
			return fmt.Sprintf(".%s.%s", objName, propName)
		}
	}
	objPropPattern := regexp.MustCompile(`\$(\w+)->(\w+)`)
	if objPropPattern.MatchString(expr) {
		matches := objPropPattern.FindStringSubmatch(expr)
		if len(matches) > 2 {
			objName := matches[1]
			propName := matches[2]
			// fmt.Printf("197 - $%s.%s\n", objName, propName)
			return fmt.Sprintf("$%s.%s", objName, propName)
		}
	}

	// Convert array access: $array['key'] -> array.key
	arrayAccessPattern := regexp.MustCompile(`\$(\w+)\[\s*['"](\w+)['"]\s*\]`)
	expr = arrayAccessPattern.ReplaceAllString(expr, "${1}.${2}")

	// Convert numeric array index: $array[0] -> index array 0
	numArrayAccessPattern := regexp.MustCompile(`\$(\w+)\[\s*(\d+)\s*\]`)
	expr = numArrayAccessPattern.ReplaceAllString(expr, "index ${1} ${2}")

	// Step 1: Temporarily replace $$var with a placeholder like __DOLLAR_VAR_var__
	doubleDollarPattern := regexp.MustCompile(`\$\$(\w+)`)
	expr = doubleDollarPattern.ReplaceAllString(expr, `.${1}`)

	// Step 2: Convert $var to .var
	simpleVarPattern := regexp.MustCompile(`\$(\w+)`)
	expr = simpleVarPattern.ReplaceAllString(expr, "$$${1}")

	restorePlaceholderPattern := regexp.MustCompile(`@@(\w+)`)
	expr = restorePlaceholderPattern.ReplaceAllString(expr, `$$$1`)

	// Convert PHP logical/comparison operators to Go template equivalents
	// ==  -> eq
	// !=  -> ne
	// <   -> lt
	// >   -> gt
	// <=  -> le
	// >=  -> ge
	// &&  -> and
	// ||  -> or
	// !   -> not
	re := regexp.MustCompile(`(\S+)\s*(==|!=|<=|>=|<|>)\s*(\S+)`)
	opMap := map[string]string{
		"==": "eq",
		"!=": "ne",
		"<=": "le",
		">=": "ge",
		"<":  "lt",
		">":  "gt",
	}

	// Replace each match with the corresponding prefix
	return re.ReplaceAllStringFunc(expr, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) == 4 {
			left, op, right := parts[1], parts[2], parts[3]
			if newOp, ok := opMap[op]; ok {
				// fmt.Printf("%s", right)
				return fmt.Sprintf("%s %s %s", newOp, left, right)
			}
		}
		return match
	})
	// return expr
}

// convertPHPBlockToGo converts a PHP code block to Go template syntax
func convertPHPBlockToGo(code string) string {
	// Handle foreach loops
	if strings.HasPrefix(code, "foreach") {
		// fmt.Println("Found Foreach 270", code)
		foreachPattern := regexp.MustCompile(`foreach\s*\(\s*(\${1,2}[a-zA-Z_][a-zA-Z0-9_]*(?:->\w+)*)\s+as\s+(\${1,2}[a-zA-Z_][a-zA-Z0-9_]*)(?:\s*=>\s*(\${1,2}[a-zA-Z_][a-zA-Z0-9_]*))?\s*\)\s*:?`)
		matches := foreachPattern.FindStringSubmatch(code)
		if len(matches) >= 3 {
			collection := matches[1]

			// Check if this is a key-value foreach
			if len(matches) > 3 && matches[3] != "" {
				// Key-value foreach ($array as $key => $value)
				keyVar := matches[2]
				valueVar := matches[3]
				// In Go templates, the convention is to use $ prefix for variables in range
				// fmt.Printf("-- {{ range %s, %s := %s }}", convertPHPVarsToGo(keyVar), convertPHPVarsToGo(valueVar), convertPHPVarsToGo(collection))
				return fmt.Sprintf("{{ range %s, %s := %s }}", convertPHPVarsToGo(keyVar), convertPHPVarsToGo(valueVar), convertPHPVarsToGo(collection))
			} else {
				// Simple foreach ($array as $item)
				// Use the iterVar as a named variable in the range loop
				return fmt.Sprintf("{{ range %s }}", convertPHPVarsToGo(collection))
			}
		}

	}

	// Handle endforeach
	if code == "endforeach" || code == "endforeach;" || code == "endforeach :" {
		return "{{ end }}"
	}

	// Handle if statements
	if strings.HasPrefix(code, "if") {
		ifPattern := regexp.MustCompile(`if\s*\(((?:[^()]+|\((?:[^()]+|\([^()]*\))*\))*)\)\s*:??`)
		matches := ifPattern.FindStringSubmatch(code)
		if len(matches) > 1 {
			condition := convertPHPExprToGo(matches[1])
			// fmt.Printf("304 - %s - %s\n", matches[1], condition)
			return fmt.Sprintf("{{ if %s }}", condition)
		}
	}

	// Handle elseif statements
	if strings.HasPrefix(code, "elseif") || strings.HasPrefix(code, "else if") {
		elseifPattern := regexp.MustCompile(`else\s*if\s*\((.*?)\)(?:\s*:)?`)
		matches := elseifPattern.FindStringSubmatch(code)
		if len(matches) == 0 {
			elseifPattern = regexp.MustCompile(`elseif\s*\((.*?)\)(?:\s*:)?`)
			matches = elseifPattern.FindStringSubmatch(code)
		}

		if len(matches) > 1 {
			condition := convertPHPExprToGo(matches[1])
			return fmt.Sprintf("{{ else if %s }}", condition)
		}
	}

	// Handle else statements
	if code == "else" || code == "else:" {
		return "{{ else }}"
	}

	// Handle endif
	if code == "endif" || code == "endif;" || code == "endif :" {
		return "{{ end }}"
	}

	// Handle for loops
	if strings.HasPrefix(code, "for") {
		// Simple conversion - Go templates don't have direct for loop equivalents
		// This is a limitation - we assume a range loop might be appropriate
		forPattern := regexp.MustCompile(`for\s*\((.*?);(.*?);(.*?)\)(?:\s*:)?`)
		matches := forPattern.FindStringSubmatch(code)
		if len(matches) > 3 {
			// This is a naive implementation - proper conversion would need more context
			return "{{ /* FOR loop converted to range - REVIEW THIS */ }}\n{{ range i := seq X Y }}"
		}
	}

	// Handle endfor
	if code == "endfor" || code == "endfor;" || code == "endfor :" {
		return "{{ end }}"
	}

	// Handle while loops
	if strings.HasPrefix(code, "while") {
		// Go templates don't have direct while loop equivalents
		whilePattern := regexp.MustCompile(`while\s*\((.*?)\)(?:\s*:)?`)
		matches := whilePattern.FindStringSubmatch(code)
		if len(matches) > 1 {
			condition := convertPHPExprToGo(matches[1])
			return fmt.Sprintf("{{ /* WHILE loop conversion - REVIEW THIS */ }}\n{{ if %s }}", condition)
		}
	}

	// Handle endwhile
	if code == "endwhile" || code == "endwhile;" || code == "endwhile :" {
		return "{{ end }}"
	}

	// Handle comments
	commentPattern := regexp.MustCompile(`//\s*(.*?)$`)
	if commentPattern.MatchString(code) {
		comment := commentPattern.FindStringSubmatch(code)[1]
		return fmt.Sprintf("{{/* %s */}}", comment)
	}

	// Handle PHP echo statements
	if strings.HasPrefix(code, "echo") {
		echoPattern := regexp.MustCompile(`echo\s+(.*?)(;|\s*$)`)
		matches := echoPattern.FindStringSubmatch(code)
		if len(matches) > 1 {
			expr := convertPHPExprToGo(matches[1])
			return fmt.Sprintf("{{ %s }}", expr)
		}
	}

	// If we couldn't match to a known pattern, return a comment
	return fmt.Sprintf("{{/* Unconverted PHP code: %s */}}", code)
}
