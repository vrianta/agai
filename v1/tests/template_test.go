package template

import (
	"fmt"
	"testing"

	"github.com/vrianta/agai/v1/internal/template"
)

func TestPHPToGoTemplate_RemoveComments(t *testing.T) {
	phpTemplate := `
		/* Multi-line comment */
		<html>
			<body>
				<p>Hello World</p>
				<?php
					/* Inside PHP comment */
					echo "PHP code";
				?>
				<?= "Short echo" ?>
			</body>
		</html>
		/* Another comment */
	`

	fmt.Println("Original template:\n", phpTemplate)

	goTemplate := template.PHPToGoTemplate(phpTemplate)

	fmt.Println("Converted template:\n", goTemplate)
}

func TestPHPToGoTemplate_ShortEcho(t *testing.T) {
	phpTemplate := `<?= strtoupper($$name) ?>`
	fmt.Println("Original template:", phpTemplate)

	goTemplate := template.PHPToGoTemplate(phpTemplate)

	fmt.Println("Converted template:", goTemplate)
}

func TestPHPToGoTemplate_ForeachBlock(t *testing.T) {
	phpTemplate := `
	<?php foreach($$users as $user): ?>
	<p><?= $user->name ?></p>
	<?php endforeach; ?>
	`
	fmt.Println("Original template:\n", phpTemplate)

	goTemplate := template.PHPToGoTemplate(phpTemplate)

	fmt.Println("Converted template:\n", goTemplate)
}

func TestPHPToGoTemplate_IfElseBlock(t *testing.T) {
	phpTemplate := `
	<?php if($$loggedIn): ?>
	<p>Welcome, <?= $$username ?></p>
	<?php else: ?>
	<p>Please log in</p>
	<?php endif; ?>
	`
	fmt.Println("Original template:\n", phpTemplate)

	goTemplate := template.PHPToGoTemplate(phpTemplate)

	fmt.Println("Converted template:\n", goTemplate)
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
			expected: "{{ range $i, $iItem := .users }}",
		},
		{
			name:     "numeric limit",
			input:    `for($i = 0; $i < $n; $i++)`,
			expected: "{{/* FOR converted: for($i = 0; $i < $n; $i++) - verify and provide seq() if needed */}}\n{{ range $i := seq 0 .$n }}",
		},
		{
			name:     "unconvertible",
			input:    `for($i = 1; $i < $n; $i += 2)`,
			expected: "{{/* FOR loop not auto-converted: ($i = 1; $i < $n; $i += 2) - review */}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := template.ConvertPHPBlockToGo(tt.input)
			fmt.Println("\n---", tt.name, "---")
			fmt.Println("Input:", tt.input)
			fmt.Println("Output:", got)

			if got != tt.expected {
				t.Errorf("expected:\n%s\ngot:\n%s", tt.expected, got)
			}
		})
	}
}
