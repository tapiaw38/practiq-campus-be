package tenant

import "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/practiqapi"

type TenantData struct {
	ID       string `json:"id"`
	SchoolID string `json:"school_id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type TenantAdminData struct {
	ID       string  `json:"id"`
	SchoolID string  `json:"school_id"`
	Name     *string `json:"name,omitempty"`
	Status   string  `json:"status"`
	Eligible *bool   `json:"eligible,omitempty"`
}

type EligibleSchoolData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListMineOutput struct {
	Data []TenantData `json:"data"`
}

type ListOutput struct {
	Data []TenantAdminData `json:"data"`
}

type ListEligibleSchoolsOutput struct {
	Data []EligibleSchoolData `json:"data"`
}

type ActivateOutput struct {
	Data TenantAdminData `json:"data"`
}

type SetStatusOutput struct {
	Data TenantAdminData `json:"data"`
}

type ListMembersOutput struct {
	Data []practiqapi.SchoolMemberInfo `json:"data"`
}

func eligible(school practiqapi.SchoolInfo) bool {
	return school.Kind == "institution" &&
		school.Billing == "direct" &&
		school.Status == "active"
}

func schoolsByID(schools []practiqapi.SchoolInfo) map[string]practiqapi.SchoolInfo {
	index := make(map[string]practiqapi.SchoolInfo, len(schools))
	for _, school := range schools {
		index[school.ID] = school
	}
	return index
}
