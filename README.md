# API Jeux Steam

API REST en Go qui expose une base de jeux Steam stockée dans MongoDB.

> État actuel : CRUD des jeux, pagination, tri et filtres fonctionnels. Authentification JWT en cours.

## Stack

- **Go** (1.22+) : serveur HTTP avec `net/http`
- **MongoDB** : base `pa2026`, collection `games`
- **Docker Compose** : lance MongoDB et l'API

## Structure du projet

```
├── main.go           Point d'entrée, déclaration des routes
├── admin/            Handlers HTTP
├── bdd/              Connexion et requêtes MongoDB
├── models/           Structs (Game, Genre, Category, APIError)
├── auth/             JWT + bcrypt (en cours)
├── docker-compose.yml
└── steam.json        Données sources
```

## Installation

### 1. Lancer les conteneurs

```bash
docker compose up -d
```

### 2. Importer les données

```bash
docker cp steam.json mongodb:/steam.json
docker exec -it mongodb mongoimport --db pa2026 --collection games --file /steam.json --jsonArray
```

Vérifier l'import :

```bash
docker exec -it mongodb mongosh pa2026 --quiet --eval "db.games.countDocuments()"
```

### 3. Lancer l'API (hors Docker)

```bash
go run main.go
```

L'API écoute sur le port **8080**.

### Variables d'environnement

| Variable | Rôle | Défaut |
|---|---|---|
| `MONGO_URI` | Adresse de MongoDB | `mongodb://localhost:27017` |
| `DB_NAME` | Nom de la base | `pa2026` |
| `PORT` | Port de l'API | `8080` |

## Modèle de données

Un jeu est identifié par son **id Steam** (champ `_id` dans MongoDB, `id` dans l'API).

```json
{
  "id": 570,
  "title": "Dota 2",
  "required_age": 0,
  "is_free": true,
  "description": "...",
  "image": "https://...",
  "website": "http://www.dota2.com/",
  "developers": ["Valve"],
  "publishers": ["Valve"],
  "platforms": ["windows", "mac", "linux"],
  "metacritic_score": 90,
  "categories": [{ "id": 1, "description": "Multi-player" }],
  "genres": [{ "id": 1, "description": "Action" }],
  "recommendations": 14301,
  "release_date": "2013-07-09",
  "release_year": 2013
}
```

## Routes

| Méthode | Route | Description | Succès |
|---|---|---|---|
| GET | `/api/v1/games` | Liste des jeux | 200 |
| GET | `/api/v1/games/{id}` | Détail d'un jeu | 200 |
| POST | `/api/v1/games` | Créer un jeu | 201 |
| PUT | `/api/v1/games/{id}` | Remplacer un jeu | 200 |
| DELETE | `/api/v1/games/{id}` | Supprimer un jeu | 204 |

Notes :
- Le **PUT remplace le jeu entier** : les champs absents du body sont vidés.
- Pour le PUT, l'id de l'**URL** fait foi, celui du body est ignoré.

## Liste des jeux : pagination, tri et filtres

`GET /api/v1/games` accepte les paramètres suivants. Tous sont facultatifs et combinables.

| Paramètre | Rôle | Valeurs | Défaut |
|---|---|---|---|
| `page` | Numéro de la page | Entier ≥ 1 | `1` |
| `limit` | Nombre de jeux par page | Entier de 1 à 100 | `5` |
| `ordering` | Tri | `title`, `-title`, `release_year`, `-release_year` | Par id croissant |
| `title` | Titre contenant ce texte, sans tenir compte des majuscules | Texte | Aucun filtre |
| `release_year` | Année de sortie exacte | Entier | Aucun filtre |

- Un `-` devant le champ de `ordering` trie en ordre décroissant.
- En cas d'égalité sur le champ trié, les jeux sont départagés par id croissant, pour garder un ordre stable entre les pages.
- `total` et `total_pages` tiennent compte des filtres.
- Une page au-delà de la dernière renvoie **200** avec `results` vide.
- Une valeur invalide (`page=0`, `limit=500`, `release_year=abc`, `ordering=prix`…) renvoie **400**.

Format de la réponse :

```json
{
  "page": 2,
  "limit": 5,
  "total": 27075,
  "total_pages": 5415,
  "results": [ { "id": 570, "title": "Dota 2", "...": "..." } ]
}
```

## Format des erreurs

Toutes les erreurs sont renvoyées en JSON :

```json
{
  "status": 404,
  "error": "Not Found",
  "message": "jeu introuvable"
}
```

| Code | Cas |
|---|---|
| 400 | Body JSON invalide, id non numérique, champs obligatoires manquants |
| 404 | Jeu introuvable |
| 409 | Un jeu avec cet id existe déjà (POST) |
| 500 | Erreur serveur ou base de données |

## Exemples de requêtes

### Lister les jeux

```bash
# Page 1, valeurs par défaut
curl -i http://localhost:8080/api/v1/games

# Page 2, 10 jeux par page
curl "http://localhost:8080/api/v1/games?page=2&limit=10"

# Les jeux les plus récents en premier
curl "http://localhost:8080/api/v1/games?ordering=-release_year"

# Les jeux de 2020 dont le titre contient "war", par ordre alphabétique
curl "http://localhost:8080/api/v1/games?title=war&release_year=2020&ordering=title"
```

Mettre l'URL entre guillemets dans le terminal quand elle contient un `&`.

### Récupérer un jeu

```bash
curl -i http://localhost:8080/api/v1/games/570
```

### Créer un jeu

```bash
curl -i -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -d '{
    "id": 999999,
    "title": "Jeu de test",
    "required_age": 12,
    "is_free": false,
    "description": "Un jeu cree pour tester l API",
    "developers": ["ESGI"],
    "publishers": ["ESGI"],
    "platforms": ["windows", "linux"],
    "metacritic_score": 75,
    "genres": [{"id": 1, "description": "Action"}],
    "release_date": "2026-10-13",
    "release_year": 2026
  }'
```

### Modifier un jeu

```bash
curl -i -X PUT http://localhost:8080/api/v1/games/999999 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Jeu de test modifie",
    "required_age": 16,
    "is_free": true,
    "platforms": ["windows"],
    "genres": [{"id": 2, "description": "Strategy"}],
    "release_year": 2026
  }'
```

### Supprimer un jeu

```bash
curl -i -X DELETE http://localhost:8080/api/v1/games/999999
```

## Sécurité

- MongoDB n'est **pas exposé publiquement** : le port 27017 est publié uniquement sur `127.0.0.1`.
- L'API joint MongoDB par le réseau interne Docker.

## Reste à faire

- [x] Pagination (`page`, `limit`)
- [x] Tri (`ordering`)
- [x] Filtres (`title`, `release_year`)
- [ ] Erreurs au format JSON sur toutes les routes
- [ ] Authentification JWT sur POST, PUT et DELETE
- [ ] Documentation Swagger / OpenAPI
- [ ] CORS et logs