package temp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
	_ "github.com/denisenkom/go-mssqldb"

	// XpertCore from XpertSofia
	"XpertRestApi-V8/api/customer"
	"XpertRestApi-V8/api/events"

	"XpertRestApi-V8/api/metaData"
	"XpertRestApi-V8/api/settings"
)

// XpertMethodStatus represents the status of a method execution
type XpertMethodStatus int

const (
	NotSuccessful XpertMethodStatus = iota
	Successful
	// ReadConfiguration reads the configuration settings for the processor
	SETTING_TEMP_CHECK_TIMER       = "TempCheckTimer"
	SETTING_WEBSOCKET_TOPIC_NAME   = "WebSocketTopic"
	SETTING_PROXIMITY_SECONDS      = "ProximitySeconds" // Added missing constant
	SETTING_STATUS_SECONDS         = "StatusSeconds"    // Added missing constant
	DEFAULT_PROXIMITY_SECONDS      = 3600               // Default value in seconds
	DEFAULT_STATUS_SECONDS         = 3600               // Default value in seconds
	
)

// Placeholder for XpertMQTTClientProducer
type XpertMQTTClientProducer struct{}

// Placeholder for XpertEmailLibrary
type XpertEmailLibrary struct{}

// SendEmailAdvanced sends an email with the given parameters and returns a success status
func (lib *XpertEmailLibrary) SendEmailAdvanced(to, subject, body string, uris []string, contentType, priority string) bool {
	// Placeholder implementation for sending email
	log.Printf("Sending email to: %s, Subject: %s, Body: %s, URIs: %v, ContentType: %s, Priority: %s",
		to, subject, body, uris, contentType, priority)
	return true 
}

//added as a placeholder
type XpertUserModel struct {
	Email string
}

// SendEmail sends an email to a list of users and logs the process
func (processor *XpertTempCheckProcessor) SendEmail(users []XpertUserModel, details, description string, eventID int, deviceMac string, oDebugMsg *XpertDebugMessageJsonObject) XpertMethodStatus {
	const METHOD_NAME = "SendEmail"
	log.Printf("%s XpertTempCheckProcessor <> %s UsersCnt: %d", deviceMac, METHOD_NAME, len(users))

	// Add debug property for user count
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error in %s: %v", METHOD_NAME, r)
		}
	}()
	oDebugMsg.AddProperty("User Count", fmt.Sprintf("%d", len(users)))

	result := NotSuccessful

	// Validate users
	if len(users) <= 0 {
		log.Printf("Error: Invalid parameter value in %s: emailTo is null or empty", METHOD_NAME)
		return result
	}

	// Initialize email library
	oEmailLibrary := &XpertEmailLibrary{}

	// Collect email addresses
	var emails []string
	for _, user := range users {
		if user.Email != "" {
			emails = append(emails, user.Email)
		}
	}

	oDebugMsg.AddProperty("Email List", strings.Join(emails, ","))

	sendEmailResult := oEmailLibrary.SendEmailAdvanced(
		strings.Join(emails, ","),
		details,
		description,
		[]string{}, // URIs
		"HTML",     // Message content type
		"Normal",   // Message priority
	)

	if !sendEmailResult {
		log.Printf("%s XpertTempCheckProcessor <> %s Email Sent Failed", deviceMac, METHOD_NAME)
		oDebugMsg.AddProperty("Email Status", "Not Sent")
		log.Printf("Error: Failed to send email in %s", METHOD_NAME)
		return result
	}

	// Add debug property for email status
	oDebugMsg.AddProperty("Email Status", "Sent")

	// Create item event action
	createItemEventActionResult := processor.CreateItemEventAction(eventID, "Email Sent", fmt.Sprintf("Email sent to: %d", len(emails)))
	if createItemEventActionResult != Successful {
		oDebugMsg.AddProperty("Action Status", "Not Sent")
		log.Printf("Error: Failed to create item event action in %s", METHOD_NAME)
		return result
	}


	oDebugMsg.AddProperty("Action Status", "Sent")

	result = Successful
	log.Printf("Exiting method: %s", METHOD_NAME)
	return result
}

// XpertDebugMessageJsonObject represents a debug message object with properties, placeholder
type XpertDebugMessageJsonObject struct {
	Properties []DebugProperty
}

