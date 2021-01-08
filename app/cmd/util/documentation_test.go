package util

import "testing"

func TestExamples(t *testing.T) {
	tests := []struct {
		desc string
		in   string
		out  string
	}{
		{
			desc: "Empty input produces empty output",
			in:   "",
			out:  "",
		},
		{
			desc: "Text is indented with a couple of spaces",
			in:   "\tLine1\n\tLine2",
			out:  "  Line1\n  Line2",
		},
		{
			desc: "Text is stripped of leading and trailing spaces",
			in:   "\t\n\tLine1\t  \n\tLine2\t  \n\t\n\n",
			out:  "  Line1\n  Line2",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := Examples(test.in)
			if got != test.out {
				t.Errorf("expected(%d):\n%s\n=====\ngot(%d):\n%s\n", len(test.out), test.out, len(got), got)
			}
		})
	}
}

func TestLongDesc(t *testing.T) {
	tests := []struct {
		desc string
		in   string
		out  string
	}{
		{
			desc: "Empty input produces empty output",
			in:   "",
			out:  "",
		},
		{
			desc: "Single line text is preserved as is",
			in:   "Some text",
			out:  "Some text",
		},
		{
			desc: "Consecutive new lines are combined into a single paragraph",
			in:   "Line1\nLine2",
			out:  "Line1 Line2",
		},
		{
			desc: "Leading and trailing spaces are stripped (single line)",
			in:   "\t  \nThe text line  \n  \t",
			out:  "The text line",
		},
		{
			desc: "Leading and trailing spaces are stripped (multi line)",
			in:   "\t  \nLine1\nLine2  \n  \t",
			out:  "Line1 Line2",
		},
		{
			desc: "Multiple paragraphs are separated by a single empty line",
			in:   "Paragraph1\n\nParagraph2\n\n\nParagraph3",
			out:  "Paragraph1\n\nParagraph2\n\nParagraph3",
		},
		{
			desc: "Indentation is not preserved",
			in:   "\tParagraph1Line1\n\tParagraph1Line2\n\n    Paragraph2Line1\n    Paragraph2Line2",
			out:  "Paragraph1Line1 Paragraph1Line2\n\nParagraph2Line1 Paragraph2Line2",
		},
		{
			desc: "Double spaced line breaks",
			in:   "Line1  \nLine2",
			out:  "Line1\nLine2",
		},
		{
			desc: "Double spaced line breaks don't preserve indentation",
			in:   "\tLine1  \n\tLine2",
			out:  "Line1\nLine2",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := LongDesc(test.in)
			if got != test.out {
				t.Errorf("expected(%d):\n%s\n=====\ngot(%d):\n%s\n", len(test.out), test.out, len(got), got)
			}
		})
	}
}