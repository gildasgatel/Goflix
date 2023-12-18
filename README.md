# Goflix

Goflix est une API en Go inspirée par Netflix, conçue pour fournir des fonctionnalités similaires dans le contexte du streaming de vidéos et de contenu multimédia.

![This is an image](https://github.com/gildasgatel/Goflix/blob/master/goflix.jpg)

## Fonctionnalités

- **Gestion des utilisateurs :** Création de comptes, authentification, gestion des profils.
- **Catalogue de contenu :** Ajout, suppression et gestion du contenu multimédia (films, séries, documentaires, etc.).
- **Gestion des favoris :** Permet aux utilisateurs de sauvegarder leurs contenus préférés.

## Installation

1. Assurez-vous d'avoir Go installé sur votre machine.
2. Clonez ce dépôt : `git clone https://github.com/gildasgatel/Goflix.git`
3. Accédez au répertoire du projet : `cd goflix`
4. Installez les dépendances : `go mod tidy`
5. Lancez l'API : `go run main.go`

## Utilisation

1. Après avoir lancé l'API, accédez à l'URL suivante : `http://localhost:4123` (ou une autre si spécifiée).
2. Utilisez les endpoints pour effectuer des requêtes (par exemple : `/users`, `/movies`, `/recommendations`, etc.).
3. Consultez la documentation fournie pour une utilisation détaillée de chaque endpoint.

## Endpoints

-**Gestion des utilisateurs :**
   
    - POST /users : Créer un nouvel utilisateur.

    - POST /login : S'identifer et retourne un token.
   
    - GET /users/{userID} : Récupérer les informations d'un utilisateur spécifique.
  
    - PUT /users/{userID} : Mettre à jour les informations d'un utilisateur.
   
    - DELETE /users/{userID} : Supprimer un utilisateur.

-**Catalogue de contenu :**
    
    - GET /movies : Récupérer la liste des films disponibles.
    
    - GET /series : Récupérer la liste des séries disponibles.
    
    - GET /movies/{movieID} : Obtenir les détails d'un film spécifique.
    
    - POST /movies : Ajouter un nouveau film au catalogue.
    
    - DELETE /movies/{movieID} : Supprimer un film du catalogue.

-**Système de recommandations :**
    
    - POST /ratings : Ajouter une évaluation d'utilisateur pour un film ou une série.
    
    - GET /ratings/{userID} : Obtenir les évaluations d'un utilisateur spécifique.

-**Gestion des favoris :**
    
    - POST /favorites : Ajouter un élément aux favoris d'un utilisateur.
    
    - GET /favorites/{userID} : Obtenir la liste des éléments favoris d'un utilisateur.
    
    - DELETE /users/{userID}/favorites/{favoriteID} : Supprimer un élément des favoris d'un utilisateur.


## Licence

Ce projet est sous licence [MIT](https://choosealicense.com/licenses/mit/).
