# Diario progetto G-Gen

#### Data : 29 marzo 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Continuato codice per passare il gcode al visualizzatore
- Impostati i cookie di javascript e Go in modo che scadano in 3 giorni. E' importante che siano sincronizzati perchè quando il client controlla se ci sono determinati cookie, se ci sono cookie creati prima e dopo e quelli prima sono scaduti, possono esserci delle discrepanze dei dati disponibili. 
  In particolare se è memorizzato es. l'id del job e il nome dell'immagine, devono esistere tutti e due altrimenti verrebbe caricata l'immagine ma quando il client fa richiesta per altre informazioni al server, il jobId risulta vuoto e ci sono degli errori, anche se l'utente vede l'immagine e pensa che i suoi dati precedenti siano rimasti memorizzati (magari l'immagine si, ma altri dati necessari salvati nei cookie no se non scadono allo stesso momento assieme al session cookie)

## Problemi riscontrati e soluzioni

\-

## Punto di situazione del lavoro

Leggermente indietro

## Programma per la prossima volta

Sistemare ui.js da riga 150