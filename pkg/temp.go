package temp

import (
	"log"
	"time"
	"fmt"
	"strings"
	"strconv"
	_ "github.com/denisenkom/go-mssqldb"

	core "Sofia/XpertCore"
	
)

//post
type XpertMethodStatus int

const (
	NotSuccessful XpertMethodStatus = iota
	Successful
)

// Placeholder for XpertMQTTClientProducer
type XpertMQTTClientProducer struct {}

// Placeholder for XpertEventService
type XpertEventService struct {}

func (service *XpertEventService) CloseEvent(customerID int, eventID int) *XpertResultObject {
	// Implement CloseEvent stored procedure call
	core.CurrentTime() //test for core
	return &XpertResultObject{}
}

func (service *XpertEventService) InsertEvent(event *XpertEventModel) *XpertResultObject {
	// Implement InsertEvent stored procedure call
	return &XpertResultObject{}
}

func (service *XpertEventService) InsertEventAction(action *XpertEventActionModel) *XpertResultObject {
	// Implement InsertEventAction stored procedure call
	return &XpertResultObject{}
}

func (service *XpertEventService) UpdateEvent(event *XpertEventModel) *XpertResultObject {
	// Implement UpdateEvent stored procedure call
	return &XpertResultObject{}
}

// Placeholder for XpertCustomerService
type XpertCustomerService struct {}

func (service *XpertCustomerService) GetAllCustomers(param int) *XpertCustomersModel {
	// Implement GetAllCustomers stored procedure call
	return &XpertCustomersModel{}
}

// Placeholder for XpertMetaDataService
type XpertMetaDataService struct {}

func (service *XpertMetaDataService) GetGroupsResponsibleFor(param1 int, param2 int) *XpertGroupsModel {
	// Implement GetGroupsResponsibleFor stored procedure call
	return &XpertGroupsModel{}
}

func (service *XpertMetaDataService) GetTemperatureCheckpoints(param1, param2, param3, param4 int, param5, param6 time.Time, param7 int) *XpertTemperatureCheckpointsModel {
	// Implement GetTemperatureCheckpoints stored procedure call
	return &XpertTemperatureCheckpointsModel{}
}

func (service *XpertMetaDataService) UpdateTemperatureCheckpoint(param1 int, param2 string, param3 int) {
	// Implement UpdateTemperatureCheckpoint stored procedure call
}

func (service *XpertMetaDataService) CreateTemperatureCheckpoint(param1, param2 int, param3 bool, param4 string, param5 int) {
	// Implement CreateTemperatureCheckpoint stored procedure call
}

// Placeholder for XpertStaffService
type XpertStaffService struct {}

func (service *XpertStaffService) GetStaffPersonById(customerID, itemID int) *XpertStaffPersonModel {
	// Implement GetStaffPersonById stored procedure call
	return &XpertStaffPersonModel{}
}

// Placeholder for XpertDeviceService
type XpertDeviceService struct {}

func (service *XpertDeviceService) GetDevice(customerID, deviceID int) *XpertDeviceModel {
	// Implement GetDevice stored procedure call
	return &XpertDeviceModel{}
}

func (service *XpertDeviceService) GetConfiguration(customerID, configID int) *Configuration {
	// Implement GetConfiguration stored procedure call
	return &Configuration{}
}

func (service *XpertDeviceService) GetDevicesByPendingTempConfig(param int) *XpertDevicesModel {
	// Implement GetDevicesByPendingTempConfig stored procedure call
	return &XpertDevicesModel{}
}

func (service *XpertDeviceService) SetDeviceConfigs(devices *XpertDevicesModel, pendingConfigID, param3, customerID int) {
	// Implement SetDeviceConfigs stored procedure call
}

// Placeholder for XpertSettingsService
type XpertSettingsService struct {}

func (service *XpertSettingsService) GetNotificationSettingsForSDCT(customerID int, param string) *XpertSDCTSettingModel {
	// Implement GetNotificationSettingsForSDCT stored procedure call
	return &XpertSDCTSettingModel{}
}

// Placeholder for XpertEmailLibrary
type XpertEmailLibrary struct {}

func (library *XpertEmailLibrary) SendEmail_Advanced(emails, details, description string, uris []string, contentType, priority int) bool {
	// Implement SendEmail_Advanced function
	return true
}

type XpertTempCheckProcessor struct {
	IsHealthy       bool
	IsTestMode      bool
	ProximitySeconds int
	TempCheckTimer  int
	StatusSeconds   int
	WebSocketTopic  string
	MProducer       *XpertMQTTClientProducer
}
//post
type XpertUserModel struct {
	Email string
	Name  string
}

type XpertResultObject struct {
	ErrorMessage string
	HasError     bool
}

type XpertEventModel struct {
	CustomerId       int
	SystemName       string
	Description      string
	MinValue         int
	MaxValue         int
	StartDateTime    time.Time
	EndDateTime      time.Time
	DateUpdated      time.Time
	DateCreated      time.Time
	ClosedDateTime   time.Time
	AllowedValueRange string
	PlanId           int
	DeviceId         int
	DateTimeToBeArchived time.Time
	UseCase          int
	ViolationValue   string
	Name             string
	DisplayName      string
	RuleName         string
	ItemId           int
	//post
	Details          string
}

