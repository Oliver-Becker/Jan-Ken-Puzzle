/*
	Algorítmos Avançados
	Trabalho 1 - Jan-Ken-Puzzle
	Bruno Delmonde - 10262818
	Óliver S Becker - 10284890
*/
package main

import (
	"fmt"
	"sort"
//	"strconv"
)

type coord struct {
	x, y int
}

type end struct {
	x, y, tipo int
}

type solucao struct {
	total, difs int
	resps []end
}

type ByPos []end

var mem map[uint64]int

var dfsVet = [4]int{1, -1, 0, 0}

func rdp(tab [][]int, results *solucao) bool {
	a, b := mem[createKey(tab)]
	results.total += a
	return b
}

func temIlhas(x, y, p int, tab [][]int) int {

	if tab[x][y] * p <= 0 {
		return 0
	}
	tab[x][y] = -tab[x][y]

	cont := 1

	if x < len(tab) - 1 {
		cont += temIlhas(x + 1, y, p, tab)
	}

	if x > 0 {
		cont += temIlhas(x - 1, y, p, tab)
	}

	if y < len(tab[0]) - 1 {
		cont += temIlhas(x, y + 1, p, tab)
	}

	if y > 0 {
		cont += temIlhas(x, y - 1, p, tab)
	}

	return cont
}

func busca(tab [][]int, control []coord, R, C int, results *solucao, key uint64) {

			//fmt.Printf("\tEntro\n")
	tam := len(control)
	aux := 0
	if rdp(tab, results) {
			//fmt.Printf("###DP = %d e vorto\n", mem[createKey(tab)])
		return
	}

	if tam == 1 {
		results.resps = append(results.resps, end{control[0].x + 1, control[0].y + 1, tab[control[0].x][control[0].y]})
		results.total++
		results.difs++
		mem[createKey(tab)] = 1
			//fmt.Printf("***Conto e vorto\n")
		return
	}

	if tam != temIlhas(control[0].x, control[0].y, 1, tab) {
		//fmt.Printf("Ilha e vorto\n")
		temIlhas(control[0].x, control[0].y, -1, tab)
		mem[createKey(tab)] = 0
		return
	}
	temIlhas(control[0].x, control[0].y, -1, tab)

	for i := range control {
			//fmt.Printf("Iterando i=%d\n", i)
		curX := control[i].x
		curY := control[i].y
		tipo := tab[curX][curY]
			//fmt.Printf("t: %d\t(%d;%d)\n", tipo, curX, curY)
		comida := tipo % 3 + 1
			//fmt.Printf("c: %d\n", comida)

		// Come pra baixo
		if curX < R - 1 && tab[curX + 1][curY] != 0 && comida == tab[curX + 1][curY] {
			//fmt.Printf("Baixo\n")

			control, control[i], control[tam - 1], tab[curX][curY], tab[curX + 1][curY] = control[:tam - 1], control[tam - 1], control[i], 0, tab[curX][curY]

			//key[

			busca(tab, control, R, C, results, key)
			//fmt.Println(createKey(tab))
			aux += mem[createKey(tab)]

			control = control[:tam]
			control[i], control[tam - 1], tab[curX][curY], tab[curX + 1][curY] = control[tam - 1], control[i], tab[curX + 1][curY], comida
		}

		// Come pra cima
		if curX > 0 && tab[curX - 1][curY] != 0 && comida == tab[curX - 1][curY] {
			//fmt.Printf("Cima\n")

			control, control[i], control[tam - 1], tab[curX][curY], tab[curX - 1][curY] = control[:tam - 1], control[tam - 1], control[i], 0, tab[curX][curY]

			busca(tab, control, R, C, results, key)
			aux += mem[createKey(tab)]

			control = control[:tam]
			control[i], control[tam - 1], tab[curX][curY], tab[curX - 1][curY] = control[tam - 1], control[i], tab[curX - 1][curY], comida
		}

		// Come pra direita
		if curY+1 < C && tab[curX][curY+1] != 0 && comida == tab[curX][curY+1] {
			//fmt.Printf("Dir\n")

			control, control[i], control[tam - 1], tab[curX][curY], tab[curX][curY + 1] = control[:tam - 1], control[tam - 1], control[i], 0, tab[curX][curY]
			busca(tab, control, R, C, results, key)
			aux += mem[createKey(tab)]

			control = control[:tam]
			control[i], control[tam - 1], tab[curX][curY], tab[curX][curY + 1] = control[tam - 1], control[i], tab[curX][curY + 1], comida
		}

		// Come pra esquerda
		if curY > 0 && tab[curX][curY-1] != 0 && comida == tab[curX][curY-1] {
			//fmt.Printf("Esq\n")

			control, control[i], control[tam - 1], tab[curX][curY], tab[curX][curY - 1] = control[:tam - 1], control[tam - 1], control[i], 0, tab[curX][curY]

			busca(tab, control, R, C, results, key)
			aux += mem[createKey(tab)]

			control = control[:tam]
			control[i], control[tam - 1], tab[curX][curY], tab[curX][curY - 1] = control[tam - 1], control[i], tab[curX][curY - 1], comida
		}

	}
		mem[createKey(tab)] = aux

		//fmt.Printf("Cabo e vorto\n")

	return
}

func (a ByPos) Len() int { //Len do Sort
	return len(a)
}

func (a ByPos) Swap(i, j int) { //Swap do Sort
	a[i], a[j] = a[j], a[i]
}

func (a ByPos) Less(i, j int) bool {	// Função de comparar para ordenar.
	if a[i].x != a[j].x {
		return a[i].x < a[j].x
	}

	if a[i].y != a[j].y {
		return a[i].y < a[j].y
	}

	return a[i].tipo < a[j].tipo
}

func main() {
	var R, C int
	var control []coord
	results := solucao{0, 0, make([]end, 0)}
	fmt.Scanf("%d %d", &R, &C)

	tab := make([][]int, R)

	for i := 0; i < R; i++ {	// Leitura da matriz do tabuleiro
		tab[i] = make([]int, C)
		for j := 0; j < C; j++ {
			fmt.Scanf("%d", &tab[i][j])
			if tab[i][j] != 0 {
				control = append(control, coord{i, j})
			}
		}
	}

	mem = make(map[uint64]int)
	key := createKey(tab)

	busca(tab, control, R, C, &results, key)	// Algorítmo backtracking
	//fmt.Printf("Fimm %d %d\n\n", results.total, len(results.resps))

	sort.Sort(ByPos(results.resps))	// Ordena o vetor de saídas

	fmt.Printf("%d\n%d\n", results.total, results.difs)	// Imprime a saída
	for i := range results.resps {
		fmt.Printf("%d %d %d\n", results.resps[i].x, results.resps[i].y, results.resps[i].tipo)
	}
}

func createKey(tab [][]int) uint64 {
	var key uint64
	key = 0

	for i := 0; i < len(tab); i++ {
		for j := 0; j < len(tab[i]); j++ {
			key += uint64(tab[i][j]) << uint64(2 * (i * len(tab[i]) + j))
		}
	}

	return key
}
