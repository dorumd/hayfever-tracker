package self_assessment

import "time"

type NoseFeeling string
type EyesFeeling string
type ThroatFeeling string

type ReportInput struct {
	ReportedAt    time.Time     `json:"reportedAt"`
	NoseFeeling   NoseFeeling   `json:"noseFeeling"`
	EyesFeeling   EyesFeeling   `json:"eyesFeeling"`
	ThroatFeeling ThroatFeeling `json:"throatFeeling"`
}

type SelfAssessmentService interface {
	SubmitReport(input ReportInput)
}

type selfAssessmentService struct {
	repository SelfAssessmentReportRepository
}

func (s selfAssessmentService) SubmitReport(input ReportInput) {
	report := SelfAssessmentReport{
		ReportedAt:    input.ReportedAt,
		NoseFeeling:   input.NoseFeeling,
		EyesFeeling:   input.EyesFeeling,
		ThroatFeeling: input.ThroatFeeling,
	}
	s.repository.Store(report)
}

func New(repository SelfAssessmentReportRepository) SelfAssessmentService {
	return selfAssessmentService{
		repository: repository,
	}
}
