package library

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ============================================= //
// === Ce fichier gère les saisies du joueur === //
// ============================================= //

// Unique lecteur du clavier, partagé par tout le jeu //
// En créer plusieurs ferait perdre des saisies gardées en mémoire //
var Reader = bufio.NewReader(os.Stdin)

// Affiche une question et renvoie la ligne tapée, sans espaces au début ni à la fin //
func ReadLine(prompt string) string {
	fmt.Print(prompt)
	line, err := Reader.ReadString('\n')
	if err != nil && line == "" {
		fmt.Println("\nFin de la saisie, à bientôt !")
		os.Exit(0)
	}
	return strings.TrimSpace(line)
}

// Lit une ligne et la convertit en nombre //
// Renvoie false si la saisie n'est pas un nombre //
func ReadChoice(prompt string) (int, bool) {
	choice, err := strconv.Atoi(ReadLine(prompt))
	return choice, err == nil
}

// Redemande tant que la saisie n'est pas un nombre entre 0 et maximum //
// Renvoie toujours un choix valide //
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
