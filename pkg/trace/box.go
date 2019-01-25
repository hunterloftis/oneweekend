package trace

import "github.com/hunterloftis/oneweekend/pkg/geom"

type Box struct {
	List
}

func NewBox(min, max geom.Vec, m Material) *Box {
	return &Box{List: *NewList(
		NewRect(geom.NewVec(min.X(), min.Y(), max.Z()), geom.NewVec(max.X(), max.Y(), max.Z()), m),
		NewFlipped(NewRect(geom.NewVec(min.X(), min.Y(), min.Z()), geom.NewVec(max.X(), max.Y(), min.Z()), m)),
		NewRect(geom.NewVec(min.X(), max.Y(), min.Z()), geom.NewVec(max.X(), max.Y(), max.Z()), m),
		NewFlipped(NewRect(geom.NewVec(min.X(), min.Y(), min.Z()), geom.NewVec(max.X(), min.Y(), max.Z()), m)),
		NewRect(geom.NewVec(max.X(), min.Y(), min.Z()), geom.NewVec(max.X(), max.Y(), max.Z()), m),
		NewFlipped(NewRect(geom.NewVec(min.X(), min.Y(), min.Z()), geom.NewVec(min.X(), max.Y(), max.Z()), m)),
	)}
}
