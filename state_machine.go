package smp

type smpState interface {
	receiveMessage1(*Protocol, SMP1) (smpState, Message, error)
	receiveMessage2(*Protocol, SMP2) (smpState, Message, error)
	receiveMessage3(*Protocol, SMP3) (smpState, Message, error)
	receiveMessage4(*Protocol, SMP4) (smpState, Message, error)
}

type smpStateBase struct{}
type smpStateExpect1 struct{ smpStateBase }
type smpStateExpect2 struct{ smpStateBase }
type smpStateExpect3 struct{ smpStateBase }
type smpStateExpect4 struct{ smpStateBase }

func abortState(e error) (smpState, Message, error) {
	//TODO: should wipe all the intermediate s1..s4 values
	return smpStateExpect1{}, SMPAbort{}, e
}

func sendSMPAbortAndRestartStateMachine() (smpState, Message, error) {
	//must return nil error otherwise the abort message will be ignored
	return abortState(nil)
}

func abortStateMachineAndNotifyCheated(p *Protocol) (smpState, Message, error) {
	p.event(Cheated)
	return sendSMPAbortAndRestartStateMachine()
}

func abortStateMachineAndNotifyError(p *Protocol) (smpState, Message, error) {
	p.event(Error)
	return sendSMPAbortAndRestartStateMachine()
}

func (smpStateBase) receiveMessage1(p *Protocol, m SMP1) (smpState, Message, error) {
	return abortStateMachineAndNotifyError(p)
}

func (smpStateBase) receiveMessage2(p *Protocol, m SMP2) (smpState, Message, error) {
	return abortStateMachineAndNotifyError(p)
}

func (smpStateBase) receiveMessage3(p *Protocol, m SMP3) (smpState, Message, error) {
	return abortStateMachineAndNotifyError(p)
}

func (smpStateBase) receiveMessage4(p *Protocol, m SMP4) (smpState, Message, error) {
	return abortStateMachineAndNotifyError(p)
}

func (smpStateExpect1) receiveMessage1(p *Protocol, m SMP1) (smpState, Message, error) {
	err := p.verifySMP1(m)
	if err != nil {
		return abortStateMachineAndNotifyCheated(p)
	}

	m2, err := p.newSMP2Message(m)
	if err != nil {
		return abortStateMachineAndNotifyCheated(p)
	}

	p.event(InProgress)
	return smpStateExpect3{}, m2, nil
}

func (smpStateExpect2) receiveMessage2(p *Protocol, m SMP2) (smpState, Message, error) {
	err := p.verifySMP2(m)
	if err != nil {
		return abortStateMachineAndNotifyCheated(p)
	}

	m3, err := p.newSMP3Message(m)
	if err != nil {
		return abortStateMachineAndNotifyCheated(p)
	}

	p.event(InProgress)
	return smpStateExpect4{}, m3, nil
}

func (smpStateExpect3) receiveMessage3(p *Protocol, m SMP3) (smpState, Message, error) {
	err := p.verifySMP3(m)
	if err != nil {
		return abortStateMachineAndNotifyCheated(p)
	}

	err = p.verifySMP3ProtocolSuccess(m)
	if err != nil {
		p.event(Failure)
		return sendSMPAbortAndRestartStateMachine()
	}

	msg, err := p.newSMP4Message(m)
	if err != nil {
		return abortStateMachineAndNotifyCheated(p)
	}

	p.event(Success)
	return smpStateExpect1{}, msg, nil
}

func (smpStateExpect4) receiveMessage4(p *Protocol, m SMP4) (smpState, Message, error) {
	err := p.verifySMP4(m)
	if err != nil {
		return abortStateMachineAndNotifyCheated(p)
	}

	err = p.verifySMP4ProtocolSuccess(m)
	if err != nil {
		p.event(Failure)
		return sendSMPAbortAndRestartStateMachine()
	}

	p.event(Success)
	return smpStateExpect1{}, nil, nil
}

func (m SMP1) received(p *Protocol) (ret Message, err error) {
	p.smpState, ret, err = p.smpState.receiveMessage1(p, m)
	return
}

func (m SMP2) received(p *Protocol) (ret Message, err error) {
	p.smpState, ret, err = p.smpState.receiveMessage2(p, m)
	return
}

func (m SMP3) received(p *Protocol) (ret Message, err error) {
	p.smpState, ret, err = p.smpState.receiveMessage3(p, m)
	return
}

func (m SMP4) received(p *Protocol) (ret Message, err error) {
	p.smpState, ret, err = p.smpState.receiveMessage4(p, m)
	return
}

func (m SMPAbort) received(p *Protocol) (ret Message, err error) {
	p.smpState = smpStateExpect1{}
	p.event(Abort)
	return
}
