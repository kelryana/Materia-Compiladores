/*
Exemplo 5: Condicionais (If e Switch)
  - Arquivo: `05_condicionais_switch.go`
  - Integrante: Ana Kelry
  - Tópicos: `if` com inicialização curta, `switch` (sem break implícito) e `type switch`.
  - Ideia de Implementação: Um "Analisador de Permissões de Usuário" que verifica os papéis
    (Admin, User, Guest) com um `switch`, e usa o `type switch` para descobrir dinamicamente
    se o payload recebido é uma string ou um int.
*/
package main

import "fmt"

// Usuario representa um usuário do sistema com seu papel (role).
type Usuario struct {
	Nome string
	Role string
}

// analisarPermissao demonstra o `switch` de Go:
// - não precisa de `break` (cada case termina implicitamente);
// - pode ter `fallthrough` quando desejado;
// - o `switch` sem expressão funciona como cadeia de if/else-if.
func analisarPermissao(u Usuario) {
	fmt.Printf("\n--- Analisando permissões de %s ---\n", u.Nome)

	switch u.Role {
	case "admin":
		fmt.Println("Acesso total: pode criar, ler, atualizar e excluir.")
	case "user":
		fmt.Println("Acesso padrão: pode criar e ler seus próprios dados.")
		// sem `break` explícito: o case já encerra aqui automaticamente
	case "guest":
		fmt.Println("Acesso somente leitura de conteúdo público.")
	default:
		fmt.Printf("Papel desconhecido (%q): acesso negado.\n", u.Role)
	}

	// Switch sem expressão: cada case é uma condição booleana.
	switch {
	case u.Role == "admin" && u.Nome == "root":
		fmt.Println("Modo superusuário ativado.")
	default:
		fmt.Println("Modo seguro padrão ativo.")
	}
}

// processarPayload demonstra o `type switch`: descoberta dinâmica do tipo
// de um valor `interface{}` em tempo de execução.
func processarPayload(payload interface{}) {
	fmt.Printf("Payload recebido: %v\n", payload)

	switch v := payload.(type) {
	case string:
		fmt.Printf("  -> É uma string de %d caracteres: %q\n", len(v), v)
	case int:
		fmt.Printf("  -> É um inteiro: dobro = %d\n", v*2)
	case float64:
		fmt.Printf("  -> É um float64: arredondado = %.0f\n", v)
	case []string:
		fmt.Printf("  -> É uma lista de strings com %d itens\n", len(v))
	default:
		fmt.Printf("  -> Tipo não suportado: %T\n", v)
	}
}

// validarUsuario retorna um erro usando `if` com inicialização curta.
func validarUsuario(u Usuario) error {
	if nome := u.Nome; nome == "" {
		return fmt.Errorf("usuário sem nome (role=%q)", u.Role)
	}
	return nil
}

func main() {
	// `if` com inicialização curta: a variável `erro` só existe dentro do bloco if/else.
	if erro := validarUsuario(Usuario{"", "guest"}); erro != nil {
		fmt.Println("Falha na validação:", erro)
	} else {
		fmt.Println("Usuário válido.")
	}

	if erro := validarUsuario(Usuario{"Ana Kelry", "user"}); erro != nil {
		fmt.Println("Falha na validação:", erro)
	} else {
		fmt.Println("Usuário válido.")
	}

	// Demonstração do switch comum com diferentes papéis.
	usuarios := []Usuario{
		{"alice", "admin"},
		{"bob", "user"},
		{"carol", "guest"},
		{"mallory", "hacker"},
	}
	for _, u := range usuarios {
		analisarPermissao(u)
	}

	// Demonstração do type switch com payloads heterogêneos.
	fmt.Println("\n=== Processamento de payloads ===")
	payloads := []interface{}{
		"comando: reiniciar servidor",
		42,
		3.14159,
		[]string{"tags", "go", "compiladores"},
		true,
	}
	for _, p := range payloads {
		processarPayload(p)
	}
}