// DebugProperty represents a single debug property, placeholder
type DebugProperty struct {
	IsList bool
	Name   string
	Value  string
}

// AddProperty adds a debug property to the debug message object
func (oDebugMsg *XpertDebugMessageJsonObject) AddProperty(name, value string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error adding property %s: %v", name, r)
		}
	}()
	oDebugMsg.Properties = append(oDebugMsg.Properties, DebugProperty{
		IsList: false,
		Name:   name,
		Value:  value,
	})
}

// XpertTempCheckProcessor represents the temperature check processor
type XpertTempCheckProcessor struct {
	HealthStatus     bool
	IsTestMode       bool
	ProximitySeconds int
	TempCheckTimer   int
	StatusSeconds    int
	WebSocketTopic   string
	MProducer        *XpertMQTTClientProducer
}

// NewXpertTempCheckProcessor initializes a new instance of XpertTempCheckProcessor
func NewXpertTempCheckProcessor() *XpertTempCheckProcessor {
	return &XpertTempCheckProcessor{
		HealthStatus:     true,
		IsTestMode:       false,
		ProximitySeconds: 3600, // Default value
		TempCheckTimer:   1000, // Default value in milliseconds
		StatusSeconds:    3600, // Default value
		WebSocketTopic:   "",
		MProducer:        nil,
	}
}

// CheckTempCheckpoints checks temperature checkpoints using the metaData API
func (processor *XpertTempCheckProcessor) CheckTempCheckpoints() {
	offset := int(time.Now().UTC().Hour() - time.Now().Hour()) // Corrected offset calculation
	oCustomers := customer.GetAllCustomers()

	for _, customer := range oCustomers.List {
		oGroups := metaData.GetGroupsResponsibleFor(0, customer.ID)
		for _, group := range oGroups.List {
			if group.CheckTimeFrames == "" {
				continue
			}
			groupTimeFrames := strings.Split(group.CheckTimeFrames, ";")
			oTempChecks := metaData.GetTemperatureCheckpoints(0, group.ID, 1, 100, time.Now().AddDate(0, 0, -1), time.Now(), customer.ID)
			for _, timeframe := range groupTimeFrames {
				timeframeParts := strings.Split(timeframe, ",")
				startHour, _ := strconv.Atoi(strings.Split(timeframeParts[0], ":")[0])
				startMinute, _ := strconv.Atoi(strings.Split(timeframeParts[0], ":")[1])
				endHour, _ := strconv.Atoi(strings.Split(timeframeParts[1], ":")[0])
				endMinute, _ := strconv.Atoi(strings.Split(timeframeParts[1], ":")[1])
				startDateTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), startHour-offset, startMinute, 0, 0, time.UTC)
				endDateTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), endHour-offset, endMinute, 0, 0, time.UTC)

				checkPerformed := false
				for _, tempCheck := range oTempChecks.List {
					if tempCheck.DateCreated.After(startDateTime) && tempCheck.DateCreated.Before(endDateTime) {
						checkPerformed = true
						break
					}
				}

				if checkPerformed {
					for _, tempCheck := range oTempChecks.List {
						if tempCheck.DateCreated.After(startDateTime) && tempCheck.DateCreated.Before(endDateTime) && tempCheck.Compliance == "" {
							metaData.UpdateTemperatureCheckpoint(group.ID, "check", customer.ID)
							break
						}
					}
				} else if !checkPerformed && time.Now().After(endDateTime) {
					var latestTempCheck *metaData.TemperatureCheckpoint
					for _, tempCheck := range oTempChecks.List {
						if latestTempCheck == nil || tempCheck.DateCreated.After(latestTempCheck.DateCreated) {
							latestTempCheck = &tempCheck
						}
					}
					if latestTempCheck == nil || latestTempCheck.Compliance != "miss" {
						log.Printf("No check found for time range, marking as missed")
						metaData.CreateTemperatureCheckpoint(0, group.ID, false, "", customer.ID)
						metaData.UpdateTemperatureCheckpoint(group.ID, "miss", customer.ID)
					}
				}
			}
		}
	}
}

// InfraCounts represents the infrastructure ID and its visit count
type InfraCounts struct {
	InfraID int
	Count   int
}

