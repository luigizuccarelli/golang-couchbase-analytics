package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"lmzsoftware.com/lzuccarelli/golang-couchbase-analytics/pkg/connectors"
	"lmzsoftware.com/lzuccarelli/golang-couchbase-analytics/pkg/schema"
	//"github.com/gorilla/mux"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

func SankeyChartHandler(w http.ResponseWriter, r *http.Request, conn connectors.Clients) {
	var response schema.Response
	//vars := mux.Vars(r)

	addHeaders(w, r)

	// remove leading htpps:// and trailing journey_id for uniquness
	query := `SELECT	page.referrerName AS source,
										page.urlName AS destination,
										COUNT(message_id) AS count
						FROM PAGEEVENTS
						GROUP BY	page.referrerName as source,
											page.urlName AS destination`

	ar, err := conn.AnalyticsQuery(query, nil)
	if err != nil {
		response = schema.Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not execute analytics query from couchbase %v", err)}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.MarshalIndent(response, "", "	")
		fmt.Fprintf(w, string(b))
		return
	}

	d, err := processSankeyResults(ar)
	if err != nil {
		response = schema.Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not process analytics results from couchbase %v", err)}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.MarshalIndent(response, "", "	")
		fmt.Fprintf(w, string(b))
		return
	}

	conn.Debug(fmt.Sprintf("Analytics result from couchbase  %v \n", d))
	response = schema.Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Data retrieved succesfully", Sankey: d}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"version\" : \"1.0.2\" , \"name\": \""+os.Getenv("NAME")+"\" }")
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func processSankeyResults(res connectors.AnalyticsResult) ([]schema.SankeyData, error) {
	var row map[string]interface{}
	var data []schema.SankeyData

	for res.Next() {
		res.Row(&row)
		if count, ok := row["count"]; ok {
			x, _ := strconv.Atoi(fmt.Sprintf("%v", count))
			d := schema.SankeyData{Value: x, From: fmt.Sprintf("%v", row["source"]), To: fmt.Sprintf("%v", row["destination"])}
			data = append(data, d)
		}
	}
	if err := res.Close(); err != nil {
		return nil, err
	}

	return data, nil
}
