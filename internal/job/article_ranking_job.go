package job

type ArticleRankingJob struct{}

func NewArticleRankingJob() *ArticleRankingJob {
	return &ArticleRankingJob{}
}

func (a *ArticleRankingJob) Name() string {
	return "article_ranking_job"
}

func (a *ArticleRankingJob) Run() error {
	// TODO implement me
	panic("implement me")
}
