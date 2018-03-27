module 04d238497ebfbc02c457096e08bf9fcb()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) 04d238497ebfbc02c457096e08bf9fcb_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 2]) {
	04d238497ebfbc02c457096e08bf9fcb();
}