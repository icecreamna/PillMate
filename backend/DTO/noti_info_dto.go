package dto

type NotiInfoResp struct {
    ID            uint     `json:"id"`
    MyMedicineID  *uint    `json:"my_medicine_id,omitempty"`
    GroupID       *uint    `json:"group_id,omitempty"`
    NotiFormatID  uint     `json:"format_id"`
    StartDate     string   `json:"start_date"`
    EndDate       string   `json:"end_date"`
    IntervalHours *int     `json:"interval_hours,omitempty"`
    IntervalDay   *int     `json:"interval_day,omitempty"`
    Times         []string `json:"times,omitempty"`
    CyclePattern  []int    `json:"cycle_pattern,omitempty"`
}