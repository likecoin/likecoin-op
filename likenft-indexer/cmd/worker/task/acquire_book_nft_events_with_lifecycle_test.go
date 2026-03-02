package task

import "testing"

func TestAcquireBookNFTEventsTaskPayloadWithLifecyclePayload_GetAddresses(t *testing.T) {
	tests := []struct {
		name     string
		payload  AcquireBookNFTEventsTaskPayloadWithLifecyclePayload
		expected []string
	}{
		{
			name: "new field takes precedence",
			payload: AcquireBookNFTEventsTaskPayloadWithLifecyclePayload{
				ContractAddress:   "0xOLD",
				ContractAddresses: []string{"0xA", "0xB"},
			},
			expected: []string{"0xA", "0xB"},
		},
		{
			name: "falls back to old single address",
			payload: AcquireBookNFTEventsTaskPayloadWithLifecyclePayload{
				ContractAddress: "0xOLD",
			},
			expected: []string{"0xOLD"},
		},
		{
			name:     "returns nil when both empty",
			payload:  AcquireBookNFTEventsTaskPayloadWithLifecyclePayload{},
			expected: nil,
		},
		{
			name: "empty slice falls back to old address",
			payload: AcquireBookNFTEventsTaskPayloadWithLifecyclePayload{
				ContractAddress:   "0xOLD",
				ContractAddresses: []string{},
			},
			expected: []string{"0xOLD"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.payload.GetAddresses()
			if len(got) != len(tt.expected) {
				t.Fatalf("got %v, want %v", got, tt.expected)
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Fatalf("got[%d] = %q, want %q", i, got[i], tt.expected[i])
				}
			}
		})
	}
}
