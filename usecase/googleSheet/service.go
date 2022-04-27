package googleSheet

import (
	entity "Atlan_Collect_Ingest/enitity"
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) AddToSheet(formId int8, oAuthCode string) ([]*entity.Responses, error) {
	resp, err := s.repo.Extract(formId)
	qIds, questions, err := s.repo.QuesIdExtract(formId)
	fmt.Printf("%v \n %v\n", qIds, questions)
	if err != nil {
		return nil, err
	}
	f := excelize.NewFile()

	for i, _ := range qIds {
		col := excelize.ToAlphaString(i)
		f.SetCellValue("Sheet1", col+strconv.Itoa(1), questions[i])
	}
	fmt.Printf("%v", len(resp))
	var counter = 0
	for i, _ := range resp {
		for _, e := range resp[i].Response {
			col := excelize.ToAlphaString(counter)
			f.SetCellValue("Sheet1", col+strconv.Itoa(i+2), e)
			counter++
		}
		counter = 0
	}

	if err = f.SaveAs("response.xlsx"); err != nil {
		return nil, err
	}

	err = UploadToDrive(oAuthCode)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//HELPER FUNCTION FOR G DRIVE

func getClient(ctx context.Context, config *oauth2.Config, code string) *http.Client {
	tok := getTokenFromWeb(config, code)
	return config.Client(ctx, tok)
}

func getTokenFromWeb(config *oauth2.Config, code string) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	//var code = "4/0AX4XfWgyl66VyJXM3LtQSpJX2qsV6ghvxeSiT3tBenLwvfIZMGU2icD-_XmiNp_6qXBJfw"
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func UploadToDrive(code string) error {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
		return err
	}
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v", err)
		return err
	}
	client := getClient(ctx, config, code)
	srv, err := drive.New(client)
	if err != nil {
		log.Printf("Unable to retrieve drive Client %v", err)
		return err
	}

	// Upload CSV and convert to Spreadsheet
	filename := "response.xlsx"                                                         // File you want to upload
	baseMimeType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" // mimeType of file you want to upload
	convertedMimeType := "application/vnd.google-apps.spreadsheet"                      // mimeType of file you want to convert on Google Drive

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	defer file.Close()
	f := &drive.File{
		Name:     filename,
		MimeType: convertedMimeType,
	}
	res, err := srv.Files.Create(f).Media(file, googleapi.ContentType(baseMimeType)).Do()
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	fmt.Printf("%s, %s, %s\n", res.Name, res.Id, res.MimeType)

	// Modify permissions
	permissiondata := &drive.Permission{
		Type:               "domain",
		Role:               "writer",
		Domain:             "google.com",
		AllowFileDiscovery: true,
	}
	pres, err := srv.Permissions.Create(res.Id, permissiondata).Do()
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	fmt.Printf("%s, %s\n", pres.Type, pres.Role)
	return nil
}
