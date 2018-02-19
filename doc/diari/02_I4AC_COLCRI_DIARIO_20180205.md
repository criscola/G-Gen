# Diario progetto G-Gen

#### Data : 05 febbraio 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Continuato ad analizzare i tool `trace2scad`, `openscad`, `slic3r` e `GCode Viewer`. Trace2scad funziona bene con immagini semplici e piatte, mentre openscad dà degli errori a volte per via del modello non manifold (con delle intersecazioni che non vanno bene al motore di render CGAL). Slic3r funziona molto bene e non ho trovato problemi. Anche GCode Viewer è praticamente perfetto e non vedo problemi, una volta datogli in pasto il g-code lui disegna così com'è. L'unico problema è appunto trace2scad che ha bisogno di immagini senza trasparenza e possibilmente con una buona risoluzione. Ho trovato anche che conviene utilizzare il parametro `-e <n>`. La pagina ufficiale dice: 

  > When the model is converted to line segments, very often there will be cases where a sequence of multiple line segments could be converted into a single line segment with minimal error. This option sets the maximum distance, in pixels, that a discarded point can be from the line that replaces it.

  In pratica serve a semplificare i poligoni che vengono generati, senza compromettere comunque il risultato finale (che a detta del creatore potrebbe migliorare in determinati casi). Devo studiare meglio l'opzione `-f <n>` per generare dei layer in base alla luminosità dei pixel per creare un effetto 3D, perchè sembra dare alcuni problemi (rende i modelli non manifold...). Devo infine studiarmi ancora gli altri parametri

- Continuato a informarmi sul go e su come creare il backend per l'applicazione, in particolare questo minilibro sembra ben fatto https://astaxie.gitbooks.io/build-web-application-with-golang/en/

- Guardato un po' di roba sulle animazioni dei siti web. Mi servirà perchè ho delle idee su come strutturare la grafica del sito (rendendola scorrevole, es. step 1 fai questo, step 2 fai quest'altro, step 3 visualizza il risultato accetta/rifai) https://www.jqueryscript.net/animation/Basic-Cross-platform-One-Page-Scroll-Plugin-jQuery-fullpage.html

## Problemi riscontrati e soluzioni

- Ho avuto un problema con un'immagine che aveva sfondo trasparente, riempiendo lo sfondo con il secchiello, l'immagine si è generata correttamente. Probabilmente a trace2scad non piace la trasparenza, quindi dovrò fare altri test e vedere ev. come eliminare riempire il canale alfa di es. bianco in modo da eliminare la trasparenza delle mie immagini

## Punto di situazione del lavoro

Studio dei tool/studio sul sito

## Programma per la prossima volta

- Continuare a testare i tool, lavorare sul sito