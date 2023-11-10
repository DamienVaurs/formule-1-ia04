package types

import "math/rand"

type Circuit struct {
	Id                string    // Circuit ID
	Name              string    // Circuit name
	Country           string    // Circuit country
	Portions          []Portion // Portions
	MeteoDistribution []int     // Distribution of meteo conditions
}

/*MeteoDistribution est construit de la sorte : [40,70].
* Ce cas signifie que 40% du temps, le temps sera RAINY, 30% DRY, 30% HEAT
 */
func (c *Circuit) GenerateMeteo() Meteo {
	var dice int = rand.Intn(100) //génère nombre entre 0 et 99
	if dice < c.MeteoDistribution[0] {
		return RAINY
	}
	if dice < c.MeteoDistribution[1] {
		return DRY
	}
	return HEAT
}
