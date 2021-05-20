package main

import (
    "flag"
    "fmt"
    "encoding/json"
)

func main() {
    var (
        dashMap             = make(map[string]map[string]string)
        metricSearchResults = make(map[string][]string)

        dashDetails         = flag.String("dash-details", "", "String: Get details for a specific dashboard by ID")
        dashList            = flag.Bool("dash-list", false, "Bool: will list all DataDog dashboard IDs")
        metricFind          = flag.String("find-metric", "", "String: Search our DataDog dashboards for the provided metric name to see if it is in use")

    )

    flag.Parse()

    if *metricFind != "" {
        fmt.Println("\nSearching for your metric across all dashboards.")
        metricSearchResults = FindMetric(*metricFind)
        fmt.Println("\n4. Feast your eyes upon the dashboards using your metric.")

        jsonMetrics, err := json.MarshalIndent(metricSearchResults, "", "  ")
        if err != nil {
            fmt.Println("ERROR: \n  ", err)
        }

        fmt.Println(string(jsonMetrics))

        var metricNames  []string
        for dashId, _ := range metricSearchResults {
            metricNames = append(metricNames, dashId)
        }

        details := make(map[string]string)
        for _, name := range metricNames {
            details = GetDashboardDetails(name)
            jsonDetails, _ := json.MarshalIndent(details, "", "  ")
            fmt.Println(string(jsonDetails))
        }


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
    }

}