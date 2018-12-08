package models

import (
	"time"

	"yunion.io/x/onecloud/pkg/appctx"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
)

const (
	BILLING_TYPE_POSTPAID = "postpaid"
	BILLING_TYPE_PREPAID  = "prepaid"
)

type SBillingResourceBase struct {
	BillingType string    `width:"36" charset:"ascii" nullable:"true" default:"postpaid" list:"user" create:"optional"`
	ExpiredAt   time.Time `nullable:"true" list:"user" create:"optional"`
}

func (self *SBillingResourceBase) GetChargeType() string {
	if len(self.BillingType) > 0 {
		return self.BillingType
	} else {
		return BILLING_TYPE_POSTPAID
	}
}

type SCloudBillingInfo struct {
	Provider            string
	Account             string
	AccountId           string
	SubAccount          string
	SubAccountId        string
	SubAccountProject   string
	SubAccountProjectId string
	Region              string
	RegionId            string
	RegionExtId         string
	Zone                string
	ZoneId              string
	ZoneExtId           string
	PriceKey            string
	ChargeType          string
	InternetChargeType  string
}

func MakeCloudBillingInfo(region *SCloudregion, zone *SZone, provider *SCloudprovider) SCloudBillingInfo {
	info := SCloudBillingInfo{}

	if zone != nil {
		info.Zone = zone.GetName()
		info.ZoneId = zone.GetId()
	}

	if region != nil {
		info.Region = region.GetName()
		info.RegionId = region.GetId()
	}

	if provider != nil {
		info.SubAccount = provider.GetName()
		info.SubAccountId = provider.GetId()

		if len(provider.ProjectId) > 0 {
			info.SubAccountProjectId = provider.ProjectId
			tc, err := db.TenantCacheManager.FetchTenantById(appctx.Background, provider.ProjectId)
			if err == nil {
				info.SubAccountProject = tc.GetName()
			}
		}

		account := provider.GetCloudaccount()
		info.Account = account.GetName()
		info.AccountId = account.GetId()

		driver, err := provider.GetDriver()

		if err == nil {
			info.Provider = driver.GetId()

			if region != nil {
				iregion, err := driver.GetIRegionById(region.ExternalId)
				if err == nil {
					info.RegionExtId = iregion.GetId()
					if zone != nil {
						izone, err := iregion.GetIZoneById(zone.ExternalId)
						if err == nil {
							info.ZoneExtId = izone.GetId()
						}
					}
				}
			}
		}
	}

	return info
}
