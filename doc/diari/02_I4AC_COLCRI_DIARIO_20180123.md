# Diario progetto G-Gen

#### Data : 22 gennaio 2018

####Autore : Cristiano Colangelo

#### Luogo: SAM Trevano

## Lavori svolti

- Analisi del dominio/dei mezzi
  - Continuata documentazione inerente al capitolo

  - Ricercato il formato .dxf (per ev. trasformare da raster → .dxf → gcode) https://it.wikipedia.org/wiki/AutoCAD_DXF

  - Ricercato il formato .dxf (per ev. trasformare da raster → .stl → gcode) https://it.wikipedia.org/wiki/STL_(formato_di_file)

  - https://github.com/thearn/stl_tools (**tool in python** per convertire da .png a STL, non sembra il massimo perchè bisogna specificare alcune opzioni arcane di difficile comprensione, e devo tenere in conto che l'utente non deve indicare troppi parametri arcani difficili da comprendere)

  - Sarebbe non male poter utilizzare **OpenScad** per trasformare da raster a STL, anche perchè è perfettamente utilizzabile da riga di comando. Ho testato il suo utilizzo: è semplicissimo ed efficace.  Adesso devo capire come aggiustare l'immagine a mio piacimento.

    - https://en.wikibooks.org/wiki/OpenSCAD_User_Manual/Using_OpenSCAD_in_a_command_line_environment (uso da CLI) 
    - https://en.wikibooks.org/wiki/OpenSCAD_User_Manual/STL_Import_and_Export (STL import & export)

  - Ho capito che comunque devo leggere una heightmap per creare immagini 3D. Per le 2D basterà appiattire la **heightmap**, praticamente il plotter stamperebbe solo il contorno. Mi conviene quindi fare prima la versione 3D in modo da far uscire la 2D in modo più naturale.
    https://en.wikipedia.org/wiki/Heightmap

  - Ho deciso che utilizzerò probabilmente **Slic3r** da riga di comando per trasformare da STL a G-Code dato che è uno dei migliori software in circolazioni open-source e free, resta il problema della trasformazione da raster a STL.

    ​

## Problemi riscontrati e soluzioni

\-

## Punto di situazione del lavoro

Analisi dei mezzi/dominio e Gantt

## Programma per la prossima volta

- Continuare il Gantt preventivo
- Continuare analisi del dominio/dei mezzi