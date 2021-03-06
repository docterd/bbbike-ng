/**
 * User: Dennis Oberhoff
 * To change this template use File | Settings | File Templates.
 */
package bbbikeng

import (
	"encoding/json"
	"log"
	"strings"
	"strconv"
)

const X0 = -780761.760862528
const X1 = 67978.2421158527
const X2 = -2285.59137120724
const Y0 = -5844741.03397902
const Y1 = 1214.24447469596
const Y2 = 111217.945663725

func ConvertStandardToWGS84(x float64, y float64) (xLat float64, yLat float64) {

	yLat = ((x-X0)*Y2 - ((y - Y0) * X2)) / (X1*Y2 - Y1*X2)
	xLat = ((x-X0)*Y1 - (y-Y0)*X1) / (X2*Y1 - X1*Y2)
	return xLat, yLat

}

func ConvertLatinToUTF8(iso8859_1_buf []byte) string {

	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)

}

func ConvertGeoJSONtoPoint(jsonInput string) (point Point) {

	points := ConvertGeoJSONtoPath(jsonInput)

	if len(points) > 0 {
		point = points[0]
	} else {
		log.Fatal("Error Converting Json", jsonInput)
	}

	return point

}

func ConvertGeoJSONtoPath(jsonInput string) (path []Point) {

	var f interface{}
	err := json.Unmarshal([]byte(jsonInput), &f)
	if err != nil {
		log.Fatal("JSON Unmarshal error:", err)
	}

	m := f.(map[string]interface{})
	dataType := m["type"]

	if dataType == "LineString" {
		var coordinates GeoJSON
		err := json.Unmarshal([]byte(jsonInput), &coordinates)
		if err != nil {
			log.Fatal("JSON Unmarshal error:", err)
		}
		for _, coord := range coordinates.Coordinates {
			path = append(path, MakeNewPoint(coord[1], coord[0]))
		}
	} else if dataType == "Point" {

		var coordinates GeoJSONPoint
		err := json.Unmarshal([]byte(jsonInput), &coordinates)
		if err != nil {
			log.Fatal("JSON Unmarshal error:", err)
		}

		point := MakeNewPoint(coordinates.Coordinates[1], coordinates.Coordinates[0])
		path = append(path, point)


	}
	return path
}

func ConvertPathToGeoJSON(path []Point)(jsonOutput string) {

	var jsonData []byte
	var err error

	if len(path) == 1 {
		var newJson GeoJSONPoint
		newJson.Type = "Point"
		newJson.Coordinates[1] = path[0].Lat
		newJson.Coordinates[0] = path[0].Lng
		jsonData, err = json.Marshal(newJson)

	} else {

		var newJson GeoJSON
		newJson.Type = "LineString"
		for _, point := range path {
			var newCoordinates [2]float64
			newCoordinates[1] = point.Lat
			newCoordinates[0] = point.Lng
			newJson.Coordinates = append(newJson.Coordinates, newCoordinates)
		}
		jsonData, err = json.Marshal(newJson)
	}

	if err != nil {
		log.Fatal("Failed to Convert Path to GeoJSON: %s", err.Error())
	}

	return string(jsonData)
}

func ConvertStringToIntArray(stringList string) (list []int) {

	stringList = strings.Replace(stringList, "{", "", -1)
	stringList = strings.Replace(stringList, "}", "", -1)
	stringList = strings.Replace(stringList, "NULL", "", -1)
	streetsSplitted := strings.Split(stringList, ",")


	for _, string := range streetsSplitted {
		converted, err := strconv.Atoi(string)
		if err == nil {
			list = append(list, converted)
		}
	}

	return list

}

func geoJsonInsert(geoJson string) (statement string) {

	return ("ST_TRANSFORM(ST_SetSRID(ST_GeomFromGeoJSON('"+ geoJson + "'), '4326'),4326)")

}

func ParseAttributes(raw string) (attributes []AttributeInterface){

	var genericAttribute []map[string]interface{}
	err := json.Unmarshal([]byte(raw), &genericAttribute)

	if err != nil {
		log.Fatal("Error while unmarshaling Attributes:", err)
		return nil
	}

	for _, entry := range genericAttribute {

		var newAttribute AttributeInterface
		switch entry["category"] {
				case "CA":
					newAttribute = new(CyclepathAttribute)
				case "GA":
					newAttribute = new(GreenwayAttribute)
				case "QA":
					newAttribute = new(QualityAttribute)
				case "UA":
					newAttribute = new(UnlitAttribute)
				case "TA":
					newAttribute = new(TrafficLightAttribute)
				case "HA":
					newAttribute = new(HandicapAttribute)
		}
		newAttribute.SetType(entry["type"].(string))
		newAttribute.SetPathFromGeoJSON(entry["geometry"])
		attributes = append(attributes, newAttribute)

	}
	return attributes
}

func (this *Attribute) SetPathFormGeoJSONString(rawjson string) () {
	var geojsonMap interface {}
	err := json.Unmarshal([]byte(rawjson), &geojsonMap)
	if err != nil {
		this.SetPathFromGeoJSON(geojsonMap)
	}
}

func (this *Attribute) SetPathFromGeoJSON(jsonInput interface {}) () {

	assertedMap := jsonInput.(map[string]interface{})
	geometryType := assertedMap["type"]

	if geometryType != "LineString" && geometryType != "Point" && geometryType != "MultiLineString" {

		log.Fatal("Error parsing Geometry:", assertedMap)
		return
	}

	geometryData := assertedMap["coordinates"].([]interface {})
	switch geometryType {
		case "Point" : {
			longitude := geometryData[0].(float64)
			latitude := geometryData[1].(float64)
			this.geometry = append(this.geometry, MakeNewPoint(latitude, longitude))
		}

		case "LineString": {
			for _, point := range geometryData {
				convertedInterface := point.([]interface {})
				longitude := convertedInterface[0].(float64)
				latitude := convertedInterface[1].(float64)
				this.geometry = append(this.geometry,  MakeNewPoint(latitude, longitude))
			}
		}
		case "MultiLineString": {
			interfaceLenght := len(geometryData)-1
			for i := 0; i <= interfaceLenght; i++ {
				convertedInnerPoint := geometryData[i].([]interface {})
				if i == 0 {
					for _, point := range convertedInnerPoint {
						convertedInterface := point.([]interface {})
						longitude := convertedInterface[0].(float64)
						latitude := convertedInterface[1].(float64)
						this.geometry = append(this.geometry, MakeNewPoint(latitude, longitude))
					}
				} else {
					convertedInterface := convertedInnerPoint[1].([]interface {})
					longitude := convertedInterface[0].(float64)
					latitude := convertedInterface[1].(float64)
					this.geometry = append(this.geometry, MakeNewPoint(latitude, longitude))
				}
			}
		}
	}

}

func (this *Attribute) GetGeoJSON() (jsonOutput string) {
	return ConvertPathToGeoJSON(this.geometry)
}