type XpertEventActionModel struct {
	ItemEventId     int
	ActionTypeId    int
	ActionType      string
	ActionDateTime  time.Time
	Description     string
	ActionUserId    int
	DateCreated     time.Time
	DateUpdated     time.Time
}

type XpertStaffPersonModel struct {
	EnableAlerts bool
	Name         string
	DeviceID     int
	Id           int
	CustomerId   int
}

type XpertCustomersModel struct {
	List []XpertCustomerModel
}

type XpertCustomerModel struct {
	Id int
}

type XpertGroupsModel struct {
	List []XpertGroupModel
}

type XpertGroupModel struct {
	Id             int
	CheckTimeFrames string
}

type XpertTemperatureCheckpointsModel struct {
	List []XpertTemperatureCheckpointModel
}

type XpertTemperatureCheckpointModel struct {
	DateCreated time.Time
	Compliance  string
}

type Configuration struct {
	ConfigDef string
}

type XpertSDCTSettingModel struct {
	SettingJson string
}

type XpertDeviceModel struct {
	ModelName       string
	UniqueId        string
	PendingConfigId int
	ConfigId        int
	CustomerId      int
	IntegrationConfigId int
	Devices []XpertDeviceModel
}

type XpertDevicesModel struct {
	List []XpertDeviceModel
}

type InfraCounts struct {
	InfraId int
	Count   int
}

//post
type XpertDebugMessageJsonObject struct {
	Message string
}

//post
type XpertUseCaseModel struct {
	UseCase   string
	UseCaseId int
}

type XpertInfrastructureModel struct {
	//post
	Id   int
	Name string
}

func NewXpertTempCheckProcessor() *XpertTempCheckProcessor {
	return &XpertTempCheckProcessor{
		IsHealthy: true,
	}
}

func (processor *XpertTempCheckProcessor) CloseEvent(eventID int, oEventService *XpertEventService, customerID int) *XpertResultObject {
	return oEventService.CloseEvent(customerID, eventID)
}

func (processor *XpertTempCheckProcessor) CreateEvent(deviceMac string, oItem *XpertStaffPersonModel,
	oProximityInfrastructure *XpertInfrastructureModel, routeID int, startDate, endDate time.Time,
	duration float64, useCase, zoneID int, oEventModel *XpertEventModel, details string, oEventService *XpertEventService) *XpertResultObject {

	if !oItem.EnableAlerts {
		return &XpertResultObject{ErrorMessage: "Item does not have events enabled, alert generation cancelled"}
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
	oEventModel.PlanId = oProximityInfrastructure.Id
	oEventModel.DeviceId = oItem.DeviceID
	oEventModel.DateTimeToBeArchived = oEventModel.DateCreated.AddDate(0, 0, 60)
	oEventModel.UseCase = useCase

	return oEventService.InsertEvent(oEventModel)
}

func CreateItemEventAction(itemEventID int, actionType, actionDetails string) XpertMethodStatus {
	if itemEventID <= 0 {
		log.Printf("Invalid itemEventID: %d", itemEventID)
		return NotSuccessful
	}

	if strings.TrimSpace(actionType) == "" {
		log.Printf("Invalid actionType: %s", actionType)
		return NotSuccessful
	}

	if strings.TrimSpace(actionDetails) == "" {
		log.Printf("Invalid actionDetails: %s", actionDetails)
		return NotSuccessful
	}

	eventTime := time.Now().UTC()
	oModel := &XpertEventActionModel{
		ItemEventId:    itemEventID,
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
		log.Printf("InsertEventAction failed: %v", oResult)
		return NotSuccessful
	}

	return Successful
}

func CheckAllInfrasVisited(routeDef interface{}, numVisits int, allVisits []*XpertEventModel) []InfraCounts {
	var unseenInfras []InfraCounts

	counts := make(map[int]int)
	infrasInRoute := []int{}

	for _, infra := range routeDef.([]interface{}) {
		infraID := int(infra.(map[string]interface{})["Id"].(float64))
		infrasInRoute = append(infrasInRoute, infraID)
		counts[infraID] = 0
	}

	for _, visit := range allVisits {
		if _, exists := counts[visit.PlanId]; exists {
			counts[visit.PlanId]++
		} else {
			counts[visit.PlanId] = 1
		}
	}

	for _, infra := range infrasInRoute {
		if counts[infra] < numVisits {
			unseenInfras = append(unseenInfras, InfraCounts{InfraId: infra, Count: counts[infra]})
		}
	}

	return unseenInfras
}

func (processor *XpertTempCheckProcessor) CheckTempCheckpoints() {
	offset := time.Now().UTC().Sub(time.Now()).Hours()
	oCustomerService := &XpertCustomerService{}
	oCustomers := oCustomerService.GetAllCustomers(0)

	for _, customer := range oCustomers.List {
		oMetaService := &XpertMetaDataService{}
		oGroups := oMetaService.GetGroupsResponsibleFor(0, customer.Id)
		for _, group := range oGroups.List {
			if group.CheckTimeFrames == "" {
				continue
			}
			groupTimeFrames := strings.Split(group.CheckTimeFrames, ";")
			oTempChecks := oMetaService.GetTemperatureCheckpoints(0, group.Id, 1, 100, time.Now().AddDate(0, 0, -1), time.Now(), customer.Id)
			for _, timeframe := range groupTimeFrames {
				timeframeParts := strings.Split(timeframe, ",")
				startHour, _ := strconv.Atoi(strings.Split(timeframeParts[0], ":")[0])
				startMinute, _ := strconv.Atoi(strings.Split(timeframeParts[0], ":")[1])
				endHour, _ := strconv.Atoi(strings.Split(timeframeParts[1], ":")[0])
				endMinute, _ := strconv.Atoi(strings.Split(timeframeParts[1], ":")[1])
				startDateTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), startHour-int(offset), startMinute, 0, 0, time.UTC)
				endDateTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), endHour-int(offset), endMinute, 0, 0, time.UTC)

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
							oMetaService.UpdateTemperatureCheckpoint(group.Id, "check", customer.Id)
							break
						}
					}
				} else if !checkPerformed && time.Now().After(endDateTime) {
					var latestTempCheck *XpertTemperatureCheckpointModel
					for _, tempCheck := range oTempChecks.List {
						if latestTempCheck == nil || tempCheck.DateCreated.After(latestTempCheck.DateCreated) {
							latestTempCheck = &tempCheck
						}
					}
					if latestTempCheck == nil || latestTempCheck.Compliance != "miss" {
						log.Printf("No check found for time range, marking as missed")
						oMetaService.CreateTemperatureCheckpoint(0, group.Id, false, "", customer.Id)
						oMetaService.UpdateTemperatureCheckpoint(group.Id, "miss", customer.Id)
					}
				}
			}
		}
	}
}