// XpertUseCaseModel represents a use case with its name and ID, placehplder
type XpertUseCaseModel struct {
	UseCase   string
	UseCaseID int
}

// CheckAllInfrasVisited checks if all infrastructures in a route have been visited the required number of times
func CheckAllInfrasVisited(routeDef map[string]interface{}, numVisits int, allVisits []events.XpertEventModel) ([]InfraCounts, []InfraCounts) {
	const METHOD_NAME = "CheckAllInfrasVisited"
	log.Printf("Entering method: %s", METHOD_NAME)
	

	// Initialize unseenInfras and counts
	unseenInfras := []InfraCounts{}
	counts := make(map[int]int)
	infrasInRoute := []int{}

	// Try-catch equivalent in Go using defer and recover
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error in %s: %v", METHOD_NAME, r)
		}
	}()

	// Populate infrasInRoute and counts from routeDef
	if infrastructures, ok := routeDef["infrastructures"].([]interface{}); ok {
		for _, infra := range infrastructures {
			if infraMap, ok := infra.(map[string]interface{}); ok {
				if id, ok := infraMap["Id"].(float64); ok { // Assuming Id is a float64 in the dynamic structure
					infraID := int(id)
					infrasInRoute = append(infrasInRoute, infraID)
					counts[infraID] = 0
				}
			}
		}
	}

	// Count visits for each infrastructure
	for _, visit := range allVisits {
		if _, exists := counts[visit.PlanID]; exists {
			counts[visit.PlanID]++
		} else {
			counts[visit.PlanID] = 1
		}
	}

	// Check if each infrastructure in the route has been visited the required number of times
	for _, infra := range infrasInRoute {
		if counts[infra] < numVisits {
			unseenInfras = append(unseenInfras, InfraCounts{
				InfraID: infra,
				Count:   counts[infra],
			})
		}
	}

	log.Printf("Exiting method: %s", METHOD_NAME)
	return unseenInfras, unseenInfras
}

// CheckTempGraces checks for temperature grace periods and handles alerts
func (processor *XpertTempCheckProcessor) CheckTempGraces() {
	const METHOD_NAME = "CheckTempGraces"
	log.Printf("Entering method: %s", METHOD_NAME)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error in %s: %v", METHOD_NAME, r)
		}
	}()

	gracePeriodUpdateTime := time.Time{}
	lastAlertUpdateTime := time.Time{}

	oStaffService := &customer.XpertStaffService{}
	oEventService := &events.XpertEventService{}
	oDeviceService := &metaData.XpertDeviceService{}
	oSettingsService := &settings.XpertSettingsService{}

	oEventModel := events.XpertEventModel{
		SystemName: "ERROR - must be re-assigned.",
		Name:       "TemperatureMonitoring",
	}

	oEvents := oEventService.GetOpenTempEvents("TEMPERATURE_GRACE")

	for _, tempEvent := range oEvents.List {
		if tempEvent.DeviceID == 0 || tempEvent.CustomerID == 0 {
			continue
		}

		oDevice := oDeviceService.GetDevice(tempEvent.CustomerID, tempEvent.DeviceID)
		if strings.ToLower(oDevice.ModelName) != "ts1" && strings.ToLower(oDevice.ModelName) != "ts2" && strings.ToLower(oDevice.ModelName) == "hs1" {
			log.Printf("DEVICE TYPE %s NOT HANDLED BY THIS PROCESSOR", oDevice.ModelName)
			continue
		}

		gracePeriodUpdateTime = tempEvent.DateUpdated
		oEventModel.DeviceID = tempEvent.DeviceID
		oEventModel.CustomerID = tempEvent.CustomerID
		oEventModel.DeviceUniqueID = oDevice.UniqueID

		var oConfig map[string]interface{}
		if oDevice.PendingConfigID > 0 {
			oConfig = oDeviceService.GetConfiguration(oDevice.CustomerID, oDevice.PendingConfigID).ConfigDef
		} else if oDevice.ConfigID > 0 {
			oConfig = oDeviceService.GetConfiguration(oDevice.CustomerID, oDevice.ConfigID).ConfigDef
		} else {
			continue
		}

		oSettings := oSettingsService.GetNotificationSettingsForSDCT(tempEvent.CustomerID, "Temp")
		oSettingsJson := oSettings.SettingJson

		if time.Now().UTC().After(gracePeriodUpdateTime.Add(time.Second * time.Duration(oConfig["hightime1"].(float64)))) {
			oLastAlertEvent := oEventService.GetEventBySystemName(0, tempEvent.ItemID, "TEMPERATURE_ALERT", false)
			lastAlertUpdateTime = oLastAlertEvent.DateUpdated

			renotifyPeriod := 999999
			if repeatInterval, ok := oSettingsJson["RepeatIntervalTimeValue"].(float64); ok {
				renotifyPeriod = int(repeatInterval)
				if strings.ToLower(oSettingsJson["RepeatIntervalTimeUnit"].(string)) == "minutes" {
					renotifyPeriod *= 60
				} else if strings.ToLower(oSettingsJson["RepeatIntervalTimeUnit"].(string)) == "hours" {
					renotifyPeriod *= 3600
				}
			}

			if time.Now().UTC().After(lastAlertUpdateTime.Add(time.Second * time.Duration(renotifyPeriod))) {
				oItem := oStaffService.GetStaffPersonByID(tempEvent.CustomerID, oDevice.ItemID)
				processor.CreateTempEvent(oItem, time.Now().UTC(), time.Now().UTC(), 1, &oEventModel, oEventService, "TEMPERATURE_ALERT",
					oConfig["lowvalue1"].(float64), oConfig["highvalue1"].(float64), tempEvent.ViolationValue, "Temperature Range Exceeded.", "0")
			}
		}
	}

	log.Printf("Exiting method: %s", METHOD_NAME)
}

