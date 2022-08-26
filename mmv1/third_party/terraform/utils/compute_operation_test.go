package google

import (
	"strings"
	"testing"

	"google.golang.org/api/compute/v1"
)

var omitAlways = []string{
	"LocalizedMessage2",
	"Help1Link2 Description",
	"https://help1.com/link2",
	"Help2Link1 Description",
	"https://help2.com/link1",
	"Help2Link2 Description",
	"https://help2.com/link2",
}

func TestComputeOperationError_Error(t *testing.T) {
	testCases := []struct {
		name           string
		input          compute.OperationError
		expectContains []string
		expectOmits    []string
	}{
		{
			name: "MessageOnly",
			input: compute.OperationError{
				Errors: []*compute.OperationErrorErrors{
					&compute.OperationErrorErrors{
						Message: "Top-level message.",
					},
				},
			},
			expectContains: []string{
				"Top-level",
			},
			expectOmits: append(omitAlways, []string{
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			}...),
		},
		{
			name: "WithLocalizedMessageAndNoHelp",
			input: compute.OperationError{
				Errors: []*compute.OperationErrorErrors{
					&compute.OperationErrorErrors{
						Message: "Top-level message.",
						ErrorDetails: []*compute.OperationErrorErrorsErrorDetails{
							&compute.OperationErrorErrorsErrorDetails{
								LocalizedMessage: &compute.LocalizedMessage{
									Locale:  "en-US",
									Message: "LocalizedMessage1 message",
								},
							},
						},
					},
				},
			},
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
			},
			expectOmits: append(omitAlways, []string{
				"Help1Link1 Description",
				"https://help1.com/link1",
			}...),
		},
		{
			name: "WithLocalizedMessageAndHelp",
			input: compute.OperationError{
				Errors: []*compute.OperationErrorErrors{
					&compute.OperationErrorErrors{
						Message: "Top-level message.",
						ErrorDetails: []*compute.OperationErrorErrorsErrorDetails{
							&compute.OperationErrorErrorsErrorDetails{
								LocalizedMessage: &compute.LocalizedMessage{
									Locale:  "en-US",
									Message: "LocalizedMessage1 message",
								},
							},
							&compute.OperationErrorErrorsErrorDetails{
								Help: &compute.Help{
									Links: []*compute.HelpLink{
										&compute.HelpLink{
											Description: "Help1Link1 Description",
											Url:         "https://help1.com/link1",
										},
									},
								},
							},
						},
					},
				},
			},
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitAlways, []string{}...),
		},
		{
			name: "WithNoLocalizedMessageAndHelp",
			input: compute.OperationError{
				Errors: []*compute.OperationErrorErrors{
					&compute.OperationErrorErrors{
						Message: "Top-level message.",
						ErrorDetails: []*compute.OperationErrorErrorsErrorDetails{
							&compute.OperationErrorErrorsErrorDetails{
								Help: &compute.Help{
									Links: []*compute.HelpLink{
										&compute.HelpLink{
											Description: "Help1Link1 Description",
											Url:         "https://help1.com/link1",
										},
									},
								},
							},
						},
					},
				},
			},
			expectContains: []string{
				"Top-level",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitAlways, []string{
				"LocalizedMessage1",
			}...),
		},
		{
			name: "WithLocalizedMessageAndHelpWithTwoLinks",
			input: compute.OperationError{
				Errors: []*compute.OperationErrorErrors{
					&compute.OperationErrorErrors{
						Message: "Top-level message.",
						ErrorDetails: []*compute.OperationErrorErrorsErrorDetails{
							&compute.OperationErrorErrorsErrorDetails{
								LocalizedMessage: &compute.LocalizedMessage{
									Locale:  "en-US",
									Message: "LocalizedMessage1 message",
								},
							},
							&compute.OperationErrorErrorsErrorDetails{
								Help: &compute.Help{
									Links: []*compute.HelpLink{
										&compute.HelpLink{
											Description: "Help1Link1 Description",
											Url:         "https://help1.com/link1",
										},
										&compute.HelpLink{
											Description: "Help1Link2 Description",
											Url:         "https://help1.com/link2",
										},
									},
								},
							},
						},
					},
				},
			},
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitAlways, []string{}...),
		},
		// The case below should never happen because the server should just send multiple links
		// but the protobuf defition would allow it, so testing anyway.
		{
			name: "WithLocalizedMessageAndTwoHelpsWithTwoLinks",
			input: compute.OperationError{
				Errors: []*compute.OperationErrorErrors{
					&compute.OperationErrorErrors{
						Message: "Top-level message.",
						ErrorDetails: []*compute.OperationErrorErrorsErrorDetails{
							&compute.OperationErrorErrorsErrorDetails{
								LocalizedMessage: &compute.LocalizedMessage{
									Locale:  "en-US",
									Message: "LocalizedMessage1 message",
								},
							},
							&compute.OperationErrorErrorsErrorDetails{
								Help: &compute.Help{
									Links: []*compute.HelpLink{
										&compute.HelpLink{
											Description: "Help1Link1 Description",
											Url:         "https://help1.com/link1",
										},
										&compute.HelpLink{
											Description: "Help1Link2 Description",
											Url:         "https://help1.com/link2",
										},
									},
								},
							},
							&compute.OperationErrorErrorsErrorDetails{
								Help: &compute.Help{
									Links: []*compute.HelpLink{
										&compute.HelpLink{
											Description: "Help2Link1 Description",
											Url:         "https://help2.com/link1",
										},
										&compute.HelpLink{
											Description: "Help2Link2 Description",
											Url:         "https://help2.com/link2",
										},
									},
								},
							},
						},
					},
				},
			},
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitAlways, []string{}...),
		},
		// This should never happen because the server should never respond with the messages for
		// two locales at once, but should rather take the locale as input to the API and serve
		// the appropriate message for that locale. However, the protobuf defition would allow it,
		// so we'll test for it. The second message in the list would be ignored.
		{
			name: "WithTwoLocalizedMessageAndHelp",
			input: compute.OperationError{
				Errors: []*compute.OperationErrorErrors{
					&compute.OperationErrorErrors{
						Message: "Top-level message.",
						ErrorDetails: []*compute.OperationErrorErrorsErrorDetails{
							&compute.OperationErrorErrorsErrorDetails{
								LocalizedMessage: &compute.LocalizedMessage{
									Locale:  "en-US",
									Message: "LocalizedMessage1 message",
								},
							},
							&compute.OperationErrorErrorsErrorDetails{
								LocalizedMessage: &compute.LocalizedMessage{
									Locale:  "en-ES",
									Message: "LocalizedMessage2 message",
								},
							},
							&compute.OperationErrorErrorsErrorDetails{
								Help: &compute.Help{
									Links: []*compute.HelpLink{
										&compute.HelpLink{
											Description: "Help1Link1 Description",
											Url:         "https://help1.com/link1",
										},
									},
								},
							},
						},
					},
				},
			},
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitAlways, []string{}...),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ComputeOperationError(tc.input)
			str := err.Error()

			for _, contains := range tc.expectContains {
				if !strings.Contains(str, contains) {
					t.Errorf("expected\n%s\nto contain, %q, and did not", str, contains)
				}
			}

			for _, omits := range tc.expectOmits {
				if strings.Contains(str, omits) {
					t.Errorf("expected\n%s\nnot to contain, %q, and did not", str, omits)
				}
			}
		})
	}
}
