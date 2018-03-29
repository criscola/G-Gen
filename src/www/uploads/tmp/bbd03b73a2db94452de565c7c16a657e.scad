module bbd03b73a2db94452de565c7c16a657e()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) bbd03b73a2db94452de565c7c16a657e_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 2]) {
	bbd03b73a2db94452de565c7c16a657e();
}