func (processor *XpertTempCheckProcessor) CheckTempGraces() {
	// Implementation of CheckTempGraces
	// This would involve interactions with various services and handling the grace period logic as in C# code
}

func (processor *XpertTempCheckProcessor) CheckPendingTempConfigs() {
	// Implementation of CheckPendingTempConfigs
	// This would involve retrieving pending configurations and validating them as in C# code
}

func GetCustomerUseCases(customerApps interface{}) []XpertUseCaseModel {
	var useCases []XpertUseCaseModel
	if customerApps.(map[string]interface{})["AssetTracking"].(bool) {
		useCases = append(useCases, XpertUseCaseModel{UseCase: "AssetTracking", UseCaseId: 1})
	}
	if customerApps.(map[string]interface{})["MELT"].(bool) {
		useCases = append(useCases, XpertUseCaseModel{UseCase: "MELT", UseCaseId: 2})
	}
	if customerApps.(map[string]interface{})["PatientFlow"].(bool) {
		useCases = append(useCases, XpertUseCaseModel{UseCase: "PatientFlow", UseCaseId: 3})
	}
	if customerApps.(map[string]interface{})["SDCT"].(bool) {
		useCases = append(useCases, XpertUseCaseModel{UseCase: "SDCT", UseCaseId: 4})
	}
	if customerApps.(map[string]interface{})["StaffSafety"].(bool) {
		useCases = append(useCases, XpertUseCaseModel{UseCase: "StaffSafety", UseCaseId: 5})
	}
	return useCases
}

func (processor *XpertTempCheckProcessor) ReadConfiguration() {
	// This involves reading configuration parameters and setting instance variables
}

func (processor *XpertTempCheckProcessor) SendEmail(users []*XpertUserModel, details, description string, eventID int, deviceMac string, oDebugMsg *XpertDebugMessageJsonObject) XpertMethodStatus {
	// Implementation of SendEmail
	// This involves sending an email and creating an item event action as in C# code
	return Successful
}

func (processor *XpertTempCheckProcessor) Start(isTestMode bool) bool {
	processor.ReadConfiguration()

	timer := time.NewTicker(time.Duration(processor.TempCheckTimer) * time.Millisecond)

	processor.MProducer = &XpertMQTTClientProducer{}

	go func() {
		for range timer.C {
			processor.CheckTempGraces()
			processor.CheckPendingTempConfigs()
			processor.CheckTempCheckpoints()
		}
	}()

	return true
}

func (processor *XpertTempCheckProcessor) StartMainThread() {
	processor.Start(false)
}

func (processor *XpertTempCheckProcessor) StopMainThread() {
	// Implementation of StopMainThread
	// This involves stopping the main thread and cleaning up resources
}

func TaskSchedulerOnUnobservedTaskException() {
	// Implementation of TaskSchedulerOnUnobservedTaskException
	// This involves handling unobserved task exceptions as in C# code
}

func main() {
	processor := NewXpertTempCheckProcessor()
	processor.StartMainThread()

	// Keep the main function running
	select {}
}