package domain

type Cliente struct {
	ID           int64  `json:"id"`
	Nome         string `json:"nome"`
	Telefone     string `json:"telefone"`
	DataCadastro string `json:"data_cadastro"` // YYYY-MM-DD
}
