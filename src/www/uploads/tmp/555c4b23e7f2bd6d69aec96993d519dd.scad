module 555c4b23e7f2bd6d69aec96993d519dd()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) 555c4b23e7f2bd6d69aec96993d519dd_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 2]) {
	555c4b23e7f2bd6d69aec96993d519dd();
}