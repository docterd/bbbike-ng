package bbbikeng

import (
	"math"
	"strconv"
)

const RADIUS = 6371

type Point struct {
	Lat float64
	Lng float64
}

type GeoJSON struct {
	Type        string
	Coordinates [][2]float64
}

type GeoJSONPoint struct {
	Type        string
	Coordinates [2]float64
}


func (f *Point) SetCoordinates(lat float64, lng float64) {

	f.Lat = lat;
	f.Lng = lng

}

func (f *Point) Coordinates()(lat float64, lng float64) {
	return f.Lat, f.Lng

}

func MakeNewPoint(lat float64, lng float64) (newPoint Point) {

	newPoint.SetCoordinates(lat, lng)
	return newPoint

}

func MakeNewPointFromString(lat string, lng string) (newPoint Point) {

	xPath, err := strconv.ParseFloat(lng, 64)
	yPath, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		panic(err)
	}
	return MakeNewPoint(xPath, yPath)
}

func (f *Point) PointIsValid() bool {
	return f.Lat != 0.0 && f.Lng != 0.0
}

func (f *Point) Compare(comparePoint Point) (equal bool) {
	thresholdLat := math.Abs(math.Abs(f.Lat) - math.Abs(comparePoint.Lat))
	thresholdLng := math.Abs(math.Abs(f.Lng) - math.Abs(comparePoint.Lng))
	return (thresholdLat <= 0.0000001 && thresholdLng <= 0.0000001)
}

func (f *Point) LatitudeLongitudeAsString() (lat string, lng string) {
	lat = strconv.FormatFloat(f.Lat, 'f', 6, 64)
	lng = strconv.FormatFloat(f.Lng, 'f', 6, 64)
	return lat, lng
}

func (this *Path) CorrectPath(parentNode *Node)() {

	isForwardConnected := parentNode.NodeGeometry.Compare(this.Path[0])
	if (!isForwardConnected && !parentNode.flippedDirection) || (isForwardConnected && parentNode.flippedDirection) {
		var flippedPath []Point
		for i := len(this.Path)-1; i >= 0; i-- {
			point := this.Path[i]
			flippedPath = append(flippedPath, point)
		}
		this.Path = flippedPath
	}
}

func DistanceFromPointToPoint(firstPoint Point, secondPoint Point) (meters int) {

	dLat := degreeToRadians(secondPoint.Lat - firstPoint.Lat)
	dLng := degreeToRadians(secondPoint.Lng - firstPoint.Lng)

	lat1 := degreeToRadians(firstPoint.Lat)
	lat2 := degreeToRadians(secondPoint.Lat)

	a := math.Sin(dLat/2) * math.Sin(dLat/2) + math.Sin(dLng/2) * math.Sin(dLng/2) * math.Cos(lat1) * math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := RADIUS * c;

	return int(d * 1000)

}

func magnitude(firstPoint Point, secondPoint Point) (magnitude float64) {

	var newPoint Point
	newPoint.Lat = secondPoint.Lat - firstPoint.Lat
	newPoint.Lng = secondPoint.Lng - secondPoint.Lng

	return math.Sqrt(math.Pow(newPoint.Lat, 2) + math.Pow(newPoint.Lng, 2))

}

func BearingBetweenPoints(firstSegment Point, secondSegment Point) (angle float64) {

	_, dLng := pointDifference(firstSegment, secondSegment)

	y := math.Sin(dLng) * math.Cos(secondSegment.Lat)
	x := math.Cos(firstSegment.Lat) * math.Sin(secondSegment.Lat) - math.Sin(firstSegment.Lat) * math.Cos(secondSegment.Lat) * math.Cos(dLng)

	return radiansToDegrees(math.Atan2(y, x))
}


func DistanceFromLinePoint(points []Point) (distance int) {

	for i := 0; i < len(points)-1; i++ {
		firstPoint := points[i]
		secondPoint := points[i+1]
		distance += DistanceFromPointToPoint(firstPoint, secondPoint)
	}
	return distance
}

func DistanceFromPointToPath(point Point, path []Point) (distance int) {

	distance = -1.0
	for i := 0; i < len(path)-1; i++ {

		firstPoint := path[i]
		secondPoint := path[i+1]
		magnitude := magnitude(secondPoint, firstPoint)

		U := (((point.Lat - firstPoint.Lat) * (secondPoint.Lat - firstPoint.Lat)) * ((point.Lng - firstPoint.Lng) * (secondPoint.Lng - firstPoint.Lng))) / math.Pow(magnitude, 2)
		if U > 0.0 || U < 1.0 {
			var newIntersection Point
			newIntersection.Lat = firstPoint.Lat + U*(secondPoint.Lat-firstPoint.Lat)
			newIntersection.Lng = firstPoint.Lng + U*(secondPoint.Lng-firstPoint.Lng)

			lastDistance := DistanceFromPointToPoint(point, newIntersection)
			if lastDistance > distance || lastDistance <= 0 {
				distance = lastDistance
			}
		}

	}

	return distance

}


func pointDifference(firstPoint Point, secondPoint Point) (dLat float64, dLon float64) {


	dLat = secondPoint.Lat - firstPoint.Lat
	dLon = secondPoint.Lat - firstPoint.Lat

	return dLat, dLon

}

func degreeToRadians(degree float64) (radians float64){

	return (degree * math.Pi / 180)
}

func radiansToDegrees(radians float64) ( degrees float64) {

	return (radians * 180 / math.Pi)

}
