module 5c16e7cede887b90ef1f6a50944af7bf()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) 5c16e7cede887b90ef1f6a50944af7bf_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 2]) {
	5c16e7cede887b90ef1f6a50944af7bf();
}