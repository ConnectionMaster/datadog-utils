package main

import (
    "strings"
    datadog "github.com/DataDog/datadog-api-client-go/api/v1/datadog"
)


// There are some WidgetDefinition types that don't even have a 
// Request type in their struct (or don't have Q in the Request struct). We want to ignore those, since
// those won't contain a metric anyway.
func CheckIgnoredWidgetDefs(object datadog.WidgetDefinition) bool {

    switch {
    case object.AlertGraphWidgetDefinition != nil:
        return true
    case object.AlertValueWidgetDefinition != nil:
        return true
    case object.ChangeWidgetDefinition != nil:
        return true
    case object.CheckStatusWidgetDefinition != nil:
        return true
    case object.EventStreamWidgetDefinition != nil:
        return true
    case object.EventTimelineWidgetDefinition != nil:
        return true
    case object.FreeTextWidgetDefinition != nil:
        return true
    case object.GroupWidgetDefinition != nil:
        return true
    case object.HostMapWidgetDefinition != nil:
        return true
    case object.IFrameWidgetDefinition != nil:
        return true
    case object.ImageWidgetDefinition != nil:
        return true
    case object.LogStreamWidgetDefinition != nil:
        return true
    case object.MonitorSummaryWidgetDefinition != nil:
        return true
    case object.NoteWidgetDefinition != nil:
        return true
    case object.SLOWidgetDefinition != nil:
        return true
    case object.ScatterPlotWidgetDefinition != nil:
        return true
    case object.ServiceMapWidgetDefinition != nil:
        return true
    case object.ServiceSummaryWidgetDefinition != nil:
        return true
    }

    return false

}



func SearchExpressionsArray(expressionsArray *[]string, searchStr string) []string{
    var resultArray     []string

    for _, expression := range *expressionsArray {
        if strings.Contains(expression, searchStr) {
            resultArray = append(resultArray, expression)
        }
    }

    return resultArray
}
