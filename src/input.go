package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// reader est l'unique lecteur de l'entrée standard. Toutes les saisies
// passent par lui : le mélanger avec les fonctions Scan du package fmt
// ferait perdre des saisies, car bufio garde en mémoire une partie de ce
// qui a déjà été tapé.
var reader = bufio.NewReader(os.Stdin)

// readLine affiche l'invite puis lit la ligne entière, jusqu'à Entrée.
// Les espaces en début et fin de ligne sont retirés, ainsi que le \r
// ajouté par Windows.
func readLine(prompt string) string {
	fmt.Print(prompt)
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		// Entrée fermée (Ctrl+D, Ctrl+Z ou fin d'un fichier redirigé) :
		// plus rien ne pourra être lu, on quitte au lieu de boucler.
		fmt.Println("\nFin de la saisie, à bientôt !")
		os.Exit(0)
	}
	return strings.TrimSpace(line)
}

// readChoice lit une ligne et la convertit en nombre entier.
// Le booléen vaut false si la ligne n'est pas exactement un nombre.
func readChoice(prompt string) (int, bool) {
	choice, err := strconv.Atoi(readLine(prompt))
	return choice, err == nil
}
