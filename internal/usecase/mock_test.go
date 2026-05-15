package usecase

import (
	"dungeon-challenge/internal/domain"
	"io"
)

type MockEventWriter struct{}

func (m *MockEventWriter) WriteEvent(domain.EventType, domain.Event) {}

func (m *MockEventWriter) WriteImpossibleMove(domain.EventType, domain.Event, string) {}

func (m *MockEventWriter) WriteDeadUser(domain.EventType, domain.Event) {}

type MockReportWriter struct{}

func (m *MockReportWriter) WriteReport(map[int]*domain.User) {}

type MockEventReader struct {
	Events []domain.Event
	Index  int
}

func (m *MockEventReader) ReadEvent() (domain.Event, error) {
	if m.Index >= len(m.Events) {
		return domain.Event{}, io.EOF
	}

	event := m.Events[m.Index]
	m.Index++

	return event, nil
}
