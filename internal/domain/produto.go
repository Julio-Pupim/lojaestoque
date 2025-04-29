package domain

import (
	"errors"
	"fmt"
)

// Produto representa um produto no sistema
type Produto struct {
	ID                int64      `json:"id"`
	Nome              string     `json:"nome"`
	Fornecedor        Fornecedor `json:"fornecedor"`
	CodigoFornecedor  string     `json:"codigo_fornecedor"`
	QuantidadeEstoque int64      `json:"quantidade_estoque"`
	Preco             Decimal    `json:"preco"`
}

// NewProduto cria uma nova instância de Produto com validação
func NewProduto(
	id int64,
	nome string,
	fornecedor Fornecedor,
	codigoFornecedor string,
	quantidadeEstoque int64,
	precoStr string, // recebemos o preço como string
) (*Produto, error) {
	// Utilizando o wrapper Decimal
	var dec Decimal
	if err := dec.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, precoStr))); err != nil {
		return nil, fmt.Errorf("preço inválido: %w", err)
	}

	p := &Produto{
		ID:                id,
		Nome:              nome,
		Fornecedor:        fornecedor,
		CodigoFornecedor:  codigoFornecedor,
		QuantidadeEstoque: quantidadeEstoque,
		Preco:             dec,
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

// Validate valida os dados do produto
func (p *Produto) Validate() error {
	if p.Nome == "" {
		return errors.New("nome não pode ser vazio")
	}
	if p.Fornecedor.Id <= 0 {
		return errors.New("fornecedor inválido")
	}
	if p.CodigoFornecedor == "" {
		return errors.New("código do fornecedor não pode ser vazio")
	}
	if p.QuantidadeEstoque < 0 {
		return errors.New("estoque não pode ser negativo")
	}
	if p.Preco.Decimal == nil || p.Preco.Sign() < 0 {
		return errors.New("preço não pode ser negativo")
	}
	return nil
}

// SetNome altera o nome do produto com validação
func (p *Produto) SetNome(nome string) error {
	if nome == "" {
		return errors.New("nome do produto não pode ser vazio")
	}
	p.Nome = nome
	return nil
}

// SetQuantidadeEstoque altera a quantidade em estoque com validação
func (p *Produto) SetQuantidadeEstoque(quantidade int64) error {
	if quantidade < 0 {
		return errors.New("quantidade em estoque não pode ser negativa")
	}
	p.QuantidadeEstoque = quantidade
	return nil
}

// SetPreco altera o preço do produto com validação
func (p *Produto) SetPreco(precoStr string) error {
	var dec Decimal
	if err := dec.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, precoStr))); err != nil {
		return fmt.Errorf("preço inválido: %w", err)
	}
	if dec.Sign() < 0 {
		return errors.New("preço não pode ser negativo")
	}
	p.Preco = dec
	return nil
}

// ValidateAndUpdate valida o produto após qualquer atualização
func (p *Produto) ValidateAndUpdate() error {
	return p.Validate()
}