// CreateItemEventAction creates an item event action and returns the method status
func (processor *XpertTempCheckProcessor) CreateItemEventAction(itemEventID int, actionType, actionDetails string) XpertMethodStatus {
	const METHOD_NAME = "CreateItemEventAction"
	log.Printf("Entering method: %s", METHOD_NAME)

	// Default result is NotSuccessful
	result := NotSuccessful

	// Validate parameters
	if itemEventID <= 0 {
		log.Printf("Error: Invalid parameter value in %s: itemEventID is null or invalid", METHOD_NAME)
		return result
	}

	if strings.TrimSpace(actionType) == "" {
		log.Printf("Error: Invalid parameter value in %s: actionType is null or empty", METHOD_NAME)
		return result
	}

	if strings.TrimSpace(actionDetails) == "" {
		log.Printf("Error: Invalid parameter value in %s: actionDetails is null or empty", METHOD_NAME)
		return result
	}

	// Create the event action model
	eventTime := time.Now().UTC()
	oModel := events.XpertEventActionModel{
		ItemEventID:    itemEventID,
		ActionTypeID:   0,
		ActionType:     actionType,
		ActionDateTime: eventTime,
		Description:    actionDetails,
		ActionUserID:   0,
		DateCreated:    eventTime,
		DateUpdated:    eventTime,
	}

	// Insert the event action using the events API
	oEventService := &events.XpertEventService{}
	oResult := oEventService.InsertEventAction(oModel)

	// Check for errors in the result
	if oResult.HasError {
		log.Printf("Error: Failed to insert event action in %s", METHOD_NAME)
		return result
	}

	// If successful, update the result
	result = Successful

	log.Printf("Exiting method: %s", METHOD_NAME)
	return result
}

