package go_textee

import (
	`reflect`
	`testing`

	gematria `github.com/andreimerlescu/go-gematria`
)

func TestTextee_ParseString(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want map[string]int32
	}{
		{
			name: "test 1",
			args: args{
				input: "All right let's move from this point on 16 March 84, let's move in time to our second location which is a specific building near where you are now. Are you ready? Just a minute. All. right, I will wait. All right, move now from this area to the front ground level of the building known as the Menara Building, to the front of, on the ground, the Menara Building.",
			},
			want: map[string]int32{
				// Populate this map with the expected substrings and their counts.
				"all right lets": 1,
				// ... more substrings and counts ...
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tt := NewTextee()
			tt.ParseString(tc.args.input)

			for substring, expectedCount := range tc.want {
				if count, ok := tt.Substrings[substring]; !ok || count.Load() != expectedCount {
					t.Errorf("Substring count for %q = %v, want %v", substring, count.Load(), expectedCount)
				}
			}
		})
	}
}

func TestTextee_CalculateGematria(t *testing.T) {
	t.Run("testing calculate gematria", func(t *testing.T) {
		got := NewTextee("manifesting three six nine").CalculateGematria()
		want := gematria.Gematria{
			Jewish:  337,
			English: 702,
			Simple:  117,
		}
		if !reflect.DeepEqual(got.Gematrias["manifesting"], want) {
			t.Errorf("CalculateGematria() = %v, want %v", got, want)
		}
	})
}
