package review

import "github.com/xanzy/go-gitlab"

type ReviewStatus int

const (
	ReviewNeeded ReviewStatus = iota
	ReviewShipIt
)

type Review struct {
	MergeRequest *gitlab.MergeRequest
	LastPipeline *gitlab.PipelineInfo
}

func New(git *gitlab.Client, r *gitlab.MergeRequest) Review {
	return Review{
		MergeRequest: r,
		LastPipeline: fetchLastPipeline(git, r),
	}
}

func (r Review) Status() ReviewStatus {
	if r.MergeRequest.Upvotes > 0 {
		return ReviewShipIt
	}

	return ReviewNeeded
}

func (r Review) CheckStatus() string {
	if r.LastPipeline != nil {
		return r.LastPipeline.Status
	} else {
		return "unknown"
	}
}

func fetchLastPipeline(git *gitlab.Client, r *gitlab.MergeRequest) *gitlab.PipelineInfo {
	pipelines, _, err := git.MergeRequests.ListMergeRequestPipelines(r.ProjectID, r.IID)
	if err != nil {
		return nil
	}
	return pipelines[0]
}
