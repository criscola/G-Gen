# Diario progetto G-Gen

#### Data : 26 marzo 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Creato server Ubuntu Server
  - Impostato virtual drive sul Windows per poter sviluppare su Windows https://kb.vmware.com/s/article/1022525
  - Impostato SSH per potersi collegare dall'esterno
  - Impostato programmi, variabili d'ambiente ecc.
- Lavorato sulla nuova grafica del sito

## Problemi riscontrati e soluzioni

Ho lasciato perdere docker perchè è troppo complicato creare una configurazione un minimo funzionante... ho creato una vm con ubuntu server. Ho avuto solo un problema con l'ssh che non funzionava perchè a quanto pare l'interfaccia virtuale su windows si era buggata. Ho impostato un altro network sulla NAT ed ha magicamente funzionato, anche l'interfaccia buggata si è "collegata" (prima era grigia anche se era abilitata).

## Punto di situazione del lavoro

Indietro

## Programma per la prossima volta

Capire perchè l'upload delle immagini si è danneggiato