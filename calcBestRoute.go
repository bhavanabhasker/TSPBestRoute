package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
)

type Estimate struct {
	Routes   []string
	Costs    int
	Duration int
	Distance float64
}

type Estimates []Estimate
type LocList []string

func (est Estimates) Less(i, j int) bool {
	return est[i].Costs < est[j].Costs
}

type CostCache struct {
	HighEstimate int
	Duration     int
	Distance     float64
}

type PriceEstimateCache struct {
	Start_latitude  float64
	Start_longitude float64
	End_latitude    float64
	End_longitude   float64
	CachedCost      CostCache
}

var PriceEstimateCacheSet = make([]PriceEstimateCache, 0)

func LookUpPriceCache(start_latitude float64, start_longitude float64, end_latitude float64,
	end_longitude float64) int {

	for i := 0; i < len(PriceEstimateCacheSet); i++ {
		if PriceEstimateCacheSet[i].Start_latitude == start_latitude &&
			PriceEstimateCacheSet[i].Start_longitude == start_longitude &&
			PriceEstimateCacheSet[i].End_latitude == end_latitude &&
			PriceEstimateCacheSet[i].End_longitude == end_longitude {
			return i
		}
	}
	return -1
}

var Chargesarrang = make(Estimates, 0)

func createbestrouteArray(arrang []string) Estimate {

	var routes Estimate
	var chargesarrang Estimate
	var pEst PriceEstimates

	for i := 0; i < len(arrang); i++ {
		if i+1 < len(arrang) {
			startaddress, _ := strconv.Atoi(arrang[i])
			startresult := FindinLocationService(startaddress)
			startlatitude := startresult.Coordinates.Lat
			startlongitude := startresult.Coordinates.Long
			endaddress, _ := strconv.Atoi(arrang[i+1])
			endresult := FindinLocationService(endaddress)
			endlatitude := endresult.Coordinates.Lat
			endlongitude := endresult.Coordinates.Long
			//fmt.Println(" Querying ", arrang[i+1])

			found := LookUpPriceCache(startlatitude, startlongitude, endlatitude, endlongitude)
			if found != -1 {
				pEst.Prices[0].Distance = PriceEstimateCacheSet[found].CachedCost.Distance
				pEst.Prices[0].Duration = PriceEstimateCacheSet[found].CachedCost.Duration
				pEst.Prices[0].HighEstimate = PriceEstimateCacheSet[found].CachedCost.HighEstimate
			} else {

				data := getResponse(startlatitude, startlongitude, endlatitude, endlongitude)
				if e := json.Unmarshal(data, &pEst); e != nil {
					log.Fatal(e)
				}

				var cc CostCache
				cc.HighEstimate = pEst.Prices[0].HighEstimate
				cc.Duration = pEst.Prices[0].Duration
				cc.Distance = pEst.Prices[0].Distance

				var pCache PriceEstimateCache
				pCache.Start_latitude = startlatitude
				pCache.Start_longitude = startlongitude
				pCache.End_latitude = endlatitude
				pCache.End_longitude = endlongitude
				pCache.CachedCost = cc

			}
			routes.Costs += pEst.Prices[0].HighEstimate
			routes.Duration += pEst.Prices[0].Duration
			routes.Distance += pEst.Prices[0].Distance
		}
	}
	chargesarrang.Routes = arrang
	chargesarrang.Distance = routes.Distance
	chargesarrang.Duration = routes.Duration
	chargesarrang.Costs = routes.Costs
	return chargesarrang
}

func calcBestRoute(startid string, locations []string) Estimate {
	perms := Perm(locations, startid)
	fmt.Print("Querying route: ", perms[0])
	bestPriceEstimate := createbestrouteArray(perms[0])
	fmt.Println(" : price ", bestPriceEstimate.Costs, " <- Best so far")
	for i := 1; i < len(perms); i++ {
		fmt.Print("Querying route: ", perms[i])
		local := createbestrouteArray(perms[i])
		fmt.Print(" : price ", local.Costs)
		if local.Costs < bestPriceEstimate.Costs {
			fmt.Print(" <- Best so far")
			bestPriceEstimate = local
		}
		fmt.Println("")
	}
	return bestPriceEstimate
}
func Permute(data sort.Interface) bool {
	i := data.Len() - 2
	for ; ; i-- { // find last ith location that is sorted
		if i < 0 {
			return false
		}

		if data.Less(i, i+1) {
			break
		}
	}

	// find jth location that's not sorted
	j := 0
	for j = data.Len() - 1; !data.Less(i, j); j-- {
	}

	data.Swap(i, j) // swap sorted and unsorted
	for i, k := i+1, data.Len()-1; i < k; i++ {
		data.Swap(i, k)
		k--
	}
	return true
}
func (slice LocList) Len() int {
	return len(slice)
}

func (slice LocList) Less(i, j int) bool {
	return slice[i] < slice[j]
}

func (slice LocList) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func PrePostPend(ls []string, st string) []string {
	pppend := make([]string, len(ls)+2)
	pppend[0] = st
	for i := 0; i < len(ls); i++ {
		pppend[i+1] = ls[i]
	}
	pppend[len(ls)+1] = st
	return pppend
}

func Perm(locations []string, startid string) [][]string {
	perm := make([][]string, 0)
	sort.Sort(LocList(locations))
	perm = append(perm, PrePostPend(locations, startid))
	for {
		if Permute(LocList(locations)) {

			cache := PrePostPend(locations, startid)
			perm = append(perm, cache)
		} else {
			return perm
		}
	}

}
