package main

/*
This API will update ticket payload with weather
risk and fruad risk data

*/
import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	//FORMS
	STORM_FORM_ID = 114093996871
	TV_FORM_ID    = 114093998312
	//FIELDS
	WIND_SPEED_FIELD_ID      = 114100596852
	TV_MODEL_FIELD_ID        = 114099896612
	CLAIM_TYPE_FIELD_ID      = 114099964311 // TV , Storm Surge
	WEATHER_RISK_FIELD_ID    = 114100712171
	FRAUD_RISK_FIELD_ID      = 114100657672
	FRAUD_CHECK_REQ_FIELD_ID = 114101778331

	PHONE_FIELD_ID         = 114100658992
	EMAIL_FIELD_ID         = 114100659172
	DAMAGE_IMAGE_URL_1     = 114101729171
	DAMAGE_IMAGE_URL_2     = 114101729171
	RECEIPT_OR_EST_IMG_URL = 114101729331
	DAMAGE_VIDEO_URL       = 114101655952
	SETTLEMENT_AMOUNT      = 114101736271
)

func Handler(w http.ResponseWriter, r *http.Request) {
	mediaFields := [4]int64{114101729171, 114101729331, 114101655952, 114101729171}

	var composedResult = ComposedResult{}
	err := json.NewDecoder(r.Body).Decode(&composedResult)
	if err != nil {
		createErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// populate ticket payload with custom field
	ticketDetails := TicketDetails{}
	var claimType, riskDesc, windSpeed string
	if composedResult.TranformedData.Status == 200 {
		ticketDetails = composedResult.TranformedData.TicketDetails
		if strings.Contains(ticketDetails.Ticket.Subject, "Storm") {
			claimType = "Storm Surge"
			ticketDetails.Ticket.TicketFormID = STORM_FORM_ID
		} else {
			claimType = "TV"
			ticketDetails.Ticket.TicketFormID = TV_FORM_ID
		}
		println("Submission ID ---> " + ticketDetails.Ticket.EventID)
	} else {
		createErrorResponse(w, "Ticket Details not recieved", http.StatusBadRequest)
	}

	if composedResult.WeatherRisk.Status == 200 {
		weatherRisk := composedResult.WeatherRisk
		riskDesc = strconv.Itoa(weatherRisk.RiskScore) + " : " + weatherRisk.Description
		ticketDetails.addCustomFields(CustomFields{WEATHER_RISK_FIELD_ID, riskDesc})
		println("riskDesc", riskDesc)
		if weatherRisk.RiskScore <= 20 {
			ticketDetails.addCustomFields(CustomFields{FRAUD_CHECK_REQ_FIELD_ID, "true"})
			ticketDetails.addCustomFields(CustomFields{FRAUD_RISK_FIELD_ID, "HIGH"})
		} else if weatherRisk.RiskScore > 20 && weatherRisk.RiskScore <= 60 {
			ticketDetails.addCustomFields(CustomFields{FRAUD_CHECK_REQ_FIELD_ID, "true"})
			ticketDetails.addCustomFields(CustomFields{FRAUD_RISK_FIELD_ID, "MEDIUM"})
		} else {
			ticketDetails.addCustomFields(CustomFields{FRAUD_CHECK_REQ_FIELD_ID, "false"})
			ticketDetails.addCustomFields(CustomFields{FRAUD_RISK_FIELD_ID, "LOW"})
		}
	}
	if composedResult.WeatherData.Status == 200 {
		weatherData := composedResult.WeatherData
		windSpeed = weatherData.History.Dailysummary[0].Maxwspdm
		ticketDetails.addCustomFields(CustomFields{WIND_SPEED_FIELD_ID, windSpeed})
		println("windSpeed", windSpeed)
	}

	//add Custom Fields
	//common fields
	ticketDetails.addCustomFields(CustomFields{CLAIM_TYPE_FIELD_ID, claimType})
	ticketDetails.addCustomFields(CustomFields{PHONE_FIELD_ID, ticketDetails.Ticket.Requester.Phone})
	ticketDetails.addCustomFields(CustomFields{EMAIL_FIELD_ID, ticketDetails.Ticket.Requester.Email})
	ticketDetails.addCustomFields(CustomFields{SETTLEMENT_AMOUNT, "0"})

	//if images stored in Bucket
	mediaBucket := composedResult.MediaBucket
	if mediaBucket.Status == 200 {
		for i, media := range mediaBucket.Media {
			ticketDetails.addCustomFields(CustomFields{mediaFields[i], media.MediaLink})
		}
	}

	if (TVClaimData{}) != composedResult.TranformedData.TVClaimData {
		tvClaimData := composedResult.TranformedData.TVClaimData
		mediaBucket := composedResult.MediaBucket
		//TV Field
		ticketDetails.addCustomFields(CustomFields{TV_MODEL_FIELD_ID, tvClaimData.TVModelNo})
		if mediaBucket.Status != 200 {
			//typeform media files
			ticketDetails.addCustomFields(CustomFields{DAMAGE_IMAGE_URL_1, tvClaimData.DamageImageURL1})
			ticketDetails.addCustomFields(CustomFields{DAMAGE_IMAGE_URL_2, tvClaimData.DamageImageURL2})
			ticketDetails.addCustomFields(CustomFields{RECEIPT_OR_EST_IMG_URL, tvClaimData.TVReceiptImage})
			ticketDetails.addCustomFields(CustomFields{DAMAGE_VIDEO_URL, tvClaimData.DamageVideoURL})
		}
	}
	if (StromClaimData{}) != composedResult.TranformedData.StromClaimData {
		stormClaimData := composedResult.TranformedData.StromClaimData
		//Storm Fields
		if mediaBucket.Status != 200 {
			ticketDetails.addCustomFields(CustomFields{DAMAGE_IMAGE_URL_1, stormClaimData.DamageImageURL1})
			ticketDetails.addCustomFields(CustomFields{DAMAGE_IMAGE_URL_2, stormClaimData.DamageImageURL2})
			ticketDetails.addCustomFields(CustomFields{RECEIPT_OR_EST_IMG_URL, stormClaimData.RepairEstimateImage})
			ticketDetails.addCustomFields(CustomFields{DAMAGE_VIDEO_URL, stormClaimData.DamageVideoURL})
		}
	}

	//add status 200 if no error
	ticketDetails.Status = 200

	//marshal to JSON
	ticketDetailsJSON, err := json.Marshal(ticketDetails)
	if err != nil {
		createErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	println("After updating Risk data, Image dats etc. to ticket -----> ", string(ticketDetailsJSON))
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(string(ticketDetailsJSON)))

}

func (ticketDetails *TicketDetails) addCustomFields(customField CustomFields) {
	ticketDetails.Ticket.CustomFields = append(ticketDetails.Ticket.CustomFields, customField)
}

type MediaBucket struct {
	Status int     `json:"status,omitempty"`
	Media  []Media `json:"media,omitempty"`
}

type Media struct {
	Bucket       string `json:"bucket,omitempty"`
	Name         string `json:"name,omitempty"`
	Size         int64  `json:"size,omitempty"`
	MediaLink    string `json:"media-link,omitempty"`
	OriginalLink string `json:"original-link,omitempty"`
}

func main() {
	println("staritng app..")
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8084", nil)
}

func createErrorResponse(w http.ResponseWriter, message string, status int) {
	errorJSON, _ := json.Marshal(&Error{
		Status:  status,
		Message: message})
	//Send custom error message to caller
	w.WriteHeader(status)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(errorJSON))
}

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ComposedResult struct {
	MediaBucket    MediaBucket    `json:"media-bucket"`
	TranformedData TranformedData `json:"tranformed-data"`
	WeatherData    WeatherData    `json:"weather-data"`
	WeatherRisk    struct {
		Status      int    `json:"status"`
		Description string `json:"description"`
		RiskScore   int    `json:"riskScore"`
	} `json:"weather-risk"`
}

