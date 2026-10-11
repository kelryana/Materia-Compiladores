/*
Exemplo 6: Tratamento de Erros (Error é Valor)
  - Arquivo: `06_tratamento_erros.go`
  - Integrante: Ana Kelry
  - Tópicos: Interface `error`, formatação com `%w` e inspeção com `errors.Is` ou `errors.As`.
  - Ideia de Implementação: Um "Validador de Cadastro". Criar uma struct de erro customizada
    (`ErroValidacao`) que informa qual campo falhou e o motivo, testando o erro no fluxo principal.
*/
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErroSentinela: erro "sentinela" usado com errors.Is para comparação por identidade.
var ErrCadastroExpirado = errors.New("token de cadastro expirado")

// ErroValidacao é um erro customizado (struct que implementa a interface `error`).
// Ele carrega contexto estruturado: qual campo falhou e o motivo da falha.
type ErroValidacao struct {
	Campo  string
	Motivo string
}

// Implementação da interface error: qualquer tipo com o método
// Error() string pode ser usado como erro em Go ("error é valor").
func (e *ErroValidacao) Error() string {
	return fmt.Sprintf("campo %q inválido: %s", e.Campo, e.Motivo)
}

// Cadastro representa os dados enviados em um formulário de registro.
type Cadastro struct {
	Nome     string
	Email    string
	Idade    int
	Token    string
	ExpiraEm time.Time
}

// validarNome retorna um erro customizado se o nome for inválido.
func validarNome(c Cadastro) error {
	if strings.TrimSpace(c.Nome) == "" {
		return &ErroValidacao{Campo: "nome", Motivo: "não pode ser vazio"}
	}
	if len(c.Nome) < 3 {
		return &ErroValidacao{Campo: "nome", Motivo: "mínimo de 3 caracteres"}
	}
	return nil
}

// validarEmail demonstra wrap de erro com %w: o erro original fica
// encapsulado e recuperável posteriormente via errors.Is/errors.As.
func validarEmail(c Cadastro) error {
	if !strings.Contains(c.Email, "@") {
		base := &ErroValidacao{Campo: "email", Motivo: "formato inválido"}
		return fmt.Errorf("falha ao validar e-mail %q: %w", c.Email, base)
	}
	return nil
}

// validarIdade retorna outro erro customizado para o campo idade.
func validarIdade(c Cadastro) error {
	if c.Idade < 18 || c.Idade > 120 {
		return &ErroValidacao{Campo: "idade", Motivo: fmt.Sprintf("%d fora do intervalo [18, 120]", c.Idade)}
	}
	return nil
}

// validarToken demonstra wrap de um erro sentinela com %w.
func validarToken(c Cadastro, agora time.Time) error {
	if c.Token == "" {
		return fmt.Errorf("autenticação do cadastro falhou: %w", ErrCadastroExpirado)
	}
	if agora.After(c.ExpiraEm) {
		return fmt.Errorf("validação de token (%s): %w", c.Token, ErrCadastroExpirado)
	}
	return nil
}

// multiErro agrega vários erros em um só. Como cada erro individual já foi
// embrulhado com %w pelos validadores, errors.Is/errors.As continuam
// funcionando ao percorrer a cadeia de cada item da lista.
type multiErro []error

func (m multiErro) Error() string {
	var partes []string
	for _, e := range m {
		partes = append(partes, e.Error())
	}
	return strings.Join(partes, "\n  - ")
}

// cadastrar executa o pipeline de validação e retorna os erros encontrados.
// Cada erro individual mantém sua cadeia de wrap (%w), preservando a
// semântica de "error é valor" que será inspecionada em analisarErro.
func cadastrar(c Cadastro, agora time.Time) multiErro {
	validadores := []func(Cadastro) error{
		validarNome,
		validarEmail,
		validarIdade,
		func(cc Cadastro) error { return validarToken(cc, agora) },
	}

	var erros multiErro
	for _, v := range validadores {
		if err := v(c); err != nil {
			erros = append(erros, err)
		}
	}
	return erros
}

// analisarErro demonstra a INSPEÇÃO de erros encapsulados:
//   - errors.Is: percorre a cadeia de wraps procurando equivalência com um erro sentinela;
//   - errors.As: extrai o erro concreto (*ErroValidacao) da cadeia para ler seus campos.
func analisarErro(erros multiErro) {
	fmt.Println(">> Análises dos erros individuais retornados pelo fluxo principal:")

	for _, e := range erros {
		fmt.Printf("   mensagem: %v\n", e)

		// errors.Is: percorre a cadeia de wraps (%w) procurando equivalência
		// com o erro sentinela ErrCadastroExpirado.
		if errors.Is(e, ErrCadastroExpirado) {
			fmt.Println("     [errors.Is] contém ErrCadastroExpirado -> solicite um novo token ao usuário.")
		}

		// errors.As: recupera o ponteiro da struct customizada e acessa Campo/Motivo.
		var ev *ErroValidacao
		if errors.As(e, &ev) {
			fmt.Printf("     [errors.As] erro de validação -> campo=%s motivo=%s\n", ev.Campo, ev.Motivo)
		}
		fmt.Println()
	}
}

func main() {
	agora := time.Date(2026, 10, 11, 12, 0, 0, 0, time.UTC)

	casos := []struct {
		desc     string
		cadastro Cadastro
	}{
		{
			desc: "Cadastro válido",
			cadastro: Cadastro{
				Nome:     "Ana Kelry",
				Email:    "ana.kelry@uern.br",
				Idade:    21,
				Token:    "abc123",
				ExpiraEm: agora.Add(time.Hour),
			},
		},
		{
			desc: "Vários campos inválidos",
			cadastro: Cadastro{
				Nome:     "Al",
				Email:    "email-sem-arroba",
				Idade:    15,
				Token:    "",
				ExpiraEm: agora.Add(time.Hour),
			},
		},
		{
			desc: "Token expirado",
			cadastro: Cadastro{
				Nome:     "Bob Silva",
				Email:    "bob@uern.br",
				Idade:    30,
				Token:    "tk-antigo",
				ExpiraEm: agora.Add(-time.Minute), // expirou antes de "agora"
			},
		},
	}

	for _, caso := range casos {
		fmt.Printf("=== Caso: %s ===\n", caso.desc)

		// "Error é valor": o resultado da validação é tratado como dados,
		// podendo ser armazenado, iterado e inspecionado normalmente.
		erros := cadastrar(caso.cadastro, agora)
		if len(erros) > 0 {
			fmt.Printf("Cadastro de %q rejeitado (%d erro(s)):\n", caso.cadastro.Nome, len(erros))
			analisarErro(erros)
		} else {
			fmt.Printf("Usuário %q cadastrado com sucesso!\n\n", caso.cadastro.Nome)
		}
	}
}
