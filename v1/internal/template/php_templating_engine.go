package template

import (
	"fmt"
	"regexp"
	"strings"
)

var phpBlockPattern = regexp.MustCompile(`<\?php([\s\S]*?)\?>`)

// / removing comments
var longComments = regexp.MustCompile(`/\*[\s\S]*?\*/`)
var oneLineComment = regexp.MustCompile(`//.*(\r?\n|$)`)

var foreachPattern = regexp.MustCompile(`foreach\s*\(\s*(\${1,2}[a-zA-Z_][a-zA-Z0-9_]*(?:->\w+)*)\s+as\s+(\${1,2}[a-zA-Z_][a-zA-Z0-9_]*)(?:\s*=>\s*(\${1,2}[a-zA-Z_][a-zA-Z0-9_]*))?\s*\)\s*:?`)
var endforeachPattern = regexp.MustCompile(`(?m)\bendforeach\b\s*[:;]?`)

var ifPattern = regexp.MustCompile(`if\s*\(((?:[^()]+|\((?:[^()]+|\([^()]*\))*\))*)\)\s*:??`)
var elseifPattern = regexp.MustCompile(`(?m)(?:else\s*if|elseif)\s*\(([\s\S]*?)\)\s*:?`)
var elsePattern = regexp.MustCompile(`(?m)\belse\b\s*[:;]?`)
var endifPattern = regexp.MustCompile(`(?m)\bendif\b\s*[:;]?`)

var forPattern = regexp.MustCompile(`(?m)for\s*\(\s*([^;]+)\s*;\s*([^;]+)\s*;\s*([^)]+)\)\s*:?`)

// PHPToGoTemplate converts PHP template code to Go html/template syntax
func PHPToGoTemplate(phpTemplate string) string {

	// // Step 1: Convert short echo tags <?= ... ?> to Go template syntax
	shortEchoPattern := regexp.MustCompile(`<\?=\s*(.*?)\s*\?>`)
	phpTemplate = shortEchoPattern.ReplaceAllStringFunc(phpTemplate, func(match string) string {
		// removing commnet
		match = oneLineComment.ReplaceAllString(match, "")
		// removing long commands /**/
		match = longComments.ReplaceAllString(match, "")

		expr := shortEchoPattern.FindStringSubmatch(match)[1]
		return "{{ " + convertPHPExprToGo(expr) + " }}"
	})

	// Handling php Blocks <?php ?>
	phpTemplate = phpBlockPattern.ReplaceAllStringFunc(phpTemplate, func(match string) string {

		match = oneLineComment.ReplaceAllString(match, "") // replacing all the // comments
		match = longComments.ReplaceAllString(match, "")   // remove all /**/

		m := phpBlockPattern.FindStringSubmatch(match)
		if len(m) < 2 {
			return match
		}
		code := strings.TrimSpace(m[1]) // handles one-liners cleanly
		return "{{" + ConvertPHPBlockToGo(code) + "}}"
	})

	// fmt.Println(phpTemplate)
	return phpTemplate
}

func convertPHPExprToGo(expr string) string {
	// Step 1: Replace -> with . for property access
	// expr = strings.ReplaceAll(expr, "->", ".")

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
	// expr = strings.ReplaceAll(expr, "&&", "and")
	// expr = strings.ReplaceAll(expr, "||", "or")
	// expr = strings.ReplaceAll(expr, "!", "not ")

	return expr
}

