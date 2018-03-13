module c1b63d86cc421185fdcc68e1ba23a0e4_1()
{
 /* Generated by trace2scad version 20150415
    http://aggregate.org/MAKE/TRACE2SCAD/
    Optimized model has 387/4796 original points
 */
 color([0.48, 0.48, 0.48])
 assign(minx=0) /* polygon minx=39 */
 assign(miny=0) /* polygon miny=45 */
 assign(maxx=12880) /* polygon maxx=12880 */
 assign(maxy=20000) /* polygon maxy=20001 */
 assign(dx=maxx-minx)
 assign(dy=maxy-miny)
 assign(maxd=((dx>dy)?dx:dy))
 scale([1/maxd, 1/maxd, 1])
 translate([-minx-dx/2, -miny-dy/2, 0])
 linear_extrude(height=1, convexity=387)
 union() {
  union() {
   polygon([[5647,19985],[5644,19969],[5475,18318],[5410,17711],[5320,16915],[5234,16197],[5145,15482],[5059,14834],[4964,14157],[4874,13537],[4665,12218],[4484,11188],[4405,10838],[4333,10585],[4155,10079],[3886,9456],[3559,8845],[3490,8733],[5668,5528],[7840,8857],[7509,9705],[7155,10715],[7017,11150],[6811,11873],[6705,12311],[6662,12506],[6486,13394],[6396,13917],[6240,14978],[6116,15944],[6030,16678],[5880,18058],[5770,19043],[5680,19778],[5651,20001]], convexity=6);
   polygon([[8050,8518],[8050,8512],[7430,7616],[6660,6416],[5797,5038],[5749,4969],[6258,4990],[6691,4975],[7115,4919],[7708,4781],[8437,4599],[8938,4500],[9218,4458],[10128,4471],[10645,4527],[10988,4595],[11423,4720],[11815,4892],[12106,5066],[12370,5276],[12585,5498],[12815,5794],[12880,5888],[12293,5834],[11717,5819],[11326,5841],[10943,5900],[10710,5959],[10446,6045],[10110,6198],[9746,6423],[9473,6645],[9085,7036],[8687,7540],[8380,7998],[8059,8525],[8050,8525]], convexity=9);
   polygon([[3242,8321],[3137,8213],[1763,6925],[1120,6243],[684,5682],[470,5308],[321,4901],[247,4578],[215,4293],[215,3847],[225,3728],[301,3250],[403,2910],[449,2810],[650,3118],[937,3473],[1358,3893],[1691,4164],[2085,4429],[2388,4599],[2782,4778],[3195,4920],[3503,4998],[3783,5050],[4071,5088],[4393,5110],[4618,5115],[5095,5090],[5503,5042],[4975,5923],[3983,7454],[3390,8363],[3347,8430]], convexity=7);
   polygon([[485,1606],[472,1593],[427,1481],[410,1372],[441,1223],[550,1033],[876,588],[930,488],[953,412],[948,360],[843,306],[728,248],[701,135],[710,85],[950,101],[1046,125],[1163,170],[1274,230],[1370,347],[1370,428],[1191,698],[957,929],[665,1200],[595,1386],[546,1593],[509,1620],[497,1620]], convexity=8);
   polygon([[5255,1608],[5244,1596],[5191,1478],[5187,1328],[5202,1241],[5314,1037],[5670,547],[5730,391],[5706,349],[5614,307],[5490,243],[5465,96],[5470,85],[5709,100],[5813,124],[5933,171],[6043,230],[6135,344],[6123,458],[6055,565],[5890,777],[5449,1192],[5365,1365],[5325,1544],[5301,1607],[5267,1620]], convexity=8);
   polygon([[6904,1596],[6511,1595],[6443,1516],[6418,1448],[6429,1020],[6456,847],[6470,725],[6540,364],[6601,85],[6626,85],[6725,268],[6845,410],[6881,447],[6825,591],[6789,699],[6765,785],[6722,1086],[6711,1309],[6776,1336],[7053,1326],[7284,1307],[7313,1513],[7307,1590],[7298,1598]], convexity=10);
   polygon([[7576,1588],[7568,1579],[7603,1512],[7692,1363],[7680,1304],[7570,1181],[7550,1148],[7560,1045],[7641,930],[7732,883],[7813,859],[7853,887],[7835,933],[7745,1032],[7755,1074],[7832,1213],[7834,1398],[7790,1494],[7703,1575],[7584,1598]], convexity=10);
   polygon([[11000,1593],[10995,1588],[11009,1538],[11110,1381],[11113,1324],[10994,1180],[10975,1149],[10985,1044],[11059,940],[11153,883],[11254,863],[11280,904],[11172,1041],[11225,1152],[11259,1238],[11265,1338],[11235,1463],[11123,1575],[11005,1598]], convexity=10);
   polygon([[1569,1581],[1557,1568],[1588,1527],[1670,1493],[1902,1418],[2041,1338],[2124,1243],[2147,1183],[2141,1047],[2085,838],[2063,784],[2059,635],[2170,448],[2440,184],[2578,110],[2667,74],[2710,71],[2695,113],[2435,475],[2346,630],[2340,766],[2346,796],[2453,879],[2601,930],[2610,949],[2487,1132],[2436,1247],[2440,1304],[2492,1326],[2633,1290],[2669,1268],[2690,1270],[2635,1458],[2626,1588],[1582,1595]], convexity=13);
   polygon([[8149,1591],[8135,1589],[8118,1564],[8140,1481],[8179,1303],[8242,1041],[8310,814],[8420,523],[8551,258],[8702,45],[8735,45],[8728,105],[8686,273],[8625,524],[8545,905],[8466,1292],[8468,1313],[8693,1326],[8864,1350],[9148,1416],[9267,1450],[9368,1485],[9445,1531],[9412,1575],[9319,1595],[8163,1594]], convexity=9);
   polygon([[9854,1589],[9846,1583],[9680,1319],[9634,1139],[9742,930],[9978,603],[10253,304],[10434,153],[10547,90],[10663,82],[10708,78],[10683,126],[10375,452],[10235,682],[10217,722],[10197,762],[10363,799],[10520,877],[10643,980],[10737,1085],[10670,1472],[10648,1590],[10358,1593],[10354,1571],[10430,1469],[10450,1308],[10419,1132],[10353,1055],[10290,1035],[10093,1045],[10033,1065],[9965,1132],[9945,1295],[9975,1393],[10097,1540],[10134,1583],[9862,1595]], convexity=13);
   polygon([[3300,1580],[3288,1575],[3116,1429],[2933,1209],[2853,1070],[2862,1038],[3330,436],[3451,243],[3472,170],[4141,189],[4142,206],[4003,250],[3856,320],[3404,778],[3345,860],[3388,883],[3750,891],[3660,1071],[3642,1101],[3207,1105],[3255,1193],[3362,1290],[3412,1290],[3455,1248],[3521,1259],[3579,1297],[3605,1339],[3398,1539],[3313,1584]], convexity=12);
   polygon([[433,852],[418,844],[152,641],[60,466],[41,319],[54,254],[113,172],[230,110],[308,91],[545,92],[551,256],[336,265],[255,321],[235,361],[245,479],[341,647],[465,844],[448,860]], convexity=8);
   polygon([[5175,835],[5148,815],[4978,698],[4897,610],[4850,520],[4808,378],[4819,273],[4870,183],[5011,110],[5170,88],[5307,86],[5329,248],[5094,266],[5022,333],[5003,366],[5010,465],[5055,573],[5205,796],[5227,836],[5229,855],[5203,854]], convexity=8);
  }
 }
}

module c1b63d86cc421185fdcc68e1ba23a0e4()
{
 /* all layers combined, scaled to within a 1mm cube */
 scale([1, 1, 1/1])
 difference() {
  union() {
   scale([1,1,2]) translate([0,0,-0.5]) c1b63d86cc421185fdcc68e1ba23a0e4_1();
  }
  translate([0,0,-2]) cube([2,2,4],center=true);
 }
}

scale([125, 125, 4]) {
	c1b63d86cc421185fdcc68e1ba23a0e4();
}