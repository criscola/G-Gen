module fa5688689a85ed5d4a3e8a964126ba41()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) fa5688689a85ed5d4a3e8a964126ba41_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 2]) {
	fa5688689a85ed5d4a3e8a964126ba41();
}