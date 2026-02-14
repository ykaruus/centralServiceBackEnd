package enums

type RegionFlag int

const (
	REGION_SOUTH     = 1287391
	REGION_NORTH     = 8708123
	REGION_SOUTHEAST = 2343249
	REGION_NORTHEAST = 1347233
	REGION_MIDWEST   = 1283829
)

var MapperRegionFlags map[RegionFlag]string = map[RegionFlag]string{
	REGION_MIDWEST:   "centro-oeste",
	REGION_SOUTH:     "sul",
	REGION_SOUTHEAST: "sudeste",
}

var RemapperRegionFlags map[string]RegionFlag = map[string]RegionFlag{
	"centro-oeste": REGION_MIDWEST,
	"sul":          REGION_SOUTH,
	"sudeste":      REGION_SOUTHEAST,
}
