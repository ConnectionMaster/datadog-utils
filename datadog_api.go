package main

import (
    "context"
    "fmt"
    "os"
    "strconv"
    "sync"
    datadog "github.com/DataDog/datadog-api-client-go/api/v1/datadog"
)


var (
    DD_API_CLIENT *datadog.APIClient
    CTX context.Context
)

func init() {
    DD_API_CLIENT, CTX = createApiClient()
}


func createApiClient() (*datadog.APIClient, context.Context){
    ctx := datadog.NewDefaultContext(context.Background())
    configuration := datadog.NewConfiguration()
    apiClient := datadog.NewAPIClient(configuration)

    return apiClient, ctx
}


func FindMetric(searchMetric string) map[string][]string{
    var (
        expressionWaitGrp   = new(sync.WaitGroup)
        searchResults       = make(map[string][]string)
        waitGrp             = new(sync.WaitGroup)
        widgetPtr             *[]datadog.Widget
        widgetExpressions   = make(map[string]map[int64][]string)
    )
    fmt.Println("\n1. Fetch all dashboards...")
    fmt.Println("\n2. Fetch the widgets from each dashboard, and retrieve every widget's expressions... asynchronously.")
    for id, _ := range GetDashboards() {
        // Part 1 - fetch the widgets array from each dashboard object
        widgetPtr = getDashboardWidgets(id)

        if len(*widgetPtr) < 1 {
            // fmt.Printf("\nLength of widget array for dashboard: %s ---- was zero. \n", id)
            continue
        } else {
            // fmt.Printf("\nLength of widget array for dashboard: %s ---- was %v\n", id, len(*widgetPtr) )
            waitGrp.Add(1)
            go func() {
                // Part 2 - Check against every set of widgets from
                // Dashboard & fetch the expression strings
                widgetExpressions[id] = GetWidgetExpressions(widgetPtr)
                fmt.Printf("-")
                waitGrp.Done()
            }()
        }
    waitGrp.Wait()
    }
    fmt.Println("----------------------------------------------------------------------------------------------")
    fmt.Printf("\n3. Search for metric --- %s --- across each expression, for every widget, in every dashboard.\n", searchMetric)
    fmt.Println("----------------------------------------------------------------------------------------------")

    for dash, expressionsMap := range widgetExpressions {
        expressionWaitGrp.Add(1)
        go func() {
            for _, expressionArray := range expressionsMap {
                 expressionHits := SearchWidgetExpressionsArray(&expressionArray, searchMetric)
                 if len(expressionHits) > 1 {
                    searchResults[dash] = expressionHits
                 } else {
                     continue
                 }
            }
            expressionWaitGrp.Done()
        }()
    expressionWaitGrp.Wait()
    }

    return searchResults
}


func GetDashboards() map[string]map[string]string{
    var dashMap = make(map[string]map[string]string)

    resp, r, err := DD_API_CLIENT.DashboardsApi.ListDashboards(CTX)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DashboardsApi.ListDashboards`: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }

    for _, dash := range *resp.Dashboards {
        dashMap[*dash.Id] = map[string]string{
            "Title": *dash.Title,
            "Author": *dash.AuthorHandle,
            "Last Modified": dash.ModifiedAt.String(),
        }
    }

    return dashMap
}


func GetDashboardDetails(dashboardID string) map[string]string {
    var (
        detailsResult     = make(map[string]string)
    )

    dashDetails, _, err := DD_API_CLIENT.DashboardsApi.GetDashboard(CTX, dashboardID)

    if err != nil {
        fmt.Println("ERROR: \n ", err)
    }
    widgetCount := strconv.Itoa(len(dashDetails.Widgets))
    detailsResult = map[string]string{
        "Author": *dashDetails.AuthorHandle,
        "ID": *dashDetails.Id,
        "URL": *dashDetails.Url,
        "Widgets": widgetCount,
    }

    return detailsResult
}

// A simple fetcher for the widgets array of a dashboard object.
// We just return a pointer to the array so we aren't passing an array
// all over the place wherever we may need this.
func getDashboardWidgets(dashboardID string) *[]datadog.Widget {
    var(
        widgetsArray    []datadog.Widget
    )

    dashDetails, _, err := DD_API_CLIENT.DashboardsApi.GetDashboard(CTX, dashboardID)
    if err != nil {
        fmt.Println("\nERROR: %v\n ", err)
    }
    widgetsArray = dashDetails.Widgets
    return &widgetsArray
}


func GetWidgetExpressions(widgetsArray *[]datadog.Widget) map[int64][]string{
    var (
        widgetExpressions = make(map[int64][]string)
        widgetDefinition  datadog.WidgetDefinition
        widgetId  int64
    )

    for _, widget := range *widgetsArray {
        widgetDefinition = widget.GetDefinition()
        widgetId = widget.GetId()

        if CheckIgnoredWidgetDefs(widgetDefinition) {
            continue
        } else {

            switch {
            case widgetDefinition.TimeseriesWidgetDefinition != nil:
                for _, request := range widget.Definition.TimeseriesWidgetDefinition.Requests {
                    widgetExpressions[widgetId] = append(widgetExpressions[widgetId], request.GetQ())
                }

            case widgetDefinition.ToplistWidgetDefinition != nil:
                for _, request := range widget.Definition.ToplistWidgetDefinition.Requests {
                    widgetExpressions[widgetId] = append(widgetExpressions[widgetId], request.GetQ())
                }

            case widgetDefinition.TableWidgetDefinition != nil:
                for _, request := range widget.Definition.TableWidgetDefinition.Requests {
                    widgetExpressions[widgetId] = append(widgetExpressions[widgetId], request.GetQ())
                }

            case widgetDefinition.HeatMapWidgetDefinition != nil:
                for _, request := range widget.Definition.HeatMapWidgetDefinition.Requests {
                    widgetExpressions[widgetId] = append(widgetExpressions[widgetId], request.GetQ())
                }

            case widgetDefinition.DistributionWidgetDefinition != nil:
                for _, request := range widget.Definition.DistributionWidgetDefinition.Requests {
                    widgetExpressions[widgetId] = append(widgetExpressions[widgetId], request.GetQ())
                }

            case widgetDefinition.GeomapWidgetDefinition != nil:
                for _, request := range widget.Definition.GeomapWidgetDefinition.Requests {
                    widgetExpressions[widgetId] = append(widgetExpressions[widgetId], request.GetQ())
                }
            }
        }
    }

    return widgetExpressions
}
