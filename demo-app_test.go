package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server *httptest.Server
	//Test Data TV
	userJson1 = `{"media-bucket":{"status":200,"media":[{"bucket":"artifacts-image","name":"iTWFtJEMlX","size":153491,"media-link":"https://www.googleapis.com/download/storage/v1/b/artifacts-image/o/iTWFtJEMlX?generation=1510874906620336&alt=media","original-link":"https://admin.typeform.com/form/results/file/download/H8mm3s/63907299/30ecfd753d05-22405371_888332421323094_6861338905885136899_n.jpg"},{"bucket":"artifacts-image","name":"DWkkKfoMgP","size":153491,"media-link":"https://www.googleapis.com/download/storage/v1/b/artifacts-image/o/DWkkKfoMgP?generation=1510874909997537&alt=media","original-link":"https://admin.typeform.com/form/results/file/download/H8mm3s/j79cNctIvogK/ab33a3a7524b-22405371_888332421323094_6861338905885136899_n.jpg"},{"bucket":"artifacts-image","name":"INYLpyoFGu","size":153491,"media-link":"https://www.googleapis.com/download/storage/v1/b/artifacts-image/o/INYLpyoFGu?generation=1510874913508920&alt=media","original-link":"https://admin.typeform.com/form/results/file/download/H8mm3s/63907004/a870839fc865-22405371_888332421323094_6861338905885136899_n.jpg"}]},
	"tranformed-data":{"status":200,"ticket_details":{"ticket":{"event_id": "paRX3DyrSU","type":"incident","subject":"Storm surge claim","priority":"normal","status":"new","comment":{"html_body":"<p><b>If you have any other insurance or warranties covering your home, please advise us of the company name.</b> : sd</p><hr><p><b>In as much detail as possible, please use the text box below to describe the full extent of the damage to your home and how you discovered it.</b> : sdf</p><hr><p><b>Please describe the details of the condition of your home prior to discovering the damage</b> : sd</p><hr><p><b>If there has been any recent maintenance carried out on your home, please describe it</b> : sdf</p><hr><p><b>Would you like to upload more images?</b> : </p><hr><p><b>If you have been provided with a repair estimate from a contractor or tradesman, you can upload this by providing a clear photo of the document or upload an existing file.</b> <a href='https://admin.typeform.com/form/results/file/download/H8mm3s/wFpTHm7AZVNO/913579797420-22405371_888332421323094_6861338905885136899_n.jpg'>https://admin.typeform.com/form/results/file/download/H8mm3s/wFpTHm7AZVNO/913579797420-22405371_888332421323094_6861338905885136899_n.jpg</a></p><hr><p><b>Where did the incident happen? (City/town name)</b> : London</p><hr><p><b>Would you like to upload more images?</b> : </p><hr><p><b>Are you aware of anything else relevant to your claim that you would like to advise us of at this stage?</b> : sadsa</p><hr><p><b>If it is safe and possible to do so, please provide images of the damage to both the outside and the inside of your home.</b> <a href='https://admin.typeform.com/form/results/file/download/H8mm3s/63907299/30ecfd753d05-22405371_888332421323094_6861338905885136899_n.jpg'>https://admin.typeform.com/form/results/file/download/H8mm3s/63907299/30ecfd753d05-22405371_888332421323094_6861338905885136899_n.jpg</a></p><hr><p><b>Are you still have possession of the damage items (i.e. damaged guttering)?</b> : </p><hr><p><b>When did the incident happen?</b> : 2017-10-12</p><hr><p><b>If it is safe and possible to do so, please provide images of the damage to both the outside and the inside of your home.</b> <a href='https://admin.typeform.com/form/results/file/download/H8mm3s/j79cNctIvogK/ab33a3a7524b-22405371_888332421323094_6861338905885136899_n.jpg'>https://admin.typeform.com/form/results/file/download/H8mm3s/j79cNctIvogK/ab33a3a7524b-22405371_888332421323094_6861338905885136899_n.jpg</a></p><hr><p><b>If it is safe and possible to do so, please take a short video to include any areas of damage. Hold the camera in landscape orientation and include a voice narration if you can to help explain the situation.</b> <a href='https://admin.typeform.com/form/results/file/download/H8mm3s/63907004/a870839fc865-22405371_888332421323094_6861338905885136899_n.jpg'>https://admin.typeform.com/form/results/file/download/H8mm3s/63907004/a870839fc865-22405371_888332421323094_6861338905885136899_n.jpg</a></p><hr><p><b>We have made the following assumptions about your property, you and anyone living with you</b> : </p><hr>"},"requester":{"locale_id":1,"name":"Amit","email":"amitkumarvarman@gmail.com","phone":"99999999999","policy_number":"DUSSS2323232"}}},"weather_api_input":{"city":"London","date":"20171012"},"tv_claim_data":{"tv_reciept_image_url":""},"storm_claim_data":{"incident_place":"London","incident_date":"2017-10-12","damage_image_url_1":"https://admin.typeform.com/form/results/file/download/H8mm3s/63907299/30ecfd753d05-22405371_888332421323094_6861338905885136899_n.jpg","damage_image_url_2":"https://admin.typeform.com/form/results/file/download/H8mm3s/j79cNctIvogK/ab33a3a7524b-22405371_888332421323094_6861338905885136899_n.jpg","estimate_image_url":"https://admin.typeform.com/form/results/file/download/H8mm3s/wFpTHm7AZVNO/913579797420-22405371_888332421323094_6861338905885136899_n.jpg","damage_video_url":"https://admin.typeform.com/form/results/file/download/H8mm3s/63907004/a870839fc865-22405371_888332421323094_6861338905885136899_n.jpg"}},
			"weather-data":{"status":200,"history":{"dailysummary":[{"fog":"0","maxpressurem":"1025","maxtempm":"7","maxwspdm":"28","minpressurem":"1014","mintempm":"0","minwspdm":"7","rain":"1","tornado":"0"}]},"response":{"version":"0.1"}},
			"weather-risk":{"status":200,"description":"Possible Stormy weather","riskScore":50}}`
	userJson2 = `{"media-bucket":{"status":200,"media":[{"bucket":"artifacts-image","name":"iTWFtJEMlX","size":153491,"media-link":"https://www.googleapis.com/download/storage/v1/b/artifacts-image/o/iTWFtJEMlX?generation=1510874906620336&alt=media","original-link":"https://admin.typeform.com/form/results/file/download/H8mm3s/63907299/30ecfd753d05-22405371_888332421323094_6861338905885136899_n.jpg"},{"bucket":"artifacts-image","name":"DWkkKfoMgP","size":153491,"media-link":"https://www.googleapis.com/download/storage/v1/b/artifacts-image/o/DWkkKfoMgP?generation=1510874909997537&alt=media","original-link":"https://admin.typeform.com/form/results/file/download/H8mm3s/j79cNctIvogK/ab33a3a7524b-22405371_888332421323094_6861338905885136899_n.jpg"},{"bucket":"artifacts-image","name":"INYLpyoFGu","size":153491,"media-link":"https://www.googleapis.com/download/storage/v1/b/artifacts-image/o/INYLpyoFGu?generation=1510874913508920&alt=media","original-link":"https://admin.typeform.com/form/results/file/download/H8mm3s/63907004/a870839fc865-22405371_888332421323094_6861338905885136899_n.jpg"}]},"tranformed-data":{"status":200,"ticket_details":{"ticket":{"comment":{"html_body":"\u003cp\u003e\u003cb\u003eIf there has been any recent maintenance carried out on your home, please describe it\u003c/b\u003e : No maintenance carried out\u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eIf you have any other insurance or warranties covering your home, please advise us of the company name.\u003c/b\u003e : No\u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eWe have made the following assumptions about your property, you and anyone living with you\u003c/b\u003e : \u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eWhen did the incident happen?\u003c/b\u003e : 2017-01-01\u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eAre you still have possession of the damage items (i.e. damaged guttering)?\u003c/b\u003e : \u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eAre you aware of anything else relevant to your claim that you would like to advise us of at this stage?\u003c/b\u003e : I would need the vendors contact for repairing the roof\u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eWould you like to upload more images?\u003c/b\u003e : \u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eWhere did the incident happen? (City/town name)\u003c/b\u003e : birmingham\u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eIn as much detail as possible, please use the text box below to describe the full extent of the damage to your home and how you discovered it.\u003c/b\u003e : Roof Damaged\u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003ePlease describe the details of the condition of your home prior to discovering the damage\u003c/b\u003e : Tiles blown away\u003c/p\u003e\u003chr\u003e"},"custom_fields":null,"email":"amitkumarvarman@gmail.com","phone":"09876512345","priority":"normal","requester":{"email":"amitkumarvarman@gmail.com","locale_id":1,"name":"Amit Varman"},"status":"new","subject":"Storm surge risk data","type":"incident"}},"weather-api-input":{"city":"birmingham","country":"","date":"20170101"}},"weather-data":{"status":200,"history":{"dailysummary":[{"fog":"0","maxpressurem":"1025","maxtempm":"7","maxwspdm":"28","minpressurem":"1014","mintempm":"0","minwspdm":"7","rain":"1","tornado":"0"}]},"response":{"version":"0.1"}},"weatherRisk":{"description":"Possible Stormy weather","riskScore":50}}`
	// ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
)

func TestHandler(t *testing.T) {
	//Convert string to reader and
	//Create request with JSON body
	req1, err := http.NewRequest("POST", "", strings.NewReader(userJson1))
	req2, err := http.NewRequest("POST", "", strings.NewReader(userJson2))
	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{"Test Data -1", args{rr, req1}},
		{"Test Data -2", args{rr, req2}},
	}
	for _, tt := range tests {
		// call ServeHTTP method
		// directly and pass Request and ResponseRecorder.
		handler := http.HandlerFunc(Handler)
		handler.ServeHTTP(tt.args.w, tt.args.r)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		//check content type
		if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
			t.Errorf("content type header does not match: got %v want %v",
				ctype, "application/json")
		}
	}
}
