package core

type Consensus struct {
		TypeConsensus int //1 pow | 2 pos | 3 pbft | 4 bftraft | 5 ripple NAO DA PRA USAR TYPE, E PALAVRA RESERVADA
		Difficulty int
}

func ConsensusCreate(typeConsensus int, difficulty int) (Consensus) {
		consensus := Consensus{}
		consensus = Consensus{typeConsensus, difficulty}
		return consensus
}

func (consensus Consensus) GetType() (int) {
		return consensus.TypeConsensus
}

func (consensus Consensus) GetDifficulty() (int) {
		return consensus.Difficulty
}
