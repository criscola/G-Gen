# Diario progetto G-Gen

#### Data : 21 febbraio 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Cercato un po' come generare il modello 3D per vedere se infine farlo o no. Sono riuscito a scaricare https://www.thingiverse.com/thing:15276 questa applicazione scritta in Java per convertire le heightmap in modelli .stl, il che sarebbe perfetto. Ho scaricato un'immagine di esempio di Super Mario da Google e l'ho caricata in Gimp per poi renderla in bianco e nero, invertire i colori e giocare un po' con contrasto e luminosità per rendere più accentuate le differenze di colore (e quindi più dettagliato il modello 3D), ovviamente sono stato attento che ogni operazione eseguita possa venire riprodotta da riga di comando. Stavo pensando che poi possa mostrare all'utente tramite delle immagini esplicative come dovrebbe giocare con luminosità e contrasto nel caso decida di creare un'immagine 3D. Il programma ha funzionato e dopo aver giocato con i 2 parametri che accetta (altezza del modello e altezza della base) sono riuscito a creare un modello 3D decente. L'unico difetto è che anche se l'immagine png è trasparente, crea lo stesso una base sotto il modello. Per ovviare a questo problema stavo pensando di eventualmente chiedere all'utente se desideri una base oppure no. Se non la desidera, importare l'.stl generato in OpenScad e rimuoverla tramite un rettangolo "eliminante" posto sotto l'immagine. Questo sarebbe fattibile con un po' più di lavoro ma chiaramente non so nemmeno se avrò abbastanza tempo per implementarlo, comunque si può fare
- Guardato come fare l'upload drag&drop dell'immagine
  - https://dev.to/gcdcoder/how-to-upload-files-with-golang-and-ajax
  - https://css-tricks.com/drag-and-drop-file-uploading/
- Continuato a lavorare sul sito
  - Finito slider per form. Adesso si può andare avanti e indietro. Ho fatto in modo che alla pressione dei bottoni (next, back) venga lanciato una metodo che fa selezionare la slide a un certo indice. Potrò aggiungere quindi i vari controlli etc. prima di andare avanti e indietro con le slides
- Ho deciso che userò FabricJS perchè glfx.js non conteneva una funzione per invertire le immagini http://fabricjs.com/ che servirà per il 3D


## Problemi riscontrati e soluzioni

\-

## Punto di situazione del lavoro

Stesura sito web

## Programma per la prossima volta

Implementare http://evanw.github.io/glfx.js/ e continuare con la stesura del sito