// converting php codes : you have to pass code inside <?php ?> it will convert them one by one
func ConvertPHPBlockToGo(code string) string {

	code = foreachPattern.ReplaceAllStringFunc(code, func(match string) string {
		m := foreachPattern.FindStringSubmatch(match)
		if len(m) < 4 {
			return match
		}

		collection := convertPHPVarsToGo(m[1])
		k := strings.TrimSpace(m[2]) // first var after `as`
		v := strings.TrimSpace(m[3]) // second var if key=>value, else ""

		switch {
		case v != "":
			// foreach ($arr as $k => $v)
			return fmt.Sprintf("range %s, %s := %s",
				convertPHPVarsToGo(k), convertPHPVarsToGo(v), collection)

		default:
			// foreach ($arr as $v)  -> {{ range $v := .collection }}
			return fmt.Sprintf("range %s := %s", convertPHPVarsToGo(k), collection)
		}
	})

	// Replace all matches with Go template end
	code = endforeachPattern.ReplaceAllString(code, "end")

	ifPattern.ReplaceAllStringFunc(code, func(match string) string {
		matches := ifPattern.FindStringSubmatch(match)

		condition := convertPHPExprToGo(matches[1])
		return fmt.Sprintf("if %s", condition)
	})

	// Handle elseif / else if blocks in multi-line code

	code = elseifPattern.ReplaceAllStringFunc(code, func(match string) string {
		matches := elseifPattern.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match // fallback if regex fails
		}
		condition := convertPHPExprToGo(matches[1])
		return fmt.Sprintf("else if %s", condition)
	})

	// Handle else statements

	code = elsePattern.ReplaceAllString(code, "else")

	// Handle endif statements

	code = endifPattern.ReplaceAllString(code, "end")

	// Handle for loops (auto-convert common patterns)

	code = forPattern.ReplaceAllStringFunc(code, func(match string) string {
		m := forPattern.FindStringSubmatch(match)
		if len(m) < 4 {
			return match
		}
		init := strings.TrimSpace(m[1])
		cond := strings.TrimSpace(m[2])
		post := strings.TrimSpace(m[3])

		// case 1: for($i = 0; $i < count($arr); $i++)
		if im := regexp.MustCompile(`^\$([a-zA-Z_]\w*)\s*=\s*0$`).FindStringSubmatch(init); im != nil {
			idx := im[1]
			if cm := regexp.MustCompile(`^\$` + idx + `\s*<\s*count\s*\(\s*\$([a-zA-Z_]\w*)\s*\)\s*$`).FindStringSubmatch(cond); cm != nil {
				arr := cm[1]
				// $item name: use idx + "Item" to avoid collisions
				itemVar := idx + "Item"
				return fmt.Sprintf("{{ range $%s, $%s := .%s }}", idx, itemVar, arr)
			}
		}

		// case 2: for($i = 0; $i < $n; $i++) -> best-effort placeholder using seq helper
		if im := regexp.MustCompile(`^\$([a-zA-Z_]\w*)\s*=\s*0$`).FindStringSubmatch(init); im != nil {
			idx := im[1]
			if cm := regexp.MustCompile(`^\$` + idx + `\s*<\s*\$?([a-zA-Z_]\w*)\s*$`).FindStringSubmatch(cond); cm != nil {
				bound := cm[1]
				// Note: this assumes you provide a "seq" or equivalent template func that yields 0..bound-1
				// If you don't have seq, this is a documentation placeholder.
				return fmt.Sprintf("{{/* FOR converted: for(%s; %s; %s) - verify and provide seq() if needed */}}\n{{ range $%s := seq 0 .%s }}", init, cond, post, idx, bound)
			}
		}

		// fallback: cannot auto-convert safely
		return fmt.Sprintf("{{/* FOR loop not auto-converted: (%s; %s; %s) - review */}}", init, cond, post)
	})

	// Handle endfor
	if code == "endfor" || code == "endfor;" || code == "endfor :" {
		return "end"
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
		return "end"
	}

	// // Handle comments
	// commentPattern := regexp.MustCompile(`//\s*(.*?)$`)
	// if commentPattern.MatchString(code) {
	// 	comment := commentPattern.FindStringSubmatch(code)[1]
	// 	return fmt.Sprintf("{{/* %s */}}", comment)
	// }

	// Handle PHP echo statements
	if strings.HasPrefix(code, "echo") {
		echoPattern := regexp.MustCompile(`echo\s+(.*?)(;|\s*$)`)
		matches := echoPattern.FindStringSubmatch(code)
		if len(matches) > 1 {
			expr := convertPHPExprToGo(matches[1])
			return fmt.Sprintf(" %s ", expr)
		}
	}

	return code

	// If we couldn't match to a known pattern, return a comment
	// return fmt.Sprintf("{{/* Unconverted PHP code: %s */}}", code)
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
