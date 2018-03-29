module 2b8c6fa99871fc2cbd6a9c75c96795c6()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) 2b8c6fa99871fc2cbd6a9c75c96795c6_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 2]) {
	2b8c6fa99871fc2cbd6a9c75c96795c6();
}