// CreateTempEvent creates a new temperature event using the events API
func (processor *XpertTempCheckProcessor) CreateTempEvent(
	oItem *settings.StaffPerson,
	startDate, endDate time.Time,
	duration float64,
	oEventModel *events.XpertEventModel,
	oXpertEventService *events.XpertEventService,
	systemName string,
	minValue, maxValue float64,
	violationValue, description, logPeriod string,
) *events.ResultObject {
	const METHOD_NAME = "CreateTempEvent"
	log.Printf("Entering method: %s", METHOD_NAME)

	result := &events.ResultObject{}

	if !oItem.EnableAlerts {
		result.ErrorMessage = "Item does not have events enabled, alert generation cancelled"
		return result
	}

	oEventModel.StartDateTime = startDate
	oEventModel.MinValue = minValue
	oEventModel.MaxValue = maxValue
	oEventModel.ViolationValue = violationValue
	oEventModel.Description = description
	oEventModel.SystemName = systemName
	oEventModel.Name = description
	oEventModel.DisplayName = description
	oEventModel.RuleName = logPeriod
	oEventModel.EndDateTime = endDate
	oEventModel.DateUpdated = time.Now().UTC()
	oEventModel.DateCreated = time.Now().UTC()
	oEventModel.ItemID = oItem.ID
	oEventModel.ClosedDateTime = oEventModel.EndDateTime
	oEventModel.AllowedValueRange = fmt.Sprintf("%f", duration)
	oEventModel.CustomerID = oItem.CustomerID
	oEventModel.DateTimeToBeArchived = oEventModel.DateCreated.AddDate(0, 0, 60)
	oEventModel.UseCase = 6

	if err := oXpertEventService.InsertEvent(oEventModel); err != nil {
		log.Printf("Error inserting event in %s: %v", METHOD_NAME, err)
	}

	log.Printf("Exiting method: %s", METHOD_NAME)
	return result
}

// CloseEvent closes an event using the events API
func (processor *XpertTempCheckProcessor) CloseEvent(eventID int, oEventService *events.XpertEventService, customerID int) *events.ResultObject {
	const METHOD_NAME = "CloseEvent"
	log.Printf("Entering method: %s", METHOD_NAME)

	// Call the CloseEvent method from the events API
	result := oEventService.CloseEvent(customerID, eventID)

	log.Printf("Exiting method: %s", METHOD_NAME)
	return result
}

// CreateEvent creates a new event using the events API
func (processor *XpertTempCheckProcessor) CreateEvent(deviceMac string, oItem *settings.StaffPerson,
	oProximityInfrastructure *metaData.Infrastructure, routeID int, startDate, endDate time.Time,
	duration float64, useCase, zoneID int, oEventModel *events.Event, details string) *events.ResultObject {

	if !oItem.EnableAlerts {
		return &events.ResultObject{ErrorMessage: "Item does not have events enabled, alert generation cancelled"} // Ensure ResultObject is defined in the events package
	}

	oEventModel.Details = details
	if strings.Contains(strings.ToLower(oEventModel.SystemName), "alert") || strings.Contains(strings.ToLower(oEventModel.SystemName), "warning") {
		oEventModel.Description = details
	} else {
		oEventModel.Description = fmt.Sprintf("%s %s %s", oItem.Name, deviceMac, oProximityInfrastructure.Name)
	}
	oEventModel.MinValue = routeID
	oEventModel.MaxValue = zoneID
	oEventModel.StartDateTime = startDate
	oEventModel.EndDateTime = endDate
	oEventModel.DateUpdated = time.Now().UTC()
	oEventModel.DateCreated = time.Now().UTC()
	oEventModel.ClosedDateTime = oEventModel.EndDateTime
	oEventModel.AllowedValueRange = fmt.Sprintf("%f", duration)
	oEventModel.PlanID = oProximityInfrastructure.ID
	oEventModel.DeviceID = oItem.DeviceID
	oEventModel.DateTimeToBeArchived = oEventModel.DateCreated.AddDate(0, 0, 60)
	oEventModel.UseCase = useCase

	return events.Insert(oEventModel)
}

// UpdateEvent updates an existing event using the events API
func (processor *XpertTempCheckProcessor) UpdateEvent(oEvent *events.XpertEventModel, oEventService *events.XpertEventService) *events.ResultObject {
	const METHOD_NAME = "UpdateEvent"
	log.Printf("Entering method: %s", METHOD_NAME)

	// Call the UpdateEvent method from the events API
	result := oEventService.UpdateEvent(oEvent)

	log.Printf("Exiting method: %s", METHOD_NAME)
	return result
}

