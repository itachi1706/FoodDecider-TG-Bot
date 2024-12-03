package services

import (
	"FoodDecider-TG-Bot/utils"
	"context"
	"errors"
	"googlemaps.github.io/maps"
	"log"
)

type GeocodingAPI interface {
	GetLocationFromPlusCode(plusCode string) (float64, float64, error)
	GetAddressFromLocation(latitude, longitude float64) (*maps.GeocodingResult, error)
}

type GeocodingAPIImpl struct {
	client *maps.Client
}

func NewGeocodingAPI() GeocodingAPI {
	key := utils.GetEnvDefault("GMAPS_API_KEY", "")
	c, err := maps.NewClient(maps.WithAPIKey(key))
	if err != nil {
		log.Fatalf("failed to create new google maps client: %v", err)
	}
	return &GeocodingAPIImpl{client: c}
}

func (g *GeocodingAPIImpl) GetLocationFromPlusCode(plusCode string) (float64, float64, error) {
	log.Println("Getting Location from Plus Code: " + plusCode)

	req := &maps.GeocodingRequest{
		Address: plusCode,
	}

	resp, err := g.client.Geocode(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to get location from plus code: %v", err)
		return 0, 0, err
	}

	if len(resp) == 0 {
		log.Println("No location found")
		return 0, 0, errors.New("no location found")
	}

	// Get first location
	location := resp[0].Geometry.Location

	return location.Lat, location.Lng, nil
}

func (g *GeocodingAPIImpl) GetAddressFromLocation(latitude, longitude float64) (*maps.GeocodingResult, error) {
	log.Printf("Getting Address from Coordinate: %f, %f\n", longitude, latitude)

	req := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: latitude,
			Lng: longitude,
		},
	}

	resp, err := g.client.ReverseGeocode(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to get address from location: %v", err)
		return nil, err
	}

	if len(resp) == 0 {
		log.Println("No address found")
		return nil, errors.New("no address found")
	}

	// Get first address
	address := resp[0]

	return &address, nil
}
