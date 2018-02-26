# Diario progetto G-Gen

#### Data : 26 febbraio 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Lavorato sull'invio/risposta fra server e client per quanto riguarda l'upload dell'immagine. Ho creato un sistema che crea una sessione quando l'utente inizia a generare il gcode e vi ho memorizzato il nome del file caricato da lui. Ho fatto in modo che venga generato sempre un file singolo per utente in modo che non si occupi spazio inutilmente. Naturalmente dovrebbe esserci anche un cronjob che spazza via il contenuto di quella cartella quando il timestamp dei file è > di un certo tempo. 

  Ho implementato anche la risposta dal server in modo che, quando l'immagine è stata caricata con successo, il backend ritorni il nome dell'immagine. Questo nome viene preso quindi da javascript per caricare l'immagine dal server in FabricJS


## Problemi riscontrati e soluzioni

\-

## Punto di situazione del lavoro

Stesura sito web

## Programma per la prossima volta

Continuare a implementare fabricjs con il caricamento dell'immagine