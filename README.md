# API Jeux Steam

API REST en Go qui expose une base de jeux Steam stockée dans MongoDB.

> État actuel : CRUD des jeux, pagination, tri, filtres, authentification JWT et documentation Swagger fonctionnels.

## Stack

- **Go** (1.25+) : serveur HTTP avec `net/http`
- **MongoDB** : base `pa2026`, collection `games`
- **JWT + bcrypt** : authentification des routes d'écriture
- **Swagger UI / OpenAPI 3** : documentation interactive
- **Docker Compose** : lance MongoDB et l'API

## Structure du projet

```
├── main.go           Point d'entrée
├── route/
│   ├── routes.go     Déclaration des routes
│   └── web/          Fichier de tests HTTP (games.http)
├── admin/            Handlers HTTP des jeux
├── auth/             Login, middleware JWT, création d'utilisateurs
├── bdd/              Connexion et requêtes MongoDB
├── models/           Structs (Game, Genre, Category, User, LoginRequest)
├── docs/             Contrat OpenAPI et page Swagger UI
├── docker-compose.yml
├── users.json        Utilisateurs de l'API (mots de passe hachés)
└── steam.json        Données sources
```

## Installation

### 1. Lancer les conteneurs

```bash
docker compose up -d --build
```

L'API écoute sur le port **8080**, et MongoDB sur `127.0.0.1:27017`.

### 2. Importer les données

```bash
docker cp steam.json mongodb:/steam.json
docker exec -it mongodb mongoimport --db pa2026 --collection games --file /steam.json --jsonArray
```

Vérifier l'import (51 822 jeux attendus) :

```bash
docker exec -it mongodb mongosh pa2026 --quiet --eval "db.games.countDocuments()"
```

### 3. Créer un utilisateur

```bash
go run . u monutilisateur -p monmotdepasse
```

L'utilisateur est ajouté à `users.json`, monté dans le conteneur de l'API : aucun redémarrage n'est nécessaire.

### Lancer l'API hors Docker (facultatif)

Arrêter d'abord le conteneur de l'API, qui occupe le port 8080 :

```bash
docker compose stop webapi
JWT_SECRET=supersecretkey123 go run .
```

### Variables d'environnement

| Variable | Rôle | Défaut |
|---|---|---|
| `MONGO_URI` | Adresse de MongoDB | `mongodb://localhost:27017` |
| `DB_NAME` | Nom de la base | `pa2026` |
| `JWT_SECRET` | Clé de signature des tokens, **obligatoire** pour les routes protégées | Aucun (défini dans `docker-compose.yml`) |

## Documentation Swagger

Une fois l'API lancée, la documentation interactive est disponible sur **http://localhost:8080/swagger**. Le contrat OpenAPI brut est sur `/swagger/openapi.yaml`.

Pour appeler une route protégée :

1. Exécuter `POST /login` et copier le `token` reçu.
2. Cliquer sur **Authorize** et coller le token, sans le préfixe `Bearer`.

La page charge Swagger UI depuis un CDN : une connexion internet est nécessaire pour l'afficher.

## Tester l'API

Le fichier [`route/web/games.http`](route/web/games.http) contient des requêtes pour toutes les routes, avec le statut attendu pour chacune. Il s'utilise avec l'extension VS Code **REST Client**.

Il repose sur l'utilisateur `testuser`, déjà présent dans `users.json`. Exécuter les requêtes dans l'ordre : le login fournit le token des suivantes.

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

| Méthode | Route | Description | Token | Succès |
|---|---|---|---|---|
| POST | `/api/v1/login` | Obtenir un token JWT | Non | 200 |
| GET | `/api/v1/games` | Liste des jeux | Non | 200 |
| GET | `/api/v1/games/{id}` | Détail d'un jeu | Non | 200 |
| POST | `/api/v1/games` | Créer un jeu | **Oui** | 201 |
| PUT | `/api/v1/games/{id}` | Remplacer un jeu | **Oui** | 200 |
| DELETE | `/api/v1/games/{id}` | Supprimer un jeu | **Oui** | 200 |
| GET | `/swagger` | Documentation Swagger UI | Non | 200 |
| GET | `/swagger/openapi.yaml` | Contrat OpenAPI | Non | 200 |

Notes :
- Le **PUT remplace le jeu entier** : les champs absents du body sont vidés.
- Pour le PUT, l'id de l'**URL** fait foi, celui du body est ignoré.

## Authentification

`POST /api/v1/login` vérifie le nom d'utilisateur et le mot de passe, puis renvoie un token JWT valable **24 h** :

```json
{ "token": "eyJhbGciOiJIUzI1NiIs..." }
```

Les routes protégées attendent ce token dans l'en-tête `Authorization: Bearer <token>`. Un token absent, invalide ou expiré renvoie **401**.

## Liste des jeux : pagination, tri et filtres

`GET /api/v1/games` accepte les paramètres suivants. Tous sont facultatifs et combinables.

| Paramètre | Rôle | Valeurs | Défaut |
|---|---|---|---|
| `page` | Numéro de la page | Entier ≥ 1 | `1` |
| `limit` | Nombre de jeux par page | Entier de 1 à 100 | `5` |
| `ordering` | Tri | `title`, `-title`, `release_year`, `-release_year` | Par id croissant |
| `title` | Titre contenant ce texte, sans tenir compte des majuscules | Texte (interprété comme une expression régulière) | Aucun filtre |
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
  "total": 51822,
  "total_pages": 10365,
  "results": [ { "id": 570, "title": "Dota 2", "...": "..." } ]
}
```

## Format des erreurs

Les erreurs sont renvoyées en texte brut, avec le code HTTP correspondant.

| Code | Cas |
|---|---|
| 400 | Body JSON invalide, id non numérique, paramètre de liste invalide |
| 401 | Identifiants incorrects, ou token absent, invalide ou expiré |
| 404 | Jeu introuvable (GET) |
| 500 | Erreur serveur ou base de données, y compris id déjà utilisé (POST) et jeu inexistant (PUT, DELETE) |

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

### Se connecter

Le token est gardé dans la variable `TOKEN`, utilisée par les exemples suivants :

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "monmotdepasse"}' \
  | sed -E 's/.*"token":"([^"]+)".*/\1/')
```

### Créer un jeu

```bash
curl -i -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
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
  -H "Authorization: Bearer $TOKEN" \
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
curl -i -X DELETE http://localhost:8080/api/v1/games/999999 \
  -H "Authorization: Bearer $TOKEN"
```

## Sécurité

- MongoDB n'est **pas exposé publiquement** : le port 27017 est publié uniquement sur `127.0.0.1`.
- L'API joint MongoDB par le réseau interne Docker.
- Les mots de passe sont hachés avec bcrypt dans `users.json`.
- Les routes `POST`, `PUT` et `DELETE` des jeux exigent un token JWT.
- `testuser` est un compte de test dont le mot de passe figure dans le dépôt : **à retirer de `users.json` en production**, et `JWT_SECRET` est à remplacer.

## Reste à faire

- [x] Pagination (`page`, `limit`)
- [x] Tri (`ordering`)
- [x] Filtres (`title`, `release_year`)
- [x] Authentification JWT sur POST, PUT et DELETE
- [x] Documentation Swagger / OpenAPI
- [ ] Erreurs au format JSON sur toutes les routes
- [ ] Codes 404 sur PUT et DELETE, 409 sur POST si l'id existe déjà
- [ ] Validation des champs obligatoires à la création d'un jeu
- [ ] CORS et logs
