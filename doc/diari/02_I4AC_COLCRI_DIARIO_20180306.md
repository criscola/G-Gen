# Diario progetto G-Gen

#### Data : 06 marzo 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Completata la barra di caricamento e interazione col server, iniziata stesura codice generazione gcode, una volta finito questo, mi rimarrà solo da implementare la visualizzazione, modificare un attimo i controlli e aggiungere il download del gcode (relativamente semplice). In pratica memorizzo nella sessione 2 struct, `GeneratorJob` che detiene le informazioni inerenti alla generazione del gcode a livello di processo, mentre `GeneratorParams` detiene i parametri della generazione. Questi 2 struct sono poi passati a una goroutine (thread) che genera il gcode, al termine produce un file .gcode. Il client manda ogni 1,5 secondi una richiesta GET per chiedere lo stato della generazione, il server ritorna una cifra che sarebbe la percentuale di completamento. Il client aggiorna la barra con questa cifra se essa è cambiata. Al raggiungimento del 100% il client smette di mandare richieste GET e manda una richiesta GET per avere il link di download. 


## Problemi riscontrati e soluzioni

- Non riuscivo a memorizzare nella sessione una map contenente degli struct, cercando ho trovato che mancava una riga `gob.Register(map[string]*GeneratorJob{})`. In pratica è necessario registrare che tipo di dato di vuole memorizzare nelle sessioni perchè le sessioni vengono serializzate con un certo encording

## Punto di situazione del lavoro

Stesura generazione GCode e barra di caricamento

## Programma per la prossima volta

Continuare generazione GCode