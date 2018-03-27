module d5318f9f4b4a4d939c7d9688f99df570()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) d5318f9f4b4a4d939c7d9688f99df570_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 2]) {
	d5318f9f4b4a4d939c7d9688f99df570();
}