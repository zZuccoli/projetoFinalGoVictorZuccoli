package main

type Aluno struct {
	Matricula string `json:"matricula"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	Ativo     bool   `json:"ativo"`
}

type Sala struct {
	ID         string   `json:"id"`
	Nome       string   `json:"nome"`
	Capacidade int      `json:"capacidade"`
	Recursos   []string `json:"recursos"`
	Ativa      bool     `json:"ativa"`
}

type Alocacao struct {
	SalaID         string `json:"sala_id"`
	DiaSemana      string `json:"dia_semana"`
	HorarioInicio  string `json:"horario_inicio"`
	HorarioTermino string `json:"horario_termino"`
}

type Turma struct {
	ID               string    `json:"id"`
	Nome             string    `json:"nome"`
	Disciplina       string    `json:"disciplina"`
	Professor        string    `json:"professor"`
	QuantidadeAlunos int       `json:"quantidade_alunos"`
	Alunos           []string  `json:"alunos"`
	StatusAlocacao   string    `json:"status_alocacao"`
	Alocacao         *Alocacao `json:"alocacao,omitempty"`
	Ativa            bool      `json:"ativa"`
}

type MatriculaRequest struct {
	AlunoID string `json:"aluno_id"`
}

type AlocacaoRequest struct {
	SalaID         string `json:"sala_id"`
	DiaSemana      string `json:"dia_semana"`
	HorarioInicio  string `json:"horario_inicio"`
	HorarioTermino string `json:"horario_termino"`
}

type UsoSala struct {
	TurmaID        string `json:"turma_id"`
	TurmaNome      string `json:"turma_nome"`
	DiaSemana      string `json:"dia_semana"`
	HorarioInicio  string `json:"horario_inicio"`
	HorarioTermino string `json:"horario_termino"`
}
