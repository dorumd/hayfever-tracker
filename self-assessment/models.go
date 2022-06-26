package self_assessment

import "time"

type SelfAssessmentReport struct {
	ReportedAt    time.Time
	NoseFeeling   NoseFeeling
	EyesFeeling   EyesFeeling
	ThroatFeeling ThroatFeeling
}

type SelfAssessmentReportRepository interface {
	Store(report SelfAssessmentReport)
}

type selfAssessmentReportRepositoryUsingMySQL struct {
}

func (s selfAssessmentReportRepositoryUsingMySQL) Store(report SelfAssessmentReport) {

}

func NewMySQLRepository() SelfAssessmentReportRepository {
	return selfAssessmentReportRepositoryUsingMySQL{}
}
