# Programmation Répartie
## 🚪 Introduction
Ce travail est réalisé dans le cadre d'un laboratoire du module SDR du Bachelor Informatique est Systèmes de Communication.

L'objectif était de réaliser une application client/serveur permettant la répartition de bénévoles pour l'organisation d'évenements.

Le cahier des charges détaillé est disponible [ici](docs/Labo_2_SDR.pdf).
### 🧍🏻‍♂️🧍🏽‍♂️ Auteurs
* Nelson Jeanrenaud
* Vincent Peer
## 📚 Guide d'utilisation
### 💾 Installation des ressources
Commencez par cloner notre repository dans le dossier de votre choix, la commande
git est la suivante :  
>git clone https://github.com/VincentPeer/SDR_Labo1.git
  
Une fois effectué, vous disposez du projet et ne reste plus qu'à mettre en service 
le serveur et le(s) client(s).

## Lancement de l'application
Pour démarrer un serveur ou un client, il faut premièrement se rendre en ligne de commande
dans le dossier main :
> SDR_Labo1\src\main 

On y trouve les dossiers mainClient, le mainServer et mainMultiServer qui contiennent chacun un 
fichier main.go qui sera celui à exécuter pour un comportement souhaité.
### Lancement d'un serveur simple
Il s'agit de la version du laboratoire 1.
Pour lancer le serveur, il suffit de se rendre dans le dossier `mainServer` et de lancer la commande suivante :
>go run . [option]

Où `option` est facultatif et propose :  
* `-P` ou `--port` : permet de spécifier le port sur lequel le serveur doit écouter (par défaut 3333)  
* `-H` ou `--host` : permet de spécifier l'adresse sur laquelle le serveur doit écouter (par défaut localhost)
* `-C` ou `--config` : permet de spécifier le dossier dans lequel le serveur doit chercher les fichiers de configuration (par défaut ./config)
* `-D` ou `--debug` : permet d'activer le mode debug (par défaut false)

### Lancement de plusieurs serveurs
Pour le mode multi-serveur, il faut se rendre dans le dossier `mainMultiServer`. Nous avons implémenté 3 serveurs 
possédant des ids de 0 à 2. Comme défini dans le fichier `networkConfig.json` Il faut les lancer dans l'ordre croissant de leur id. Pour lancer un serveur, il faut entrer la commande suivante :
> go run . [serverId]

Où `[serverId]` est obligatoire et indique l'id du serveur à lancer.
Les log concernant le trafic réseau sont affichés, que ce soit des requêtes/réponses avec un client ou des messages de synchronisation entre serveurs.

### Lancement d'un client
Pour lancer un client avec le mode multi-serveur, il suffit de se rendre dans le dossier `mainClient` et d'entrer la commande suivante :
>go run . [clientName] -I [serverId] [option] 

Où `[clientName]` est le nom du client à lancer, obligatoire dans la version multi-serveur.  
`[serverId]` est l'id du serveur auquel il doit se connecter.  
`-I serverId`, ou `--id serverId` est optionnel, sans id entré, un id aléatoire sera choisi parmi la liste des serveurs (id 0 à 2).  

`option` est facultatif et propose :
* `-H` ou `--host` : permet de spécifier l'adresse sur laquelle le client doit se connecter (par défaut localhost)
* `-D` ou `--config` : permet d'activer le mode debug (par défaut false)

Si aucun argument n'est fourni, c'est-à-dire uniquement `go run .`, nous établissons une connexion telle qu'au laboratoire 1 avec le
serveur simple.
### 🦟 Mode debug 
Le mode debug permet de voir les messages échangés entre le serveur et le client.
Pour l'activer, il suffit de lancer le serveur avec l'argument `-D` ou `--debug`.

Pour tester les races conditions, il suffit de lancer le client avec l'argument `-d` ou `--debug` également.
L'accès au ressources par des clients lancés de cette manière est bloqué pendant 5 secondes, permettant de tester le conditions de concurrence.

### Test de la concurrence
Voici une marche à suivre permettant de tester les conditions de concurrence.   
Commencez par lancer les 3 serveurs dans l'ordre croissant de leur id dans le dossier `mainMultiServer` :
> go run . 0  
> go run . 1  
> go run . 2  

Lancez ensuite 2 clients, un en mode debug et un en mode normal :
> go run . client1 -I 0 -D  
> go run . client2 -I 0   

A présent, on peut vérifier la concurrence en profitant du client en debug qui est bloqué pendant 5 secondes pour chaque commande, la suite du test consiste à
faire un test de lecture puis un test en écriture.
#### Test en lecture
Dans le client en mode debug, entrez 3 afin de lister les événements, puis entrez aussitôt 3 dans le client normal. Entrées  à saisir :
> clientDebug : 3  
> clientNormal : 3

