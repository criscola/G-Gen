# Diario progetto G-Gen

#### Data : 10 aprile 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Implementata generazione codice con Ultimaker
- Sistemato un bug con i Cookie

## Problemi riscontrati e soluzioni

Ho riscontrato un problema con Cura per la generazione del codice Ultimaker, in pratica cercavo di modificare alcuni parametri passandoli da riga di comando ma non cambiava nulla nel codice finale, allora trovando una risposta su Github ho capito che si possono modificare solo le impostazioni "foglie" del JSON (cioè voci che non hanno figli) mentre quelle "padre" sono usate solo dalla GUI per mostrare menu e cose del genere. Ho quindi modificato l'impostazione "foglia" ed ha funzionato

## Punto di situazione del lavoro

In linea

## Programma per la prossima volta

Iniziare doc, fare test 