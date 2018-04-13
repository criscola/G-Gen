# Diario progetto G-Gen

#### Data : 11 aprile 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Corretto bug generazione G-Code con Ultimaker, avevo dimenticato un argomento nell'invocazione da riga di comando
- Corretto bug download G-Code, si aprivano finestre di download pari al numero di volte che l'utente rifaceva la procedura di generazione. Avevo aggiunto il listener dell'evento del click del bottone in un punto sbagliato del programma, mi è bastato spostare il listener fuori da quella posizione
- Aggiunto controllo in caso l'immagine non fosse disponibile, se non lo fosse l'utente potrebbe ora caricare un'altra immagine. Ciò è essenziale in caso sia memorizzato ancora il Cookie contenente il percorso dell'immagine ma l'immagine venisse eliminata dal server (ciò è normale in quanto il server non deve memorizzare le immagini degli utenti persistentemente)
- Parlato con il committente per chiedergli se volesse aggiungere ulteriori parametri, egli mi ha detto che desidererebbe poter definire il punto di partenza dell'ugello della stampante 3D. Cercherò di implementare questa cosa se ci sarà ancora tempo
- Iniziata la documentazione dell'implementazione

## Problemi riscontrati e soluzioni

-

## Punto di situazione del lavoro

In linea

## Programma per la prossima volta

Continuare doc implementazione