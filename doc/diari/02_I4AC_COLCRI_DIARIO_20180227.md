# Diario progetto G-Gen

#### Data : 27 febbraio 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Finito upload dell'immagine con tanto di controlli dell'estensione e possibilità di eliminare l'immagine. Se l'immagine viene eliminata dal client, manda una richiesta DELETE al server per eliminare l'immagine anche da là. 


## Problemi riscontrati e soluzioni

- Ho avuto un problema con la cache del browser. In pratica eliminavo sia il file fisicamente dal server che dal client con javascript, tuttavia ricaricava sempre la stessa immagine anche se ne caricavo un'altra, allora ho scoperto che quando Fabricjs cercava di caricare l'immagine nel canvas, avendo lo stesso nome, la cache locale dava la stessa immagine. Ho quindi modificato un po' il sistema per fare in modo di cambiare sempre il nome dell'immagine, questo mi ha anche permesso di semplificare un po' il codice

## Punto di situazione del lavoro

Stesura sito web

## Programma per la prossima volta

Continuare a implementare luminosità/contrasto di fabricjs