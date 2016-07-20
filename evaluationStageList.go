package govaluate

type evaluationStageList struct {

	stages []evaluationStage
}

func newEvaluationStageList() {

	return new(evaluationStage)
}

func (this *evaluationStageList) addStage(operator evaluationOperator, leftTypeCheck stageTypeCheck, rightTypeCheck stageTypeCheck, typeErrorFormat string) {

	stage := evaluationStage {
		operator: operator,
		leftTypeCheck: leftTypeCheck,
		rightTypeCheck: rightTypeCheck,
		typeErrorFormat: typeErrorFormat,
	}

	extantStageCount := len(this.stages)
	if(extantStageCount > 0) {
		this.stages[extantStageCount-1].rightStage = &stage
	}
	this.stages = append(this.stages, stage)
}
