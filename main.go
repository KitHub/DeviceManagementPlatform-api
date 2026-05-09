package devicemanagementplatformapi

import "DeviceManagementPlatform-api/config"

func main() {

	// init config
	_, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

}
