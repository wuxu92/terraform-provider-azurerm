package util

import (
	"log"
	"strings"
)

func GetRPOwner(rp string) string {
	rp = strings.ReplaceAll(strings.ToUpper(rp), " ", "")
	if val, ok := serviceMappingAlt[rp]; ok {
		return val
	}
	//
	log.Printf("no owner for rp: %v", rp)
	return ""
}

var serviceMappingAlt = map[string]string{
	// yunliu1
	"AZURESTACKHCI":   "yunliu1",
	"BATCH":           "yunliu1",
	"DNS":             "yunliu1",
	"KEYVAULT":        "yunliu1",
	"KUSTO":           "yunliu1",
	"LOADTEST":        "yunliu1",
	"LOADTESTSERVICE": "yunliu1",
	"PRIVATEDNS":      "yunliu1",
	// v-elenaxin
	"ANALYSISSERVICES":            "v-elenaxin",
	"DESKTOPVIRTUALIZATION":       "v-elenaxin",
	"ELASTIC":                     "v-elenaxin",
	"MARIADB":                     "v-elenaxin",
	"MEDIA":                       "v-elenaxin",
	"MIXEDREALITY":                "v-elenaxin",
	"MSSQL":                       "v-elenaxin",
	"MYSQL":                       "v-elenaxin",
	"POSTGRES":                    "v-elenaxin",
	"SQL":                         "v-elenaxin",
	"MICROSOFTSQLSERVER/AZURESQL": "v-elenaxin",
	"VIDEOANALYZER":               "v-elenaxin",
	// henglu
	"COGNITIVE":         "henglu",
	"CONTAINER":         "henglu",
	"CONTAINERSERVICES": "henglu",
	"DATADOG":           "henglu",
	"DATAFACTORY":       "henglu",
	"PURVIEW":           "henglu",
	"SPRINGCLOUD":       "henglu",
	"SYNAPSE":           "henglu",
	// jiaweitao
	"ATTESTATION":      "jiaweitao",
	"CONNECTIONS":      "jiaweitao",
	"SERVICECONNECTOR": "jiaweitao",
	"DEVTESTLABS":      "jiaweitao",
	"DEVTEST":          "jiaweitao",
	"HDINSIGHT":        "jiaweitao",
	"STREAMANALYTICS":  "jiaweitao",
	"TRAFFICMANAGER":   "jiaweitao",
	"VMWARE":           "jiaweitao",
	// v-cheye
	"BOT":                  "v-cheye",
	"COMMUNICATION":        "v-cheye",
	"CONFIDENTIALLEDGER":   "v-cheye",
	"COSMOS":               "v-cheye",
	"COSMOSDB":             "v-cheye",
	"COSTMANAGEMENT":       "v-cheye",
	"CUSTOMPROVIDERS":      "v-cheye",
	"DATABOXEDGE":          "v-cheye",
	"HSM":                  "v-cheye",
	"LOGANALYTICS":         "v-cheye",
	"MAPS":                 "v-cheye",
	"NETAPP":               "v-cheye",
	"PORTAL":               "v-cheye",
	"POWERBI":              "v-cheye",
	"SERVICEFABRIC":        "v-cheye",
	"SERVICEFABRICMANAGED": "v-cheye",
	// wangta
	"APPCONFIGURATION":    "wangta",
	"APPLICATIONINSIGHTS": "wangta",
	"DATAPROTECTION":      "wangta",
	"MONITOR":             "wangta",
	"MSI":                 "wangta",
	"RELAY":               "wangta",
	// xiaxin.yi
	"APPSERVICE":      "xiaxin.yi",
	"DATABRICKS":      "xiaxin.yi",
	"EVENTGRID":       "xiaxin.yi",
	"EVENTHUB":        "xiaxin.yi",
	"HEALTHCARE":      "xiaxin.yi",
	"LOADBALANCER":    "xiaxin.yi",
	"NOTIFICATIONHUB": "xiaxin.yi",
	"SEARCH":          "xiaxin.yi",
	"SERVICEBUS":      "xiaxin.yi",
	"SIGNALR":         "xiaxin.yi",
	"WEB":             "xiaxin.yi",
	// xuwu1
	"AUTOMATION": "xuwu1",
	"FIREWALL":   "xuwu1",
	"FLUIDRELAY": "xuwu1",
	"NETWORK":    "xuwu1",
	"POLICY":     "xuwu1",
	"ORBITAL":    "xuwu1",
	// xuzhang3
	"MACHINELEARNING": "xuzhang3",
	"APIMANAGEMENT":   "xuzhang3",
	"REDIS":           "xuzhang3",
	"REDISENTERPRISE": "xuzhang3",
	// yicma
	"COMPUTE":               "yicma",
	"DIGITALTWINS":          "yicma",
	"DISKS":                 "yicma",
	"IOTCENTRAL":            "yicma",
	"IOTHUB":                "yicma",
	"IOTTIMESERIESINSIGHTS": "yicma",
	"TIMESERIESINSIGHTS":    "yicma",
	"LEGACY":                "yicma",
	// zhaoting.weng
	"AADB2C":            "zhaoting.weng",
	"AUTHORIZATION":     "zhaoting.weng",
	"CDN":               "zhaoting.weng",
	"DATABASEMIGRATION": "zhaoting.weng",
	"DOMAINSERVICES":    "zhaoting.weng",
	"FRONTDOOR":         "zhaoting.weng",
	"STORAGE":           "zhaoting.weng",
	// zhenteng
	"DATASHARE":        "zhenteng",
	"HPCCACHE":         "zhenteng",
	"LOGIC":            "zhenteng",
	"SECURITYCENTER":   "zhenteng",
	"SENTINEL":         "zhenteng",
	"RECOVERYSERVICES": "zhenteng",
	// zhhu
	"ADVISOR":             "zhhu",
	"BILLING":             "zhhu",
	"BLUEPRINTS":          "zhhu",
	"CONSUMPTION":         "zhhu",
	"LIGHTHOUSE":          "zhhu",
	"LOGZ":                "zhhu",
	"MAINTENANCE":         "zhhu",
	"MANAGEDAPPLICATIONS": "zhhu",
	"MANAGEMENTGROUP":     "zhhu",
	"RESOURCE":            "zhhu",
	"RESOURCES":           "zhhu",
	"SUBSCRIPTION":        "zhhu",
}
