package main

import (
	"strconv"
	"strings"
)

func indiceAluno(matricula string) int {
	for i, aluno := range alunos {
		if aluno.Matricula == matricula {
			return i
		}
	}
	return -1
}

func indiceSala(id string) int {
	for i, sala := range salas {
		if sala.ID == id {
			return i
		}
	}
	return -1
}

func indiceTurma(id string) int {
	for i, turma := range turmas {
		if turma.ID == id {
			return i
		}
	}
	return -1
}

func alunoEstaNaTurma(turma Turma, matricula string) bool {
	for _, alunoID := range turma.Alunos {
		if alunoID == matricula {
			return true
		}
	}

	return false
}

func normalizarDiaSemana(dia string) (string, bool) {
	dia = strings.ToLower(strings.TrimSpace(dia))

	switch dia {
	case "segunda", "segunda-feira":
		return "segunda-feira", true

	case "terca", "terça", "terca-feira", "terça-feira":
		return "terça-feira", true

	case "quarta", "quarta-feira":
		return "quarta-feira", true

	case "quinta", "quinta-feira":
		return "quinta-feira", true

	case "sexta", "sexta-feira":
		return "sexta-feira", true

	case "sabado", "sábado":
		return "sábado", true

	case "domingo":
		return "domingo", true

	default:
		return "", false
	}
}

func horarioParaMinutos(horario string) (int, bool) {
	partes := strings.Split(horario, ":")

	if len(partes) != 2 || len(partes[0]) != 2 || len(partes[1]) != 2 {
		return 0, false
	}

	hora, errHora := strconv.Atoi(partes[0])
	minuto, errMinuto := strconv.Atoi(partes[1])

	if errHora != nil || errMinuto != nil {
		return 0, false
	}

	if hora < 0 || hora > 23 || minuto < 0 || minuto > 59 {
		return 0, false
	}

	return hora*60 + minuto, true
}

func horariosSobrepostos(inicioA, fimA, inicioB, fimB string) bool {
	inicioAMin, okA := horarioParaMinutos(inicioA)
	fimAMin, okB := horarioParaMinutos(fimA)
	inicioBMin, okC := horarioParaMinutos(inicioB)
	fimBMin, okD := horarioParaMinutos(fimB)

	if !okA || !okB || !okC || !okD {
		return false
	}

	return inicioAMin < fimBMin && fimAMin > inicioBMin
}

func conflitoSala(salaID, turmaIgnorada string, nova Alocacao) (string, bool) {
	for _, turma := range turmas {
		if !turma.Ativa ||
			turma.ID == turmaIgnorada ||
			turma.Alocacao == nil {
			continue
		}

		if turma.Alocacao.SalaID == salaID &&
			turma.Alocacao.DiaSemana == nova.DiaSemana &&
			horariosSobrepostos(
				nova.HorarioInicio,
				nova.HorarioTermino,
				turma.Alocacao.HorarioInicio,
				turma.Alocacao.HorarioTermino,
			) {

			return turma.ID, true
		}
	}

	return "", false
}

func conflitoAgendaAluno(
	matricula string,
	turmaIgnorada string,
	nova Alocacao,
) (string, bool) {

	for _, turma := range turmas {
		if !turma.Ativa ||
			turma.ID == turmaIgnorada ||
			turma.Alocacao == nil {
			continue
		}

		if !alunoEstaNaTurma(turma, matricula) {
			continue
		}

		if turma.Alocacao.DiaSemana == nova.DiaSemana &&
			horariosSobrepostos(
				nova.HorarioInicio,
				nova.HorarioTermino,
				turma.Alocacao.HorarioInicio,
				turma.Alocacao.HorarioTermino,
			) {

			return turma.ID, true
		}
	}

	return "", false
}