// CheckPendingTempConfigs checks for pending temperature configurations and processes them
func (processor *XpertTempCheckProcessor) CheckPendingTempConfigs() {
	const METHOD_NAME = "CheckPendingTempConfigs"
	log.Printf("Entering method: %s", METHOD_NAME)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error in %s: %v", METHOD_NAME, r)
		}
	}()

	oDeviceService := &metaData.XpertDeviceService{}
	oCustService := &customer.XpertCustomerService{}

	// Get devices with pending temperature configurations
	oDevices := oDeviceService.GetDevicesByPendingTempConfig(1)

	for _, oDevice := range oDevices.List {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Error processing device %d in %s: %v", oDevice.ID, METHOD_NAME, r)
				}
			}()

			var ipAddress, username, password string
			cmdMatch := true
			var oIntegration *customer.XpertIntegrationModel

			// Get integrations for the customer
			oIntegrations := oCustService.GetIntegrations(oDevice.CustomerID)
			for _, custInt := range oIntegrations.List {
				if strings.TrimSpace(strings.ToLower(custInt.Type)) == "arc" {
					oIntegration = &custInt
					break
				}
			}

			if oIntegration == nil {
				log.Printf("No ARC integration found for device %d", oDevice.ID)
				return
			}

			// Parse integration JSON
			var oArcMessage map[string]interface{}
			if err := json.Unmarshal([]byte(oIntegration.JSON), &oArcMessage); err != nil {
				log.Printf("Error parsing integration JSON for device %d: %v", oDevice.ID, err)
				return
			}

			// Extract IP address, username, and password
			if ip, ok := oArcMessage["ErcIPAddress"].(string); ok && ip != "" {
				ipAddress = ip
			} else if ip, ok := oArcMessage["IPAddress"].(string); ok {
				ipAddress = ip
			}
			if user, ok := oArcMessage["Username"].(string); ok {
				username = user
			}
			if pass, ok := oArcMessage["Password"].(string); ok {
				password = pass
			}

			// Check if the associated configuration exists
			configDef, err := processor.fetchConfig(ipAddress, username, password, oDevice.IntegrationConfigID, "tagconfiglist")
			if err != nil {
				log.Printf("Error fetching configuration for device %d: %v", oDevice.ID, err)
				return
			}

			// Parse configuration XML to JSON
			configJSON, err := processor.parseXMLToJSON(configDef)
			if err != nil {
				log.Printf("Error parsing configuration XML for device %d: %v", oDevice.ID, err)
				return
			}

			// Check if the associated tag commands exist
			tagCommands, err := processor.fetchConfig(ipAddress, username, password, oDevice.UniqueID, "tagconfigdump")
			if err != nil {
				log.Printf("Error fetching tag commands for device %d: %v", oDevice.ID, err)
				return
			}

			// Parse tag commands XML to JSON
			commandJSON, err := processor.parseXMLToJSON(tagCommands)
			if err != nil {
				log.Printf("Error parsing tag commands XML for device %d: %v", oDevice.ID, err)
				return
			}

			// Compare commands and configurations
			for _, cmd := range commandJSON["response"].(map[string]interface{})["TAGCONFIG"].(map[string]interface{})["item"].([]interface{}) {
				found := false
				for _, config := range configJSON["response"].(map[string]interface{})["CONFIG"].(map[string]interface{})["cmd"].([]interface{}) {
					if processor.compareCommands(cmd.(string), config.(string)) {
						found = true
						break
					}
				}
				if !found {
					cmdMatch = false
					break
				}
			}

			// If commands match, set device configurations
			if cmdMatch {
				oDeviceService.SetDeviceConfigs([]metaData.XpertDeviceModel{oDevice}, oDevice.PendingConfigID, 0, oDevice.CustomerID)
			}
		}()
	}

	log.Printf("Exiting method: %s", METHOD_NAME)
}

// GetCustomerUseCases retrieves a list of use cases based on customer applications
func GetCustomerUseCases(customerApps map[string]bool) []XpertUseCaseModel {
	const METHOD_NAME = "GetCustomerUseCases"
	log.Printf("Entering method: %s", METHOD_NAME)

	useCases := []XpertUseCaseModel{}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error in %s: %v", METHOD_NAME, r)
		}
	}()

	if customerApps["AssetTracking"] {
		useCases = append(useCases, XpertUseCaseModel{UseCase: "AssetTracking", UseCaseID: 1})
	}
	if customerApps["MELT"] {
		useCases = append(useCases, XpertUseCaseModel{UseCase: "MELT", UseCaseID: 2})
	}
	if customerApps["PatientFlow"] {
		useCases = append(useCases, XpertUseCaseModel{UseCase: "PatientFlow", UseCaseID: 3})
	}
	if customerApps["SDCT"] {
		useCases = append(useCases, XpertUseCaseModel{UseCase: "SDCT", UseCaseID: 4})
	}
	if customerApps["StaffSafety"] {
		useCases = append(useCases, XpertUseCaseModel{UseCase: "StaffSafety", UseCaseID: 5})
	}

	log.Printf("Exiting method: %s", METHOD_NAME)
	return useCases
}

