module 724256b28f7b59385b3f4163180fc0dc()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) 724256b28f7b59385b3f4163180fc0dc_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 2]) {
	724256b28f7b59385b3f4163180fc0dc();
}