type WeatherData struct {
	Status  int `json:"status"`
	History struct {
		Dailysummary []struct {
			Fog          string `json:"fog"`
			Maxpressurem string `json:"maxpressurem"`
			Maxtempm     string `json:"maxtempm"`
			Maxwspdm     string `json:"maxwspdm"`
			Minpressurem string `json:"minpressurem"`
			Mintempm     string `json:"mintempm"`
			Minwspdm     string `json:"minwspdm"`
			Rain         string `json:"rain"`
			Tornado      string `json:"tornado"`
		} `json:"dailysummary"`
	} `json:"history"`
	Response struct {
		Version string `json:"version"`
	} `json:"response"`
}
type CustomFields struct {
	ID    int64  `json:"id"`
	Value string `json:"value"`
}

type TranformedData struct {
	Status          int             `json:"status,omitempty"`
	TicketDetails   TicketDetails   `json:"ticket_details,omitempty"`
	WeatherAPIInput WeatherAPIInput `json:"weather_api_input,omitempty"`
	TVClaimData     TVClaimData     `json:"tv_claim_data,omitempty"`
	StromClaimData  StromClaimData  `json:"storm_claim_data,omitempty"`
}

type TVClaimData struct {
	TVPrice         string `json:"tv_price,omitempty"`
	CrimeRef        string `json:"crime_ref,omitempty"`
	IncidentDate    string `json:"incident_date,omitempty"`
	TVModelNo       string `json:"tv_model_no,omitempty"`
	TVMake          string `json:"tv_make,omitempty"`
	TVSerialNo      string `json:"tv_serial_no,omitempty"`
	DamageImageURL1 string `json:"damage_image_url_1,omitempty"`
	DamageImageURL2 string `json:"damage_image_url_2,omitempty"`
	TVReceiptImage  string `json:"tv_reciept_image_url,omitempty"`
	DamageVideoURL  string `json:"damage_video_url,omitempty"`
}

type StromClaimData struct {
	IncidentPlace       string `json:"incident_place,omitempty"`
	IncidentDate        string `json:"incident_date,omitempty"`
	DamageImageURL1     string `json:"damage_image_url_1,omitempty"`
	DamageImageURL2     string `json:"damage_image_url_2,omitempty"`
	RepairEstimateImage string `json:"estimate_image_url,omitempty"`
	DamageVideoURL      string `json:"damage_video_url,omitempty"`
}

type TicketDetails struct {
	Status int `json:"status"`
	Ticket struct {
		Type     string `json:"type"`
		Subject  string `json:"subject"`
		Priority string `json:"priority"`
		Status   string `json:"status"`
		Comment  struct {
			HTMLBody string   `json:"html_body"`
			Uploads  []string `json:"uploads,omitempty"`
		} `json:"comment"`
		CustomFields []CustomFields `json:"custom_fields,omitempty"`
		Requester    struct {
			LocaleID     int    `json:"locale_id"`
			Name         string `json:"name"`
			Email        string `json:"email"`
			Phone        string `json:"phone"`
			PolicyNumber string `json:"policy_number"`
		} `json:"requester"`
		TicketFormID int64     `json:"ticket_form_id"`
		EventID      string    `json:"event_id"`
		Token        string    `json:"token"`
		SubmittedAt  time.Time `json:"submitted_at"`
	} `json:"ticket"`
}

type WeatherAPIInput struct {
	City    string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	Date    string `json:"date,omitempty"` //YYYYMMDD
}
