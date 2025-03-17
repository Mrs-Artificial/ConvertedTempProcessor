package tempold


/*

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/robfig/cron/v3"
)

var db *sqlx.DB
var err error	

// Struct definitions
type XpertTempCheckProcessor struct {
	IsHealthy        bool
	IsTestMode       bool
	ProximitySeconds int
	TempCheckTimer   int
	StatusSeconds    int
	WebSocketTopic   string
	MQTTClient       mqtt.Client
	timer            *cron.Cron
	tokenSource      context.CancelFunc
}

type XpertEventService struct{}

func (s *XpertEventService) GetEventBySystemName(customerID, itemID int, systemName string, includeClosed bool) XpertEventModel {
	// Implement the logic to get event by system name
	return XpertEventModel{}
}

func (s *XpertEventService) GetOpenTempEvents(eventType string) XpertEventsModel {
	// Implement the logic to get open temperature events
	return XpertEventsModel{}
}

func (s *XpertEventService) CloseEvent(customerID, eventID int) XpertResultObject {
	// Implement the logic to close the event
	return XpertResultObject{}
}

func (s *XpertEventService) InsertEvent(event *XpertEventModel) XpertResultObject {
	// Implement the logic to insert the event
	return XpertResultObject{}
}

func (s *XpertEventService) InsertEventAction(action *XpertEventActionModel) XpertResultObject {
	// Implement the logic to insert the event action
	return XpertResultObject{}
}

type XpertStaffService struct{}

func (s *XpertStaffService) GetStaffPersonById(customerId, itemId int) *XpertStaffPersonModel {
	// Implement logic to get staff person by customerId and itemId
	return &XpertStaffPersonModel{}
}

type XpertDeviceService struct{}

func (s *XpertDeviceService) GetDevice(customerId, deviceId int) XpertDeviceModel {
	// Implement logic to get a device by customerId and deviceId
	return XpertDeviceModel{}
}

func (s *XpertDeviceService) GetDevicesByPendingTempConfig(configId int) XpertDevicesModel {
	// Implement logic to get devices by pending temp config
	return XpertDevicesModel{}
}

type XpertSettingsService struct{}

func (s *XpertSettingsService) GetNotificationSettingsForSDCT(customerId int, settingType string) XpertSDCTSettingModel {
	// Implement logic to get notification settings for SDCT
	return XpertSDCTSettingModel{}
}

type XpertCustomerService struct{}

type XpertCustomersModel struct {
	List []XpertCustomerModel
}

func (s *XpertCustomerService) GetAllCustomers(dummyParam int) XpertCustomersModel {
	// Implement logic to get all customers
	return XpertCustomersModel{}
}

type XpertMetaDataService struct{}

func (s *XpertMetaDataService) UpdateTemperatureCheckpoint(groupId int, compliance string, customerId int) {
	// Implement logic to update temperature checkpoint
}

func (s *XpertMetaDataService) CreateTemperatureCheckpoint(dummyParam, groupId int, param1 bool, param2 string, customerId int) {
	// Implement logic to create temperature checkpoint
}

func (s *XpertMetaDataService) GetTemperatureCheckpoints(dummyParam, groupId, param1, param2 int, startDate, endDate time.Time, customerId int) XpertTemperatureCheckpointsModel {
	// Implement logic to get temperature checkpoints
	return XpertTemperatureCheckpointsModel{}
}

func (s *XpertMetaDataService) GetGroupsResponsibleFor(dummyParam, customerId int) XpertGroupsModel {
	// Implement logic to get groups responsible for a customer
	return XpertGroupsModel{}
}

type XpertGroupsModel struct {
	List []XpertGroupModel
}

type XpertLogStrings struct{}
type XpertResultObject struct {
	HasError     bool
	ErrorMessage string
}

type XpertEventModel struct {
	Details              string
	SystemName           string
	Description          string
	MinValue             int
	MaxValue             int
	StartDateTime        time.Time
	EndDateTime          time.Time
	DateUpdated          time.Time
	DateCreated          time.Time
	ClosedDateTime       time.Time
	AllowedValueRange    string
	PlanId               int
	DeviceId             int
	DateTimeToBeArchived time.Time
	UseCase              int
	CustomerId           int
	DeviceUniqueId       string
	ViolationValue       float64
}

type XpertStaffPersonModel struct {
	EnableAlerts bool
	Name         string
	DeviceID     int
	CustomerId   int
	Id           int
}

type XpertInfrastructureModel struct {
	Name string
	Id   int
}

type XpertCustomerModel struct {
	Id int
}

type XpertGroupModel struct {
	CheckTimeFrames string
	Id              int
}

type XpertTemperatureCheckpointsModel struct {
	List []XpertTemperatureCheckpointModel
}

type XpertTemperatureCheckpointModel struct {
	DateCreated time.Time
	Compliance  string
}

type XpertEventActionModel struct {
	ItemEventId    int
	ActionTypeId   int
	ActionType     string
	ActionDateTime time.Time
	Description    string
	ActionUserId   int
	DateCreated    time.Time
	DateUpdated    time.Time
}

type InfraCounts struct {
	infraId int
	count   int
}

type XpertDebugMessageJsonObject struct {
	Properties []DebugProperty
}

type DebugProperty struct {
	isList bool
	Name   string
	Value  string
}

type XpertUserModel struct {
	Email string
}

type XpertEmailLibrary struct{}

func (e *XpertEmailLibrary) SendEmail_Advanced(emails, details, description string, uris []string, contentType, priority int) bool {
	// Simulate sending email
	return true
}

type XpertEnums struct{}

const (
	XpertMethodStatusSuccessful    = 0
	XpertMethodStatusNotSuccessful = 1
)

type XpertIntegrationModel struct {
	JSON string
}

type XpertIntegrationsModel struct {
	List []XpertIntegrationModel
}

type XpertDeviceModel struct {
	ModelName           string
	CustomerId          int
	UniqueId            string
	PendingConfigId     int
	ConfigId            int
	IntegrationConfigId int
	ItemId              int
}

type XpertDevicesModel struct {
	List []XpertDeviceModel
}

type XpertEventsModel struct {
	List []XpertEventModel
}

type XpertSDCTSettingModel struct {
	SettingJson string
}

type XpertUseCaseModel struct {
	UseCase   string
	UseCaseId int
}

func NewXpertTempCheckProcessor() *XpertTempCheckProcessor {
	_, cancel := context.WithCancel(context.Background())
	processor := &XpertTempCheckProcessor{
		IsHealthy:        true,
		tokenSource:      cancel,
		ProximitySeconds: 3600,
		TempCheckTimer:   3600,
		StatusSeconds:    3600,
	}
	processor.timer = cron.New()
	return processor
}

func (p *XpertTempCheckProcessor) CloseEvent(eventID int, oEventService *XpertEventService, customerID int) XpertResultObject {
	//the C# method CloseEvent
	var result XpertResultObject
	result = oEventService.CloseEvent(customerID, eventID)
	return result
}

func (p *XpertTempCheckProcessor) CreateEvent(deviceMac string, oItem *XpertStaffPersonModel, oProximityInfrastructure *XpertInfrastructureModel, routeId int, startDate, endDate time.Time, duration float64, useCase, zoneId int, oEventModel *XpertEventModel, details string) XpertResultObject {
	//the C# method CreateEvent
	var result XpertResultObject
	if !oItem.EnableAlerts {
		result.ErrorMessage = "Item does not have events enabled, alert generation cancelled"
		return result
	}
	oEventModel.Details = details
	if strings.Contains(strings.ToLower(oEventModel.SystemName), "alert") || strings.Contains(strings.ToLower(oEventModel.SystemName), "warning") {
		oEventModel.Description = details
	} else {
		oEventModel.Description = fmt.Sprintf("%s %s %s", oItem.Name, deviceMac, oProximityInfrastructure.Name)
	}
	oEventModel.MinValue = routeId
	oEventModel.MaxValue = zoneId
	oEventModel.StartDateTime = startDate
	oEventModel.EndDateTime = endDate
	oEventModel.DateUpdated = time.Now().UTC()
	oEventModel.DateCreated = time.Now().UTC()
	oEventModel.ClosedDateTime = endDate
	oEventModel.AllowedValueRange = fmt.Sprintf("%f", duration)
	oEventModel.PlanId = oProximityInfrastructure.Id
	oEventModel.DeviceId = oItem.DeviceID
	oEventModel.DateTimeToBeArchived = oEventModel.DateCreated.Add(60 * 24 * time.Hour)
	oEventModel.UseCase = useCase
	eventService := &XpertEventService{}
	result = eventService.InsertEvent(oEventModel)
	return result
}

func (p *XpertTempCheckProcessor) CreateItemEventAction(itemEventId int, actionType, actionDetails string) int {
	//the C# method CreateItemEventAction
	const methodName = "CreateItemEventAction"
	var result int
	result = XpertMethodStatusNotSuccessful

	if itemEventId <= 0 || actionType == "" || actionDetails == "" {
		log.Printf("Invalid parameters for %s", methodName)
		return result
	}

	eventTime := time.Now().UTC()

	oModel := &XpertEventActionModel{
		ItemEventId:    itemEventId,
		ActionTypeId:   0,
		ActionType:     actionType,
		ActionDateTime: eventTime,
		Description:    actionDetails,
		ActionUserId:   0,
		DateCreated:    eventTime,
		DateUpdated:    eventTime,
	}

	oEventService := &XpertEventService{}
	oResult := oEventService.InsertEventAction(oModel)
	if oResult.HasError {
		log.Println(fmt.Sprintf("Error inserting event action in %s", methodName))
		return result
	}

	result = XpertMethodStatusSuccessful
	return result
}

func CheckAllInfrasVisited(routeDef interface{}, numVisits int, allVisits []XpertEventModel) ([]InfraCounts, []InfraCounts) {
	//the C# method CheckAllInfrasVisited
	const methodName = "CheckAllInfrasVisited"
	var unseenInfras []InfraCounts
	counts := make(map[int]int)
	var infrasInRoute []int

	for _, infra := range routeDef.([]interface{}) {
		infraId := infra.(map[string]interface{})["Id"].(int)
		infrasInRoute = append(infrasInRoute, infraId)
		counts[infraId] = 0
	}

	for _, visit := range allVisits {
		if _, ok := counts[visit.PlanId]; ok {
			counts[visit.PlanId]++
		}
	}

	for _, infra := range infrasInRoute {
		if counts[infra] < numVisits {
			unseenInfras = append(unseenInfras, InfraCounts{infra, counts[infra]})
		}
	}

	return unseenInfras, unseenInfras
}

func (p *XpertTempCheckProcessor) CheckTempCheckpoints() {
	//the C# method CheckTempCheckpoints
	offset := time.Now().UTC().Hour() - time.Now().Hour()

	customerService := &XpertCustomerService{}
	oCustomers := customerService.GetAllCustomers(0)
	for _, customer := range oCustomers.List {
		metaService := &XpertMetaDataService{}
		oGroups := metaService.GetGroupsResponsibleFor(0, customer.Id)
		for _, group := range oGroups.List {
			if group.CheckTimeFrames == "" {
				continue
			}
			groupTimeFrames := strings.Split(group.CheckTimeFrames, ";")
			oTempChecks := metaService.GetTemperatureCheckpoints(0, group.Id, 1, 100, time.Now().AddDate(0, 0, -1), time.Now(), customer.Id)
			for _, timeframe := range groupTimeFrames {
				timeFrameParts := strings.Split(timeframe, ",")
				startDateTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), convertToInt(timeFrameParts[0])+(-1*offset), convertToInt(timeFrameParts[1]), 0, 0, time.UTC)
				endDateTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), convertToInt(timeFrameParts[2])+(-1*offset), convertToInt(timeFrameParts[3]), 0, 0, time.UTC)
				checkPerformed := false
				for _, tempCheck := range oTempChecks.List {
					if endDateTime.After(tempCheck.DateCreated) && tempCheck.DateCreated.After(startDateTime) {
						checkPerformed = true
						break
					}
				}
				if checkPerformed {
					// Update temperature checkpoint if compliance is empty
					checkpoint := oTempChecks.List[0]
					if checkpoint.Compliance == "" {
						metaService.UpdateTemperatureCheckpoint(group.Id, "check", customer.Id)
					}
				} else if !checkPerformed && time.Now().After(endDateTime) {
					checkpoint := oTempChecks.List[0]
					if checkpoint.Compliance != "miss" {
						log.Println("No check found for time range, marking as missed")
						metaService.CreateTemperatureCheckpoint(0, group.Id, false, "", customer.Id)
						metaService.UpdateTemperatureCheckpoint(group.Id, "miss", customer.Id)
					}
				}
			}
		}
	}
}

func (p *XpertTempCheckProcessor) CreateTempEvent(item *XpertStaffPersonModel, startDate, endDate time.Time, duration float64, eventModel *XpertEventModel, eventService *XpertEventService, eventType string, minValue, maxValue, violationValue float64, description, details string) XpertResultObject {
	//create a temperature event
	eventModel.Details = details
	eventModel.Description = description
	eventModel.MinValue = int(minValue)
	eventModel.MaxValue = int(maxValue)
	eventModel.StartDateTime = startDate
	eventModel.EndDateTime = endDate
	eventModel.DateUpdated = time.Now().UTC()
	eventModel.DateCreated = time.Now().UTC()
	eventModel.ClosedDateTime = endDate
	eventModel.AllowedValueRange = fmt.Sprintf("%f", duration)
	eventModel.PlanId = item.Id
	eventModel.DeviceId = item.DeviceID
	eventModel.DateTimeToBeArchived = eventModel.DateCreated.Add(60 * 24 * time.Hour)
	eventModel.UseCase = 1 // Example use case
	return eventService.InsertEvent(eventModel)
}

func (p *XpertTempCheckProcessor) CheckTempGraces() {
	//the C# method CheckTempGraces
	log.Printf("XpertTempCheckProcessor CheckTempGraces Start")

	staffService := &XpertStaffService{}
	eventService := &XpertEventService{}
	deviceService := &XpertDeviceService{}
	settingsService := &XpertSettingsService{}
	eventModel := XpertEventModel{
		SystemName: "ERROR - must be re-assigned.",
	}

	events := eventService.GetOpenTempEvents("TEMPERATURE_GRACE")
	for _, tempEvent := range events.List {
		if tempEvent.DeviceId == 0 || tempEvent.CustomerId == 0 {
			continue
		}
		device := deviceService.GetDevice(tempEvent.CustomerId, tempEvent.DeviceId)
		if device.ModelName != "ts1" && device.ModelName != "ts2" && device.ModelName == "hs1" {
			log.Printf("XpertTempCheckProcessor CheckTempGraces DEVICE TYPE %s NOT HANDLED BY THIS PROCESSOR", device.ModelName)
			continue
		}

		eventModel.DeviceId = tempEvent.DeviceId
		eventModel.CustomerId = tempEvent.CustomerId
		eventModel.DeviceUniqueId = device.UniqueId

		config := getDeviceConfiguration(deviceService, &device)
		if config == nil {
			continue
		}

		settings := settingsService.GetNotificationSettingsForSDCT(tempEvent.CustomerId, "Temp")
		settingsJson := parseSettingsJson(settings.SettingJson)
		if time.Now().After(tempEvent.DateUpdated.Add(time.Duration(config["hightime1"].(int)) * time.Second)) {
			lastAlertEvent := eventService.GetEventBySystemName(0, tempEvent.DeviceId, "TEMPERATURE_ALERT", false)
			if time.Now().After(lastAlertEvent.DateUpdated.Add(time.Duration(getRenotifyPeriod(settingsJson)) * time.Second)) {
				item := staffService.GetStaffPersonById(tempEvent.CustomerId, device.ItemId)
				p.CreateTempEvent(item, time.Now(), time.Now(), 1, &eventModel, eventService, "TEMPERATURE_ALERT", config["lowvalue1"].(float64), config["highvalue1"].(float64), tempEvent.ViolationValue, "Temperature Range Exceeded.", "0")
			}
		}
	}
}

func (p *XpertTempCheckProcessor) CheckPendingTempConfigs() {
	//the C# method CheckPendingTempConfigs
	log.Printf("XpertTempCheckProcessor CheckPendingTempConfigs Start")

	deviceService := &XpertDeviceService{}
	customerService := &XpertCustomerService{}

	devices := deviceService.GetDevicesByPendingTempConfig(1)
	for _, device := range devices.List {
		tryProcessDevice(&device, deviceService, customerService)
	}
}

func (p *XpertTempCheckProcessor) Start(isTestMode bool) bool {
	//the C# method Start
	log.Printf("Entering Start")

	p.ReadConfiguration()

	p.timer.AddFunc(fmt.Sprintf("@every %ds", p.TempCheckTimer), func() {
		p.CheckTempGraces()
		p.CheckPendingTempConfigs()
		p.CheckTempCheckpoints()
	})
	p.timer.Start()

	p.MQTTClient = createMQTTClient()

	go p.CheckTempGraces()
	go p.CheckPendingTempConfigs()
	go p.CheckTempCheckpoints()

	return true
}

func (p *XpertTempCheckProcessor) Stop() {
	//the C# method Stop
	p.timer.Stop()
	p.tokenSource()
}

func (p *XpertTempCheckProcessor) ReadConfiguration() {
	//the C# method ReadConfiguration
	log.Printf("Reading configuration in ReadConfiguration")

	// Implementation to read configuration values
	// Example:
	p.TempCheckTimer = getConfigValue("TempCheckTimer", 3600)
	p.WebSocketTopic = getConfigValueString("WebSocketTopic", "")
	p.ProximitySeconds = getConfigValue("ProximitySeconds", 3600)
	p.StatusSeconds = getConfigValue("StatusSeconds", 3600)
}
/*
func main() {
	processor := NewXpertTempCheckProcessor()
	processor.Start(false)
}


func convertToInt(value string) int {
	i, _ := strconv.Atoi(value)
	return i
}

func getConfigValue(key string, defaultValue int) int {
	// Implement logic to read configuration values
	return defaultValue
}

func getConfigValueString(key string, defaultValue string) string {
	// Implement logic to read configuration values
	return defaultValue
}

func createMQTTClient() mqtt.Client {
	// Implement logic to create and return MQTT client
	return nil
}

func getDeviceConfiguration(deviceService *XpertDeviceService, device *XpertDeviceModel) map[string]any {
	// Implement logic to get device configuration
	return nil
}

func parseSettingsJson(settingJson string) map[string]any {
	// Implement logic to parse settings JSON
	return nil
}

func getRenotifyPeriod(settingsJson map[string]any) int {
	// Implement logic to get renotify period from settings JSON
	return 0
}

func tryProcessDevice(device *XpertDeviceModel, deviceService *XpertDeviceService, customerService *XpertCustomerService) {
	const methodName = "tryProcessDevice"
	log.Printf("Processing device %s in %s", device.UniqueId, methodName)
	// Implement the rest of the logic
}

*/
// The above code is a Go implementation of a temperature check processor.