package library

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/* Ce fichier centralise la lecture et la conversion des saisies entrées par le joueur. */

/* La variable Reader est l'unique lecteur utilisé pour récupérer les saisies du joueur. */
var Reader = bufio.NewReader(os.Stdin)

/* La fonction ReadLine affiche une demande, lit une ligne entière et retire les espaces inutiles. */
func ReadLine(prompt string) string {
	fmt.Print(prompt)
	line, err := Reader.ReadString('\n')
	if err != nil && line == "" {
		fmt.Println("\nFin de la saisie, à bientôt !")
		os.Exit(0)
	}
	return strings.TrimSpace(line)
}

/* La fonction ReadChoice lit une saisie et indique si elle peut être convertie en nombre entier. */
func ReadChoice(prompt string) (int, bool) {
	choice, err := strconv.Atoi(ReadLine(prompt))
	return choice, err == nil
}

func ReadChoiceEntre(prompt string, maximum int) int {
	for {
		choix, ok := ReadChoice(prompt)
		if !ok || choix < 0 || choix > maximum {
			fmt.Printf("Choix invalide. Entrez un nombre entre 0 et %d\n", maximum)
			continue
		}
		return choix
	}
}
