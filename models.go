package main

type TrainScheduleReqBody struct {
	Org        string `json:"org"`
	Dest       string `json:"dest"`
	DepartDate string `json:"depart_date"`
	TrainId    string `json:"train_id"`
	TotPsg     int    `json:"tot_psg"`
}

type UpdateTrainScheduleReqBody struct {
	TrainScheduleReqBody
	Stock int `json:"stock"`
}
