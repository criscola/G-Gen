module 7b4c16b7cdb155c6f7bc9e0874013615()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) 7b4c16b7cdb155c6f7bc9e0874013615_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 2]) {
	7b4c16b7cdb155c6f7bc9e0874013615();
}