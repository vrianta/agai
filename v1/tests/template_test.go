package template

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vrianta/agai/v1/internal/template"
)

func TestPHPToGoTemplate_ShortEcho(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "strtoupper",
			input:    `<?= strtoupper($$name) ?>`,
			expected: `{{ Upper .name }}`,
		},
		{
			name:     "strtolower",
			input:    `<?= strtolower($$name) ?>`,
			expected: `{{ lower .name }}`,
		},
		{
			name:     "strlen",
			input:    `<?= strlen($$text) ?>`,
			expected: `{{ (Strlen .text) }}`,
		},
		{
			name:     "count",
			input:    `<?= count($$items) ?>`,
			expected: `{{ (Len .items) }}`,
		},
		{
			name:     "htmlspecialchars",
			input:    `<?= htmlspecialchars($$raw) ?>`,
			expected: `{{ Html .raw }}`,
		},
		{
			name:     "isset",
			input:    `<?= isset($$email) ?>`,
			expected: `{{ ne .email nil }}`,
		},
		{
			name:     "empty",
			input:    `<?= empty($$email) ?>`,
			expected: `{{ eq .email "" }}`,
		},
		{
			name:     "print",
			input:    `<?= print($$value) ?>`,
			expected: `{{ print .value }}`,
		},
		{
			name:     "unknown function fallback",
			input:    `<?= foo($$x, $$y) ?>`,
			expected: `{{ call .foo .x .y }}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := template.PHPToGoTemplate(tt.input)

			fmt.Println("\n---", tt.name, "---")
			fmt.Println("Input:", tt.input)
			fmt.Println("Output:", got)

			if got != tt.expected {
				t.Errorf("\nexpected:\n%s\ngot:\n%s", tt.expected, got)
			}
		})
	}
}

func TestConvertPHPBlockToGo_ForeachTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// --- plain tokens ---
		{
			name:     "endforeach token with semicolon",
			input:    `endforeach;`,
			expected: `{{ end }}`,
		},
		{
			name:     "endforeach token without semicolon",
			input:    `endforeach`,
			expected: `{{ end }}`,
		},

		// --- FOREACH variants ---
		{
			name:     "value only over slice",
			input:    `foreach($$users as $user):`,
			expected: `{{ range $user := .users }}`,
		},
		{
			name:     "key => value over slice/map",
			input:    `foreach($$items as $id => $item):`,
			expected: `{{ range $id, $item := .items }}`,
		},
		{
			name:     "index and value from object property",
			input:    `foreach($$order->lines as $i => $line):`,
			expected: `{{ range $i, $line := .order.lines }}`,
		},
		{
			name:     "scalar values",
			input:    `foreach($$tags as $tag):`,
			expected: `{{ range $tag := .tags }}`,
		},
		{
			name:     "nested property value only",
			input:    `foreach($$user->emails as $email):`,
			expected: `{{ range $email := .user.emails }}`,
		},
		{
			name:     "key only (discard value)",
			input:    `foreach($$map as $k => $_):`,
			expected: `{{ range $k, $_ := .map }}`,
		},
		{
			name:     "value only with underscore key in source",
			input:    `foreach($$list as $_ => $v):`,
			expected: `{{ range $_, $v := .list }}`,
		},

		// --- Additional high-value cases ---
		{
			name:     "array index with string key",
			input:    `foreach($$data['rows'] as $row):`,
			expected: `{{ range $row := .data.rows }}`,
		},
		{
			name:     "array index with variable key",
			input:    `foreach($$data[$i] as $row):`,
			expected: `{{ range $row := index .data $i }}`,
		},
		{
			name:     "array_keys helper",
			input:    `foreach(array_keys($$map) as $k):`,
			expected: `{{ range $k := keys .map }}`,
		},
		{
			name:     "array_values helper",
			input:    `foreach(array_values($$map) as $v):`,
			expected: `{{ range $v := values .map }}`,
		},
		{
			name:     "function call as source",
			input:    `foreach(getUsers($$org) as $u):`,
			expected: `{{ range $u := call .getUsers .org }}`,
		},
		{
			name:     "by-reference value",
			input:    `foreach($$nums as &$n):`,
			expected: `{{ range $n := .nums }}`,
		},
		{
			name:     "extra spaces tolerated",
			input:    `foreach (  $$items   as   $x  ):`,
			expected: `{{ range $x := .items }}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := template.ConvertPHPBlockToGo(tt.input)

			if got != tt.expected {
				t.Errorf("\nexpected:\n%s\ngot:\n%s", tt.expected, got)
			}
		})
	}
}

