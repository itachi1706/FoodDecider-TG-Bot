package utils

import "math"

// Vincenty formula constants
const (
	a = 6378137.0         // Semi-major axis of the Earth (meters)
	f = 1 / 298.257223563 // Flattening of the Earth
	b = (1 - f) * a       // Semi-minor axis of the Earth (meters)
)

// VincentyDistance calculates the distance between two points (latitude and longitude) using the Vincenty formula
func VincentyDistance(lat1, lon1, lat2, lon2 float64) float64 {
	sq := func(x float64) float64 { return x * x }
	degToRad := func(x float64) float64 { return x * math.Pi / 180 }

	var lambda, tmp, q, p float64 = 0, 0, 0, 0
	var sigma, sinSigma, cosSigma float64 = 0, 0, 0
	var sinAlpha, cos2Alpha, cos2Sigma float64 = 0, 0, 0
	var c float64 = 0

	C := (sq(a) - sq(b)) / sq(b)

	uX := math.Atan((1 - f) * math.Tan(degToRad(lat1)))
	sinUX := math.Sin(uX)
	cosUX := math.Cos(uX)

	uY := math.Atan((1 - f) * math.Tan(degToRad(lat2)))
	sinUY := math.Sin(uY)
	cosUY := math.Cos(uY)

	l := degToRad(lon2) - degToRad(lon1)

	lambda = l

	for i := 0; i < 10; i++ {

		tmp = math.Cos(lambda)
		q = cosUY * math.Sin(lambda)
		p = cosUX*sinUY - sinUX*cosUY*tmp

		sinSigma = math.Sqrt(q*q + p*p)
		cosSigma = sinUX*sinUY + cosUX*cosUY*tmp
		sigma = math.Atan2(sinSigma, cosSigma)

		sinAlpha = (cosUX * cosUY * math.Sin(lambda)) / sinSigma
		cos2Alpha = 1 - sq(sinAlpha)
		cos2Sigma = cosSigma - (2*sinUX*sinUY)/cos2Alpha

		c = f / 16.0 * cos2Alpha * (4 + f*(4-3*cos2Alpha))
		tmp = lambda
		lambda = (l + (1-c)*f*sinAlpha*(sigma+c*sinSigma*(cos2Sigma+c*cosSigma*(-1+2*cos2Sigma*cos2Sigma))))

		if math.Abs(lambda-tmp) < 0.00000001 {
			break
		}
	}

	uu := cos2Alpha * C
	a := 1 + uu/16384*(4096+uu*(-768+uu*(320-175*uu)))
	b := uu / 1024 * (256 + uu*(-128+uu*(74-47*uu)))

	deltaSigma := (b * sinSigma * (cos2Sigma + 1.0/4.0*b*(cosSigma*(-1+2*sq(cos2Sigma))*(-3+4*sq(sinSigma))*(-3+4*sq(cos2Sigma)))))

	return b * a * (sigma - deltaSigma)
}
