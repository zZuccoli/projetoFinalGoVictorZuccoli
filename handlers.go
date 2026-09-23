package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AlunoHandler struct{}
type SalaHandler struct{}
type TurmaHandler struct{}

func (AlunoHandler) CriarAluno(c *gin.Context) {
	var aluno Aluno

	if err := c.ShouldBindJSON(&aluno); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "JSON inválido",
		})
		return
	}

	aluno.Matricula = strings.TrimSpace(aluno.Matricula)
	aluno.Nome = strings.TrimSpace(aluno.Nome)
	aluno.Email = strings.TrimSpace(aluno.Email)

	if aluno.Matricula == "" ||
		aluno.Nome == "" ||
		aluno.Email == "" {

		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "matricula, nome e email são obrigatórios",
		})
		return
	}

	if indiceAluno(aluno.Matricula) != -1 {
		c.JSON(http.StatusConflict, gin.H{
			"erro": "aluno já cadastrado",
		})
		return
	}

	aluno.Ativo = true

	alunos = append(alunos, aluno)

	c.JSON(http.StatusCreated, aluno)
}

func (AlunoHandler) ListarAlunos(c *gin.Context) {
	c.JSON(http.StatusOK, alunos)
}

func (AlunoHandler) BuscarAluno(c *gin.Context) {
	indice := indiceAluno(c.Param("id"))

	if indice == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "aluno não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, alunos[indice])
}

func (SalaHandler) CriarSala(c *gin.Context) {
	var sala Sala

	if err := c.ShouldBindJSON(&sala); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "JSON inválido",
		})
		return
	}

	sala.ID = strings.TrimSpace(sala.ID)
	sala.Nome = strings.TrimSpace(sala.Nome)

	if sala.ID == "" || sala.Nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "id e nome são obrigatórios",
		})
		return
	}

	if sala.Capacidade <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "capacidade deve ser maior que zero",
		})
		return
	}

	if indiceSala(sala.ID) != -1 {
		c.JSON(http.StatusConflict, gin.H{
			"erro": "sala já cadastrada",
		})
		return
	}

	if sala.Recursos == nil {
		sala.Recursos = []string{}
	}

	sala.Ativa = true

	salas = append(salas, sala)

	c.JSON(http.StatusCreated, sala)
}

func (SalaHandler) ListarSalas(c *gin.Context) {
	c.JSON(http.StatusOK, salas)
}

func (SalaHandler) ConsultarGrade(c *gin.Context) {
	indice := indiceSala(c.Param("id"))

	if indice == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "sala não encontrada",
		})
		return
	}

	grade := []UsoSala{}

	for _, turma := range turmas {
		if turma.Ativa &&
			turma.Alocacao != nil &&
			turma.Alocacao.SalaID == salas[indice].ID {

			grade = append(grade, UsoSala{
				TurmaID:        turma.ID,
				TurmaNome:      turma.Nome,
				DiaSemana:      turma.Alocacao.DiaSemana,
				HorarioInicio:  turma.Alocacao.HorarioInicio,
				HorarioTermino: turma.Alocacao.HorarioTermino,
			})
		}
	}

	dia := c.Query("dia_semana")
	inicio := c.Query("inicio")
	fim := c.Query("fim")

	if dia == "" && inicio == "" && fim == "" {
		c.JSON(http.StatusOK, gin.H{
			"sala":  salas[indice],
			"grade": grade,
		})
		return
	}

	if dia == "" || inicio == "" || fim == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "informe dia_semana, inicio e fim juntos",
		})
		return
	}

	diaNormalizado, ok := normalizarDiaSemana(dia)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "dia da semana inválido",
		})
		return
	}

	inicioMin, inicioValido := horarioParaMinutos(inicio)
	fimMin, fimValido := horarioParaMinutos(fim)

	if !inicioValido || !fimValido || inicioMin >= fimMin {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "horário inválido; use HH:MM e início menor que término",
		})
		return
	}

	conflitos := []UsoSala{}

	for _, uso := range grade {
		if uso.DiaSemana == diaNormalizado &&
			horariosSobrepostos(
				inicio,
				fim,
				uso.HorarioInicio,
				uso.HorarioTermino,
			) {

			conflitos = append(conflitos, uso)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"sala_id":    salas[indice].ID,
		"dia_semana": diaNormalizado,
		"inicio":     inicio,
		"fim":        fim,
		"disponivel": len(conflitos) == 0,
		"conflitos":  conflitos,
	})
}

func (TurmaHandler) CriarTurma(c *gin.Context) {
	var turma Turma

	if err := c.ShouldBindJSON(&turma); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "JSON inválido",
		})
		return
	}

	turma.ID = strings.TrimSpace(turma.ID)
	turma.Nome = strings.TrimSpace(turma.Nome)
	turma.Disciplina = strings.TrimSpace(turma.Disciplina)
	turma.Professor = strings.TrimSpace(turma.Professor)

	if turma.ID == "" ||
		turma.Nome == "" ||
		turma.Disciplina == "" ||
		turma.Professor == "" {

		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "id, nome, disciplina e professor são obrigatórios",
		})
		return
	}

	if indiceTurma(turma.ID) != -1 {
		c.JSON(http.StatusConflict, gin.H{
			"erro": "turma já cadastrada",
		})
		return
	}

	turma.Alunos = []string{}
	turma.QuantidadeAlunos = 0
	turma.StatusAlocacao = "nao_alocada"
	turma.Alocacao = nil
	turma.Ativa = true

	turmas = append(turmas, turma)

	c.JSON(http.StatusCreated, turma)
}

