module a19a12ba3cca29b9ee51b6d476657b0c()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) a19a12ba3cca29b9ee51b6d476657b0c_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 2]) {
	a19a12ba3cca29b9ee51b6d476657b0c();
}