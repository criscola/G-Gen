module 66be130d793d1f2f40493f00316c9324()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) 66be130d793d1f2f40493f00316c9324_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 2]) {
	66be130d793d1f2f40493f00316c9324();
}