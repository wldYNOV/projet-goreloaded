# Go Reloaded

Un formatteur de texte en ligne de commande qui implémente la spécification du projet jointe.
Nécessite Go 1.20 ou version ultérieure. Utilise uniquement la bibliothèque standard de Go ; aucune dépendance.

## Exécution

Extrayez l'archive, ouvrez un terminal dans `go-reloaded`, puis exécutez :

```sh
go run . sample.txt result.txt
cat result.txt
```

Texte attendu :

```text
It was an amazing day! Simply add 66 and 2. 'Hello World'...
```

N'importe quel chemin de fichier d'entrée et de sortie fonctionne :

```sh
go run . input.txt output.txt
```

## Règles

- `(hex)` / `(bin)` convertissent le mot précédent en décimal, y compris les entiers plus grands que 64 bits.
- `(up)` / `(low)` / `(cap)` modifient la casse du mot précédent.
- `(up, N)` / `(low, N)` / `(cap, N)` modifient les N mots précédents, en ignorant la ponctuation et les guillemets. Si N dépasse le nombre de mots disponibles, tous les mots disponibles sont modifiés.
- Les commandes s'exécutent de gauche à droite et sont supprimées du résultat.
- La capitalisation met en majuscule la première lettre et en minuscule les lettres restantes, avec prise en charge Unicode.
- `.,!?:;` se collent au texte précédent ; la ponctuation adjacente forme un groupe. Un mot suivant est séparé par un espace.
- Les guillemets simples appariés se collent au texte qu'ils encadrent. Les apostrophes dans les contractions comme `don't` restent intactes.
- `a` / `A` devient `an` / `An` avant un mot commençant par une voyelle ou un h, selon la règle littérale de l'énoncé.
- Tous les espaces, y compris les retours à la ligne et les tabulations, sont normalisés en espaces simples ; aucune ligne vide supplémentaire n'est ajoutée.

Les nombres mal formés, les mots précédents manquants, les comptes nuls ou trop grands, ainsi que les comptes sur les commandes de nombres, renvoient des erreurs. Les commandes non reconnues restent du texte ordinaire.
La sortie n'est écrite qu'après un formatage réussi. Les fichiers d'entrée et de sortie peuvent être identiques.

## Vérification

```sh
gofmt -w *.go
go test ./...
go vet ./...
go run . sample.txt result.txt
```

Les tests couvrent les exemples de spécification, les commandes chainées, la ponctuation, les guillemets, les articles, l'Unicode, les grands entiers, les entrées invalides et la gestion des fichiers.

## Fichiers

- `main.go` : arguments de ligne de commande et E/S de fichiers.
- `formatter.go` : tokenisation, transformations et rendu.
- `formatter_test.go` : tests de formatage pilotés par tableau.
- `main_test.go` : tests de gestion des fichiers et des erreurs.
- `sample.txt` / `expected.txt` : exemple d'entrée et de sortie attendue.
