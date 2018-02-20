# Diario progetto G-Gen

#### Data : 20 febbraio 2018

#### Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Letto https://blue-jay.github.io/views/#basic-usage, probabilmente il miglior esempio di sito web funzionate fatto in MVC per go. Aggiunge molte cose che non sono necessarie per il mio progetto ma è comunque utile quando non parla di componenti esterni
- https://github.com/julienschmidt/httprouter sarebbe non male utilizzare un router esterno per gestire le connessioni HTTP in quanto quello incluso con la libreria standard è molto semplice, ad esempio non è presente un comportamento di default per routes non riconosciute (es se scrivo `localhost/ciao` nella barra di navigazione del browser e non vi è effettivamente un handler per questa route, essa viene gestita dalla route di default, in poche parole ritorna sulla homepage invece di ritornare un errore 404 pagina non trovata)
- Creato un file `base.tmpl` contenente la base di tutte le pagine del sito. E' popolato con alcune tag che cambiano sempre come il titolo della pagina, l'head (le info base come charset, percorsi assets, javascript ecc. sono comunque già presenti, ma si può ampliare se una determinata pagina necessita ad esempio di uno script in più), il contenuto della pagina. Ho aggiunto anche la possibilità di modificare il footer ma in quanto rimarrà sempre uguale penso lo toglierò. I contenuti delle pagine sono contenute nella cartella `/templates`, es `/templates/home`
- Scremato un po' di codice inutile
- Implementato https://github.com/kenwheeler/slick/ per il form. Saranno da creare i form in html dei vari step e modificare i bottoni avanti/indietro in modo da gestire i controlli. 


## Problemi riscontrati e soluzioni

- Non riuscivo ad applicare l'approccio di blue-jay del suo esempio con un file base.tmpl e vari file che contenevano i vari contenuti (es. `home/index.tmpl` contiene i dati della homepage che viene passata al base.tmpl per il parsing), leggendo questo ho capito https://stackoverflow.com/questions/36617949/how-to-use-base-template-file-for-golang-html-template 

## Punto di situazione del lavoro

Stesura sito web

## Programma per la prossima volta

Implementare http://evanw.github.io/glfx.js/ e continuare con la stesura del sito