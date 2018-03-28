# Diario progetto G-Gen

#### Data : 28 marzo 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Creato bottone per il download del GCode
- Iniziata l'integrazione e la modifica del visualizzatore
  - Stralciati menu e statistiche FPS
  - Modificati controlli zoom, distanza massima e minima zoom..
  - Iniziato a creare il codice per passare gcode dalla pagina all'iframe del viewer

## Problemi riscontrati e soluzioni

Ho lasciato perdere docker perchè è troppo complicato creare una configurazione un minimo funzionante... ho creato una vm con ubuntu server. Ho avuto solo un problema con l'ssh che non funzionava perchè a quanto pare l'interfaccia virtuale su windows si era buggata. Ho impostato un altro network sulla NAT ed ha magicamente funzionato, anche l'interfaccia buggata si è "collegata" (prima era grigia anche se era abilitata).

## Punto di situazione del lavoro

Leggermente indietro

## Programma per la prossima volta

Completare il caricamento del gcode nel viewer