func TestConvertPHPBlockToGo_IfElseTokens_WithConditions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// --- plain tokens ---
		{
			name:     "else token",
			input:    `else:`,
			expected: `{{ else }}`,
		},
		{
			name:     "endif token",
			input:    `endif;`,
			expected: `{{ end }}`,
		},

		// --- IF cases ---
		{
			name:     "if simple var",
			input:    `if($$loggedIn):`,
			expected: `{{ if .loggedIn }}`,
		},
		{
			name:     "if equality string",
			input:    `if($$role == "admin"):`,
			expected: `{{ if eq .role "admin" }}`,
		},
		{
			name:     "if numeric compare and logical AND",
			input:    `if($$age >= 18 && $$country == "IN"):`,
			expected: `{{ if and (ge .age 18) (eq .country "IN") }}`,
		},
		{
			name:     "if function count",
			input:    `if(count($$users) > 0):`,
			expected: `{{ if gt (Len .users) 0 }}`,
		},
		{
			name:     "if unary not",
			input:    `if(!$$disabled):`,
			expected: `{{ if not .disabled }}`,
		},
		{
			name:     "if nested parens with or",
			input:    `if(($$a && ($$b || $$c))):`,
			expected: `{{ if and .a (or .b .c) }}`,
		},
		{
			name:     "if isset",
			input:    `if(isset($$email)):`,
			expected: `{{ if ne .email nil }}`,
		},
		{
			name:     "if empty",
			input:    `if(empty($$email)):`,
			expected: `{{ if eq .email "" }}`,
		},
		{
			name:     "if object property",
			input:    `if($$user->active):`,
			expected: `{{ if .user.active }}`,
		},

		// --- ELSEIF variants ---
		{
			name:     "elseif simple var (elseif)",
			input:    `elseif($$isAdmin):`,
			expected: `{{ else if .isAdmin }}`,
		},
		{
			name:     "elseif simple var (else if spacing)",
			input:    `else if($$isOwner):`,
			expected: `{{ else if .isOwner }}`,
		},
		{
			name:     "elseif with string equality",
			input:    `elseif($$role == "editor"):`,
			expected: `{{ else if eq .role "editor" }}`,
		},
		{
			name:     "elseif with NOT",
			input:    `elseif(!$$disabled):`,
			expected: `{{ else if not .disabled }}`,
		},
		{
			name:     "elseif with AND of comparisons",
			input:    `elseif($$score > 80 && $$passed == true):`,
			expected: `{{ else if and (gt .score 80) (eq .passed true) }}`,
		},
		{
			name:     "elseif using count and empty",
			input:    `elseif(count($$items) == 0 || empty($$query)):`,
			expected: `{{ else if or (eq (Len .items) 0) (eq .query "") }}`,
		},
		{
			name:     "elseif object property",
			input:    `elseif($$user->verified):`,
			expected: `{{ else if .user.verified }}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := template.ConvertPHPBlockToGo(tt.input)

			fmt.Println("\n---", tt.name, "---")
			fmt.Println("Input:", tt.input)
			fmt.Println("Output:", got)

			if got != tt.expected {
				t.Errorf("\nexpected:\n%s\ngot:\n%s", tt.expected, got)
			}
		})
	}
}

func TestConvertPHPBlockToGo_ForLoop(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "array foreach style",
			input:    `for($i = 0; $i < count($users); $i++)`,
			expected: "range $i, $iItem := .users",
		},
		{
			name:     "numeric limit",
			input:    `for($i = 0; $i < $n; $i++)`,
			expected: "/* FOR converted: for($i = 0; $i < $n; $i++) - verify and provide seq() if needed */ range $i := seq 0 .n",
		},
		{
			name:     "unconvertible",
			input:    `for($i = 1; $i < $n; $i += 2)`,
			expected: "/* FOR loop not auto-converted: ($i = 1; $i < $n; $i += 2) - review */",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := template.ConvertPHPBlockToGo(tt.input)
			fmt.Println("\n---", tt.name, "---")
			fmt.Println("Input:", tt.input)
			fmt.Println("Output:", got)

			if got != tt.expected {
				t.Errorf("\nexpected:\n%s\ngot:\n%s", tt.expected, got)
			}
		})
	}
}

func TestPHPToGoTemplate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// --- comment removal anywhere ---
		{
			name: "remove multi-line comments outside PHP",
			input: `
<?php /* top */ ?>
<div>ok</div>
<?php /* bottom */ ?>`,
			expected: "{{}}\n<div>ok</div>\n{{}}",
		},
		{
			name: "remove one-line comments outside PHP",
			input: `<?php // outside ?>
<div>ok</div> <?php // tail ?>`,
			expected: "{{}}\n<div>ok</div> {{}}",
		},
		{
			name: "remove comment-only PHP block",
			input: `
<?php /* inside php comment */ ?>
<span>ok</span>`,
			expected: "{{}}\n<span>ok</span>",
		},

		// --- raw short echo ---
		{
			name:     "short echo literal string",
			input:    `Before <?= "X" ?> After`,
			expected: `Before {{ "X" }} After`,
		},
		{
			name:     "short echo variable",
			input:    `<p>Hello, <?= $$name ?></p>`,
			expected: `<p>Hello, {{ .name }}</p>`,
		},
		{
			name:     "short echo in attribute",
			input:    `<img src="<?= $$url ?>">`,
			expected: `<img src="{{ .url }}">`,
		},

		// --- long-form php echo ---
		{
			name:     "php echo variable",
			input:    `<b><?php echo $$title; ?></b>`,
			expected: `<b>{{ .title }}</b>`,
		},
		{
			name:     "php echo function on var",
			input:    `<?php echo strtoupper($$name); ?>`,
			expected: `{{ Upper .name }}`,
		},

		// --- supported function mapping in short echo ---
		{
			name:     "strtoupper",
			input:    `<?= strtoupper($$name) ?>`,
			expected: `{{ Upper .name }}`,
		},
		{
			name:     "strtolower",
			input:    `<?= strtolower($$name) ?>`,
			expected: `{{ lower .name }}`,
		},
		{
			name:     "strlen",
			input:    `<?= strlen($$text) ?>`,
			expected: `{{ (Strlen .text) }}`,
		},
		{
			name:     "count",
			input:    `<?= count($$items) ?>`,
			expected: `{{ (Len .items) }}`,
		},
		{
			name:     "htmlspecialchars",
			input:    `<?= htmlspecialchars($$raw) ?>`,
			expected: `{{ Html .raw }}`,
		},
		{
			name:     "isset",
			input:    `<?= isset($$email) ?>`,
			expected: `{{ ne .email nil }}`,
		},
		{
			name:     "empty",
			input:    `<?= empty($$query) ?>`,
			expected: `{{ eq .query "" }}`,
		},
		{
			name:     "print passthrough",
			input:    `<?= print($$value) ?>`,
			expected: `{{ print .value }}`,
		},

		// --- fallback and spacing robustness ---
		{
			name:     "unknown function fallback",
			input:    `<?= foo($$x, $$y) ?>`,
			expected: `{{ call .foo .x .y }}`,
		},
		{
			name:     "extra spaces tolerated",
			input:    `<?php   echo    strtolower( $$Name )   ;   ?>`,
			expected: `{{ lower .Name }}`,
		},
		{
			name: "strip comments around mixed content",
			input: `
<?php /* head */ ?>
<html>
<?php // stray ?>
<body>
<?= strtoupper($$name) /* tail */ ?>
</body>
</html>`,
			expected: "{{}}\n<html>\n{{}}\n<body>\n{{ Upper .name }}\n</body>\n</html>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strings.TrimSpace(template.PHPToGoTemplate(tt.input))
			want := strings.TrimSpace(tt.expected)

			fmt.Println("\n---", tt.name, "---")
			fmt.Println("Input:\n", tt.input)
			fmt.Println("Output:\n", got)

			if got != want {
				t.Errorf("\nexpected:\n%s\n\ngot:\n%s", want, got)
			}
		})
	}
}