// IsHealthy checks the health status of the processor
func (processor *XpertTempCheckProcessor) IsHealthy(diagnosticData string) bool {
	return processor.HealthStatus
}


func (processor *XpertTempCheckProcessor) ReadConfiguration() {
	const METHOD_NAME = "ReadConfiguration"
	log.Printf("Entering method: %s", METHOD_NAME)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error in %s: %v", METHOD_NAME, r)
		}
	}()

	oConfigReader := NewXpertConfigurationReader()
	if err := oConfigReader.OpenConfigurationFile(); err != nil {
		log.Printf("Error opening configuration file: %v", err)
		return
	}

	// Read configuration parameters
	if tempCheckTimer, err := strconv.Atoi(oConfigReader.GetConfigurationParameter(SETTING_TEMP_CHECK_TIMER)); err == nil {
		processor.TempCheckTimer = tempCheckTimer
	} else {
		log.Printf("Error reading TempCheckTimer: %v", err)
	}

	processor.WebSocketTopic = oConfigReader.GetConfigurationParameter(SETTING_WEBSOCKET_TOPIC_NAME)

	if proximitySeconds, err := strconv.Atoi(oConfigReader.GetConfigurationParameter(SETTING_PROXIMITY_SECONDS)); err == nil {
		processor.ProximitySeconds = proximitySeconds
	} else {
		processor.ProximitySeconds = DEFAULT_PROXIMITY_SECONDS
	}

	if statusSeconds, err := strconv.Atoi(oConfigReader.GetConfigurationParameter(SETTING_STATUS_SECONDS)); err == nil {
		processor.StatusSeconds = statusSeconds
	} else {
		processor.StatusSeconds = DEFAULT_STATUS_SECONDS
	}

	log.Printf("Exiting method: %s", METHOD_NAME)
}

// fetchConfig fetches configuration or tag commands from the specified URI
func (processor *XpertTempCheckProcessor) fetchConfig(ipAddress, username, password string, id interface{}, endpoint string) (string, error) {
	uri := fmt.Sprintf("%s/epe/cfg/%s?configid=%v", ipAddress, endpoint, id)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// parseXMLToJSON parses an XML string into a JSON object
func (processor *XpertTempCheckProcessor) parseXMLToJSON(xmlData string) (map[string]interface{}, error) {
	if xmlData == "" {
		return nil, fmt.Errorf("xmlData is empty")
	}

	var jsonData map[string]interface{}
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlData); err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(doc.Root())
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}

// compareCommands compares two commands for equivalence
func (processor *XpertTempCheckProcessor) compareCommands(cmd, config string) bool {
	// Example logic: Compare trimmed strings
	return strings.TrimSpace(cmd) == strings.TrimSpace(config)
}

// Start initializes the processor and starts the main loop
func (processor *XpertTempCheckProcessor) Start(isTestMode bool) bool {
	processor.ReadConfiguration()

	if processor.TempCheckTimer <= 0 {
		log.Printf("Invalid TempCheckTimer value: %d", processor.TempCheckTimer)
		return false
	}

	timer := time.NewTicker(time.Duration(processor.TempCheckTimer) * time.Millisecond)
	processor.MProducer = &XpertMQTTClientProducer{}

	go func() {
		for range timer.C {
			processor.CheckTempCheckpoints()
		}
	}()

	log.Println("Processor started successfully.")
	return true
}

// Infrastructure represents a proximity infrastructure
type Infrastructure struct {
	ID   int
	Name string
}

func main() {
	processor := NewXpertTempCheckProcessor()
	processor.Start(false)

	// Keep the main function running
	select {}
}
