# Diario progetto G-Gen

#### Data : 23 gennaio 2018

####Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Analisi del dominio/dei mezzi
  - Convertire SVG a .scad https://github.com/Spiritdude/SVG2SCAD
  - Convertire immagini in .scad + lista di tools http://aggregate.org/MAKE/TRACE2SCAD/
  - Tool python per convertire immagini in .scad https://github.com/l0b0/img2scad

- Creata VM di ubuntu per testare trace2scad. Ho seguito il sito originale cercando di ricreare da me il logo klingon e ci sono riuscito. In pratica ho preso il codice (uno script bash) dal sito e l'ho incollato nelal mia VM in un file. Poi ho scaricato il logo klingon e ho fatto partire il comando :

  `$ ./trace2scad.sh -o fd.scad klingon_dondewi.png`

  Lo script mi ha generato un file .scad da cui ho potuto prendere il codice generato (allocato in un `module`) e scalarlo in questo modo in OpenScad:

  `scale([150,150,3]) klingon_dondewi();`

  Infine ho renderizzato l'immagine. Ho provato anche con un'immagine per cercare di fare una litografia di essa con il comando:

  `./trace2scad,sh -f 0 -l 8 testface.jpg`

  Ma non sono riuscito a generare un'immagine renderizzabile per il momento

## Problemi riscontrati e soluzioni

- Problema nel setup di una macchina virtuale Ubuntu per fare dei test

Al momento della creazione della vm mi è stato chiesto di disabilitare il Device Guard, quindi ho seguito questa guida https://kb.vmware.com/s/article/2146361

## Punto di situazione del lavoro

Analisi dei mezzi/dominio e Gantt

## Programma per la prossima volta

- Continuare il Gantt preventivo
- Continuare analisi del dominio/dei mezzi

