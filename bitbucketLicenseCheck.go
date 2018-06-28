package main

import (
	"log"
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strconv"
	"crypto/tls"
	"github.com/spf13/viper"
	"path/filepath"
	"os"
)

type StatusJSON struct {
	CurrentNumberOfUsers int `json: currentNumberOfUsers`
	ServerId string `json: serverId`
}

type LicenseJSON struct {
	CreationDate int `json: creationDate`
	PurchaseDate int `json: purchaseDate`
	ExpiryDate int `json: expiryDate`
	MaintenanceExpiryDate int `json: maintenanceExpiryDate`
	NumberOfDaysBeforeMaintenanceExpiry int `json: NumberOfDaysBeforeMaintenanceExpiry`
	GracePeriodEndDate int `json: gracePeriodEndDate`
	NumberofDaysBeforeGracePeriodExpiry int `json: numberOfDaysBeforeGracePeriodExpiry`
	MaximumNumberOfUsers int `json: maximumNumberOfUsers`
    UnlimitedNumberOfUsers bool `json: unlimitedNumberOfUsers`
	ServerId string `json: serverId`
	SupportEntitlementNumber string `json: supportEntitlementNumber`
	License string `json: license`
	Status StatusJSON `json: status`
}

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	BaseUrl string `json:"baseurl"`
}

func main() {

	//Read in Config file
	var appConf Config
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil{
		log.Fatalf("Error getting working directory info. Error: %v", err)
	}

	viper.SetConfigName("settings")
	viper.AddConfigPath(currentDir)
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading 'settings.conf' located in %s. \nError: %v", currentDir, err)
	}

	err = viper.Unmarshal(&appConf)
	if err != nil {
		log.Fatalf("Error parsing data from 'settings.conf. \nError: %v", err)
	}

	//URL to get the license info
	licenseUrl := appConf.BaseUrl + "admin/license"

	//Disable security check of the http certificate
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	//Create http client and set up request
	client := http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest(http.MethodGet, licenseUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(appConf.Username, appConf.Password)

	//Perform the GET
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error getting response from url. \nError: %v", err)
	}

	//get body of the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body from url. \nError: %v", err)
	}

	//Parse the json body
	licenseInfo := LicenseJSON{}
	err = json.Unmarshal(body, &licenseInfo)
	if err != nil {
		log.Fatalf("Error reading json. Error: %v", err)
	}

	//convert int to string and print out in a format for PRTG
	numLicensesString := strconv.Itoa(licenseInfo.Status.CurrentNumberOfUsers)
	fmt.Println("<prtg> <result> <channel>Bitbucket License Count</channel> <value>" + numLicensesString + "</value> </result> </prtg>")

}