On observe que le client normal doit attendre que le client debug ait fini sa requête pour afficher les évênements. En effet, le client normal n'a pas directement accès à la section critique tant le client debug n'a pas terminé.
#### Test en écriture
Dans le client en mode debug, saisissez les entrées suivantes afin de modifier la base de donnée en ajoutant un nouvel évenement :  
(Ne pas oublier que certaine étape demande au client debug d'attendre 5 secondes).
> clientDebug : 1   (création d'un événement)  
> clientDebug : a   (username)  
> clientDebug : 1   (password)  
> clientDebug : manif   (nom manifestation)  
> clientDebug : cuisine     (nom job)  
> clientDebug : 3   (nombre de bénévole)  
> clientDebug : STOP    (fin de la saisie)  
> clientNormal : 3   (affichage des événements)  

On observe que le client normal doit attendre que le client debug ait fini sa requête afin d'afficher la liste des événements
et que le nouvel événement `manif` a bien été ajouté.  

Le même test peut être réalisé en connectant les clients sur des serveurs différents. On peut également lancer plus de 2 clients.


### 👨🏽‍⚕️ Utilitaire godoc
Afin d'avoir une documentation claire de nos packages, fonctions et l'ensemble
de notre projet, il est possible de générer un fichier html contenant les
commentaires précisés pour chaque entité. Ce fichier peut ensuite être visualisé
en lancant un serveur local.
Pour cela, il faut installer godoc sur votre machine en tapant la commande :  
`go install -v golang.org/x/tools/cmd/godoc@latest`  
Ensuite, à partir d'un terminal dans le dossier *SDR_Labo1*, tapez la commande  
`godoc -http=:6060`  
Dans votre navigateur, entrez l'URL  
`http://localhost:6060/pkg/SDR_Labo1/`  
Vous pouvez alors parcourir notre documentation. Les packages main ne sont pas visibles, ainsi que les fonctions
non exportées.

### 🔎 Détails d'implémentation
* Lorsqu'une saisie concerne l'id d'une manifestation ou d'un poste, l'indice commence à 0.
* Lorsque l'utilisateur doit se loguer, il ne peut plus revenir en arrière et n'a pas d'autre choix que de réussir le log in.
* Les alignements des colonnes pour les affichages de manifestation, poste et bénévole fonctionnent tant que l'utilisateur n'entre pas de données extrêmement longues.
* Le client (`mainClient/main.go`) contient en dur les ports et les id des serveurs disponibles. Si le nombre de serveur vient à évoluer, il faudra mettre à jour ce fichier en ajoutant l'id du nouveau serveur et son port, le reste de l'implémentation prendra en compte ces changements.
#### Diagramme de classes
![Diagramme de classes](docs/uml.svg)
### Implémentation de lamport
On crée une go routine pour chaque serveur, qui va écouter sur le port défini dans le fichier de configuration.
Lorsqu'un client envoie une requête, le client envoie la requête à la go routine qui gère la base de données du serveur, assurant ainsi l'exclusion mutuelle sur la base de données pour les clients.

Pour l'exclusion mutuelle entre les serveurs, on utilise un algorithme de Lamport. Chaque serveur possède une estampille et enregistre le dernier message qu'il a envoyé. Il enregistre également, pour chaque serveurs le dernier message reçu et l'estampille associée. 
Lorsque le serveur veut accéder à la section critique il envoit via des canaux, une requête d'accès à la section critique au go routines communiquant avec les autres serveurs. Il attend ensuite la réponse de ces derniers qui est renvoyée à la go routine qui gère la base de données via un autre canal. Lorsque la go routine a reçu la réponse de tous les serveurs, elle peut accéder à la section critique. Lorsque la section critique est terminée, le serveur envoit un message de libération de la section critique aux autres serveurs.

Avant de quitter la section critique, si l'action sur la base de données la modifiée (ajout, suppression, modification), le serveur envoit un message de synchronisation à tous les autres serveurs. Ce message contient l'action à effectuer sur la base de données du serveur qui a reçu le message et l'estampille du serveur qui a effectué l'action. Lorsque le serveur reçoit un message de synchronisation, il effectue l'action sur sa base de données et met à jour son estampille.
### Tests automatisés
Une série de tests automatisés ont pour but de tester les fonctions utilisées pour la communication entre le client et
le serveur. Dans le projet, ils se trouvent dans les sources dans le dossier `integrationTests`. Emplacement :
> SDR_Labo1\src\integrationTests

Pour lancer les tests, il faut d'abord lancer un serveur, le plus simple est de se rendre dans `SDR_Labo1\src\integrationTests` et d'entrer `go run .` pour lancer un serveur simple.
Ensuite, il faut se rendre dans le dossier `integrationTests` et effectuer la commande suivante :
> go test

Le résultat des tests est simple, on obtient un `PASS` si tout a passé avec succès, et un `FAIL` sinon.
Pour avoir plus de détail sur le noms des tests effetués, l'option -v peut être ajouté à la commande :
>go test -v

Important: Comme les tests modifient les données, une fois la commande `go test`lancée, on ne peut pas la relancer une seconde fois car il y aura des erreurs. En effet, nous avons par exemple testé qu'il est possible pour un bénévole de s'inscrire à un poste qui requière encore 1 personne. Or, cette place ne sera plus disponible si nous relançons les tests une deuxième fois et lèvera une erreur.  
Pour relancer, il suffit de terminer l'exécution et de la relancer.



## 📖 Protocole
### ⬅ Format des messages
Les paramètres sont séparés par des virgules, les messages sont séparés par des points-virgules.
Le premier paramètre est toujours le type du message. Les différents types de messages sont les suivants :
* `LOGIN` : demande de connexion au serveur
* `CREATE_EVENT` : demande de création d'une manifestation
* `CLOSE_EVENT` : demande de fermeture d'une manifestation
* `EVENT_REG` : demande d'inscription à une manifestation
* `GET_EVENTS` : demande de récupération des manifestations
  * Si le paramètre 2 est spécifié, alors le serveur envoie les informations de la manifestation correspondant à l'id. Sinon, il envoie toutes les manifestations.
* `GET_JOBS` : demande de récupération des postes d'une manifestation spécifique
* `STOP` : demande d'arrêt de la connexion
* `DEBUG` : demande de lancement du mode debug

### ➡ Format des réponses
Le serveur répond au client avec des messages formatés de la même manière que les messages envoyés par le client. Mais les types de messages sont différents :
* `OK` : réponse positive
* `NOTOK` : réponse négative

