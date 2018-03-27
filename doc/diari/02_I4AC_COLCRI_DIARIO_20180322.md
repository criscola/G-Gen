# Diario progetto G-Gen

#### Data : 22 marzo 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Ho iniziato a cambiare la grafica del sito ma ho riscontrato un problema grave spiegato nella prossima sezione

## Problemi riscontrati e soluzioni

Ho avuto un problema per quanto riguarda il server. Ogni volta che aggiornavo i miei file .css, il browser continuava a caricare i vecchi. Non era un problema di cache perchè continuavo a ricaricarle. Dopo 2 ore o giù di lì ho trovato che l'errore era in Vagrant, leggendo su internet si tratta di un bug nella sincronizzazione delle cartelle, ovvero quando si modificano i file sulla macchina host, le modifiche non vengono replicate correttamente nella macchina virtuale. In effetti facendo un `destroy` della macchina e ricaricandola da capo, le modifiche venivano apportate, e inoltre sulla vecchia VM di ubuntu su VMWare non c'era questo problema. Per ovviare a questo problema ho trovato alcune soluzioni che si applicavano solo a webserver quali nginx e apache e molte lamentele su internet, il bug è derivato da VirtualBox (hypervisor sul quale Vagrant si appoggia). Non mi è rimasto altro che utilizzare un altro tool come Vagrant, in questo caso Docker... per lo meno il file di configurazione della macchina che avevo fatto per Vagrant potrò riutilizzarlo per docker (i comandi bash per installare le dipendenze sono uguali dato che è sempre ubuntu).

## Punto di situazione del lavoro

Indietro

## Programma per la prossima volta

Cambiare il sito, impostare docker