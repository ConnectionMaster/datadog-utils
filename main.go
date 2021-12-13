package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

func main() {
	var (
		dashMap              = make(map[string]map[string]string)
		metricSearchResults  = make(map[string][]string)
		monitorSearchResults = make(map[string][]string)
		searchWaitGrp        = new(sync.WaitGroup)

		dashDetails = flag.String("dash-details", "", "String: Get details for a specific dashboard by ID")
		dashList    = flag.Bool("dash-list", false, "Bool: Will list all DataDog dashboards.")
		metricFind  = flag.String("find-metric", "", "String: Search our DataDog dashboards for the provided metric name to see if it is in use")
		monitorList = flag.Bool("monitor-list", false, "Bool: Will list all DataDog monitors.")
	)

	flag.Parse()

	if *metricFind != "" {
		searchWaitGrp.Add(1)
		go func() {
			fmt.Println("\nSearching for your metric across all  & monitors.")
			metricSearchResults = FindMetric(*metricFind)
			fmt.Println("\n4. Feast your eyes upon the dashboards and monitors using your metric.")

			jsonMetrics, err := json.MarshalIndent(metricSearchResults, "", "  ")
			if err != nil {
				fmt.Println("ERROR: \n  ", err)
			}

			fmt.Println(string(jsonMetrics))

			var dashIds []string
			for dashId, _ := range metricSearchResults {
				dashIds = append(dashIds, dashId)
			}

			details := make(map[string]string)
			for _, name := range dashIds {
				details = GetDashboardDetails(name)
				jsonDetails, _ := json.MarshalIndent(details, "", "  ")
				fmt.Println(string(jsonDetails))
			}
			searchWaitGrp.Done()
		}()

		searchWaitGrp.Add(1)
		go func() {
			monitors := GetMonitors()
			monitorsWaitGrp := new(sync.WaitGroup)

			for _, monitor := range monitors {
				monitorsWaitGrp.Add(1)
				go func() {
					if strings.Contains(monitor.GetQuery(), *metricFind) {
						monitorSearchResults[monitor.GetName()] = append(monitorSearchResults[monitor.GetName()], monitor.GetQuery())
					}

					monitorsWaitGrp.Done()
				}()

				monitorsWaitGrp.Wait()
			}

			searchWaitGrp.Done()
		}()

		searchWaitGrp.Wait()

		jsonMonitors, err := json.MarshalIndent(monitorSearchResults, "", "  ")
		if err != nil {
			fmt.Println("ERROR: ", err)
		}

		fmt.Println(string(jsonMonitors))

	} else if *dashList != false {

		dashMap = GetDashboards()
		fmt.Println("Here is the list of dashboards in an attractive JSON format: ")
		jsonDashboards, _ := json.MarshalIndent(dashMap, "", "  ")
		fmt.Println(string(jsonDashboards))

	} else if *dashDetails != "" {

		details := GetDashboardDetails(*dashDetails)
		jsonDetails, err := json.MarshalIndent(details, "", "  ")
		if err != nil {
			fmt.Println("ERROR: \n  ", err)
		}

		fmt.Println("Here are the high-level details: ")
		fmt.Println(string(jsonDetails))

		// fmt.Println("Widget ID's and expressions are: ")
		// fmt.Println(string(expressionJson))
	} else if *monitorList != false {
		var (
			monitorMap = make(map[string]map[string]string)
			monitorId  string
		)
		monitors := GetMonitors()

		for _, monitor := range monitors {
			monitorId = strconv.FormatInt(monitor.GetId(), 10)
			monitorCreatorName := monitor.Creator.GetName()
			monitorMap[monitor.GetName()] = map[string]string{
				"ID":      monitorId,
				"Creator": monitorCreatorName,
				"Query":   monitor.GetQuery(),
			}
		}

		jsonMonitors, err := json.MarshalIndent(monitorMap, "", "  ")
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}

		fmt.Println("Here are all of the monitors & some info: \n")
		fmt.Println(string(jsonMonitors))
	}

}