func (TurmaHandler) ListarTurmas(c *gin.Context) {
	c.JSON(http.StatusOK, turmas)
}

func (TurmaHandler) MatricularAluno(c *gin.Context) {
	indiceDaTurma := indiceTurma(c.Param("id"))

	if indiceDaTurma == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "turma não encontrada",
		})
		return
	}

	var requisicao MatriculaRequest

	if err := c.ShouldBindJSON(&requisicao); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "JSON inválido",
		})
		return
	}

	requisicao.AlunoID = strings.TrimSpace(requisicao.AlunoID)

	indiceDoAluno := indiceAluno(requisicao.AlunoID)

	if indiceDoAluno == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "aluno não encontrado",
		})
		return
	}

	turma := &turmas[indiceDaTurma]
	aluno := alunos[indiceDoAluno]

	if !turma.Ativa || !aluno.Ativo {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"erro": "turma ou aluno inativo",
		})
		return
	}

	if alunoEstaNaTurma(*turma, aluno.Matricula) {
		c.JSON(http.StatusConflict, gin.H{
			"erro": "aluno já está matriculado nesta turma",
		})
		return
	}

	if turma.Alocacao != nil {
		indiceDaSala := indiceSala(turma.Alocacao.SalaID)

		if indiceDaSala == -1 || !salas[indiceDaSala].Ativa {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"erro": "sala alocada não está disponível",
			})
			return
		}

		if len(turma.Alunos)+1 > salas[indiceDaSala].Capacidade {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"erro": "capacidade insuficiente na sala alocada",
			})
			return
		}

		if turmaConflitante, conflito :=
			conflitoAgendaAluno(
				aluno.Matricula,
				turma.ID,
				*turma.Alocacao,
			); conflito {

			c.JSON(http.StatusConflict, gin.H{
				"erro":              "conflito de agenda do aluno",
				"turma_conflitante": turmaConflitante,
			})
			return
		}
	}

	turma.Alunos = append(turma.Alunos, aluno.Matricula)
	turma.QuantidadeAlunos = len(turma.Alunos)

	c.JSON(http.StatusCreated, turma)
}

func (TurmaHandler) ListarAlunosTurma(c *gin.Context) {
	indiceDaTurma := indiceTurma(c.Param("id"))

	if indiceDaTurma == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "turma não encontrada",
		})
		return
	}

	alunosDaTurma := []Aluno{}

	for _, matricula := range turmas[indiceDaTurma].Alunos {
		indice := indiceAluno(matricula)

		if indice != -1 {
			alunosDaTurma = append(
				alunosDaTurma,
				alunos[indice],
			)
		}
	}

	c.JSON(http.StatusOK, alunosDaTurma)
}

func (TurmaHandler) AlocarSala(c *gin.Context) {
	indiceDaTurma := indiceTurma(c.Param("id"))

	if indiceDaTurma == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "turma não encontrada",
		})
		return
	}

	var requisicao AlocacaoRequest

	if err := c.ShouldBindJSON(&requisicao); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "JSON inválido",
		})
		return
	}

	requisicao.SalaID = strings.TrimSpace(requisicao.SalaID)

	indiceDaSala := indiceSala(requisicao.SalaID)

	if indiceDaSala == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "sala não encontrada",
		})
		return
	}

	diaNormalizado, ok :=
		normalizarDiaSemana(requisicao.DiaSemana)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "dia da semana inválido",
		})
		return
	}

	inicioMin, inicioValido :=
		horarioParaMinutos(requisicao.HorarioInicio)

	fimMin, fimValido :=
		horarioParaMinutos(requisicao.HorarioTermino)

	if !inicioValido ||
		!fimValido ||
		inicioMin >= fimMin {

		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "horário inválido; use HH:MM e início menor que término",
		})
		return
	}

	turma := &turmas[indiceDaTurma]
	sala := salas[indiceDaSala]

	if !turma.Ativa || !sala.Ativa {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"erro": "turma ou sala inativa",
		})
		return
	}

	if len(turma.Alunos) > sala.Capacidade {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"erro":              "capacidade insuficiente",
			"quantidade_alunos": len(turma.Alunos),
			"capacidade_sala":   sala.Capacidade,
		})
		return
	}

	novaAlocacao := Alocacao{
		SalaID:         sala.ID,
		DiaSemana:      diaNormalizado,
		HorarioInicio:  requisicao.HorarioInicio,
		HorarioTermino: requisicao.HorarioTermino,
	}

	if turmaConflitante, conflito :=
		conflitoSala(
			sala.ID,
			turma.ID,
			novaAlocacao,
		); conflito {

		c.JSON(http.StatusConflict, gin.H{
			"erro":              "conflito de agenda da sala",
			"turma_conflitante": turmaConflitante,
		})
		return
	}

	for _, matricula := range turma.Alunos {
		if turmaConflitante, conflito :=
			conflitoAgendaAluno(
				matricula,
				turma.ID,
				novaAlocacao,
			); conflito {

			c.JSON(http.StatusConflict, gin.H{
				"erro":              "conflito de agenda de aluno",
				"aluno_id":          matricula,
				"turma_conflitante": turmaConflitante,
			})
			return
		}
	}

	turma.Alocacao = &novaAlocacao
	turma.StatusAlocacao = "alocada"

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "turma alocada com sucesso",
		"turma":    turma,
	})
}
