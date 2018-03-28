module 2f97141429906c42d315520d98a75e76()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) 2f97141429906c42d315520d98a75e76_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 2]) {
	2f97141429906c42d315520d98a75e76();
}