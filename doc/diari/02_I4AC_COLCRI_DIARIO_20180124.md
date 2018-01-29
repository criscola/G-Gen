# Diario progetto G-Gen

#### Data : 24 gennaio 2018

####Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Analisi del dominio/dei mezzi
  - Provato a far funzionare png23d, ma al momento della compilazione con `make install` gcc ritorna un errore di compilazione nel file bitmap.c dicendo che la variabile `STDIN_IOFIL` non è dichiarata. Non ho trovato soluzioni su internet al problema, così ho provato a vedere da dove saltasse fuori la variabile e ho trovato che faceva parte di una determinata libreria. L'ho inclusa nel sorgente ma mi sono saltati fuori altri problemi in altri file, allora ho lasciato perdere.
  - Provato a far funzionare img2scad ma facendo partire il `setup.py` il compilatore di Python segnala un errore di sintassi, allora ho lasciato perdere dato che si vede che non capisco come l'autore possa aver distribuito sorgenti che restituiscono errori di _sintassi_...
  - Riprovato a far funzionare la "litografia" di trace2scad con l'immagine di esempio del sito ufficiale (la faccia dell'autore) ma anche questa ci impiega tantissimo. Anche se alla fine l'elaborazione restituisse un risultato, non sarebbe fattibile far attendere così tanto gli utenti sul sito dato che probabilmente lo darebbero per non funzionante dopo nemmeno 1 minuto di attesa (io credo). Almeno ho appurato che funziona bene con immagini stile loghi.

Le restanti 4 ore del pomeriggio sono stato assente.

## Problemi riscontrati e soluzioni

Vedesi lavori svolti

## Punto di situazione del lavoro

Analisi dei mezzi/dominio e Gantt

## Programma per la prossima volta

- Continuare il Gantt preventivo
- Continuare analisi del dominio/dei mezzi

