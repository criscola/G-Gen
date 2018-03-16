# Diario progetto G-Gen

#### Data : 07 marzo 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Spostato il progetto su Ubuntu, messi in ordine tool di sviluppo
- Continuato a lavorare sulla generazione del gcode


## Problemi riscontrati e soluzioni

- Ho avuto qualche problema perchè il codice che funzionava su windows non funzionava su ubuntu, ho scoperto che la thread che gestiva il lavoro non aveva il tempo di funzionare su linux, quindi ho dovuto creare un meccanismo di sincronizzazione usando i channels di golang per comunicare fra la main thread e la thread di generazione, adesso ci sto ancora lavorando/cercando come fare

## Punto di situazione del lavoro

Stesura generazione GCode e barra di caricamento

## Programma per la prossima volta

Continuare generazione GCode