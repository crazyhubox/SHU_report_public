package report

type ReportResult struct {
	result map[string]string
	date   string
}

func NewReportResult(date string) *ReportResult {
	var tMap = make(map[string]string)
	//map映射对象必须要创建，不然就会指向nil，指向nil的map无法保存任何值
	return &ReportResult{
		result: tMap,
		date:   date,
	}
}
