package policy

//CreateTapePolicy is a start record from scracth policy
//Start recording events could be a heavy processing since we have to make database dump so we should verify if a dump already was created
type CreateTapePolicy struct {
}

//DatabaseBackupAlreadyCreated verify if database dump was created
func (policy *CreateTapePolicy) DatabaseBackupAlreadyCreated(solutionID string) (bool, error) {
	return false, nil
}
