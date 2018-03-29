module e9510504141b615ae5ab5679b853d6ae()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) e9510504141b615ae5ab5679b853d6ae_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 2]) {
	e9510504141b615ae5ab5679b853d6ae();
}