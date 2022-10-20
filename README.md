# SDR_Labo1

## Guide d'utilisation
### Installation des ressources
Commencez par cloner notre repository dans le dossier de votre choix, la commande
git est la suivante :  
`git clone https://github.com/VincentPeer/SDR_Labo1.git`  
Une fois effectué, vous disposer du projet et ne reste plus qu'à mettre en service 
le serveur et le(s) client(s).

### Lancement du serveur


### Lancement d'un client

### Utilitaire godoc
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
Vous pouvez alors parcourir notre documentation.