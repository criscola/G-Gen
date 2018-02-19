# Diario progetto G-Gen

#### Data : 07 febbraio 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Continuato a leggere la guida https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/03.2.html su Go e su come usarlo per il web 
- Iniziata stesura sito web partendo da questo sito http://tutsplus.github.io/smooth-page-transition-with-history-web-api/ che provvede una buona base da cui partire e con la transizione che avevo proprio in mente (un'animazione che scatta al cambio della pagina, servirà per il processo di conversione da foto a G-Code)
- Guardato un po' di siti dove si può testare la compatibilità del proprio sito web online
  - https://www.browserling.com/ contiene i browser principali fino a molte versioni vecchie tranne che explorer che è a pagamento
  - https://netrenderer.com per testare la compatibilità con explorer
- http://tutsplus.github.io/smooth-page-transition-with-history-web-api/ Bottoni "pressabili" CSS

## Problemi riscontrati e soluzioni

- Ho avuto un problema avviando il webserver in Ubuntu, non funzionava sulla porta 80 quindi cercando su internet ho visto che in ubuntu sono necessari i privilegi di root per ascoltare su porte sotto la 1024, ho usato quindi il comando `sudo go run ggen.go` per avviarlo correttamente

## Punto di situazione del lavoro

Stesura sito web

## Programma per la prossima volta

- Continuare a lavorare sui tools